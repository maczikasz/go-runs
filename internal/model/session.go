package model

import (
	"time"
)

type Session struct {
	Runbook         RunbookRef        `json:"runbook"`
	SessionId       string            `json:"session_id"`
	Stats           SessionStatistics `json:"stats"`
	TriggeringError Error             `json:"error"`
}

type SessionStatistics struct {
	CompletedSteps map[string]time.Time `json:"completed_steps"`
}
