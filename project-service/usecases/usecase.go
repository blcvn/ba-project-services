package usecases

import (
	"context"

	"github.com/blcvn/backend/services/project-service/repository/postgres"
)

// ProjectRepo defines the interface for repository access from usecase
type ProjectRepo interface {
	Save(context.Context, *postgres.Project) (*postgres.Project, error)
	Update(context.Context, *postgres.Project) (*postgres.Project, error)
	FindByID(context.Context, string) (*postgres.Project, error)
	ListAll(context.Context, string, string, int32, int32) ([]*postgres.Project, int64, error)
	Delete(context.Context, string) error
}

type ProjectUsecase struct {
	repo ProjectRepo
}

func NewProjectUsecase(repo ProjectRepo) *ProjectUsecase {
	return &ProjectUsecase{repo: repo}
}

func (uc *ProjectUsecase) Create(ctx context.Context, p *postgres.Project) (*postgres.Project, error) {
	return uc.repo.Save(ctx, p)
}

func (uc *ProjectUsecase) Update(ctx context.Context, p *postgres.Project) (*postgres.Project, error) {
	return uc.repo.Update(ctx, p)
}

func (uc *ProjectUsecase) Get(ctx context.Context, id string) (*postgres.Project, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *ProjectUsecase) List(ctx context.Context, status, search string, page, limit int32) ([]*postgres.Project, int64, error) {
	return uc.repo.ListAll(ctx, status, search, page, limit)
}

func (uc *ProjectUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
