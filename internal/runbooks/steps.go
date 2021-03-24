package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/pkg/errors"
)

type RunbookStepLocation struct {
	LocationType string
	Ref          string
}

type Markdown struct {
	Content string
}

type RunbookStepDetailsEntity struct {
	model.RunbookStepData `bson:"inline"`
	Id                    string `json:"id" bson:"_id,omitempty"`
	Location              RunbookStepLocation
}

type RunbookStepEntityFinder interface {
	FindRunbookStepEntityById(id string) (RunbookStepDetailsEntity, error)
}

type RunbookStepEntityWriter interface {
	WriteRunbookStepEntity(entity RunbookStepDetailsEntity) (string, error)
}

type RunbookStepMarkdownWriter interface {
	WriteRunbookStepMarkdown(markdown *Markdown, storageType string) (string, error)
}
type RunbookStepMarkdownResolver interface {
	ResolveRunbookStepMarkdown(entity RunbookStepLocation) (*Markdown, error)
}

type RunbookStepDetailsFinder struct {
	RunbookStepsEntityFinder    RunbookStepEntityFinder
	RunbookStepMarkdownResolver RunbookStepMarkdownResolver
}

func (receiver RunbookStepDetailsFinder) FindRunbookStepDetailsById(id string) (model.RunbookStepDetails, error) {
	runbookStepDetails, err := receiver.RunbookStepsEntityFinder.FindRunbookStepEntityById(id)
	if err != nil {
		return model.RunbookStepDetails{}, err
	}

	markdown, err := receiver.RunbookStepMarkdownResolver.ResolveRunbookStepMarkdown(runbookStepDetails.Location)
	if err != nil {
		return model.RunbookStepDetails{}, err
	}

	return model.RunbookStepDetails{
		RunbookStepData: model.RunbookStepData{
			Id:      runbookStepDetails.Id,
			Summary: runbookStepDetails.Summary,
			Type:    runbookStepDetails.Type,
		},
		Markdown: markdown.Content,
	}, nil
}

type RunbookStepDetailsWriter struct {
	RunbookStepsEntityWriter  RunbookStepEntityWriter
	RunbookStepMarkdownWriter RunbookStepMarkdownWriter
}

func (w RunbookStepDetailsWriter) WriteRunbookStepDetails(data model.RunbookStepData, markdown Markdown, markdownLocationType string) (string, error) {
	markdownLocationId, err := w.RunbookStepMarkdownWriter.WriteRunbookStepMarkdown(&markdown, markdownLocationType)

	if err != nil {
		return "", errors.Wrap(err, "failed to save markdown")
	}

	id, err := w.RunbookStepsEntityWriter.WriteRunbookStepEntity(RunbookStepDetailsEntity{
		RunbookStepData: data,
		Location: RunbookStepLocation{
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
