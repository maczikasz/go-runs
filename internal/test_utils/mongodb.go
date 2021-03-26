package test_utils

import (
	"context"
	"fmt"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

func RunMongoDBDockerTest(testFunction func(t *testing.T, client *mongodb.MongoClient) error, t *testing.T) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}
	_ = resource.Expire(60)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		mongodburl := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodburl))
		if err != nil {
			return err
		}
		err = client.Ping(ctx, readpref.Primary())

		if err != nil {
			log.Error("Failed to connect to mongo")
			return err
		}

		mongoClient, disconnectFunction := mongodb.InitializeMongoClient(mongodburl, "local")
		defer disconnectFunction()

		return testFunction(t, mongoClient)
	}); err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err)
	}
}
