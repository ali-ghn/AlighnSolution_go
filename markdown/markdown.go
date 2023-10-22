package markdown

import (
	"context"
	"github.com/ali-ghn/AlighnSolution_go/shared"
	"github.com/ali-ghn/AlighnSolution_go/template"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const Collection = "MarkdownTemplate"

type MarkdownTemplate struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}

type IMarkdownHelper interface {
	MarkdownToHtml(input string) string
	TemplateMarkdownToHtml(templateName string) string
}

type MarkdownHelper struct {
	Client *mongo.Client
	tt     *template.TemplateText
}

func NewIMarkdownHelper(client *mongo.Client, tt *template.TemplateText) *MarkdownHelper {
	return &MarkdownHelper{
		Client: client,
		tt:     tt,
	}
}

func (mh *MarkdownHelper) MarkdownToHtml(input string) string {
	unsafe := blackfriday.Run([]byte(input))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html)
}

func (mh *MarkdownHelper) CreateTemplateMarkdown(markdownTemplate *MarkdownTemplate) (*MarkdownTemplate, error) {
	markdownTemplate.Id = primitive.NewObjectID()
	markdownTemplate.CreatedAt = time.Now().UTC()
	markdownTemplate.UpdatedAt = time.Now().UTC()
	markdownTemplate.IsActive = true
	_, err := mh.Client.Database(shared.DATABASE_NAME).Collection(Collection).InsertOne(context.TODO(), markdownTemplate)
	if err != nil {
		return nil, err
	}
	// TODO add result.InsertId
	return markdownTemplate, nil
}

func (mh *MarkdownHelper) TemplateMarkdownToHtml(templateName string, keyValuePair map[string]string) []byte {
	filter := bson.D{{"name", templateName}}
	var markdownTemplate MarkdownTemplate
	err := mh.Client.Database(shared.DATABASE_NAME).Collection(Collection).FindOne(context.TODO(), filter).Decode(&markdownTemplate)
	if err != nil {
		return nil
	}
	textTemplate := mh.tt.FromString(markdownTemplate.Content)
	textTemplate.ReplaceTemplate(keyValuePair)
	unsafe := blackfriday.Run([]byte(textTemplate.Content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return html
}
