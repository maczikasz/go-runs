package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/runbooks"
	"github.com/maczikasz/go-runs/internal/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:generate moq -out mocks/runbook_mocks.go -skip-ensure . RunbookStepDetailsFinder RunbookDetailsFinder

type RunbookStepDetailsFinder interface {
	FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error)
}

type RunbookStepWriter interface {
	WriteRunbookStepDetails(data model.RunbookStepData, markdown runbooks.Markdown, markdownLocationType string) (string, error)
}

type runbookStepDetailsHandler struct {
	runbookManager RunbookStepDetailsFinder
	stepWriter     RunbookStepWriter
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

type stepDTO struct {
	Summary  string       `json:"summary"`
	Markdown MarkdownInfo `json:"markdown"`
	Type     string       `json:"type"`
}

type MarkdownInfo struct {
	Content string
	Type    string
}

func (r runbookStepDetailsHandler) CreateNewStep(context *gin.Context) {
	var step stepDTO
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

type RunbookDetailsFinder interface {
	FindRunbookDetailsById(id string) (model.RunbookDetails, error)
}

type RunbookDetailsWriter interface {
	CreateRunbookFromStepIds(steps []string) (string, error)
}

type runbookHandler struct {
	runbookDetailsFinder     RunbookDetailsFinder
	runbookDetailsWriter     RunbookDetailsWriter
	runbookStepDetailsFinder RunbookStepDetailsFinder
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

type runbookDTO struct {
	Steps []string
}

func (r runbookHandler) CreateNewRunbook(context *gin.Context) {
	var runbook runbookDTO
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
