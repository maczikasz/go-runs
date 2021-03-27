package model

type (
	RunbookRef struct {
		Id string `json:"id"`
	}

	RunbookDetails struct {
		Steps []RunbookStepData `json:"steps"`
	}

	RunbookStepData struct {
		Id      string `json:"id,omitempty"`
		Summary string `json:"summary"`
		Type    string `json:"type"`
	}

	RunbookStepDetailsEntity struct {
		RunbookStepData
		Location RunbookStepLocation
	}

	RunbookStepLocation struct {
		LocationType string
		Ref          string
	}

	Markdown struct {
		Content string
	}
)
