package model

type Runbook struct {
	RunbookId string
}

type RunbookDetails struct {
	Steps []RunbookStepSummary
}

type RunbookStepSummary struct {
	Id       string
	Summary  string
	Type     string
}

type RunbookStepDetails struct {
	Markdown string
}