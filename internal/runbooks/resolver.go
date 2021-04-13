package runbooks

import (
	"github.com/maczikasz/go-runs/internal/model"
)

type (
	E struct {
		Key   string
		Value *MarkdownHandlers
	}

	MarkdownResolver interface {
		ResolveMarkdownFromLocationString(string) (*model.Markdown, error)
	}

	MarkdownWriter interface {
		WriteMarkdown(markdown *model.Markdown) (string, error)
		DeleteMarkdown(ref string) error
	}

	MarkdownHandlers struct {
		Resolver MarkdownResolver
		Writer   MarkdownWriter
	}

	MapRunbookMarkdownResolver struct {
		resolvers map[string]*MarkdownHandlers
	}
)

type Builder []E

func BuildNewMapRunbookMarkdownResolver(resolverList Builder) *MapRunbookMarkdownResolver {
	resolvers := make(map[string]*MarkdownHandlers)

	for _, resolver := range resolverList {
		resolvers[resolver.Key] = resolver.Value
	}

	return NewMapRunbookMarkdownResolver(resolvers)
}

func NewMapRunbookMarkdownResolver(resolvers map[string]*MarkdownHandlers) *MapRunbookMarkdownResolver {
	return &MapRunbookMarkdownResolver{resolvers: resolvers}
}

func (m MapRunbookMarkdownResolver) ResolveRunbookStepMarkdown(location model.RunbookStepLocation) (*model.Markdown, error) {
	resolver, ok := m.resolvers[location.LocationType]

	if !ok {
		return nil, model.CreateDataNotFoundError("resolver_type", location.LocationType)
	}

	return resolver.Resolver.ResolveMarkdownFromLocationString(location.Ref)
}

func (m MapRunbookMarkdownResolver) DeleteRunbookMarkdown(location model.RunbookStepLocation) error {

	resolver, ok := m.resolvers[location.LocationType]

	if !ok {
		return model.CreateDataNotFoundError("resolver_type", location.LocationType)
	}

	return resolver.Writer.DeleteMarkdown(location.Ref)
}
func (m MapRunbookMarkdownResolver) WriteRunbookStepMarkdown(markdown *model.Markdown, storageType string) (string, error) {
	resolver, ok := m.resolvers[storageType]

	if !ok {
		return "", model.CreateDataNotFoundError("resolver_type", storageType)
	}

	return resolver.Writer.WriteMarkdown(markdown)
}
