package model

type Runbook struct {
	Id string `json:"id"`
}

type RunbookDetails struct {
	Steps []RunbookStepSummary `json:"steps"`
}

type RunbookStepSummary struct {
	Id      string `json:"id"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type RunbookStepDetails struct {
	Summary  string `json:"summary"`
	Type     string `json:"type"`
	Markdown string `json:"markdown"`
}
