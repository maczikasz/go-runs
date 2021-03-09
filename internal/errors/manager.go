package errors

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/sessions"
	"github.com/pkg/errors"
)

type ErrorManager interface {
	HandleIncomingError(e model.Error) (string, error)
}

type DefaultErrorManager struct {
	session  sessions.SessionManager
	runbooks runbooks.RunbookManager
}

func WithManagers(manager *DefaultErrorManager, sessionManager sessions.SessionManager, runbookManager runbooks.RunbookManager) *DefaultErrorManager {
	manager.session = sessionManager
	manager.runbooks = runbookManager
	return manager
}

func (manager DefaultErrorManager) HandleIncomingError(e model.Error) (string, error) {
	runbook, err := manager.runbooks.FindRunbookForError(e)

	if err != nil {
		return "", errors.Wrap(err, "failed to find runbook")
	}

	sessionId := manager.session.CreateNewSessionForRunbook(runbook)

	return sessionId, nil

}
