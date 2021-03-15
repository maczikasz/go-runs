package sessions

import (
	"github.com/google/uuid"
	"github.com/maczikasz/go-runs/internal/model"
	"time"
)

type FakeSessionManager struct {
	sessions map[string]model.Session
}

func (s FakeSessionManager) CreateNewSessionForRunbook(r model.Runbook) string {
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

func (s FakeSessionManager) GetSession(sessionId string) (model.Session, error) {
	res, ok := s.sessions[sessionId]

	if !ok {
		return model.Session{}, model.CreateDataNotFoundError("session", sessionId)
	}

	return res, nil
}
