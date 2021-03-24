package gridfs

import (
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io/ioutil"
	"time"
)

type Client struct {
	Bucket *gridfs.Bucket
}

func (c Client) ReadFileContentFromLocation(s string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return "", errors.Wrap(err, "invalid ID format for mongodb")
	}
	stream, err := c.Bucket.OpenDownloadStream(objectID)
	if err != nil {
		if err == gridfs.ErrFileNotFound {
			return "", model.CreateDataNotFoundError("gridfs_markdown", s)
		}
		return "", err
	}

	err = stream.SetReadDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(stream)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c Client) WriteFileToLocation(content string, id primitive.ObjectID) error {
	uploadStream, err := c.Bucket.OpenUploadStreamWithID(id, "/markdown/"+id.Hex())

	if err != nil {
		return err
	}

	err = uploadStream.SetWriteDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		return err
	}

	data := []byte(content)
	write, err := uploadStream.Write(data)

	if err != nil {
		return err
	}

	if write != len(data) {
		log.Warnf("Did not write full file to gridfs for file with id %s; only written %d bytes", id.Hex(), write)
	}
	err = uploadStream.Close()

	if err != nil {
		log.Warnf("Could not close upload stream file with id %s", id.Hex())
		return err
	}

	return nil

}
