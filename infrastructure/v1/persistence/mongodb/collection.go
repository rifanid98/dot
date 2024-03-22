package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionImpl struct {
	*mongo.Collection
}

func (coll *collectionImpl) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) SingleResult {
	return coll.Collection.FindOne(ctx, filter, opts...)
}

func (coll *collectionImpl) InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return coll.Collection.InsertOne(ctx, document, opts...)
}

func (coll *collectionImpl) UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return coll.Collection.UpdateOne(ctx, filter, update, opts...)
}

func (coll *collectionImpl) UpdateMany(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return coll.Collection.UpdateMany(ctx, filter, update, opts...)
}

func (coll *collectionImpl) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return coll.Collection.DeleteOne(ctx, filter, opts...)
}

func (coll *collectionImpl) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (Cursor, error) {
	return coll.Collection.Find(ctx, filter, opts...)
}

func (coll *collectionImpl) InsertMany(ctx context.Context, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return coll.Collection.InsertMany(ctx, documents, opts...)
}

func (coll *collectionImpl) Aggregate(ctx context.Context, pipeline any, opts ...*options.AggregateOptions) (Cursor, error) {
	return coll.Collection.Aggregate(ctx, pipeline, opts...)
}

func (coll *collectionImpl) DeleteMany(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return coll.Collection.DeleteMany(ctx, filter, opts...)
}

func (coll *collectionImpl) CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	return coll.Collection.CountDocuments(ctx, filter, opts...)
}
