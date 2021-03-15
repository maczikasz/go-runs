package model

import (
	"time"
)

type Session struct {
	Runbook   Runbook           `json:"runbook"`
	SessionId string            `json:"session_id"`
	Stats     SessionStatistics `json:"stats"`
}

type SessionStatistics struct {
	CompletedSteps map[string]time.Time `json:"completed_steps"`
}
