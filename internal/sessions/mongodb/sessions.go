package mongodb

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	SessionManager struct {
		mongoClient *mongodb.MongoClient
	}

	SessionEntity struct {
		ID              primitive.ObjectID `bson:"_id,omitempty"`
		RunbookRef      primitive.ObjectID `bson:"runbook"`
		TriggeringError model.Error        `bson:"error"`
	}

	StatEntry struct {
		ID              primitive.ObjectID  `bson:"_id,omitempty"`
		SessionID       primitive.ObjectID  `bson:"session_id"`
		CompletedStepId primitive.ObjectID  `bson:"step_id"`
		CompletionTime  primitive.Timestamp `bson:"time"`
	}
)

func NewMongoDBSessionManager(mongoClient *mongodb.MongoClient) *SessionManager {
	return &SessionManager{mongoClient: mongoClient}
}

func (s SessionManager) CreateNewSession(runbook model.RunbookRef, error model.Error) (string, error) {
	collection, cancelFunc, ctx := s.mongoClient.Collection("sessions")
	defer cancelFunc()

	runbookObjectId, err := primitive.ObjectIDFromHex(runbook.Id)

	if err != nil {
		return "", err
	}

	insertOneResult, err := collection.InsertOne(ctx, SessionEntity{
		RunbookRef:      runbookObjectId,
		TriggeringError: error,
	})

	if err != nil {
		return "", err
	}

	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s SessionManager) GetSession(sessionId string) (model.Session, error) {
	collection, cancelFunc, ctx := s.mongoClient.Collection("sessions")
	defer cancelFunc()

	sessionObjectId, err := primitive.ObjectIDFromHex(sessionId)

	if err != nil {
		return model.Session{}, errors.Wrap(err, "Could not parse id")
	}

	var sessionEntity SessionEntity

	err = collection.FindOne(ctx, bson.M{"_id": sessionObjectId}).Decode(&sessionEntity)

	if err != nil {
		return model.Session{}, err
	}

	aggrResult, err := s.getStatsForSession(sessionObjectId)

	if err != nil {
		return model.Session{}, err
	}

	return s.mergeEntities(sessionEntity, aggrResult), nil
}

func (s SessionManager) mergeEntities(sessionEntity SessionEntity, aggrResult map[string]time.Time) model.Session {
	return model.Session{
		Runbook:         model.RunbookRef{Id: sessionEntity.RunbookRef.Hex()},
		SessionId:       sessionEntity.ID.Hex(),
		Stats:           model.SessionStatistics{CompletedSteps: aggrResult},
		TriggeringError: sessionEntity.TriggeringError,
	}
}

func (s SessionManager) getStatsForSession(sessionObjectId primitive.ObjectID) (map[string]time.Time, error) {
	statsCollection, statsCancel, statsCtx := s.mongoClient.Collection("stats")
	defer statsCancel()
	cursor, err := statsCollection.Aggregate(
		statsCtx,
		[]bson.M{
			{"$match": bson.M{"session_id": sessionObjectId}},
			{"$group": bson.M{"_id": "$step_id", "lastCompletionTime": bson.M{"$last": "$time"}}},
		},
	)

	if err != nil {
		return nil, err
	}

	var aggrResult []bson.M

	err = cursor.All(statsCtx, &aggrResult)

	if err != nil {
		return nil, err
	}

	if aggrResult == nil {
		return map[string]time.Time{}, nil
	}
	result := make(map[string]time.Time)

	for _, r := range aggrResult {
		result[r["_id"].(primitive.ObjectID).Hex()] = time.Unix(int64(r["lastCompletionTime"].(primitive.Timestamp).T), 0)
	}

	return result, nil
}

func (s SessionManager) GetAllSessions() ([]model.Session, error) {
	collection, cancelFunc, ctx := s.mongoClient.Collection("sessions")
	defer cancelFunc()

	var sessionEntities []SessionEntity

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &sessionEntities)

	if err != nil {
		return nil, err
	}

	var result []model.Session

	for _, session := range sessionEntities {
		stats, err := s.getStatsForSession(session.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, model.Session{
			Runbook:         model.RunbookRef{Id: session.RunbookRef.Hex()},
			SessionId:       session.ID.Hex(),
			Stats:           model.SessionStatistics{CompletedSteps: stats},
			TriggeringError: session.TriggeringError,
		})
	}

	if result == nil {
		return []model.Session{}, nil
	}

	return result, nil
}

func (s SessionManager) CompleteStepInSession(sessionId string, stepId string, now time.Time) error {
	statsCollection, statsCancel, statsCtx := s.mongoClient.Collection("stats")
	defer statsCancel()

	sessionObjectId, err := primitive.ObjectIDFromHex(sessionId)

	if err != nil {
		return errors.Wrap(err, "Could not parse id")
	}

	stepObjectId, err := primitive.ObjectIDFromHex(stepId)

	if err != nil {
		return errors.Wrap(err, "Could not parse id")
	}

	_, err = statsCollection.InsertOne(statsCtx, StatEntry{
		SessionID:       sessionObjectId,
		CompletedStepId: stepObjectId,
		CompletionTime:  primitive.Timestamp{T: uint32(now.Unix())},
	})

	if err != nil {
		return err
	}

	return nil
}
