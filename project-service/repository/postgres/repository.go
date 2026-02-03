package postgres

import "gorm.io/gorm"

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}
