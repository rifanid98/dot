package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"dot/core"
)

//go:generate mockery --name Client --filename client.go --output ./mocks
type Client interface {
	Database(name string, opts ...*options.DatabaseOptions) Database
	Connect(ctx context.Context) error
	StartSession(ic *core.InternalContext) (Session, error)
	Ping(ctx context.Context, rp *readpref.ReadPref) error
}

//go:generate mockery --name Database --filename database.go --output ./mocks
type Database interface {
	Collection(name string) Collection
	ListCollectionNames(ctx context.Context, filter any, opts ...*options.ListCollectionsOptions) ([]string, error)
	GetDatabase() *mongo.Database
}

//go:generate mockery --name Collection --filename collection.go --output ./mocks
type Collection interface {
	FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) SingleResult
	InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter any, opts ...*options.FindOptions) (Cursor, error)
	Aggregate(ctx context.Context, pipeline any, opts ...*options.AggregateOptions) (Cursor, error)
	DeleteMany(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)
}

//go:generate mockery --name Session --filename session.go --output ./mocks
type Session interface {
	context.Context
	Client() Client
	StartTransaction(...*options.TransactionOptions) error
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
	WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (any, error), opts ...*options.TransactionOptions) (any, error)
	EndSession(context.Context)
	ClusterTime() bson.Raw
	OperationTime() *primitive.Timestamp
	AdvanceClusterTime(bson.Raw) error
	AdvanceOperationTime(*primitive.Timestamp) error
}

//go:generate mockery --name SingleResult --filename single_result.go --output ./mocks
type SingleResult interface {
	Err() error
	Decode(v any) error
	DecodeBytes() (bson.Raw, error)
}

//go:generate mockery --name Cursor --filename cursor.go --output ./mocks
type Cursor interface {
	Decode(val any) error
	All(ctx context.Context, results any) error
	Err() error
	ID() int64
}
