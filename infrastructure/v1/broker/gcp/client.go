package gcp

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type clientImpl struct {
	*pubsub.Client
}

func NewClient(client *pubsub.Client) *clientImpl {
	return &clientImpl{Client: client}
}

func (c *clientImpl) CreateTopic(ctx context.Context, topicID string) (Topic, error) {
	topic, err := c.Client.CreateTopic(ctx, topicID)
	return newTopic(topic), err
}

func (c *clientImpl) CreateTopicWithConfig(ctx context.Context, topicID string, tc *pubsub.TopicConfig) (Topic, error) {
	topic, err := c.Client.CreateTopicWithConfig(ctx, topicID, tc)
	return newTopic(topic), err
}

func (c *clientImpl) Topic(id string) Topic {
	return newTopic(c.Client.Topic(id))
}

func (c *clientImpl) TopicInProject(id, projectID string) Topic {
	return newTopic(c.Client.TopicInProject(id, projectID))
}

func (c *clientImpl) DetachSubscription(ctx context.Context, sub string) (*pubsub.DetachSubscriptionResult, error) {
	return c.Client.DetachSubscription(ctx, sub)
}

func (c *clientImpl) Topics(ctx context.Context) *pubsub.TopicIterator {
	return c.Client.Topics(ctx)
}

func (c *clientImpl) Subscription(id string) Subscription {
	return c.Client.Subscription(id)
}

func (c *clientImpl) Subscriptions(ctx context.Context) *pubsub.SubscriptionIterator {
	return c.Client.Subscriptions(ctx)
}

func (c *clientImpl) SubscriptionInProject(id, projectID string) Subscription {
	return c.Client.SubscriptionInProject(id, projectID)
}

func (c *clientImpl) CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (Subscription, error) {
	return c.Client.CreateSubscription(ctx, id, cfg)
}
