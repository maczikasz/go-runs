package runbooks

import (
	"fmt"
	"github.com/maczikasz/go-runs/internal/model"
)

type DataNotFoundError struct {
	dataType string
	id       string
}

func CreateDataNotFoundError(dataType string, id string) error {
	return DataNotFoundError{
		dataType: dataType,
		id:       id,
	}
}

func (d DataNotFoundError) Error() string {
	return fmt.Sprintf("could not found %s with id %s", d.dataType, d.id)
}

type RunbookManager interface {
	FindRunbookForError(e model.Error) (model.Runbook, error)
	FindRunbookDetailsById(id string) (model.RunbookDetails, error)
	FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error)
}
