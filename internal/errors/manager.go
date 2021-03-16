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
	SessionCreator SessionCreator
	RunbookFinder  RunbookFinder
}

func (manager DefaultErrorManager) GetSessionForError(e model.Error) (string, error) {
	runbook, err := manager.RunbookFinder.FindRunbookForError(e)

	if err != nil {
		return "", errors.Wrap(err, "failed to find runbook")
	}

	sessionId := manager.SessionCreator.CreateNewSessionForRunbook(runbook)

	return sessionId, nil

}
