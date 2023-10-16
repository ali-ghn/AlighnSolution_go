package store

import (
	"fmt"
	"net/http"

	"github.com/ali-ghn/Coinopay_Go/auth"
	"github.com/ali-ghn/Coinopay_Go/shared"
	"github.com/ali-ghn/Coinopay_Go/user"
	"github.com/gin-gonic/gin"
)

type StoreController struct {
	auth auth.IAuth
	sr   IStoreRepository
	ur   user.IUserRepository
}

func NewStoreController(auth auth.Auth, sr StoreRepository, ur user.UserRepository) StoreController {
	return StoreController{
		auth: auth,
		sr:   sr,
		ur:   ur,
	}
}

func (sc StoreController) CreateStore(c *gin.Context) {
	var scr StoreCreateRequest
	err := c.Bind(&scr)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	if scr.Name == "" {
		c.String(http.StatusBadRequest, "Field 'name' is necessary")
		return
	}
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	user, err := sc.ur.GetUserByEmail(claims.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist")
		return
	}

	repoStore := Store{
		Name:        scr.Name,
		Description: scr.Description,
		AvatarId:    scr.AvatarId,
		OwnerId:     user.Id,
	}

	res, err := sc.sr.Create(&repoStore)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	storeResponse := StoreCreateResponse{
		Id:          res.Id.Hex(),
		Name:        res.Name,
		Description: res.Description,
		AvatarId:    res.AvatarId,
	}

	c.JSON(http.StatusOK, storeResponse)
}

func (sc StoreController) GetStore(c *gin.Context) {
	var storeGetRequest StoreGetRequest
	err := c.BindJSON(&storeGetRequest)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	user, err := sc.ur.GetUserByEmail(claims.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist")
		return
	}

	isAdmin := false

	for _, v := range user.Roles {
		if v == shared.ADMIN_ROLE {
			isAdmin = true
		}
	}

	store, err := sc.sr.GetStore(storeGetRequest.Id)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}

	storeResponse := StoreGetResponse{
		Id:          store.Id,
		Name:        store.Name,
		Description: store.Description,
		AvatarId:    store.AvatarId,
	}

	if store.OwnerId == user.Id || isAdmin {
		c.JSON(http.StatusOK, storeResponse)
		return
	}

	c.String(http.StatusForbidden, "You don't have access to this resource")
}

func (sc StoreController) GetStores(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	user, err := sc.ur.GetUserByEmail(claims.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist")
		return
	}
	stores, err := sc.sr.GetStoresByUser(user.Id.Hex())
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, stores)
}

func (sc StoreController) GetStoresByUser(c *gin.Context) {
	var sgr StoresGetRequest
	c.BindJSON(&sgr)
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := sc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	user, err := sc.ur.GetUserByEmail(claims.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "User doesn't exist")
		return
	}

	isAdmin := false

	for _, v := range user.Roles {
		if v == shared.ADMIN_ROLE {
			isAdmin = true
		}
	}

	if isAdmin {
		stores, err := sc.sr.GetStoresByUser(sgr.UserId)

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "Something went wrong, please try again")
			return
		}

		c.JSON(http.StatusOK, stores)
		return
	}

	c.String(http.StatusForbidden, "You don't have access to this resource")
}
