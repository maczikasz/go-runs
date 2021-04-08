package mongodb

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestWritingAndReadingSessionsFromDB(t *testing.T) {
	test_utils.RunMongoDBDockerTest(DoTestWritingAndReadingSessionsFromDB, t)
}

func DoTestWritingAndReadingSessionsFromDB(t *testing.T, client *mongodb.MongoClient) error {

	manager := SessionManager{mongoClient: client}

	Convey("Given session is created in the DB", t, func() {
		runbookId := primitive.NewObjectID().Hex()
		sessionId, _ := manager.CreateNewSession(model.RunbookRef{Id: runbookId}, model.Error{
			Name:    "test",
			Message: "test",
			Tags:    []string{},
		})

		Convey("When all sessions are listed the session is returned", func() {
			allSessions, err := manager.GetAllSessions()

			So(err, ShouldBeNil)

			Convey("Then session is present", func() {
				So(allSessions, test_utils.ShouldMatch, func(value interface{}) bool {
					if session, ok := value.(model.Session); ok {
						return session.SessionId == sessionId &&
							session.Stats.CompletedSteps != nil &&
							len(session.Stats.CompletedSteps) == 0 &&
							session.Runbook.Id == runbookId
					}
					return false
				})
			})
		})

		Convey("When a step is added", func() {
			stepId := primitive.NewObjectID().Hex()
			completeTime := time.Now()
			_ = manager.CompleteStepInSession(sessionId, stepId, completeTime)

			Convey("Then the stats contain the step", func() {
				session, _ := manager.GetSession(sessionId)

				So(session.Stats.CompletedSteps, ShouldContainKey, stepId)
				So(session.Stats.CompletedSteps[stepId].Unix(), ShouldEqual, completeTime.Unix())
			})
		})

		Convey("When two completions are added for the same step", func() {
			stepId := primitive.NewObjectID().Hex()
			completeTime := time.Now()
			_ = manager.CompleteStepInSession(sessionId, stepId, completeTime)

			completeTime2 := time.Now()
			_ = manager.CompleteStepInSession(sessionId, stepId, completeTime2)

			Convey("Then the stats contain the step", func() {
				session, _ := manager.GetSession(sessionId)

				So(session.Stats.CompletedSteps, ShouldContainKey, stepId)
				So(session.Stats.CompletedSteps[stepId].Unix(), ShouldEqual, completeTime2.Unix())
			})
		})
	})

	return nil
}
