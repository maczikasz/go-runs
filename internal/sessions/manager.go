package sessions

import "github.com/maczikasz/go-runs/internal/model"

func NewInMemorySessionManager() SessionManager {
	return FakeSessionManager{sessions: map[string]Session{}}
}

type SessionManager interface {
	CreateNewSessionForRunbook(r model.Runbook) string
	GetSession(sessionId string) (Session, error)
}
