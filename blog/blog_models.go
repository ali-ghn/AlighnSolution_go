package blog

import "go.mongodb.org/mongo-driver/bson/primitive"

type BlogPost struct {
	Id        primitive.ObjectID `bson:"_id"`
	Title     string
	Content   string
	AuthorId  primitive.ObjectID
	Status    string
	CreatedAt int64
	UpdatedAt int64
}

type GetPublishedBlogPostsRequest struct {
	Skip  int64
	Limit int64
}
type GetBlogPostRequest struct {
	Id string
}
type GetBlogPostResponse struct {
	Id              string
	Title           string
	Content         string
	AuthorFirstName string
	AuthorLastName  string
}
