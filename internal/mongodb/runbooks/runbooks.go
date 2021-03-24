package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RunbookDataManager struct {
	Client *mongodb.MongoClient
}

func (m RunbookDataManager) FindRunbookById(id string) (model.RunbookRef, error) {
	collection, cancelFunc, ctx := m.Client.Collection("runbooks")
	defer cancelFunc()
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.RunbookRef{}, errors.Wrap(err, "invalid ID format for mongodb")
	}
	err = collection.FindOne(ctx, bson.M{"_id": hex}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.RunbookRef{}, model.CreateDataNotFoundError("runbook", id)
		}
	}

	return model.RunbookRef{Id: id}, nil
}

type RunbookEntity struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty"`
	Steps []primitive.ObjectID `bson:"steps"`
}

func (m RunbookDataManager) CreateRunbookFromStepIds(steps []string) (string, error) {

	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("runbooks")
	defer rbCancelFunc()

	var objectIdSteps []primitive.ObjectID

	for _, stepId := range steps {
		//We can safely ignore error as steps should be verified before
		hex, _ := primitive.ObjectIDFromHex(stepId)
		objectIdSteps = append(objectIdSteps, hex)
	}

	runbookEntity := RunbookEntity{
		Steps: objectIdSteps,
	}

	insertOneResult, err := rbCollection.InsertOne(rbCtx, &runbookEntity)

	if err != nil {
		return "", errors.Wrap(err, "failed to insert runbook")
	}

	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil

}

func (m RunbookDataManager) FindRunbookDetailsById(id string) (model.RunbookDetails, error) {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("runbooks")
	defer rbCancelFunc()
	runbook := RunbookEntity{}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.RunbookDetails{}, errors.Wrap(err, "invalid ID format for mongodb")
	}

	err = rbCollection.FindOne(rbCtx, bson.M{"_id": objectId}).Decode(&runbook)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.RunbookDetails{}, model.CreateDataNotFoundError("runbook", id)
		}
	}
	stepCollection, stepCancelFunc, stepCtx := m.Client.Collection("steps")
	defer stepCancelFunc()

	cursor, err := stepCollection.Find(stepCtx, bson.M{"_id": bson.M{"$in": runbook.Steps}})

	if err != nil {
		return model.RunbookDetails{}, err
	}

	var stepSummaries []runbooks.RunbookStepDetailsEntity
	//var b []bson.M

	err = cursor.All(stepCtx, &stepSummaries)
	//panic(spew.Sdump(b))

	if err != nil {
		return model.RunbookDetails{}, err
	}

	if len(stepSummaries) == 0 {
		return model.RunbookDetails{}, model.CreateDataNotFoundError("steps", id)
	}

	if len(stepSummaries) != len(runbook.Steps) {
		log.Warnf("Did not manage to return all runbook steps for runbook %s", id)
	}

	var steps []model.RunbookStepData

	for _, summary := range stepSummaries {
		data := summary.RunbookStepData
		//TODO really?
		data.Id = summary.Id

		steps = append(steps, data)
	}
	return model.RunbookDetails{Steps: steps}, nil
}
