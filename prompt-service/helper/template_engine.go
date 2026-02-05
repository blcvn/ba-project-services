package helper

import (
	"bytes"
	"fmt"
	"text/template"
)

type TemplateEngine struct{}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

// Render processes a text template with provided variables
func (e *TemplateEngine) Render(content string, variables map[string]string) (string, error) {
	tmpl, err := template.New("prompt").Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
