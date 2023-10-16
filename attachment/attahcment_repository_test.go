package attachment

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	ar     AttachmentRepository
)

func init() {
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	ar = NewAttachmentRepository(client)
}

func TestUpload(t *testing.T) {
	testData := []byte("This is a test")
	res, err := ar.Upload(testData, "testFile.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !res {
		t.Errorf(fmt.Errorf("something went wrong").Error())
	}
}

func TestDownload(t *testing.T) {
	filename := "testFile.txt"
	res, err := ar.Download(filename)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(res))
}
