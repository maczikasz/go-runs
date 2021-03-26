package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/server/dto"
	"github.com/maczikasz/go-runs/internal/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:generate moq -out ../mocks/steps.go --skip-ensure . RunbookStepDetailsFinder RunbookStepWriter

type (
	RunbookStepDetailsFinder interface {
		FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error)
	}

	RunbookStepWriter interface {
		WriteRunbookStepDetails(data model.RunbookStepData, markdown runbooks.Markdown, markdownLocationType string) (string, error)
	}

	RunbookStepDetailsHandler struct {
		runbookManager RunbookStepDetailsFinder
		stepWriter     RunbookStepWriter
	}
)

func NewRunbookStepDetailsHandler(runbookManager RunbookStepDetailsFinder, stepWriter RunbookStepWriter) *RunbookStepDetailsHandler {
	return &RunbookStepDetailsHandler{runbookManager: runbookManager, stepWriter: stepWriter}
}

func (r RunbookStepDetailsHandler) RetrieveRunbookStepDetails(context *gin.Context) {

	stepId := context.Param("stepId")

	stepDetails, err := r.runbookManager.FindRunbookStepDetailsById(stepId)

	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.JSON(200, stepDetails)
}

func (r RunbookStepDetailsHandler) CreateNewStep(context *gin.Context) {
	var step dto.StepDTO
	err := context.BindJSON(&step)

	if err != nil {
		log.Warnf("Could not parse json %s", err)
		context.Status(http.StatusBadRequest)
		return
	}

	stepId, err := r.stepWriter.WriteRunbookStepDetails(model.RunbookStepData{
		Summary: step.Summary,
		Type:    step.Type,
	}, runbooks.Markdown{Content: step.Markdown.Content}, step.Markdown.Type)

	if err != nil {
		log.Warnf("Failed to save  %s", err)
		context.Status(http.StatusBadRequest)
		return
	}

	context.String(http.StatusOK, stepId)
}
