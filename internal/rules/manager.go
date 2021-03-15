package rules

import (
	"fmt"
	"github.com/maczikasz/go-runs/internal/model"
)

type RuleManager interface {
	FindMatch(error2 model.Error) (string, bool)
}

type RuleRunbookPair struct {
	RunbookId string
	Rule      Rule
}

func (m RuleRunbookPair) String() string {
	return fmt.Sprintf("rule: %s, runbook: %s", m.Rule, m.RunbookId)
}
