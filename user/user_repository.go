package user

import (
	"context"
	"fmt"
	"time"

	"github.com/ali-ghn/Coinopay_Go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "Users"

type IUserRepository interface {
	Create(user *User) (*User, error)
	GetUser(id string) (*User, error)
	GetUsers(filter interface{}, options *options.FindOptions) (*[]User, error)
	GetUserByEmail(email string) (*User, error)
	UserExists(email string) bool
	UpdateUser(user *User) (*User, error)
}

type UserRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) UserRepository {
	return UserRepository{
		client: client,
	}
}

func (ur UserRepository) Create(user *User) (*User, error) {
	user.Id = primitive.NewObjectID()
	user.CreatedAt = time.Now().UTC().Unix()
	user.UpdatedAt = time.Now().UTC().Unix()
	res, err := ur.client.Database(shared.DATABASE_NAME).Collection(collectionName).InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	user.Id = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (ur UserRepository) GetUser(id string) (*User, error) {
	user := User{}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	err = ur.client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	user.Id, _ = primitive.ObjectIDFromHex(id)
	return &user, nil
}

func (ur UserRepository) GetUserByEmail(email string) (*User, error) {
	user := User{}
	filter := bson.D{{Key: "email", Value: email}}
	err := ur.client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur UserRepository) GetUsers(filter interface{}, options *options.FindOptions) (*[]User, error) {
	var users []User
	c, err := ur.client.Database(shared.DATABASE_NAME).Collection(collectionName).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	err = c.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (ur UserRepository) UserExists(email string) bool {
	filter := bson.D{{Key: "email", Value: email}}
	var user User
	ur.client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&user)
	fmt.Println(user.Email)
	return user.Email != ""
}

func (ur UserRepository) UpdateUser(user *User) (*User, error) {
	filter := bson.D{{Key: "email", Value: user.Email}}
	user.UpdatedAt = time.Now().UTC().Unix()
	res := ur.client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOneAndReplace(context.TODO(), filter, user)
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
