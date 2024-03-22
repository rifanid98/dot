package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"dot/core"
)

type clientImpl struct {
	client *mongo.Client
	Client
}

func NewClient(client *mongo.Client) *clientImpl {
	return &clientImpl{client: client}
}

func (c *clientImpl) Database(name string, opts ...*options.DatabaseOptions) Database {
	return &databaseImpl{Database: c.client.Database(name, opts...)}
}

func (c *clientImpl) Connect(ctx context.Context) error {
	// mongo clientInstance does not use context on connect method. There is a ticket
	// with a request to deprecate this functionality and another one with
	// explanation why it could be useful in synchronous requests.
	// https://jira.mongodb.org/browse/GODRIVER-1031
	// https://jira.mongodb.org/browse/GODRIVER-979
	return c.client.Connect(ctx)
}

func (c *clientImpl) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.client.Ping(ctx, rp)
}

func (c *clientImpl) StartSession(ic *core.InternalContext) (Session, error) {
	session, err := c.client.StartSession()
	ctx, err := session.WithTransaction(ic.ToContext(), func(ctx mongo.SessionContext) (any, error) {
		return ctx, nil
	})
	return NewSession(ctx.(mongo.SessionContext), c), err
}
