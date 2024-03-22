package gcp

import (
	"context"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
)

type topicImpl struct {
	*pubsub.Topic
}

func newTopic(topic *pubsub.Topic) *topicImpl {
	return &topicImpl{Topic: topic}
}

func (t *topicImpl) Config(ctx context.Context) (pubsub.TopicConfig, error) {
	return t.Topic.Config(ctx)
}

func (t *topicImpl) Update(ctx context.Context, cfg pubsub.TopicConfigToUpdate) (pubsub.TopicConfig, error) {
	return t.Topic.Update(ctx, cfg)
}

func (t *topicImpl) ID() string {
	return t.Topic.ID()
}

func (t *topicImpl) String() string {
	return t.Topic.String()
}

func (t *topicImpl) Delete(ctx context.Context) error {
	return t.Topic.Delete(ctx)
}

func (t *topicImpl) Exists(ctx context.Context) (bool, error) {
	return t.Topic.Exists(ctx)
}

func (t *topicImpl) IAM() *iam.Handle {
	return t.Topic.IAM()
}

func (t *topicImpl) Subscriptions(ctx context.Context) *pubsub.SubscriptionIterator {
	return t.Topic.Subscriptions(ctx)
}

func (t *topicImpl) Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult {
	return t.Topic.Publish(ctx, msg)
}

func (t *topicImpl) Stop() {
	t.Topic.Stop()
}

func (t *topicImpl) Flush() {
	t.Topic.Flush()
}

func (t *topicImpl) ResumePublish(orderingKey string) {
	t.Topic.ResumePublish(orderingKey)
}

func (t *topicImpl) EnableMessageOrdering() Topic {
	t.Topic.EnableMessageOrdering = true
	return t
}

func (t *topicImpl) DisableMessageOrdering() Topic {
	t.Topic.EnableMessageOrdering = false
	return t
}

func (t *topicImpl) GetOriginTopic() *pubsub.Topic {
	return t.Topic
}
