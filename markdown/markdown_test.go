package markdown

import (
	"context"
	"github.com/ali-ghn/AlighnSolution_go/cors"
	"github.com/ali-ghn/AlighnSolution_go/template"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"testing"
)

var mh *MarkdownHelper

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err.Error())
	}
	tt := template.NewTemplateText()
	mh = NewIMarkdownHelper(client, tt)

}
func TestMarkdownHelper_CreateTemplateMarkdown(t *testing.T) {
	markdownTemplateText := `| {11} | {12} | {13} |
|------|------|------|
| {21} | {22} | {23} |
| {31} | {32} | {33} |`
	markdownTemplate := &MarkdownTemplate{Content: markdownTemplateText, Name: "main01"}
	markdownTemplate, err := mh.CreateTemplateMarkdown(markdownTemplate)
	if err != nil {
		t.Error(err.Error())
	}
	spew.Dump(markdownTemplate)
}

func MarkdownHelper_TemplateMarkdownToHtmlControllerFunction(c *gin.Context) {
	name := "main01"
	keyValuePair := map[string]string{}
	columnKeys := []string{"a", "b", "c"}
	var rowKeys []string
	rowKeys = columnKeys
	for _, key := range columnKeys {
		for _, rowKey := range rowKeys {
			keyValuePair[key+rowKey] = "Alighn"
		}
	}
	markdownTemplate := mh.TemplateMarkdownToHtml(name, keyValuePair)
	c.Data(http.StatusOK, "text/html; charset=utf-8", markdownTemplate)
}

func TestMarkdownHelper_TemplateMarkdownToHtml(t *testing.T) {
	r := gin.Default()
	// CORS middleware
	r.GET("/", MarkdownHelper_TemplateMarkdownToHtmlControllerFunction)
	r.Use(cors.CORSMiddleware())
	r.Run("127.0.0.1:8005")
}
