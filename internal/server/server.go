package server

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/infra"
	"github.com/maczikasz/go-runs/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"sync"
)

func StartHttpServer(wg *sync.WaitGroup, context *StartupContext) {
	defer wg.Done()
	r := setupRouter(context, []string{"localhost:3000"})

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type StartupContext struct {
	RunbookDetailsFinder     RunbookDetailsFinder
	SessionStore             SessionStore
	RunbookStepDetailsFinder RunbookStepDetailsFinder
	SessionFromErrorCreator  SessionFromErrorCreator
}

func setupRouter(context *StartupContext, acceptedOrigins []string) http.Handler {
	r := mux.NewRouter()

	errorHandler := incomingErrorHandler{errorHandler: context.SessionFromErrorCreator}
	sessionHandler := sessionHandler{sessionStore: context.SessionStore}
	runbookHandler := runbookHandler{runbookDetailsFinder: context.RunbookDetailsFinder}
	runbookStepDetailsHandler := runbookStepDetailsHandler{context.RunbookStepDetailsFinder}
	r.Handle("/errors", errorHandler).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/sessions/{sessionId}", sessionHandler).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/runbooks/{runbookId}", runbookHandler).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/details/{stepId}", runbookStepDetailsHandler).Methods(http.MethodGet, http.MethodOptions)

	r.Use(mux.CORSMethodMiddleware(r))
	middleware := infra.CORSPreflightOriginMiddleware{AcceptedOrigins: util.ToSet(acceptedOrigins)}
	r.Use(middleware.Middleware)

	n := negroni.Classic()
	n.UseHandler(r)
	return n
}
