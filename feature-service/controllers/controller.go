package controllers

import (
	"context"
	"time"

	"github.com/blcvn/backend/services/feature-service/helper"
	"github.com/blcvn/backend/services/feature-service/repository/postgres"
	"github.com/blcvn/backend/services/feature-service/usecases"
	pb "github.com/blcvn/kratos-proto/go/feature"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type FeatureController struct {
	pb.UnimplementedFeatureServiceServer
	uc        *usecases.FeatureUsecase
	transform *helper.Transform
}

func NewFeatureController(uc *usecases.FeatureUsecase, transform *helper.Transform) *FeatureController {
	return &FeatureController{uc: uc, transform: transform}
}

func (c *FeatureController) CreateFeature(ctx context.Context, req *pb.CreateFeatureRequest) (*pb.FeatureReply, error) {
	if req.Payload == nil || req.Payload.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "Project ID is required")
	}

	f := &postgres.Feature{
		ProjectID:   req.Payload.ProjectId,
		Name:        req.Payload.Name,
		Description: req.Payload.Description,
	}
	if req.Payload.ParentId != "" {
		f.ParentID = &req.Payload.ParentId
	}

	res, err := c.uc.Create(ctx, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FeatureReply{
		Result:  &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Feature created successfully"},
		Payload: convertToProto(res),
	}, nil
}

func (c *FeatureController) UpdateFeature(ctx context.Context, req *pb.UpdateFeatureRequest) (*pb.FeatureReply, error) {
	if req.Payload == nil || req.Payload.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Feature ID is required")
	}

	f, err := c.uc.Get(ctx, req.Payload.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Feature not found")
	}

	if req.Payload.Name != "" {
		f.Name = req.Payload.Name
	}
	if req.Payload.Description != "" {
		f.Description = req.Payload.Description
	}
	if req.Payload.Status != "" {
		f.Status = req.Payload.Status
	}
	if req.Payload.CurrentPhase != "" {
		f.CurrentPhase = req.Payload.CurrentPhase
	}
	if req.Payload.CurrentStep != "" {
		f.CurrentStep = req.Payload.CurrentStep
	}

	if req.Payload.Progress != nil {
		f.Progress = req.Payload.Progress.AsMap()
	}
	if req.Payload.Artifacts != nil {
		f.Artifacts = req.Payload.Artifacts.AsMap()
	}
	if req.Payload.Data != nil {
		f.Data = req.Payload.Data.AsMap()
	}

	res, err := c.uc.Update(ctx, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FeatureReply{
		Result:  &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Feature updated successfully"},
		Payload: convertToProto(res),
	}, nil
}

func (c *FeatureController) GetFeature(ctx context.Context, req *pb.GetFeatureRequest) (*pb.FeatureReply, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Feature ID is required")
	}

	res, err := c.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Feature not found")
	}

	return &pb.FeatureReply{
		Result:  &pb.Result{Code: pb.ResultCode_SUCCESS},
		Payload: convertToProto(res),
	}, nil
}

func (c *FeatureController) ListFeatures(ctx context.Context, req *pb.ListFeaturesRequest) (*pb.ListFeaturesReply, error) {
	if req.Payload == nil || req.Payload.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "Project ID is required")
	}

	var parentID *string
	if req.Payload.ParentId != "" {
		parentID = &req.Payload.ParentId
	}

	res, err := c.uc.List(ctx, req.Payload.ProjectId, parentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protos := make([]*pb.Feature, len(res))
	for i, f := range res {
		protos[i] = convertToProto(f)
	}

	return &pb.ListFeaturesReply{
		Result:     &pb.Result{Code: pb.ResultCode_SUCCESS},
		Pagination: req.Pagination,
		Payload:    protos,
	}, nil
}

func (c *FeatureController) DeleteFeature(ctx context.Context, req *pb.DeleteFeatureRequest) (*pb.FeatureReply, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Feature ID is required")
	}

	err := c.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FeatureReply{
		Result: &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Feature deleted successfully"},
	}, nil
}

func convertToProto(f *postgres.Feature) *pb.Feature {
	if f == nil {
		return nil
	}

	progress, _ := structpb.NewStruct(f.Progress)
	if progress == nil {
		progress = &structpb.Struct{Fields: make(map[string]*structpb.Value)}
	}

	artifacts, _ := structpb.NewStruct(f.Artifacts)
	if artifacts == nil {
		artifacts = &structpb.Struct{Fields: make(map[string]*structpb.Value)}
	}

	data, _ := structpb.NewStruct(f.Data)
	if data == nil {
		data = &structpb.Struct{Fields: make(map[string]*structpb.Value)}
	}

	res := &pb.Feature{
		Id:           f.ID,
		ProjectId:    f.ProjectID,
		Name:         f.Name,
		Description:  f.Description,
		Order:        f.Order,
		Status:       f.Status,
		CurrentPhase: f.CurrentPhase,
		CurrentStep:  f.CurrentStep,
		Progress:     progress,
		Artifacts:    artifacts,
		Data:         data,
		CreatedAt:    f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    f.UpdatedAt.Format(time.RFC3339),
	}
	if f.ParentID != nil {
		res.ParentId = *f.ParentID
	}

	return res
}
