package errors

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/pkg/errors"
)

type (
	SessionCreator interface {
		CreateNewSession(runbook model.RunbookRef, err error) string
	}

	RunbookFinder interface {
		FindRunbookForError(e model.Error) (model.RunbookRef, error)
	}

	DefaultErrorManager struct {
		sessionCreator SessionCreator
		runbookFinder  RunbookFinder
	}
)

func NewDefaultErrorManager(sessionCreator SessionCreator, runbookFinder RunbookFinder) *DefaultErrorManager {
	return &DefaultErrorManager{sessionCreator: sessionCreator, runbookFinder: runbookFinder}
}

func (manager DefaultErrorManager) ManageErrorWitSession(e model.Error) (string, error) {
	runbook, err := manager.runbookFinder.FindRunbookForError(e)

	if err != nil {
		return "", errors.Wrap(err, "failed to find runbook")
	}

	sessionId := manager.sessionCreator.CreateNewSession(runbook, err)

	return sessionId, nil

}
