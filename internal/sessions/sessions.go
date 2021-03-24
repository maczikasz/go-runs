package sessions

import (
	"github.com/google/uuid"
	"github.com/maczikasz/go-runs/internal/model"
	"time"
)

type FakeSessionManager struct {
	sessions map[string]model.Session
}

func (s FakeSessionManager) CreateNewSession(r model.RunbookRef, err error) string {
	sessionId := uuid.New().String()
	newSession := model.Session{
		Runbook:   r,
		SessionId: sessionId,
		Stats: model.SessionStatistics{
			CompletedSteps: map[string]time.Time{},
		},
	}
	s.sessions[sessionId] = newSession
	newSession.Stats.CompletedSteps["rbs1"] = time.Now()

	return sessionId
}

func NewInMemorySessionManager() *FakeSessionManager {
	return &FakeSessionManager{sessions: map[string]model.Session{}}
}

func (s FakeSessionManager) GetSession(sessionId string) (model.Session, error) {
	res, ok := s.sessions[sessionId]

	if !ok {
		return model.Session{}, model.CreateDataNotFoundError("session", sessionId)
	}

	return res, nil
}
