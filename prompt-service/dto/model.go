package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PromptTemplate represents the database model for prompt templates
type PromptTemplate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	Version     string    `gorm:"type:varchar(50);default:'v1'"`
	Content     string    `gorm:"type:text;not null"`
	Variables   string    `gorm:"type:jsonb;default:'[]'"` // JSON array of variables
	Tags        string    `gorm:"type:jsonb;default:'[]'"` // JSON array of tags
	Status      string    `gorm:"type:varchar(50);default:'active';index"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`
}

// TableName specifies the table name
func (PromptTemplate) TableName() string {
	return "prompt_templates"
}

// BeforeUpdate hook
func (t *PromptTemplate) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	// Logic to increment version could go here
	return nil
}

// Experiment represents the database model for experiments
type Experiment struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name             string    `gorm:"type:varchar(255);not null"`
	PromptTemplateID uuid.UUID `gorm:"type:uuid;not null;index"`
	ModelID          uuid.UUID `gorm:"type:uuid;not null;index"`
	Config           string    `gorm:"type:jsonb;default:'{}'"`
	Status           string    `gorm:"type:varchar(50);default:'active'"`
	CreatedAt        time.Time `gorm:"default:now()"`
}

// TableName specifies the table name
func (Experiment) TableName() string {
	return "prompt_experiments"
}
