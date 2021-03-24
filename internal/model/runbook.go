package model

type RunbookRef struct {
	Id string `json:"id"`
}

type RunbookDetails struct {
	Steps []RunbookStepData `json:"steps"`
}

type RunbookStepData struct {
	Id      string `json:"id,omitempty" bson:"-"`
	Summary string `json:"summary"`
	Type    string `json:"type"`
}

type RunbookStepDetails struct {
	RunbookStepData `json:"inline"`
	Markdown        string `json:"markdown"`
}
