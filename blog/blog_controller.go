package blog

import (
	"math"
	"net/http"

	"github.com/ali-ghn/Coinopay_Go/auth"
	"github.com/ali-ghn/Coinopay_Go/shared"
	"github.com/ali-ghn/Coinopay_Go/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogController struct {
	br   BlogRepository
	auth auth.IAuth
	ur   user.IUserRepository
}

func NewBlogController(br BlogRepository, auth auth.Auth, ur user.UserRepository) BlogController {
	return BlogController{
		br:   br,
		auth: auth,
		ur:   ur,
	}
}

func (bc BlogController) CreateBlogPost(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := bc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := bc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	isBlogger := false
	for _, v := range cUser.Roles {
		if v == shared.BLOGGER_ROLE || v == shared.ADMIN_ROLE {
			isBlogger = true
		}
	}
	if !isBlogger {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
	var blogPost BlogPost
	err = c.Bind(&blogPost)
	if blogPost.Content == "" || blogPost.Title == "" || blogPost.Status == "" {
		c.String(http.StatusBadRequest, "Fields 'content', 'title' and 'status' are necessary")
		return
	}
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	blogPost.AuthorId = cUser.Id
	nBlogPost, err := bc.br.CreateBlogPost(&blogPost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusCreated, nBlogPost)
}

func (bc BlogController) GetBlogPostsByAdmin(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := bc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := bc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	isBlogger := false
	for _, v := range cUser.Roles {
		if v == shared.BLOGGER_ROLE || v == shared.ADMIN_ROLE {
			isBlogger = true
		}
	}
	if !isBlogger {
		c.String(http.StatusForbidden, "You don't have access to this resource")
		return
	}
	filter := bson.D{{Key: "authorid", Value: cUser.Id}}
	blogPosts, err := bc.br.GetBlogPosts(filter, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, blogPosts)
}

func (bc BlogController) GetPublishedBlogPosts(c *gin.Context) {
	filter := bson.D{{Key: "status", Value: shared.BlogPostStatusPublished}}
	var getPublishedBlogPostsReq GetPublishedBlogPostsRequest
	err := c.BindJSON(&getPublishedBlogPostsReq)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	if getPublishedBlogPostsReq.Limit == 0 {
		getPublishedBlogPostsReq.Limit = math.MaxInt64
	}
	options := options.Find().SetSort(bson.D{{Key: "createdat", Value: -1}}).SetSkip(getPublishedBlogPostsReq.Skip).SetLimit(getPublishedBlogPostsReq.Limit)
	blogPosts, err := bc.br.GetBlogPosts(filter, options)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	c.JSON(http.StatusOK, blogPosts)
}

func (bc BlogController) GetBlogPost(c *gin.Context) {
	var blogPostRequest GetBlogPostRequest
	err := c.BindJSON(&blogPostRequest)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object")
		return
	}
	blogPost, err := bc.br.GetBlogPost(blogPostRequest.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	token := c.Request.Header.Get("Authorization")
	if blogPost.Status == shared.BlogPostStatusPublished && token == "" {
		author, err := bc.ur.GetUser(blogPost.AuthorId.Hex())
		if err != nil {
			c.String(http.StatusInternalServerError, "Something went wrong, please try again")
			return
		}
		c.JSON(http.StatusOK, GetBlogPostResponse{
			Id:              blogPost.Id.Hex(),
			Title:           blogPost.Title,
			Content:         blogPost.Content,
			AuthorFirstName: author.FirstName,
			AuthorLastName:  author.LastName,
		})
		return
	}
	if token == "" {
		c.String(http.StatusForbidden, "Authorization Failed, please login")
		return
	}
	claims, err := bc.auth.ParseToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Token is invalid")
		return
	}
	cUser, err := bc.ur.GetUserByEmail(claims.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again")
		return
	}
	if cUser.Id == blogPost.AuthorId {
		c.JSON(http.StatusOK, blogPost)
		return
	}
	c.String(http.StatusForbidden, "You don't have access to this resource.")
}
