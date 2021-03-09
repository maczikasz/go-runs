package main

import (
	"github.com/maczikasz/go-runs/internal/context"
	"github.com/maczikasz/go-runs/internal/errors"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server"
	"github.com/maczikasz/go-runs/internal/sessions"
	log "github.com/sirupsen/logrus"
	"sync"
)

func main() {
	log.SetLevel(log.TraceLevel)
	wg := sync.WaitGroup{}
	wg.Add(1)

	sessionManager := sessions.NewInMemorySessionManager()
	runbookManager := runbooks.FakeRunbookManager{}
	errorManager := errors.DefaultErrorManager{}
	startupContext := context.StartupContext{
		ErrorManager:   errors.WithManagers(&errorManager, sessionManager, runbookManager),
		SessionManager: sessionManager,
		RunbookManager: runbookManager,
	}

	startupContext.Validate()

	server.StartHttpServer(&wg, &startupContext)
	wg.Wait()

}
