package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/server/dto"
	"github.com/maczikasz/go-runs/internal/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:generate moq -out ../mocks/steps.go --skip-ensure . RunbookStepDetailsFinder RunbookStepWriter

type (
	RunbookStepDetailsFinder interface {
		FindRunbookStepDetailsById(id string) (model.RunbookStepData, *model.Markdown, error)
		ListAllSteps() ([]model.RunbookStepData, error)
	}

	RunbookStepWriter interface {
		WriteRunbookStepDetails(data model.RunbookStepData, markdown *model.Markdown, markdownLocationType string) (string, error)
		UpdateRunbookStepDetails(stepId string, data model.RunbookStepData, markdown *model.Markdown, markdownLocationType string) error
		DeleteStepDetails(id string) error
	}

	ReverseRunbookFinder interface {
		FindRunbooksByStepId(stepId string) ([]model.RunbookRef, error)
	}

	RunbookStepDetailsHandler struct {
		runbookStepDetailsFinder RunbookStepDetailsFinder
		stepWriter               RunbookStepWriter
		reverseRunbookFinder     ReverseRunbookFinder
	}
)

func NewRunbookStepDetailsHandler(runbookManager RunbookStepDetailsFinder, stepWriter RunbookStepWriter, reverseRunbookFinder ReverseRunbookFinder) *RunbookStepDetailsHandler {
	return &RunbookStepDetailsHandler{runbookStepDetailsFinder: runbookManager, stepWriter: stepWriter, reverseRunbookFinder: reverseRunbookFinder}
}

func (r RunbookStepDetailsHandler) RetrieveRunbookStepDetails(context *gin.Context) {

	stepId := context.Param("stepId")

	stepDetails, md, err := r.runbookStepDetailsFinder.FindRunbookStepDetailsById(stepId)

	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.JSON(200, dto.RunbookStepDetailDTO{
		RunbookStepData: stepDetails,
		Markdown:        md.Content,
	})
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
	}, &model.Markdown{Content: step.Markdown.Content}, step.Markdown.Type)

	if err != nil {
		log.Warnf("Failed to save  %s", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.String(http.StatusOK, stepId)
}

func (r RunbookStepDetailsHandler) ListAllSteps(context *gin.Context) {
	steps, err := r.runbookStepDetailsFinder.ListAllSteps()

	if err != nil {
		log.Warnf("Failed to save  %s", err)
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusOK, steps)
}

func (r RunbookStepDetailsHandler) DeleteRunbookStep(context *gin.Context) {
	stepId := context.Param("stepId")

	runbookRefs, err := r.reverseRunbookFinder.FindRunbooksByStepId(stepId)

	if err != nil {
		log.Warnf("Failed to check runbook with stepId %s", err)
		context.Status(http.StatusInternalServerError)
		_ = context.Error(err)
		return
	}

	if len(runbookRefs) != 0 {
		log.Warnf("The following runbooks (%s) are still using step %s", runbookRefs, stepId)
		context.String(http.StatusBadRequest, "there are runbooks using the step")
		return
	}

	err = r.stepWriter.DeleteStepDetails(stepId)

	err = util.HandleDataError(context, err)

	if err != nil {
		return
	}

	context.Status(http.StatusOK)
}

func (r RunbookStepDetailsHandler) UpdateRunbookStep(context *gin.Context) {
	var step dto.StepDTO
	err := context.BindJSON(&step)
	stepId := context.Param("stepId")

	if err != nil {
		log.Warnf("Could not parse json %s", err)
		context.Status(http.StatusBadRequest)
		return
	}

	err = r.stepWriter.UpdateRunbookStepDetails(stepId, model.RunbookStepData{
		Summary: step.Summary,
		Type:    step.Type,
	}, &model.Markdown{Content: step.Markdown.Content}, step.Markdown.Type)

	err = util.HandleDataError(context, err)

	if err != nil {
		return
	}

	err = util.HandleDataError(context, err)
	if err != nil {
		return
	}

	context.Status(http.StatusOK)

}
