package settings

import (
	"net/http"

	"github.com/ali-ghn/Coinopay_Go/auth"
	"github.com/ali-ghn/Coinopay_Go/shared"
	"github.com/ali-ghn/Coinopay_Go/user"
	"github.com/gin-gonic/gin"
)

type SettingsController struct {
	str  ISettingsRepository
	auth auth.IAuth
	ur   user.IUserRepository
}

func NewSettingsController(str SettingsRepository, auth auth.Auth, ur user.UserRepository) SettingsController {
	return SettingsController{
		str:  str,
		auth: auth,
		ur:   ur,
	}
}

func (sc SettingsController) CreateSiteSettings(c *gin.Context) {
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
	cUser, err := sc.ur.GetUserByEmail(claims.Email)
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
	var siteSettings SiteSettings
	err = c.Bind(&siteSettings)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	res, err := sc.str.CreateSettings(&siteSettings)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, res)
}

func (sc SettingsController) GetSiteSettings(c *gin.Context) {
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
	cUser, err := sc.ur.GetUserByEmail(claims.Email)
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
	settings, err := sc.str.GetLatestSettings()
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, settings)
}
