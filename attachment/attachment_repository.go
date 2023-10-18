package attachment

import (
	"bytes"

	"github.com/ali-ghn/AlighnSolution_go/shared"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type IAttachmentRepository interface {
	Upload(data []byte, filename string) (bool, error)
	Download(id string) ([]byte, error)
}

type AttachmentRepository struct {
	Client *mongo.Client
}

func NewAttachmentRepository(client *mongo.Client) AttachmentRepository {
	return AttachmentRepository{
		Client: client,
	}
}

func (a AttachmentRepository) Upload(data []byte, filename string) (bool, error) {
	bucket, err := gridfs.NewBucket(
		a.Client.Database(shared.DATABASE_NAME),
	)
	if err != nil {
		return false, err
	}
	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return false, err
	}
	defer uploadStream.Close()
	_, err = uploadStream.Write(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a AttachmentRepository) Download(id string) ([]byte, error) {
	bucket, _ := gridfs.NewBucket(
		a.Client.Database(shared.DATABASE_NAME),
	)
	var buf bytes.Buffer
	_, err := bucket.DownloadToStreamByName(id, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
