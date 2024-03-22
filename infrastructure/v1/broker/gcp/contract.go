package gcp

import (
	"context"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
)

// //go:generate mockery --name Topic --filename topic.go --output ./mocks
type Topic interface {
	Config(ctx context.Context) (pubsub.TopicConfig, error)
	Update(ctx context.Context, cfg pubsub.TopicConfigToUpdate) (pubsub.TopicConfig, error)
	ID() string
	String() string
	Delete(ctx context.Context) error
	Exists(ctx context.Context) (bool, error)
	IAM() *iam.Handle
	Subscriptions(ctx context.Context) *pubsub.SubscriptionIterator
	Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult
	Stop()
	Flush()
	ResumePublish(orderingKey string)
	EnableMessageOrdering() Topic
	DisableMessageOrdering() Topic
	GetOriginTopic() *pubsub.Topic
}

// //go:generate mockery --name Subscription --filename subscription.go --output ./mocks
type Subscription interface {
	String() string
	ID() string
	Delete(ctx context.Context) error
	Exists(ctx context.Context) (bool, error)
	Config(ctx context.Context) (pubsub.SubscriptionConfig, error)
	Update(ctx context.Context, cfg pubsub.SubscriptionConfigToUpdate) (pubsub.SubscriptionConfig, error)
	IAM() *iam.Handle
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}

//go:generate mockery --name Client --filename client.go --output ./mocks
type Client interface {
	CreateTopic(ctx context.Context, topicID string) (Topic, error)
	CreateTopicWithConfig(ctx context.Context, topicID string, tc *pubsub.TopicConfig) (Topic, error)
	Topic(id string) Topic
	TopicInProject(id, projectID string) Topic
	DetachSubscription(ctx context.Context, sub string) (*pubsub.DetachSubscriptionResult, error)
	Topics(ctx context.Context) *pubsub.TopicIterator
	Subscription(id string) Subscription
	SubscriptionInProject(id, projectID string) Subscription
	Subscriptions(ctx context.Context) *pubsub.SubscriptionIterator
	CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (Subscription, error)
}
