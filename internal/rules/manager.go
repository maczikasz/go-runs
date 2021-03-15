package rules

import (
	"fmt"
)

type RuleRunbookPair struct {
	RunbookId string
	Rule      Rule
}

func (m RuleRunbookPair) String() string {
	return fmt.Sprintf("rule: %s, runbook: %s", m.Rule, m.RunbookId)
}
