package template

import (
	"strings"
	"testing"
)

func TestFromString(t *testing.T) {
	text := "123456"
	template := NewTemplateText()
	template.FromString(text)
	if !strings.EqualFold(template.Content, text) {
		t.Errorf("FromString() doesnt work properly, expected %v got %v", text, template.Content)
	}
}

func TestReplaceTemplate(t *testing.T) {
	content := "Hello {Username}\n This is a token {Token}"
	template := NewTemplateText()
	template.Content = content
	finalText := "Hello Alighn\n This is a token 123456"
	replacement := map[string]string{
		"Username": "Alighn",
		"Token":    "123456",
	}
	template.ReplaceTemplate(replacement)
	if template.Content != finalText {
		t.Errorf("Replacement does not work properly, expected %v got %v", finalText, template.Content)
	}
	print(template)
}
