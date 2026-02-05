package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/blcvn/backend/services/prompt-service/common/errors"
	"github.com/blcvn/backend/services/prompt-service/dto"
	"github.com/blcvn/backend/services/prompt-service/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type promptRepository struct {
	db *gorm.DB
}

func NewPromptRepository(db *gorm.DB) *promptRepository {
	return &promptRepository{db: db}
}

// CreateTemplate creates a new prompt template
func (r *promptRepository) CreateTemplate(ctx context.Context, payload *entities.CreateTemplatePayload) (*entities.PromptTemplate, errors.BaseError) {
	// Check existing name
	var count int64
	r.db.Model(&dto.PromptTemplate{}).Where("name = ?", payload.Name).Count(&count)
	if count > 0 {
		return nil, errors.Conflict("template with this name already exists")
	}

	varsJSON, _ := json.Marshal(payload.Variables)
	tagsJSON, _ := json.Marshal(payload.Tags)

	dtoTemplate := &dto.PromptTemplate{
		ID:          uuid.New(),
		Name:        payload.Name,
		Description: payload.Description,
		Version:     "v1",
		Content:     payload.Content,
		Variables:   string(varsJSON),
		Tags:        string(tagsJSON),
		Status:      string(entities.TemplateStatusActive),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(dtoTemplate).Error; err != nil {
		return nil, errors.Internal(err)
	}

	return r.dtoToEntity(dtoTemplate)
}

// GetTemplate retrieves a template
func (r *promptRepository) GetTemplate(ctx context.Context, id string) (*entities.PromptTemplate, errors.BaseError) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.BadRequest("invalid id format")
	}

	var dtoTemplate dto.PromptTemplate
	if err := r.db.WithContext(ctx).Where("id = ?", uid).First(&dtoTemplate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("template not found")
		}
		return nil, errors.Internal(err)
	}

	return r.dtoToEntity(&dtoTemplate)
}

// GetTemplateByName retrieves a template by name
func (r *promptRepository) GetTemplateByName(ctx context.Context, name string) (*entities.PromptTemplate, errors.BaseError) {
	var dtoTemplate dto.PromptTemplate
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&dtoTemplate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("template not found")
		}
		return nil, errors.Internal(err)
	}

	return r.dtoToEntity(&dtoTemplate)
}

// ListTemplates lists templates
func (r *promptRepository) ListTemplates(ctx context.Context, filter *entities.TemplateFilter) ([]*entities.PromptTemplate, int64, errors.BaseError) {
	query := r.db.WithContext(ctx).Model(&dto.PromptTemplate{})

	if filter.Status != "" {
		query = query.Where("status = ?", string(filter.Status))
	}
	// TODO: Implement tag filtering (requires JSONB query)

	var total int64
	query.Count(&total)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(int(offset)).Limit(int(filter.PageSize))
	}

	var dtos []dto.PromptTemplate
	if err := query.Order("created_at DESC").Find(&dtos).Error; err != nil {
		return nil, 0, errors.Internal(err)
	}

	results := make([]*entities.PromptTemplate, 0, len(dtos))
	for _, d := range dtos {
		entity, _ := r.dtoToEntity(&d)
		results = append(results, entity)
	}

	return results, total, nil
}

// UpdateTemplate updates a template
func (r *promptRepository) UpdateTemplate(ctx context.Context, payload *entities.UpdateTemplatePayload) (*entities.PromptTemplate, errors.BaseError) {
	uid, err := uuid.Parse(payload.ID)
	if err != nil {
		return nil, errors.BadRequest("invalid id format")
	}

	updates := make(map[string]interface{})
	if payload.Content != "" {
		updates["content"] = payload.Content
	}
	if payload.Status != "" {
		updates["status"] = string(payload.Status)
	}
	if payload.Variables != nil {
		varsJSON, _ := json.Marshal(payload.Variables)
		updates["variables"] = string(varsJSON)
	}
	if payload.Tags != nil {
		tagsJSON, _ := json.Marshal(payload.Tags)
		updates["tags"] = string(tagsJSON)
	}
	updates["updated_at"] = time.Now()

	if err := r.db.WithContext(ctx).Model(&dto.PromptTemplate{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
		return nil, errors.Internal(err)
	}

	return r.GetTemplate(ctx, payload.ID)
}

// DeleteTemplate deletes a template
func (r *promptRepository) DeleteTemplate(ctx context.Context, id string) errors.BaseError {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.BadRequest("invalid id format")
	}

	if err := r.db.WithContext(ctx).Delete(&dto.PromptTemplate{}, "id = ?", uid).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *promptRepository) dtoToEntity(d *dto.PromptTemplate) (*entities.PromptTemplate, errors.BaseError) {
	var vars []entities.Variable
	_ = json.Unmarshal([]byte(d.Variables), &vars)

	var tags []string
	_ = json.Unmarshal([]byte(d.Tags), &tags)

	return &entities.PromptTemplate{
		ID:          d.ID.String(),
		Name:        d.Name,
		Description: d.Description,
		Version:     d.Version,
		Content:     d.Content,
		Variables:   vars,
		Tags:        tags,
		Status:      entities.TemplateStatus(d.Status),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}
