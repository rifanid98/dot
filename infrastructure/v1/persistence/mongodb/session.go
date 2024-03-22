package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sessionImpl struct {
	session mongo.SessionContext
	client  Client
}

func NewSession(session mongo.SessionContext, client Client) *sessionImpl {
	return &sessionImpl{session: session, client: client}
}

func (s *sessionImpl) Deadline() (deadline time.Time, ok bool) {
	return s.session.Deadline()
}

func (s *sessionImpl) Done() <-chan struct{} {
	return s.session.Done()
}

func (s *sessionImpl) Err() error {
	return s.session.Err()
}

func (s *sessionImpl) Value(key any) any {
	return s.session.Value(key)
}

func (s *sessionImpl) Client() Client {
	return s.client
}

func (s *sessionImpl) StartTransaction(options ...*options.TransactionOptions) error {
	return s.session.StartTransaction(options...)
}

func (s *sessionImpl) AbortTransaction(ctx context.Context) error {
	defer s.EndSession(ctx)
	return s.session.AbortTransaction(ctx)
}

func (s *sessionImpl) CommitTransaction(ctx context.Context) error {
	defer s.EndSession(ctx)
	return s.session.CommitTransaction(ctx)
}

func (s *sessionImpl) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (any, error), opts ...*options.TransactionOptions) (any, error) {
	return s.session.WithTransaction(ctx, fn)
}

func (s *sessionImpl) EndSession(ctx context.Context) {
	s.session.EndSession(ctx)
}

func (s *sessionImpl) ClusterTime() bson.Raw {
	return s.session.ClusterTime()
}

func (s *sessionImpl) OperationTime() *primitive.Timestamp {
	return s.session.OperationTime()
}

func (s *sessionImpl) AdvanceClusterTime(raw bson.Raw) error {
	return s.session.AdvanceClusterTime(raw)
}

func (s *sessionImpl) AdvanceOperationTime(timestamp *primitive.Timestamp) error {
	return s.session.AdvanceOperationTime(timestamp)
}
