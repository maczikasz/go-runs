package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/util"
)

//go:generate moq -out mocks/runbook_mocks.go -skip-ensure . RunbookStepDetailsFinder RunbookDetailsFinder

type RunbookStepDetailsFinder interface {
	FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error)
}

type runbookStepDetailsHandler struct {
	runbookManager RunbookStepDetailsFinder
}

func (r runbookStepDetailsHandler) RetrieveRunbookStepDetails(context *gin.Context) {

	stepId := context.Param("stepId")

	stepDetails, err := r.runbookManager.FindRunbookStepDetailsById(stepId)

	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.JSON(200, stepDetails)
}

type RunbookDetailsFinder interface {
	FindRunbookDetailsById(id string) (model.RunbookDetails, error)
}

type runbookHandler struct {
	runbookDetailsFinder RunbookDetailsFinder
}

func (r runbookHandler) RetrieveRunbook(context *gin.Context) {

	runbookId := context.Param("runbookId")

	runbookDetails, err := r.runbookDetailsFinder.FindRunbookDetailsById(runbookId)
	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.JSON(200, runbookDetails)
}
