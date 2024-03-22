package adapter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"

	"dot/config"
	"dot/core"
	"dot/core/v1/port/subscriber"
	"dot/infrastructure/v1/broker/gcp"
	"dot/pkg/util"
)

var log = util.NewLogger()

type pubsubImpl struct {
	gcp gcp.Client
	cfg *config.GcpPubsubConfig
	uc  subscriber.SubscriberUsecase
}

func NewGcpPubsub(
	gcp gcp.Client,
	cfg *config.GcpPubsubConfig,
	uc subscriber.SubscriberUsecase,
) *pubsubImpl {
	return &pubsubImpl{
		gcp: gcp,
		cfg: cfg,
		uc:  uc,
	}
}

func (p *pubsubImpl) Publish(ic *core.InternalContext, data []byte) *core.CustomError {
	ctxData := ic.GetData()
	topicName := ctxData["topic"].(string)

	topic, cerr := p.createTopic(ic, topicName)
	if cerr != nil {
		return cerr
	}

	msgAttr := map[string]string{}
	for k, v := range ctxData {
		msgAttr[k] = v.(string)
	}
	publish := topic.Publish(ic.ToContext(), &pubsub.Message{
		Data:        data,
		OrderingKey: p.cfg.OrderingKey,
		Attributes:  msgAttr,
	})

	id, err := publish.Get(ic.ToContext())
	if err != nil {
		log.Error(ic.ToContext(), "failed to get published message info", err)
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed get topic published message info",
		}
	}

	log.Info(ic.ToContext(), "message published", id, time.Now())

	return nil
}

func (p *pubsubImpl) Subscribe(ic *core.InternalContext) *core.CustomError {
	topics := p.gcp.Topics(ic.ToContext())
	for {
		pTopic, err := topics.Next()
		if err == iterator.Done {
			log.Info(ic.ToContext(), "no more topics")
			return nil
		}
		if err != nil {
			log.Error(ic.ToContext(), "failed to get next topic", err)

			return &core.CustomError{
				Code:    core.INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			}
		}

		log.Info(ic.ToContext(), "[CHECK] checking topic...")

		sTopic := strings.Split(pTopic.String(), "/topics/")
		topicName := sTopic[1]

		// create normal topic
		gTopic, cerr := p.createTopic(ic, topicName)
		if cerr != nil {
			return cerr
		}

		// create dead-letter topic
		dlTopicName := strings.ReplaceAll(topicName, "-dead-letter", "")
		dlTopicName = dlTopicName + "-dead-letter"
		dlTopic, cerr := p.createTopic(ic, dlTopicName)
		if cerr != nil {
			return cerr
		}

		subscriberName := "subscriber-" + topicName
		subs, cerr := p.createSubscription(ic, gTopic, dlTopic, subscriberName)
		if cerr != nil {
			return cerr
		}

		go func() {
			log.Info(ic.ToContext(), "[IDLE] subscriber "+subscriberName)

			err = subs.Receive(ic.ToContext(), func(ctx context.Context, msg *pubsub.Message) {
				log.Info(ic.ToContext(), subs.String()+" subscriber got message "+msg.ID+" "+string(msg.Data))

				// TODO: process received message
				if msg.DeliveryAttempt != nil {
					log.Info(ic.ToContext(), fmt.Sprintf("message: %s, delivery attempts: %d", msg.Data, *msg.DeliveryAttempt))
				}

				go p.processMessage(msg)
			})
			if err != nil {
				log.Error(ic.ToContext(), "failed to start subscriber", err)
			}

			log.Info(ic.ToContext(), "[END] subscriber "+subscriberName)
		}()
	}
}

func (p *pubsubImpl) createTopic(ic *core.InternalContext, name string) (gcp.Topic, *core.CustomError) {
	log.Info(ic.ToContext(), "[CHECK] topic "+name)

	topic := p.gcp.Topic(name)
	exists, err := topic.Exists(ic.ToContext())
	if exists {
		log.Info(ic.ToContext(), "[EXISTS] topic "+name)
		return topic.EnableMessageOrdering(), nil
	}

	log.Info(ic.ToContext(), "[CREATE] topic "+name)

	topic, err = p.gcp.CreateTopic(ic.ToContext(), name)
	if err != nil {
		log.Error(ic.ToContext(), "[FAILED] create topic "+name, err)
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to create topic",
		}
	}

	log.Info(ic.ToContext(), "[CREATED] topic "+name)
	return topic.EnableMessageOrdering(), nil
}

func (p *pubsubImpl) createSubscription(ic *core.InternalContext, topic gcp.Topic, dlTopic gcp.Topic, name string) (gcp.Subscription, *core.CustomError) {
	log.Info(ic.ToContext(), "[CHECK] subscriber "+name)

	subscriber := p.gcp.Subscription(name)
	exists, err := subscriber.Exists(ic.ToContext())
	if exists {
		log.Info(ic.ToContext(), "[EXISTS] subscriber "+name)
		return subscriber, nil
	}

	log.Info(ic.ToContext(), "[CREATE] subscriber "+name)

	subscriber, err = p.gcp.CreateSubscription(ic.ToContext(), name, pubsub.SubscriptionConfig{
		Topic:                 topic.GetOriginTopic(),
		AckDeadline:           20 * time.Second,
		EnableMessageOrdering: true,
		RetryPolicy: &pubsub.RetryPolicy{
			MinimumBackoff: time.Second * 10,
			MaximumBackoff: time.Minute * 3,
		},
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     dlTopic.String(),
			MaxDeliveryAttempts: 5,
		},
	})
	if err != nil {
		log.Error(ic.ToContext(), "[FAILED] create subscription "+name, err)
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to create subscription",
		}
	}

	log.Info(ic.ToContext(), "[CREATED] subscriber "+name)
	return subscriber, nil
}

func (p *pubsubImpl) processMessage(msg *pubsub.Message) {
	var ctxData = map[string]any{}
	for k, v := range msg.Attributes {
		ctxData[k] = v
	}

	ic := core.NewInternalContext(uuid.New().String()).InjectData(ctxData)

	cerr := p.uc.ProcessMessage(ic, msg)
	if cerr != nil {
		log.Error(ic.ToContext(), "failed to process message", cerr)
		//TODO: pesan yang gagal di proses harus diapakan?
		msg.Nack()
		return
	}

	msg.Ack()
	log.Info(ic.ToContext(), "message "+msg.ID+" acked")
	return
}
