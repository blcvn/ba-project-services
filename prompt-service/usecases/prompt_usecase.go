package usecases

import (
	"context"
	"fmt"

	"github.com/blcvn/backend/services/prompt-service/entities"
	"github.com/blcvn/backend/services/prompt-service/helper"
	"github.com/blcvn/backend/services/prompt-service/repository/postgres"
)

type PromptUsecase struct {
	repo           *postgres.PromptRepository
	templateEngine *helper.TemplateEngine
}

func NewPromptUsecase(repo *postgres.PromptRepository, engine *helper.TemplateEngine) *PromptUsecase {
	return &PromptUsecase{
		repo:           repo,
		templateEngine: engine,
	}
}

func (u *PromptUsecase) CreateTemplate(ctx context.Context, payload *entities.CreateTemplatePayload) (*entities.PromptTemplate, error) {
	// Validation
	if payload.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// 1. Create Template
	tmpl, err := u.repo.CreateTemplate(ctx, payload)
	if err != nil {
		return nil, err
	}

	// 2. Create Initial Version if provided
	if payload.Version != "" && payload.Content != "" {
		versionPayload := &entities.CreateVersionPayload{
			TemplateID: tmpl.ID.String(),
			Version:    payload.Version,
			Content:    payload.Content,
			Variables:  payload.Variables,
		}
		_, err := u.repo.CreateVersion(ctx, versionPayload)
		if err != nil {
			// Basic cleanup if version creation fails (optional, better with transaction)
			// u.repo.DeleteTemplate(ctx, tmpl.ID)
			return nil, fmt.Errorf("failed to create initial version: %w", err)
		}
		// Reload template with versions?
	}

	return tmpl, nil
}

func (u *PromptUsecase) RetrieveTemplate(ctx context.Context, name string) (*entities.PromptTemplate, error) {
	return u.repo.GetTemplate(ctx, name)
}

func (u *PromptUsecase) CreateVersion(ctx context.Context, payload *entities.CreateVersionPayload) (*entities.PromptVersion, error) {
	// Validation
	if payload.Version == "" {
		return nil, fmt.Errorf("version is required")
	}
	if payload.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// Check if template exists? (Repo handles constraint but good to check)

	return u.repo.CreateVersion(ctx, payload)
}

func (u *PromptUsecase) RenderPrompt(ctx context.Context, name string, version string, variables map[string]string) (string, string, error) {
	// 1. Get Template
	tmpl, err := u.repo.GetTemplate(ctx, name)
	if err != nil {
		return "", "", fmt.Errorf("template not found: %w", err)
	}

	var content string
	var versionUsed string

	// 2. Select Version
	if version != "" {
		// Specific version
		v, err := u.repo.GetVersion(ctx, tmpl.ID.String(), version)
		if err != nil {
			return "", "", fmt.Errorf("version not found: %w", err)
		}
		content = v.Content
		versionUsed = v.Version
	} else {
		// Latest / Current Version
		// If CurrentVersionID is set, use it. Else use latest created?
		// For simplicity, let's assume valid templates have CurrentVersionID or we pick the last one.
		// Since we didn't implement 'GetLatestVersion' in repo yet, let's look at Loaded Versions.
		if len(tmpl.Versions) > 0 {
			// Naive: take the last one. Ideally sort by created_at desc in repo.
			// assuming 'Preload("Versions")' returns them.
			latest := tmpl.Versions[len(tmpl.Versions)-1]
			content = latest.Content
			versionUsed = latest.Version
		} else {
			return "", "", fmt.Errorf("no versions available for template")
		}
	}

	// 3. Render
	rendered, err := u.templateEngine.Render(content, variables)
	if err != nil {
		return "", "", err
	}

	return rendered, versionUsed, nil
}

func (u *PromptUsecase) ListTemplates(ctx context.Context, page, pageSize int) ([]*entities.PromptTemplate, int64, error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}
	return u.repo.ListTemplates(ctx, page, pageSize)
}
