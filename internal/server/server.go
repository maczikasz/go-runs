package server

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func StartHttpServer(wg *sync.WaitGroup, context *context.StartupContext) {
	defer wg.Done()
	r := mux.NewRouter()

	errorHandler := incomingErrorHandler{errorManager: context.ErrorManager}
	sessionHandler := sessionHandler{sessionManager: context.SessionManager}
	runbookHandler := runbookHandler{runbookManager: context.RunbookManager}
	runbookStepDetailsHandler := runbookStepDetailsHandler{runbookManager: context.RunbookManager}
	r.Handle("/errors", errorHandler).Methods("POST")
	r.Handle("/sessions/{sessionId}", sessionHandler).Methods("GET")
	r.Handle("/runbooks/{runbookId}", runbookHandler).Methods("GET")
	r.Handle("/steps/{stepId}", runbookStepDetailsHandler).Methods("GET")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
