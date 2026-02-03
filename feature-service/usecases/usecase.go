package usecases

type FeatureRepo interface{}

type FeatureUsecase struct {
	repo FeatureRepo
}

func NewFeatureUsecase(repo FeatureRepo) *FeatureUsecase {
	return &FeatureUsecase{repo: repo}
}
