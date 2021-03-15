package server

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/runbooks"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type runbookStepDetailsHandler struct {
	runbookManager runbooks.RunbookManager
}

func (r runbookStepDetailsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {


	if request.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(request)
	stepId := vars["stepId"]

	stepDetails, err := r.runbookManager.FindRunbookStepDetailsById(stepId)

	err = HandleDataError(writer, request, err)
	if err != nil {
		return
	}

	err = WriteJsonResponse(writer, stepDetails)
	if err != nil {
		log.Warnf("failed to write to response %s", err)
	}
}

type runbookHandler struct {
	runbookManager runbooks.RunbookManager
}

func (r runbookHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {


	if request.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(request)
	runbookId := vars["runbookId"]

	runbookDetails, err := r.runbookManager.FindRunbookDetailsById(runbookId)
	err = HandleDataError(writer, request, err)
	if err != nil {
		return
	}

	err = WriteJsonResponse(writer, runbookDetails)

	if err != nil {
		log.Warnf("failed to write to response %s", err)
	}
}
