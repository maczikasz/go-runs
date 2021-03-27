package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RunbookStepsDataManager struct {
	Client *mongodb.MongoClient
}

type runbookStepMongoEntity struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Summary  string             `bson:"summary"`
	Type     string             `bson:"type"`
	Location model.RunbookStepLocation
}

func (m RunbookStepsDataManager) WriteRunbookStepEntity(entity model.RunbookStepDetailsEntity) (string, error) {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("steps")
	defer rbCancelFunc()
	mongoEntity := runbookStepMongoEntity{
		Summary:  entity.Summary,
		Type:     entity.Type,
		Location: entity.Location,
	}

	insertOneResult, err := rbCollection.InsertOne(rbCtx, mongoEntity)

	if err != nil {
		return "", err
	}

	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m RunbookStepsDataManager) FindRunbookStepData(id string) (model.RunbookStepDetailsEntity, error) {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("steps")
	defer rbCancelFunc()
	details := runbookStepMongoEntity{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.RunbookStepDetailsEntity{}, errors.Wrap(err, "invalid ID format for mongodb")
	}
	err = rbCollection.FindOne(rbCtx, bson.M{"_id": objectID}).Decode(&details)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.RunbookStepDetailsEntity{}, model.CreateDataNotFoundError("runbook", id)
		}
	}

	return model.RunbookStepDetailsEntity{
		RunbookStepData: model.RunbookStepData{
			Id:      details.Id.Hex(),
			Summary: details.Summary,
			Type:    details.Type,
		},
		Location: details.Location,
	}, nil
}
