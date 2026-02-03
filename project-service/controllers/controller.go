package controllers

import (
	"github.com/blcvn/backend/services/project-service/helper"
	"github.com/blcvn/backend/services/project-service/usecases"
	pb "github.com/blcvn/kratos-proto/go/project"
)

type ProjectController struct {
	pb.UnimplementedProjectServiceServer
	uc        *usecases.ProjectUsecase
	transform *helper.Transform
}

func NewProjectController(uc *usecases.ProjectUsecase, transform *helper.Transform) *ProjectController {
	return &ProjectController{uc: uc, transform: transform}
}
