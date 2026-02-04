package postgres

import (
	"context"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	// Ensure the table is migrated
	db.AutoMigrate(&Project{})
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Save(ctx context.Context, p *Project) (*Project, error) {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, p *Project) (*Project, error) {
	if err := r.db.WithContext(ctx).Model(&Project{}).Where("id = ?", p.ID).Updates(p).Error; err != nil {
		return nil, err
	}
	// Fetch updated record to return fresh data
	var updated Project
	if err := r.db.WithContext(ctx).First(&updated, "id = ?", p.ID).Error; err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id string) (*Project, error) {
	var p Project
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) ListAll(ctx context.Context, status, search string, page, limit int32) ([]*Project, int64, error) {
	var projects []*Project
	var total int64
	db := r.db.WithContext(ctx).Model(&Project{})

	if status != "" && status != "all" {
		db = db.Where("status = ?", status)
	}
	if search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := db.Offset(int(offset)).Limit(int(limit)).Order("created_at desc").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Project{}, "id = ?", id).Error
}
