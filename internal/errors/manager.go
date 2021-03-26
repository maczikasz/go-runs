package errors

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/pkg/errors"
)

type SessionCreator interface {
	CreateNewSession(runbook model.RunbookRef, err error) string
}

type RunbookFinder interface {
	FindRunbookForError(e model.Error) (model.RunbookRef, error)
}

type DefaultErrorManager struct {
	SessionCreator SessionCreator
	RunbookFinder  RunbookFinder
}

func (manager DefaultErrorManager) ManageErrorWitSession(e model.Error) (string, error) {
	runbook, err := manager.RunbookFinder.FindRunbookForError(e)

	if err != nil {
		return "", errors.Wrap(err, "failed to find runbook")
	}

	sessionId := manager.SessionCreator.CreateNewSession(runbook, err)

	return sessionId, nil

}
