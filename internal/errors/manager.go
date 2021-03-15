package errors

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/pkg/errors"
)

type SessionCreator interface {
	CreateNewSessionForRunbook(runbook model.Runbook) string
}

type RunbookFinder interface {
	FindRunbookForError(e model.Error) (model.Runbook, error)
}

type DefaultErrorManager struct {
	session  SessionCreator
	runbooks RunbookFinder
}

func WithManagers(manager *DefaultErrorManager, sessionManager SessionCreator, runbookManager RunbookFinder) *DefaultErrorManager {
	manager.session = sessionManager
	manager.runbooks = runbookManager
	return manager
}

func (manager DefaultErrorManager) GetSessionForError(e model.Error) (string, error) {
	runbook, err := manager.runbooks.FindRunbookForError(e)

	if err != nil {
		return "", errors.Wrap(err, "failed to find runbook")
	}

	sessionId := manager.session.CreateNewSessionForRunbook(runbook)

	return sessionId, nil

}
