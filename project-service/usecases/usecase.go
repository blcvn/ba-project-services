package usecases

type ProjectRepo interface{}

type ProjectUsecase struct {
	repo ProjectRepo
}

func NewProjectUsecase(repo ProjectRepo) *ProjectUsecase {
	return &ProjectUsecase{repo: repo}
}
