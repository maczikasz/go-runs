package model

type RuleEntity struct {
	RuleType    string `json:"rule_type"`
	MatcherType string `json:"matcher_type"`
	RuleContent string `json:"rule_content"`
	RunbookId   string `json:"runbook_id"`
	ID          string `json:"id" bson:"_id,omitempty"`
}
