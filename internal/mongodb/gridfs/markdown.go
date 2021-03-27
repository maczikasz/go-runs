package gridfs

import (
	"github.com/maczikasz/go-runs/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarkdownResolver struct {
	Client *Client
}

type MarkdownWriter struct {
	Client *Client
}

func (m MarkdownWriter) WriteMarkdown(markdown *model.Markdown) (string, error) {
	id := primitive.NewObjectID()
	err := m.Client.WriteFileToLocation(markdown.Content, id)

	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (m MarkdownResolver) ResolveMarkdownFromLocationString(s string) (*model.Markdown, error) {

	markdownContent, err := m.Client.ReadFileContentFromLocation(s)
	if err != nil {
		return nil, err
	}

	return &model.Markdown{Content: markdownContent}, err
}
