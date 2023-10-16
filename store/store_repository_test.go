package store

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	sr     StoreRepository
	client *mongo.Client
)

func init() {
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	sr = NewStoreRepository(client)
}

func TestCreate(t *testing.T) {
	bId, err := primitive.ObjectIDFromHex("63262acf57524bad10ab8002")
	if err != nil {
		t.Errorf(err.Error())
	}
	repoStore := Store{
		Name:        "Store Name Test",
		OwnerId:     bId,
		Description: "Description Test",
		AvatarId:    "AvatarID Test",
	}
	res, err := sr.Create(&repoStore)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(res.Id)
}

func TestGetStore(t *testing.T) {
	storeId := "632f170f887fb0f1bb317ee8"
	store, err := sr.GetStore(storeId)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(store.Name)
}

func TestGetStores(t *testing.T) {
	filter := bson.D{{Key: "name", Value: "Store Name Test"}}
	stores, err := sr.GetStores(filter)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(*stores) == 0 {
		t.Errorf(fmt.Errorf("Stores is empty").Error())
	}
	for _, v := range *stores {
		fmt.Println(v.Id)
	}
}

func TestGetStoresByUserId(t *testing.T) {
	ownerId := "OwnerId Test"
	stores, err := sr.GetStoresByUser(ownerId)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(*stores) == 0 {
		t.Errorf(fmt.Errorf("Stores is empty").Error())
	}
	for _, v := range *stores {
		fmt.Println(v.Id)
	}
}

func TestUpdateStore(t *testing.T) {
	bId, err := primitive.ObjectIDFromHex("632f16faa65c3a6bf3859838")
	if err != nil {
		t.Errorf(fmt.Errorf("Error while parsing id").Error())
	}
	store := Store{
		Id:   bId,
		Name: "New name Test",
	}
	res, err := sr.UpdateStore(&store)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(res.Name)
	fmt.Println(res.Id)
}
