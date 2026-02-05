package entities

import (
	"time"
)

// TemplateStatus represents the status of a prompt template
type TemplateStatus string

const (
	TemplateStatusActive   TemplateStatus = "active"
	TemplateStatusArchived TemplateStatus = "archived"
	TemplateStatusDraft    TemplateStatus = "draft"
)

// Variable represents a variable in a prompt template
type Variable struct {
	Name         string
	Description  string
	Type         string // string, number, boolean, json
	Required     bool
	DefaultValue string
}

// PromptTemplate represents a reusable prompt structure
type PromptTemplate struct {
	ID          string
	Name        string
	Description string
	Version     string
	Content     string
	Variables   []Variable
	Tags        []string
	Status      TemplateStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PromptExperiment represents an A/B test for prompts
type PromptExperiment struct {
	ID               string
	Name             string
	PromptTemplateID string
	ModelID          string
	Config           map[string]string
	Status           string // active, completed
	CreatedAt        time.Time
}

// RenderedPrompt represents the result of filling a template
type RenderedPrompt struct {
	Content   string
	Variables map[string]string
}

// CreateTemplatePayload payload for creating a template
type CreateTemplatePayload struct {
	Name        string
	Description string
	Content     string
	Variables   []Variable
	Tags        []string
}

// UpdateTemplatePayload payload for updating a template
type UpdateTemplatePayload struct {
	ID        string
	Content   string
	Variables []Variable
	Status    TemplateStatus
	Tags      []string
}

// TemplateFilter filter for listing templates
type TemplateFilter struct {
	Status   TemplateStatus
	Tags     []string
	Page     int32
	PageSize int32
}
