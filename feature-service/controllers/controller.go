package controllers

import (
	"github.com/blcvn/backend/services/feature-service/helper"
	"github.com/blcvn/backend/services/feature-service/usecases"
	pb "github.com/blcvn/kratos-proto/go/feature"
)

type FeatureController struct {
	pb.UnimplementedFeatureServiceServer
	uc        *usecases.FeatureUsecase
	transform *helper.Transform
}

func NewFeatureController(uc *usecases.FeatureUsecase, transform *helper.Transform) *FeatureController {
	return &FeatureController{uc: uc, transform: transform}
}
