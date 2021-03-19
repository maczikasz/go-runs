package runbooks

import (
	"fmt"
	"github.com/maczikasz/go-runs/internal/model"
	log "github.com/sirupsen/logrus"
)

type RuleMatcher interface {
	FindMatchingRunbook(e model.Error) (string, bool)
}

type FakeRunbookManager struct {
	RuleManager RuleMatcher
}

func (r FakeRunbookManager) FindRunbookForError(e model.Error) (model.Runbook, error) {

	rulebookId, found := r.RuleManager.FindMatchingRunbook(e)
	if !found {
		log.Debugf("Could not find matching rule for error %s", e)
		return model.Runbook{}, model.CreateDataNotFoundError("runbook", fmt.Sprintf("error: %s", e))
	}

	return model.Runbook{Id: rulebookId}, nil
}

func (f FakeRunbookManager) FindRunbookStepDetailsById(stepId string) (model.RunbookStepDetails, error) {
	switch stepId {
	case "rbs1":
		return model.RunbookStepDetails{
			Summary:  "Test workaround",
			Type:     "Workaround",
			Markdown: "Test MD 1",
		}, nil
	case "rbs2":
		return model.RunbookStepDetails{
			Markdown: "Test MD 2",
			Summary:  "Test investigation",
			Type:     "Investigation",
		}, nil
	case "rbs3":
		return model.RunbookStepDetails{
			Markdown: "Test MD 3",
			Summary:  "Escalate to whoever",
			Type:     "Escalation",
		}, nil

	case "rbs4":
		return model.RunbookStepDetails{
			Markdown: "Test MD 4",
			Summary:  "Escalate to whoever 2",
			Type:     "Escalation",
		}, nil

	case "rbs5":
		return model.RunbookStepDetails{
			Markdown: "Test MD 5",
			Summary:  "Test investigation 2",
			Type:     "Investigation",
		}, nil

	case "rbs6":
		return model.RunbookStepDetails{
			Markdown: "Test MD 6",
			Summary:  "Test workaround 2",
			Type:     "Workaround",
		}, nil
	default:
		return model.RunbookStepDetails{}, model.CreateDataNotFoundError("step_details", stepId)

	}

}

func (r FakeRunbookManager) FindRunbookDetailsById(id string) (model.RunbookDetails, error) {

	switch id {
	case "test-1":
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
			},
		}

		return model.RunbookDetails{
			Steps: steps,
		}, nil
	case "test-2":
		steps := []model.RunbookStepSummary{
			{
				Id:      "rbs2",
				Summary: "Test step 2",
				Type:    "Investigation",
			}, {
				Id:      "rbs5",
				Summary: "Escalate to whoever",
				Type:    "Escalation",
			},
		}

		return model.RunbookDetails{
			Steps: steps,
		}, nil
	case "test-3":
		steps := []model.RunbookStepSummary{
			{
				Id:      "rbs4",
				Summary: "Test step 1",
				Type:    "Workaround",
			},
			{
				Id:      "rbs6",
				Summary: "Test step 2",
				Type:    "Investigation",
			},
		}

		return model.RunbookDetails{
			Steps: steps,
		}, nil

	}

	return model.RunbookDetails{}, model.CreateDataNotFoundError("steps", id)
}
