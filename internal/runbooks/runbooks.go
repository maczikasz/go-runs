package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
)

type FakeRunbookManager struct {
}

func (r FakeRunbookManager) FindRunbookForError(e model.Error) (model.Runbook, error) {

	return model.Runbook{RunbookId: "1"}, nil
}

func (f FakeRunbookManager) FindRunbookStepDetailsById(stepId string) (model.RunbookStepDetails, error) {
	switch stepId {
	case "rbs1":
		return model.RunbookStepDetails{Markdown: "Test MD 1"}, nil
	case "rbs2":
		return model.RunbookStepDetails{Markdown: "Test MD 2"}, nil
	case "rbs3":
		return model.RunbookStepDetails{Markdown: "Test MD 3"}, nil
	default:
		return model.RunbookStepDetails{}, CreateDataNotFoundError("step_details", stepId)

	}

}

func (r FakeRunbookManager) FindRunbookDetailsById(id string) (model.RunbookDetails, error) {

	if id == "1" {
		steps := []model.RunbookStepSummary{
			{
				Id:      "rbs1",
				Summary: "Test step 1",
				Type:    "Workaround",
			},
			{
				Id:      "rbs2",
				Summary: "Test step 2",
				Type:    "Investigation",
			}, {
				Id:      "rbs3",
				Summary: "Escalate to whoever",
				Type:    "Escalation",
			},
		}

		return model.RunbookDetails{
			Steps: steps,
		}, nil
	}

	return model.RunbookDetails{}, CreateDataNotFoundError("steps", id)
}
