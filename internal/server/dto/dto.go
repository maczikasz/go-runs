package dto

type RuleCreateDTO struct {
	RuleType    string `json:"rule_type"`
	MatcherType string `json:"matcher_type"`
	RuleContent string `json:"rule_content"`
	RunbookId   string `json:"runbook_id"`
}

type StepDTO struct {
	Summary  string       `json:"summary"`
	Markdown MarkdownInfo `json:"markdown"`
	Type     string       `json:"type"`
}

type MarkdownInfo struct {
	Content string
	Type    string
}

type RunbookDTO struct {
	Steps []string
}

