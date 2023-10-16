package store

import "go.mongodb.org/mongo-driver/bson/primitive"

type Store struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string
	OwnerId     primitive.ObjectID
	Token       string
	Description string
	AvatarId    string
}

type StoreCreateRequest struct {
	Name        string
	Description string
	AvatarId    string
}

type StoreCreateResponse struct {
	Id          string
	Name        string
	Description string
	AvatarId    string
}

type StoreGetRequest struct {
	Id string `json:"Id"`
}

type StoreGetResponse struct {
	Id          primitive.ObjectID
	Name        string
	Description string
	AvatarId    string
}

type StoresGetRequest struct {
	UserId string
}
