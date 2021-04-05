package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RunbookDataManager struct {
	Client *mongodb.MongoClient
}

type RunbookEntity struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty"`
	Name  string               `bson:"name"`
	Steps []primitive.ObjectID `bson:"steps"`
}

func (m RunbookDataManager) FindRunbooksByStepId(stepId string) ([]model.RunbookRef, error) {
	collection, cancelFunc, ctx := m.Client.Collection("runbooks")
	defer cancelFunc()
	hex, err := primitive.ObjectIDFromHex(stepId)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ID format for mongodb")
	}

	cursor, err := collection.Find(ctx, bson.M{"steps": hex})

	if err != nil {
		return nil, err
	}

	var entities []RunbookEntity

	err = cursor.All(ctx, &entities)
	if err != nil {
		return nil, err
	}

	var result []model.RunbookRef

	for _, entity := range entities {
		result = append(result, model.RunbookRef{Id: entity.ID.Hex()})
	}

	return result, nil
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

func (m RunbookDataManager) DeleteRunbook(id string) error {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("runbooks")
	defer rbCancelFunc()

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "invalid ID format for mongodb")
	}

	result, err := rbCollection.DeleteOne(rbCtx, bson.M{"_id": hex})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return model.CreateDataNotFoundError("runbook", id)
	}

	return nil
}

func (m RunbookDataManager) UpdateRunbook(id string, name string, steps []string) error {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("runbooks")
	defer rbCancelFunc()

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "invalid ID format for mongodb")
	}

	var objectIdSteps []primitive.ObjectID

	for _, stepId := range steps {
		//We can safely ignore error as steps should be verified before
		hex, _ := primitive.ObjectIDFromHex(stepId)
		objectIdSteps = append(objectIdSteps, hex)
	}

	runbookEntity := RunbookEntity{
		ID:    hex,
		Name:  name,
		Steps: objectIdSteps,
	}
	result, err := rbCollection.ReplaceOne(rbCtx, bson.M{"_id": hex}, &runbookEntity)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return model.CreateDataNotFoundError("runbook", id)
	}

	return nil

}

func (m RunbookDataManager) CreateRunbookFromDetails(steps []string, name string) (string, error) {

	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("runbooks")
	defer rbCancelFunc()

	var objectIdSteps []primitive.ObjectID

	for _, stepId := range steps {
		//We can safely ignore error as steps should be verified before
		hex, _ := primitive.ObjectIDFromHex(stepId)
		objectIdSteps = append(objectIdSteps, hex)
	}

	runbookEntity := RunbookEntity{
		Name:  name,
		Steps: objectIdSteps,
	}

	insertOneResult, err := rbCollection.InsertOne(rbCtx, &runbookEntity)

	if err != nil {
		return "", errors.Wrap(err, "failed to insert runbook")
	}

	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil

}

func (m RunbookDataManager) ListAllRunbooks() ([]model.RunbookSummary, error) {
	rbCollection, rbCancelFunc, rbCtx := m.Client.Collection("runbooks")
	defer rbCancelFunc()
	var entities []RunbookEntity

	cursor, err := rbCollection.Find(rbCtx, bson.M{})

	if err != nil {
		return nil, err
	}

	err = cursor.All(rbCtx, &entities)

	if err != nil {
		return nil, err
	}

	var result []model.RunbookSummary

	for _, entity := range entities {
		var steps []string
		for _, step := range entity.Steps {
			steps = append(steps, step.Hex())
		}
		result = append(result, model.RunbookSummary{
			Id:    entity.ID.Hex(),
			Name:  entity.Name,
			Steps: steps,
		})
	}

	if result == nil {
		return []model.RunbookSummary{}, nil
	}

	return result, nil
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

	var stepSummaries []runbookStepMongoEntity

	err = cursor.All(stepCtx, &stepSummaries)

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
		data := model.RunbookStepData{
			Id:      summary.Id.Hex(),
			Summary: summary.Summary,
			Type:    summary.Type,
		}

		steps = append(steps, data)
	}
	return model.RunbookDetails{Name: runbook.Name, Steps: steps}, nil
}
