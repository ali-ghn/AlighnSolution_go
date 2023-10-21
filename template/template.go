package template

import (
	"regexp"
	"strings"
)

type TemplateText struct {
	Content string
}

func NewTemplateText() *TemplateText {
	return new(TemplateText)
}

func (input *TemplateText) FromString(text string) *TemplateText {
	input.Content = text
	return input
}

func (input *TemplateText) ReplaceTemplate(keyValuePair map[string]string) {
	re := regexp.MustCompile(`\{ *([a-zA-Z_]+|[a-zA-Z_]+(( +[a-zA-Z_]+)+)) *\}`)
	foundKeyReplacement := re.FindAllString(input.Content, -1)
	for _, value := range foundKeyReplacement {
		input.Content = strings.Replace(input.Content, value, keyValuePair[strings.Replace(strings.Replace(value, "{", "", -1), "}", "", -1)], -1)
	}
}
