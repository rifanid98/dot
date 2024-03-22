package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type databaseImpl struct {
	*mongo.Database
}

func (db *databaseImpl) Collection(colName string) Collection {
	return &collectionImpl{db.Database.Collection(colName)}
}

func (db *databaseImpl) ListCollectionNames(ctx context.Context, filter any, opts ...*options.ListCollectionsOptions) ([]string, error) {
	return db.Database.ListCollectionNames(ctx, filter, opts...)
}

func (db *databaseImpl) GetDatabase() *mongo.Database {
	return db.Database
}
