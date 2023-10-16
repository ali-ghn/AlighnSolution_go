package blog

import (
	"context"
	"fmt"
	"testing"

	"github.com/ali-ghn/Coinopay_Go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	br BlogRepository
)

func init() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	br = NewBlogRepository(client)
}

func TestCreateBlogPost(t *testing.T) {
	authorId, err := primitive.ObjectIDFromHex("632a13361d431376a04117ab")
	if err != nil {
		t.Error(err)
	}
	blogPost := BlogPost{
		Title:    "blog post title",
		Content:  "# This is a header",
		AuthorId: authorId,
		Status:   shared.BlogPostStatusPublished,
	}
	res, err := br.CreateBlogPost(&blogPost)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res.Content)
	fmt.Println(res.Id)
}

func TestGetBlogPost(t *testing.T) {
	id := "634183481e8a042c48b214f5"
	blogPost, err := br.GetBlogPost(id)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blogPost.Id)
	fmt.Println(blogPost.Content)
}

func TestGetBlogPosts(t *testing.T) {
	filter := bson.D{{}}
	blogPosts, err := br.GetBlogPosts(filter, nil)
	if err != nil {
		t.Error(err)
	}
	for _, v := range *blogPosts {
		fmt.Println(v.Id)
		fmt.Println(v.Content)
	}
}
