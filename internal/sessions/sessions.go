package sessions

import (
	"github.com/google/uuid"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"time"
)

type SessionStatistics struct {
	CompletedSteps map[string]time.Time
}

type Session struct {
	Runbook   model.Runbook
	SessionId string
	Stats     SessionStatistics
}

type FakeSessionManager struct {
	sessions map[string]Session
}

func (s FakeSessionManager) CreateNewSessionForRunbook(r model.Runbook) string {
	sessionId := uuid.New().String()
	newSession := Session{
		Runbook:   r,
		SessionId: sessionId,
	}
	s.sessions[sessionId] = newSession

	return sessionId
}

func (s FakeSessionManager) GetSession(sessionId string) (Session, error) {
	res, ok := s.sessions[sessionId]

	if !ok {
		return Session{}, runbooks.CreateDataNotFoundError("session", sessionId)
	}

	return res, nil
}
