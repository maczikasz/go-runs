package runbooks

import (
	"fmt"
	"github.com/maczikasz/go-runs/internal/model"
	log "github.com/sirupsen/logrus"
)

type RuleMatcher interface {
	FindMatchingRunbook(e model.Error) (string, bool)
}

type RunbookFinder interface {
	FindRunbookById(id string) (model.RunbookRef, error)
}

type RunbookManager struct {
	ruleManager   RuleMatcher
	runbookFinder RunbookFinder
}

func NewRunbookManager(ruleManager RuleMatcher, runbookFinder RunbookFinder) *RunbookManager {
	return &RunbookManager{ruleManager: ruleManager, runbookFinder: runbookFinder}
}

func (r RunbookManager) FindRunbookForError(e model.Error) (model.RunbookRef, error) {

	runbookId, found := r.ruleManager.FindMatchingRunbook(e)
	if !found {
		log.Debugf("Could not find matching rule for error %s", e)
		return model.RunbookRef{}, model.CreateDataNotFoundError("runbook", fmt.Sprintf("error: %s", e))
	}

	return r.runbookFinder.FindRunbookById(runbookId)
}
