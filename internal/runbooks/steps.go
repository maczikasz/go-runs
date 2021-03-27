package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/pkg/errors"
)

type (
	RunbookStepEntityFinder interface {
		FindRunbookStepData(id string) (model.RunbookStepDetailsEntity, error)
	}
	RunbookStepMarkdownResolver interface {
		ResolveRunbookStepMarkdown(entity model.RunbookStepLocation) (*model.Markdown, error)
	}

	RunbookStepDetailsFinder struct {
		RunbookStepsEntityFinder    RunbookStepEntityFinder
		RunbookStepMarkdownResolver RunbookStepMarkdownResolver
	}
)

func (receiver RunbookStepDetailsFinder) FindRunbookStepDetailsById(id string) (model.RunbookStepData, *model.Markdown, error) {
	runbookStepDetails, err := receiver.RunbookStepsEntityFinder.FindRunbookStepData(id)
	if err != nil {
		return model.RunbookStepData{}, nil, err
	}

	markdown, err := receiver.RunbookStepMarkdownResolver.ResolveRunbookStepMarkdown(runbookStepDetails.Location)
	if err != nil {
		return model.RunbookStepData{}, nil, err
	}

	return runbookStepDetails.RunbookStepData, markdown, nil
}

type (
	RunbookStepEntityWriter interface {
		WriteRunbookStepEntity(entity model.RunbookStepDetailsEntity) (string, error)
	}

	RunbookStepMarkdownWriter interface {
		WriteRunbookStepMarkdown(markdown *model.Markdown, storageType string) (string, error)
	}

	RunbookStepDetailsWriter struct {
		RunbookStepsEntityWriter  RunbookStepEntityWriter
		RunbookStepMarkdownWriter RunbookStepMarkdownWriter
	}
)

func (w RunbookStepDetailsWriter) WriteRunbookStepDetails(data model.RunbookStepData, markdown *model.Markdown, markdownLocationType string) (string, error) {
	markdownLocationId, err := w.RunbookStepMarkdownWriter.WriteRunbookStepMarkdown(markdown, markdownLocationType)

	if err != nil {
		return "", errors.Wrap(err, "failed to save markdown")
	}

	id, err := w.RunbookStepsEntityWriter.WriteRunbookStepEntity(model.RunbookStepDetailsEntity{
		RunbookStepData: data,
		Location: model.RunbookStepLocation{
			LocationType: markdownLocationType,
			Ref:          markdownLocationId,
		},
	})

	if err != nil {
		//TODO cleanup makrdown
		return "", errors.Wrap(err, "failed to save step")
	}

	return id, nil
}
