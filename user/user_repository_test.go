package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/ali-ghn/AlighnSolution_go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ur     UserRepository
	client *mongo.Client
)

func init() {
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	ur = NewUserRepository(client)
}

func TestCreate(t *testing.T) {
	repoUser := User{
		AvatarId: "AvatarId Test",
		Roles: []string{
			shared.USER_ROLE,
		},
		Email:             "Test@email.com",
		Password:          "Password Test",
		IsActive:          true,
		EmailConfirmation: true,
		TotpSecret:        "TOTP Secret Test",
	}
	res, err := ur.Create(&repoUser)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(res.Id)
}

func TestGet(t *testing.T) {
	userId := "632f1cd7745a0ae2490aeb33"
	user, err := ur.GetUser(userId)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(user.Id)
}

func TestGetByEmail(t *testing.T) {
	email := "Test@email.com"
	user, err := ur.GetUserByEmail(email)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(user.Id)
}

func TestGetUsers(t *testing.T) {
	filter := bson.D{{"email", "Test@email.com"}}
	users, err := ur.GetUsers(filter, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(*users) == 0 {
		t.Errorf(fmt.Errorf("Users is empty").Error())
	}
	for _, v := range *users {
		fmt.Println(v.Email)
	}
}

func TestUserExists(t *testing.T) {
	userEmail := "alighndev@protonmail.com"
	res := ur.UserExists(userEmail)
	if !res {
		t.Errorf(fmt.Errorf("expected user existence but returned false").Error())
	}
	fakeUserMail := "Fake@email.com"
	res = ur.UserExists(fakeUserMail)
	if res {
		t.Errorf(fmt.Errorf("expected user to not exist but returned true").Error())
	}
}

func TestUpdateUser(t *testing.T) {
	bId, err := primitive.ObjectIDFromHex("63262acf57524bad10ab8002")
	if err != nil {
		t.Errorf(err.Error())
	}
	user := User{
		Id:    bId,
		Email: "alighndev@protonmail.com",
		Roles: []string{
			shared.ADMIN_ROLE,
		},
		AvatarId: "New Avatar Id",
	}
	res, err := ur.UpdateUser(&user)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res.AvatarId != "New Avatar Id" {
		t.Errorf(fmt.Errorf("expected New Avatar Id as the new AvatarId").Error())
	}
	fmt.Println(res.AvatarId)
	fmt.Println(res.Roles)
}

func TestGetSpecificRoles(t *testing.T) {
	filter := bson.D{{Key: "roles", Value: shared.SUPPORT_ROLE}}
	users, err := ur.GetUsers(filter, nil)
	if err != nil {
		t.Error(err)
	}
	for _, v := range *users {
		fmt.Println(v.Id)
	}
}
