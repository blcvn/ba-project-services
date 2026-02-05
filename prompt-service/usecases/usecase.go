package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/blcvn/backend/services/prompt-service/common/errors"
	"github.com/blcvn/backend/services/prompt-service/entities"
)

type iPromptRepository interface {
	CreateTemplate(ctx context.Context, payload *entities.CreateTemplatePayload) (*entities.PromptTemplate, errors.BaseError)
	GetTemplate(ctx context.Context, id string) (*entities.PromptTemplate, errors.BaseError)
	GetTemplateByName(ctx context.Context, name string) (*entities.PromptTemplate, errors.BaseError)
	ListTemplates(ctx context.Context, filter *entities.TemplateFilter) ([]*entities.PromptTemplate, int64, errors.BaseError)
	UpdateTemplate(ctx context.Context, payload *entities.UpdateTemplatePayload) (*entities.PromptTemplate, errors.BaseError)
	DeleteTemplate(ctx context.Context, id string) errors.BaseError
}

type promptUsecase struct {
	repo iPromptRepository
}

func NewPromptUsecase(repo iPromptRepository) *promptUsecase {
	return &promptUsecase{repo: repo}
}

func (u *promptUsecase) CreateTemplate(ctx context.Context, payload *entities.CreateTemplatePayload) (*entities.PromptTemplate, errors.BaseError) {
	if payload.Name == "" || payload.Content == "" {
		return nil, errors.BadRequest("name and content are required")
	}
	return u.repo.CreateTemplate(ctx, payload)
}

func (u *promptUsecase) GetTemplate(ctx context.Context, id string) (*entities.PromptTemplate, errors.BaseError) {
	return u.repo.GetTemplate(ctx, id)
}

func (u *promptUsecase) ListTemplates(ctx context.Context, filter *entities.TemplateFilter) ([]*entities.PromptTemplate, int64, errors.BaseError) {
	return u.repo.ListTemplates(ctx, filter)
}

func (u *promptUsecase) UpdateTemplate(ctx context.Context, payload *entities.UpdateTemplatePayload) (*entities.PromptTemplate, errors.BaseError) {
	return u.repo.UpdateTemplate(ctx, payload)
}

func (u *promptUsecase) DeleteTemplate(ctx context.Context, id string) errors.BaseError {
	return u.repo.DeleteTemplate(ctx, id)
}

func (u *promptUsecase) RenderTemplate(ctx context.Context, name string, variables map[string]string) (*entities.RenderedPrompt, errors.BaseError) {
	template, err := u.repo.GetTemplateByName(ctx, name)
	if err != nil {
		return nil, err
	}

	content := template.Content

	// Validate variables and replace
	for _, v := range template.Variables {
		val, ok := variables[v.Name]
		if !ok {
			if v.Required && v.DefaultValue == "" {
				return nil, errors.BadRequest(fmt.Sprintf("missing required variable: %s", v.Name))
			}
			val = v.DefaultValue
		}

		// Simple replacement logic (can be upgraded to text/template or pongo2)
		// Syntax: {{variable}}
		placeholder := fmt.Sprintf("{{%s}}", v.Name)
		content = strings.ReplaceAll(content, placeholder, val)
	}

	return &entities.RenderedPrompt{
		Content:   content,
		Variables: variables,
	}, nil
}
