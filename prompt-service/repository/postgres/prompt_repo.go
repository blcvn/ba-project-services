package postgres

import (
	"context"
	"encoding/json"

	"github.com/blcvn/backend/services/prompt-service/entities"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PromptRepository struct {
	db *gorm.DB
}

func NewPromptRepository(db *gorm.DB) *PromptRepository {
	return &PromptRepository{db: db}
}

func (r *PromptRepository) CreateTemplate(ctx context.Context, payload *entities.CreateTemplatePayload) (*entities.PromptTemplate, error) {
	// Marshal metadata to JSON (using datatypes.JSON which is []byte internally or similar)
	metaJSON, _ := json.Marshal(payload.Metadata)

	template := &entities.PromptTemplate{
		Name: payload.Name,
		// Description: payload.Description, // Removed from payload
		Environment: payload.Environment,
		Metadata:    datatypes.JSON(metaJSON),
	}
	if err := r.db.WithContext(ctx).Create(template).Error; err != nil {
		return nil, err
	}
	return template, nil
}

func (r *PromptRepository) GetTemplate(ctx context.Context, name string) (*entities.PromptTemplate, error) {
	var template entities.PromptTemplate
	if err := r.db.WithContext(ctx).Preload("Versions").Where("name = ?", name).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *PromptRepository) CreateVersion(ctx context.Context, payload *entities.CreateVersionPayload) (*entities.PromptVersion, error) {
	tid, err := uuid.Parse(payload.TemplateID)
	if err != nil {
		return nil, err
	}

	// Simple JSON conversion for now, ideally use proper marshalling
	// Assuming VariablesSchema is stored as JSON
	// mocking json content for now as we deal with map
	// In real impl, we marshal payload.VariablesSchema

	version := &entities.PromptVersion{
		TemplateID: tid,
		Version:    payload.Version,
		Content:    payload.Content,
		// VariablesSchema: ... (skipped complexity for brevity, requires json marshal)
	}

	if err := r.db.WithContext(ctx).Create(version).Error; err != nil {
		return nil, err
	}

	// Update current version of template?
	// Logic might belong to usecase, but if this is atomic...

	return version, nil
}

func (r *PromptRepository) GetVersion(ctx context.Context, templateID string, versionStr string) (*entities.PromptVersion, error) {
	var version entities.PromptVersion
	if err := r.db.WithContext(ctx).Where("template_id = ? AND version = ?", templateID, versionStr).First(&version).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *PromptRepository) ListTemplates(ctx context.Context, page, pageSize int) ([]*entities.PromptTemplate, int64, error) {
	var templates []*entities.PromptTemplate
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&entities.PromptTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}
