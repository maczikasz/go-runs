package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
)

type RunbookManager interface {
	FindRunbookForError(e model.Error) (model.Runbook, error)
	FindRunbookDetailsById(id string) (model.RunbookDetails, error)
	FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error)
}
