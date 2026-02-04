package usecases

import (
	"context"
	"time"

	"github.com/blcvn/backend/services/feature-service/repository/postgres"
)

type FeatureRepo interface {
	Create(ctx context.Context, f *postgres.Feature) (*postgres.Feature, error)
	Update(ctx context.Context, f *postgres.Feature) (*postgres.Feature, error)
	Get(ctx context.Context, id string) (*postgres.Feature, error)
	List(ctx context.Context, projectID string, parentID *string) ([]*postgres.Feature, error)
	Delete(ctx context.Context, id string) error
}

type FeatureUsecase struct {
	repo FeatureRepo
}

func NewFeatureUsecase(repo FeatureRepo) *FeatureUsecase {
	return &FeatureUsecase{repo: repo}
}

func (uc *FeatureUsecase) Create(ctx context.Context, f *postgres.Feature) (*postgres.Feature, error) {
	if f.Status == "" {
		f.Status = "active"
	}
	return uc.repo.Create(ctx, f)
}

func (uc *FeatureUsecase) Update(ctx context.Context, f *postgres.Feature) (*postgres.Feature, error) {
	f.UpdatedAt = time.Now()
	// Fetch existing to ensure it exists
	existing, err := uc.repo.Get(ctx, f.ID)
	if err != nil {
		return nil, err
	}
	// Use existing created_at if not set
	if f.CreatedAt.IsZero() {
		f.CreatedAt = existing.CreatedAt
	}
	return uc.repo.Update(ctx, f)
}

func (uc *FeatureUsecase) Get(ctx context.Context, id string) (*postgres.Feature, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *FeatureUsecase) List(ctx context.Context, projectID string, parentID *string) ([]*postgres.Feature, error) {
	return uc.repo.List(ctx, projectID, parentID)
}

func (uc *FeatureUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
