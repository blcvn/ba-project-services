package postgres

import (
	"context"

	"gorm.io/gorm"
)

type FeatureRepository struct {
	db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) *FeatureRepository {
	return &FeatureRepository{db: db}
}

func (r *FeatureRepository) Create(ctx context.Context, f *Feature) (*Feature, error) {
	if err := r.db.WithContext(ctx).Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (r *FeatureRepository) Update(ctx context.Context, f *Feature) (*Feature, error) {
	if err := r.db.WithContext(ctx).Save(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (r *FeatureRepository) Get(ctx context.Context, id string) (*Feature, error) {
	var f Feature
	if err := r.db.WithContext(ctx).First(&f, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FeatureRepository) List(ctx context.Context, projectID string, parentID *string) ([]*Feature, error) {
	var fs []*Feature
	query := r.db.WithContext(ctx).Where("project_id = ?", projectID)
	if parentID != nil {
		if *parentID != "all" {
			query = query.Where("parent_id = ?", *parentID)
		}
		// If *parentID == "all", we don't add parent_id filter, returning all features
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.Order("sort_order ASC").Find(&fs).Error; err != nil {
		return nil, err
	}
	return fs, nil
}

func (r *FeatureRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Feature{}, "id = ?", id).Error
}
