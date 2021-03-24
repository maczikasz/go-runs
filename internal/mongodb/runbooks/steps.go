package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RunbookStepsDataManager struct {
	Client *mongodb.MongoClient
}

func (m RunbookStepsDataManager) WriteRunbookStepEntity(entity runbooks.RunbookStepDetailsEntity) (string, error) {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("steps")
	defer rbCancelFunc()
	insertOneResult, err := rbCollection.InsertOne(rbCtx, entity)

	if err != nil {
		return "", err
	}

	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m RunbookStepsDataManager) FindRunbookStepEntityById(id string) (runbooks.RunbookStepDetailsEntity, error) {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("steps")
	defer rbCancelFunc()
	details := runbooks.RunbookStepDetailsEntity{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return runbooks.RunbookStepDetailsEntity{}, errors.Wrap(err, "invalid ID format for mongodb")
	}
	err = rbCollection.FindOne(rbCtx, bson.M{"_id": objectID}).Decode(&details)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return runbooks.RunbookStepDetailsEntity{}, model.CreateDataNotFoundError("runbook", id)
		}
	}

	return details, nil
}
