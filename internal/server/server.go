package server

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/context"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"sync"
)

func StartHttpServer(wg *sync.WaitGroup, context *context.StartupContext) {
	defer wg.Done()
	r := setupRouter(context, []string{"localhost:3000"})

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRouter(context *context.StartupContext, acceptedOrigins []string) http.Handler {
	r := mux.NewRouter()

	errorHandler := incomingErrorHandler{errorManager: context.ErrorManager}
	sessionHandler := sessionHandler{sessionManager: context.SessionManager}
	runbookHandler := runbookHandler{runbookManager: context.RunbookManager}
	runbookStepDetailsHandler := runbookStepDetailsHandler{runbookManager: context.RunbookManager}
	r.Handle("/errors", errorHandler).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/sessions/{sessionId}", sessionHandler).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/runbooks/{runbookId}", runbookHandler).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/details/{stepId}", runbookStepDetailsHandler).Methods(http.MethodGet, http.MethodOptions)

	r.Use(mux.CORSMethodMiddleware(r))
	middleware := CORSPreflightOriginMiddleware{AcceptedOrigins: toSet(acceptedOrigins)}
	r.Use(middleware.Middleware)

	n := negroni.Classic()
	n.UseHandler(r)
	return n
}

func toSet(origins []string) (result map[string]bool) {
	result = make(map[string]bool)
	for _, v := range origins {
		result[v] = true
	}

	return
}
