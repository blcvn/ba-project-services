package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSONMap is a helper for storing map[string]interface{} as JSONB/JSON
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type Feature struct {
	ID           string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProjectID    string  `gorm:"type:uuid;not null;index"`
	ParentID     *string `gorm:"type:uuid;index"`
	Name         string  `gorm:"not null"`
	Description  string
	Order        int32  `gorm:"column:sort_order;default:0"`
	Status       string `gorm:"default:active"`
	CurrentPhase string
	CurrentStep  string

	Progress  JSONMap `gorm:"type:jsonb"`
	Artifacts JSONMap `gorm:"type:jsonb"`
	Data      JSONMap `gorm:"type:jsonb"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
