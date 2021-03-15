package server

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:generate moq -out mocks/runbook_mocks.go -skip-ensure . RunbookStepDetailsFinder RunbookDetailsFinder

type RunbookStepDetailsFinder interface {
	FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error)
}

type runbookStepDetailsHandler struct {
	runbookManager RunbookStepDetailsFinder
}

func (r runbookStepDetailsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(request)
	stepId := vars["stepId"]

	stepDetails, err := r.runbookManager.FindRunbookStepDetailsById(stepId)

	err = util.HandleDataError(writer, request, err)
	if err != nil {
		return
	}

	err = util.WriteJsonResponse(writer, stepDetails)
	if err != nil {
		log.Warnf("failed to write to response %s", err)
	}
}


type RunbookDetailsFinder interface {
	FindRunbookDetailsById(id string) (model.RunbookDetails, error)
}

type runbookHandler struct {
	runbookDetailsFinder RunbookDetailsFinder
}

func (r runbookHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodOptions {
		return
	}

	vars := mux.Vars(request)
	runbookId := vars["runbookId"]

	runbookDetails, err := r.runbookDetailsFinder.FindRunbookDetailsById(runbookId)
	err = util.HandleDataError(writer, request, err)
	if err != nil {
		return
	}

	err = util.WriteJsonResponse(writer, runbookDetails)

	if err != nil {
		log.Warnf("failed to write to response %s", err)
	}
}
