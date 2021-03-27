package gridfs

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/maczikasz/go-runs/internal/mongodb"
	"github.com/maczikasz/go-runs/internal/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func DoTestMarkdownIsWrittenThenResolved(t *testing.T, client *mongodb.MongoClient) error {

	Convey("Given mongodb is connected to gridfs", t, func() {

		fsClient, _ := client.NewGridFSClient()
		resolver := MarkdownResolver{Client: &Client{Bucket: fsClient}}
		writer := MarkdownWriter{Client: &Client{Bucket: fsClient}}

		const testContent = "TEST_CONTENT"

		Convey("Given data is written to GridFS", func() {
			fileId, wErr := writer.WriteMarkdown(&model.Markdown{Content: testContent})

			So(wErr, ShouldBeNil)

			So(fileId, ShouldNotBeEmpty)

			Convey("When data is read back from gridfs", func() {
				resolvedMarkdown, rErr := resolver.ResolveMarkdownFromLocationString(fileId)
				So(rErr, ShouldBeNil)

				Convey("Then content matches", func() {
					So(resolvedMarkdown, ShouldNotBeNil)
					So(resolvedMarkdown.Content, ShouldEqual, testContent)
				})
			})
		})

		Convey("When non existent fileId is read back from GridFS", func() {
			_, resolverError := resolver.ResolveMarkdownFromLocationString(primitive.NewObjectID().Hex())

			Convey("Then DataNotFoundError is returned", func() {
				So(resolverError, ShouldHaveSameTypeAs, model.DataNotFoundError{})
			})
		})
	})

	return nil
}

func TestMarkdownIsWrittenThenResolved(t *testing.T) {
	test_utils.RunMongoDBDockerTest(DoTestMarkdownIsWrittenThenResolved, t)
}
