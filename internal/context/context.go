package context

import (
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/rules"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/sessions"
)

type StartupContext struct {
	ErrorManager   errors.ErrorManager
	SessionManager sessions.SessionManager
	RunbookManager runbooks.RunbookManager
	RuleManager    rules.RuleManager
}

func (c StartupContext) Validate() {
	if c.ErrorManager == nil || c.SessionManager == nil || c.RunbookManager == nil {
		panic("startup context validation failed")
	}
}
