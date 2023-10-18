package user

import (
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"

	auth "github.com/ali-ghn/AlighnSolution_go/auth"
	"github.com/ali-ghn/AlighnSolution_go/cryptography"
	crypto "github.com/ali-ghn/AlighnSolution_go/cryptography"
	"github.com/ali-ghn/AlighnSolution_go/email"
	"github.com/ali-ghn/AlighnSolution_go/password"
	"github.com/ali-ghn/AlighnSolution_go/shared"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserController struct {
	ur   IUserRepository
	eh   crypto.IEncryptionHelper
	auth auth.IAuth
	es   email.IEmailSender
}

func NewUserController(ur UserRepository, eh crypto.EncryptionHelper, auth auth.Auth, es email.EmailSender) UserController {
	return UserController{
		ur:   ur,
		eh:   eh,
		auth: auth,
		es:   es,
	}
}

func (uc UserController) SignUp(c *gin.Context) {
	var user UserSignUpRequest
	err := c.Bind(&user)

	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}

	if user.Email == "" || user.Password == "" {
		c.String(http.StatusBadRequest, "Password or email is invalid")
		return
	}

	_, err = email.ValidateEmail(user.Email)
	if err != nil {
		c.String(http.StatusBadRequest, "Email is invalid")
		return
	}

	_, err = password.ValidatePassword(user.Password)

	if err != nil {
		c.String(http.StatusBadRequest, "Password is invalid")
		return
	}

	if uc.ur.UserExists(user.Email) {
		c.String(http.StatusBadRequest, "User already exists")
		return
	}

	hPassword := crypto.Hash([]byte(user.Password))

	encPassword, err := uc.eh.Encrypt(hPassword)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	repoUser := User{
		Email:             user.Email,
		Password:          fmt.Sprintf("%x", encPassword),
		IsActive:          true,
		Roles:             []string{shared.USER_ROLE},
		EmailConfirmation: false,
	}

	newUser, err := uc.ur.Create(&repoUser)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	resUser := UserSignUpResponse{
		Id:    newUser.Id.Hex(),
		Email: newUser.Email,
	}

	c.JSON(http.StatusCreated, resUser)

	emailConfirmationToken, err := uc.eh.Encrypt([]byte(resUser.Email))

	confirmationUrl := fmt.Sprintf("http://localhost:8081/EmailConfirmation/%x", emailConfirmationToken)

	message := fmt.Sprintf("لینک تایید ایمیل شما: %v", confirmationUrl)

	uc.es.Send("تایید ایمیل", message, shared.ACCOUNT_EMAIL, resUser.Email)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
}

func (uc UserController) SignIn(c *gin.Context) {
	var signIn UserSignInRequest
	err := c.Bind(&signIn)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	user, err := uc.ur.GetUserByEmail(signIn.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist")
		return
	}
	if !user.EmailConfirmation {
		c.String(http.StatusBadRequest, "Your email has not been verified")
		return
	}

	hashedPass := cryptography.Hash([]byte(signIn.Password))

	strPassword, err := hex.DecodeString(user.Password)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	decPassword, err := uc.eh.Decrypt(strPassword)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	if fmt.Sprintf("%x", hashedPass) != fmt.Sprintf("%x", decPassword) {
		c.String(http.StatusForbidden, "Incorrect credential")
		return
	}

	userClaims := auth.UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		Email: user.Email,
	}

	token, err := uc.auth.CreateToken(&userClaims)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	c.String(http.StatusOK, token)

}

func (uc UserController) ForgetPassword(c *gin.Context) {
	var fPasswordReq UserForgetPasswordRequest
	err := c.Bind(&fPasswordReq)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	user, err := uc.ur.GetUserByEmail(fPasswordReq.Email)
	if err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist")
		return
	}
	strPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	decPassword, err := uc.eh.Decrypt(strPassword)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	encryptedPassword, err := uc.eh.Encrypt(decPassword)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	confirmationUrl := fmt.Sprintf("http://localhost:8081/ForgetPassword/%x/%v", encryptedPassword, user.Email)
	message := fmt.Sprintf("لینک بازگردانی رمز عبور شما: %v", confirmationUrl)

	err = uc.es.Send("بازگردانی رمز عبور", message, shared.ACCOUNT_EMAIL, user.Email)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.String(http.StatusOK, "Your password recovery email has sent")
}

func (uc UserController) VerifyForgetPassword(c *gin.Context) {
	var verifyForgetPass UserVerifyForgetPasswordRequest
	err := c.Bind(&verifyForgetPass)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	strToken, err := hex.DecodeString(verifyForgetPass.Token)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	decToken, err := uc.eh.Decrypt(strToken)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	user, err := uc.ur.GetUserByEmail(verifyForgetPass.Email)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	strPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	decUserPassword, err := uc.eh.Decrypt(strPassword)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	if fmt.Sprintf("%x", decUserPassword) != fmt.Sprintf("%x", decToken) {
		c.String(http.StatusForbidden, "Token Is invalid")
		return
	}
	hPassword := crypto.Hash([]byte(verifyForgetPass.Password))

	encPassword, err := uc.eh.Encrypt(hPassword)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	user.Password = fmt.Sprintf("%x", encPassword)

	_, err = uc.ur.UpdateUser(user)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.String(http.StatusAccepted, "Password has changed")
}

func (uc UserController) VerifyEmail(c *gin.Context) {
	var confirmation UserEmailConfirmationRequest
	err := c.Bind(&confirmation)

	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}

	hexToken := confirmation.ConfirmationToken

	strToken, err := hex.DecodeString(hexToken)

	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	decEmail, err := uc.eh.Decrypt(strToken)

	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	user, err := uc.ur.GetUserByEmail(string(decEmail))

	if user.EmailConfirmation {
		c.String(http.StatusBadRequest, "Your email has already been verified")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	user.EmailConfirmation = true

	_, err = uc.ur.UpdateUser(user)

	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	c.String(http.StatusOK, "Email has confirmed.")
}

func (uc UserController) InsertUser(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := uc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := uc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	isAdmin := false
	for _, v := range cUser.Roles {
		if v == shared.ADMIN_ROLE {
			isAdmin = true
		}
	}
	if !isAdmin {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
	var createUser CreateUserRequest
	err = c.Bind(&createUser)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	hPassword := crypto.Hash([]byte(createUser.Password))

	encPassword, err := uc.eh.Encrypt(hPassword)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	repoUser := User{
		Email:             createUser.Email,
		Password:          fmt.Sprintf("%x", encPassword),
		IsActive:          true,
		Roles:             createUser.Roles,
		EmailConfirmation: createUser.EmailConfirmation,
		AvatarId:          createUser.AvatarId,
		FirstName:         createUser.FirstName,
		LastName:          createUser.LastName,
	}

	res, err := uc.ur.Create(&repoUser)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (uc UserController) GetUser(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := uc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := uc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	isAdmin := false
	for _, v := range cUser.Roles {
		if v == shared.ADMIN_ROLE {
			isAdmin = true
		}
	}
	if !isAdmin {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
	var getUser GetUserRequest
	err = c.BindJSON(&getUser)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	user, err := uc.ur.GetUser(getUser.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, user)
}
func (uc UserController) GetUsers(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := uc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := uc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	isAdmin := false
	for _, v := range cUser.Roles {
		if v == shared.ADMIN_ROLE {
			isAdmin = true
		}
	}
	if !isAdmin {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
	var getUsers GetUsersRequest
	err = c.BindJSON(&getUsers)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	options := options.Find().SetSort(bson.D{{"createdat", -1}}).SetSkip(getUsers.Skip).SetLimit(getUsers.Limit)
	filter := bson.D{{}}
	users, err := uc.ur.GetUsers(filter, options)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc UserController) GetUserByRole(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := uc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := uc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	isAdmin := false
	for _, v := range cUser.Roles {
		if v == shared.ADMIN_ROLE {
			isAdmin = true
		}
	}
	if !isAdmin {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
	var getUsersByRole GetUsersByRole
	err = c.BindJSON(&getUsersByRole)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	options := options.Find().SetSort(bson.D{{Key: "createdat", Value: -1}}).SetSkip(getUsersByRole.Skip).SetLimit(getUsersByRole.Limit)
	filter := bson.D{{Key: "roles", Value: getUsersByRole.Role}}
	users, err := uc.ur.GetUsers(filter, options)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc UserController) UpdateUser(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := uc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := uc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	isAdmin := false
	for _, v := range cUser.Roles {
		if v == shared.ADMIN_ROLE {
			isAdmin = true
		}
	}
	if !isAdmin {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
	var userUpdate UpdateUserRequest
	err = c.Bind(&userUpdate)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	user, err := uc.ur.GetUser(userUpdate.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	if userUpdate.Password != "" && userUpdate.Password != " " {
		hPassword := crypto.Hash([]byte(userUpdate.Password))
		encPassword, err := uc.eh.Encrypt(hPassword)
		if err != nil {
			c.String(http.StatusInternalServerError, "Something went wrong, please try again")
			return
		}
		user.Password = fmt.Sprintf("%x", encPassword)
	}
	user.Email = userUpdate.Email
	user.AvatarId = userUpdate.AvatarId
	user.EmailConfirmation = userUpdate.EmailConfirmation
	user.FirstName = userUpdate.FirstName
	user.LastName = userUpdate.LastName
	user.Roles = userUpdate.Roles
	user.IsActive = userUpdate.IsActive
	user.KycVerified = userUpdate.KycVerified
	res, err := uc.ur.UpdateUser(user)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, res)
}
