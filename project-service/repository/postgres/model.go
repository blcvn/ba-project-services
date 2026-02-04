package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project represents the project entity in the database
type Project struct {
	ID                     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID                 string `gorm:"type:uuid;index"`
	TenantID               string `gorm:"type:uuid;index"`
	Name                   string `gorm:"not null"`
	Description            string
	Status                 string `gorm:"default:'active'"` // active, paused, completed, archived
	ConfluenceEnabled      bool   `gorm:"default:false"`
	ConfluenceSpaceKey     string
	ConfluenceParentPageID string
	FeatureCount           int32 `gorm:"default:0"`
	Progress               int32 `gorm:"default:0"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate hook to generate UUID if not present
func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}
