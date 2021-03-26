package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/server/dto"
	"github.com/maczikasz/go-runs/internal/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:generate moq -out ../mocks/runbooks.go --skip-ensure . RunbookDetailsWriter RunbookDetailsFinder

type (
	RunbookDetailsFinder interface {
		FindRunbookDetailsById(id string) (model.RunbookDetails, error)
	}

	RunbookDetailsWriter interface {
		CreateRunbookFromStepIds(steps []string) (string, error)
	}

	RunbookHandler struct {
		runbookDetailsFinder     RunbookDetailsFinder
		runbookDetailsWriter     RunbookDetailsWriter
		runbookStepDetailsFinder RunbookStepDetailsFinder
	}
)

func NewRunbookHandler(runbookDetailsFinder RunbookDetailsFinder, runbookDetailsWriter RunbookDetailsWriter, runbookStepDetailsFinder RunbookStepDetailsFinder) *RunbookHandler {
	return &RunbookHandler{runbookDetailsFinder: runbookDetailsFinder, runbookDetailsWriter: runbookDetailsWriter, runbookStepDetailsFinder: runbookStepDetailsFinder}
}

func (r RunbookHandler) RetrieveRunbook(context *gin.Context) {

	runbookId := context.Param("runbookId")

	runbookDetails, err := r.runbookDetailsFinder.FindRunbookDetailsById(runbookId)
	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.JSON(200, runbookDetails)
}

func (r RunbookHandler) CreateNewRunbook(context *gin.Context) {
	var runbook dto.RunbookDTO
	err := context.BindJSON(&runbook)

	if err != nil {
		log.Warnf("Could not parse json %s", err)
		context.Status(http.StatusBadRequest)
		return
	}

	for _, stepId := range runbook.Steps {
		_, err := r.runbookStepDetailsFinder.FindRunbookStepDetailsById(stepId)
		if err != nil {
			log.Warnf("Could not find step with id %s", stepId)
			context.Status(http.StatusBadRequest)
			return
		}
	}

	runbookId, err := r.runbookDetailsWriter.CreateRunbookFromStepIds(runbook.Steps)

	if err != nil {
		log.Warnf("Failed to insert runbook %s", err)
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	context.String(http.StatusOK, runbookId)
}
