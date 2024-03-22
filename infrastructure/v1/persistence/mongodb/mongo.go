package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"dot/config"
	"dot/pkg/util"
)

var log = util.NewLogger()

func New(config *config.MongoDbConfig) (Database, Client, error) {
	var (
		auth options.Credential
		ctx  = context.Background()
	)

	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	if config.Production {
		uri = fmt.Sprintf(
			"mongodb+srv://%s:%s@%s/%s",
			config.Username,
			config.Password,
			config.Host,
			config.Name,
		)
	}

	auth.Username = config.Username
	auth.Password = config.Password

	opts := options.Client().
		ApplyURI(uri).
		SetAuth(auth)

	if config.Debug {
		opts.SetMonitor(&event.CommandMonitor{
			Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
				log.Info(ctx, fmt.Sprintf("%v", evt.Command))
			},
		})
	}

	c, err := mongo.NewClient(opts)
	if err != nil {
		log.Error(ctx, "failed to initiate mongo db client : %v", err)
		panic(err)
	}
	err = c.Connect(ctx)
	if err != nil {
		log.Error(ctx, "failed to connect mongodb client : %v", err)
		panic(err)
	}
	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error(ctx, "failed to ping mongodb client : %v", err)
		panic(err)
	}

	client := NewClient(c)
	return client.Database(config.Name), client, nil
}

func Index(collection *mongo.Collection, options *options.IndexOptions, data ...[]string) error {
	ctx := context.Background()

	for _, values := range data {
		document := bsonx.Doc{}
		for _, value := range values {
			document = append(document, bsonx.Elem{Key: value, Value: bsonx.Int32(1)})
		}
		_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    document,
			Options: options,
		})
		if err != nil {
			log.Error(ctx, "failed to create mongodb index : %v", err)
			panic(err)
		}
	}

	return nil
}
