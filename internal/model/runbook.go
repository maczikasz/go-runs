package model

type (
	RunbookRef struct {
		Id string `json:"id"`
	}

	RunbookSummary struct {
		Id    string   `json:"id"`
		Name  string   `json:"name"`
		Steps []string `json:"steps"`
	}

	RunbookDetails struct {
		Name  string            `json:"name"`
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
