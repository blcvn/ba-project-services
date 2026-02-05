package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PromptTemplate struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description      string         `gorm:"type:text"`
	CurrentVersionID *uuid.UUID     `gorm:"type:uuid"`
	Environment      int32          `gorm:"default:0"` // 1=DEV, 2=STAGING, 3=PROD
	Metadata         datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt        time.Time      `gorm:"default:now()"`
	UpdatedAt        time.Time      `gorm:"default:now()"`

	// Associations
	Versions []PromptVersion `gorm:"foreignKey:TemplateID"`
}

type PromptVersion struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TemplateID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Version    string         `gorm:"type:varchar(50);not null"`
	Content    string         `gorm:"type:text;not null"`
	Variables  datatypes.JSON `gorm:"type:jsonb;default:'[]'"` // Array of Variable
	IsActive   bool           `gorm:"default:true"`
	CreatedAt  time.Time      `gorm:"default:now()"`
	UpdatedAt  time.Time      `gorm:"default:now()"`
}

type CreateTemplatePayload struct {
	Name        string
	Version     string
	Content     string
	Environment int32
	Variables   interface{} // JSON serializable
	Metadata    map[string]string
}

type Variable struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Description  string `json:"description"`
	DefaultValue string `json:"default_value"`
	Required     bool   `json:"required"`
}

type CreateVersionPayload struct {
	TemplateID string
	Version    string
	Content    string
	Variables  interface{}
}
