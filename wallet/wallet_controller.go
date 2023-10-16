package wallet

import (
	"net/http"

	"github.com/ali-ghn/Coinopay_Go/auth"
	"github.com/ali-ghn/Coinopay_Go/shared"
	"github.com/ali-ghn/Coinopay_Go/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type WalletController struct {
	wh   IWalletHelper
	wr   IWalletRepository
	ur   user.IUserRepository
	auth auth.IAuth
}

func NewWalletController(wh WalletHelper, wr WalletRepository, ur user.UserRepository, auth auth.Auth) WalletController {
	return WalletController{
		wh:   wh,
		wr:   wr,
		ur:   ur,
		auth: auth,
	}
}

func (wc WalletController) GetWalletOverview(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := wc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := wc.ur.GetUserByEmail(claims.Email)
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
	var walletGetRequest GetWalletRequest
	c.BindJSON(&walletGetRequest)
	walletOverview, err := wc.wh.GetWalletOverview(walletGetRequest.CryptoCode)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, walletOverview)
}

func (wc WalletController) GetWalletAddress(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := wc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := wc.ur.GetUserByEmail(claims.Email)
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
	var walletGetRequest GetWalletRequest
	c.BindJSON(&walletGetRequest)
	walletAddress, err := wc.wh.GetWalletAddress(walletGetRequest.CryptoCode)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, walletAddress)
}

func (wc WalletController) GetWallets(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := wc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := wc.ur.GetUserByEmail(claims.Email)
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
	filter := bson.D{{}}
	wallets, err := wc.wr.GetWallets(filter, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	var getWalletsResponse []GetWalletsResponse
	for _, v := range *wallets {
		walletOverview, err := wc.wh.GetWalletOverview(v.Symbol)
		if err != nil {
			c.String(http.StatusInternalServerError, "Something went wrong, please try again")
			return
		}
		getWalletsResponse = append(getWalletsResponse, GetWalletsResponse{
			Id:         v.Id.Hex(),
			CryptoCode: v.Symbol,
			Name:       v.Name,
			Symbol:     v.Symbol,
			Balance:    walletOverview.balance,
		})
	}
	c.JSON(http.StatusOK, getWalletsResponse)
}

func (wc WalletController) CreateTransaction(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := wc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := wc.ur.GetUserByEmail(claims.Email)
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
	var createTransaction CreateTransactionRequest
	err = c.Bind(&createTransaction)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	res, err := wc.wh.CreateWalletTransaction(createTransaction.CryptoCode, &CreateWalletTransactionRequest{
		Destinations: []TransactionDestination{
			{
				Destination:        createTransaction.Destination,
				Amount:             createTransaction.Amount,
				SubtractFromAmount: createTransaction.SubtractFee,
			},
		},
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, res)
}
