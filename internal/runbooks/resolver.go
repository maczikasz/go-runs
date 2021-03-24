package runbooks

import "github.com/maczikasz/go-runs/internal/model"

type MarkdownResolver interface {
	ResolveMarkdownFromLocationString(string) (*Markdown, error)
}

type MarkdownWriter interface {
	WriteMarkdown(markdown *Markdown) (string, error)
}

type MarkdownHandlers struct {
	Resolver MarkdownResolver
	Writer   MarkdownWriter
}

type MapRunbookMarkdownResolver struct {
	Resolvers map[string]MarkdownHandlers
}

func (m MapRunbookMarkdownResolver) ResolveRunbookStepMarkdown(location RunbookStepLocation) (*Markdown, error) {
	resolver, ok := m.Resolvers[location.LocationType]

	if !ok {
		return nil, model.CreateDataNotFoundError("resolver_type", location.LocationType)
	}

	return resolver.Resolver.ResolveMarkdownFromLocationString(location.Ref)
}

func (m MapRunbookMarkdownResolver) WriteRunbookStepMarkdown(markdown *Markdown, storageType string) (string, error) {
	resolver, ok := m.Resolvers[storageType]

	if !ok {
		return "", model.CreateDataNotFoundError("resolver_type", storageType)
	}

	return resolver.Writer.WriteMarkdown(markdown)
}
