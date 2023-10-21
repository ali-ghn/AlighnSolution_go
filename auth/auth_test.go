package auth

import (
	"github.com/dgrijalva/jwt-go"
	"strings"
	"testing"
	"time"
)

func TestNewAuth(t *testing.T) {
	secret := "w9RyLLMQ6xckghjfRNs9JDbDd7pabwGJ"
	auth := NewAuth([]byte(secret))
	if auth.AuthKey == nil {
		t.Errorf("secret was not set in the initializer, expected %v got %v", secret, auth.AuthKey)
	}
}

func TestCreateToken(t *testing.T) {
	secret := "w9RyLLMQ6xckghjfRNs9JDbDd7pabwGJ"
	auth := NewAuth([]byte(secret))
	userClaim := UserClaim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "AlighnSolution",
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Id:        "alighn@duck.com or userid",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "AlighnSolution",
			Subject:   "User main token",
		},
		Email: "alighn@duck.com",
	}
	token, err := auth.CreateToken(&userClaim)
	if err != nil {
		t.Error(err.Error())
	}
	if strings.EqualFold(token, "") {
		t.Errorf("token was not created successfully")
	}
}

func TestParseToken(t *testing.T) {
	secret := "w9RyLLMQ6xckghjfRNs9JDbDd7pabwGJ"
	auth := NewAuth([]byte(secret))
	userClaim := UserClaim{
		StandardClaims: jwt.StandardClaims{
			Audience:  "AlighnSolution",
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Id:        "alighn@duck.com or userid",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "AlighnSolution",
			Subject:   "User main token",
		},
		Email: "alighn@duck.com",
	}
	token, err := auth.CreateToken(&userClaim)
	if err != nil {
		t.Error(err.Error())
	}
	parsedToken, err := auth.ParseToken(token)
	if err != nil {
		t.Error(err.Error())
	}
	result, err := parsedToken.Compare(&userClaim)
	if err != nil {
		t.Error(err.Error())
	}
	if !result {
		t.Errorf("parsed token do not match, expected %v as the primary key got %v", userClaim.Email, parsedToken.Email)
	}
}
