package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id"`
	AvatarId          string
	Roles             []string
	IsActive          bool
	Email             string
	Password          string
	EmailConfirmation bool
	FirstName         string
	LastName          string
	KycVerified       bool
	TotpSecret        string
	CreatedAt         int64
	UpdatedAt         int64
}

type UserSignUpRequest struct {
	Email    string
	Password string
}

type UserSignUpResponse struct {
	Id    string
	Email string
}

type UserSignInRequest struct {
	Email    string
	Password string
}

type UserEmailConfirmationRequest struct {
	ConfirmationToken string
}

type UserForgetPasswordRequest struct {
	Email string
}

type UserVerifyForgetPasswordRequest struct {
	Token    string
	Email    string
	Password string
}

type CreateUserRequest struct {
	Email             string
	Password          string
	AvatarId          string
	Roles             []string
	EmailConfirmation bool
	FirstName         string
	LastName          string
	KycVerified       bool
}

type GetUserRequest struct {
	Id string
}

type GetUsersRequest struct {
	Skip  int64
	Limit int64
}

type GetUsersByRole struct {
	Skip  int64
	Limit int64
	Role  string
}

type UpdateUserRequest struct {
	Id                string
	AvatarId          string
	Roles             []string
	IsActive          bool
	Email             string
	Password          string
	EmailConfirmation bool
	FirstName         string
	LastName          string
	KycVerified       bool
}
