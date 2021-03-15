package sessions

import "github.com/maczikasz/go-runs/internal/model"

func NewInMemorySessionManager() SessionManager {
	return FakeSessionManager{sessions: map[string]model.Session{}}
}

type SessionManager interface {
	CreateNewSessionForRunbook(r model.Runbook) string
	GetSession(sessionId string) (model.Session, error)
}
