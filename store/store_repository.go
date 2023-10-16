package store

import (
	"context"

	"github.com/ali-ghn/Coinopay_Go/shared"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "Stores"
)

type IStoreRepository interface {
	Create(store *Store) (*Store, error)
	GetStore(id string) (*Store, error)
	GetStores(filter interface{}) (*[]Store, error)
	GetStoresByUser(userId string) (*[]Store, error)
	GetStoreByToken(token string) (*Store, error)
	UpdateStore(store *Store) (*Store, error)
}

type StoreRepository struct {
	client *mongo.Client
}

func NewStoreRepository(client *mongo.Client) StoreRepository {
	return StoreRepository{
		client: client,
	}
}

func (sr StoreRepository) Create(store *Store) (*Store, error) {
	store.Id = primitive.NewObjectID()
	store.Token = uuid.New().String()
	res, err := sr.client.Database(shared.DATABASE_NAME).Collection(collectionName).InsertOne(context.TODO(), store)
	if err != nil {
		return nil, err
	}
	store.Id = res.InsertedID.(primitive.ObjectID)
	return store, nil
}

func (sr StoreRepository) GetStore(id string) (*Store, error) {
	store := Store{}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	err = sr.client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (sr StoreRepository) GetStoreByToken(token string) (*Store, error) {
	store := Store{}
	filter := bson.D{{Key: "token", Value: token}}
	err := sr.client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (sr StoreRepository) GetStores(filter interface{}) (*[]Store, error) {
	stores := []Store{}
	cur, err := sr.client.Database(shared.DATABASE_NAME).Collection(collectionName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &stores)
	if err != nil {
		return nil, err
	}
	return &stores, nil
}

func (sr StoreRepository) GetStoresByUser(userId string) (*[]Store, error) {
	stores := []Store{}
	bId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "ownerid", Value: bId}}
	cur, err := sr.client.Database(shared.DATABASE_NAME).Collection(collectionName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &stores)
	if err != nil {
		return nil, err
	}
	return &stores, nil
}

func (sr StoreRepository) UpdateStore(store *Store) (*Store, error) {
	filter := bson.D{{Key: "_id", Value: store.Id}}
	err := sr.client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOneAndReplace(context.TODO(), filter, store).Decode(&store)
	if err != nil {
		return nil, err
	}
	return store, err
}
