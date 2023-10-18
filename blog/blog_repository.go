package blog

import (
	"context"
	"time"

	"github.com/ali-ghn/AlighnSolution_go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "BlogPosts"
)

type IBlogRepository interface {
	CreateBlogPost(blogPost *BlogPost) (*BlogPost, error)
	GetBlogPost(id string) (*BlogPost, error)
	GetBlogPosts(filter interface{}, options options.FindOptions) (*[]BlogPost, error)
	UpdateBlogPost(blogPost *BlogPost) (*BlogPost, error)
}

type BlogRepository struct {
	Client *mongo.Client
}

func NewBlogRepository(client *mongo.Client) BlogRepository {
	return BlogRepository{
		Client: client,
	}
}

func (br BlogRepository) CreateBlogPost(blogPost *BlogPost) (*BlogPost, error) {
	blogPost.Id = primitive.NewObjectID()
	blogPost.CreatedAt = time.Now().UTC().Unix()
	blogPost.UpdatedAt = time.Now().UTC().Unix()
	_, err := br.Client.Database(shared.DATABASE_NAME).Collection(collectionName).InsertOne(context.TODO(), blogPost)
	if err != nil {
		return nil, err
	}
	return blogPost, nil
}

func (br BlogRepository) GetBlogPost(id string) (*BlogPost, error) {
	bId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", bId}}
	var blogPost BlogPost
	err = br.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&blogPost)
	if err != nil {
		return nil, err
	}
	return &blogPost, nil
}

func (br BlogRepository) GetBlogPosts(filter interface{}, options *options.FindOptions) (*[]BlogPost, error) {
	var blogPosts []BlogPost
	cur, err := br.Client.Database(shared.DATABASE_NAME).Collection(collectionName).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	cur.All(context.TODO(), &blogPosts)
	return &blogPosts, nil
}

func (br BlogRepository) UpdateBlogPost(blogPost *BlogPost) (*BlogPost, error) {
	filter := bson.D{{"_id", blogPost.Id}}
	err := br.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOneAndReplace(context.TODO(), filter, blogPost).Decode(&blogPost)
	if err != nil {
		return nil, err
	}
	return blogPost, nil
}
