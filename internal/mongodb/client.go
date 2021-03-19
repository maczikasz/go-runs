package mongodb

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	database string
	mongo    *mongo.Client
}

func (c MongoClient) Collection(s string) (*mongo.Collection, context.CancelFunc, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	database := c.mongo.Database(c.database)
	collection := database.Collection(s)

	return collection, cancel, ctx
}

type DisconnectFunction func()

func InitializeMongoClient(mongoUrl string, database string) (*MongoClient, DisconnectFunction) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

	if err != nil {
		log.Errorf("Failed to connect to mongo %s", err)
		panic("Failed to connect to mongodb")
	}

	return &MongoClient{
			database: database,
			mongo:    client,
		}, func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}
}
