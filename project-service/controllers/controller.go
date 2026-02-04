package controllers

import (
	"context"
	"fmt"

	"github.com/blcvn/backend/services/project-service/helper"
	"github.com/blcvn/backend/services/project-service/repository/postgres"
	"github.com/blcvn/backend/services/project-service/usecases"
	pb "github.com/blcvn/kratos-proto/go/project"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectController struct {
	pb.UnimplementedProjectServiceServer
	uc        *usecases.ProjectUsecase
	transform *helper.Transform
}

func NewProjectController(uc *usecases.ProjectUsecase, transform *helper.Transform) *ProjectController {
	return &ProjectController{uc: uc, transform: transform}
}

func (c *ProjectController) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.ProjectReply, error) {
	if req.Payload == nil || req.Payload.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Project name is required")
	}

	// Extract headers
	var userID, tenantID string
	if tr, ok := transport.FromServerContext(ctx); ok {
		userID = tr.RequestHeader().Get("X-User-ID")
		tenantID = tr.RequestHeader().Get("X-Tenant-ID")
		if tenantID == "" {
			tenantID = tr.RequestHeader().Get("Grpc-Metadata-tenant_id")
		}
	}
	fmt.Printf("CreateProject: UserID=%s TenantID=%s\n", userID, tenantID)

	if userID == "" || tenantID == "" {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Missing Identity headers. UserID=%s TenantID=%s", userID, tenantID))
	}

	// Map Proto to Model
	p := &postgres.Project{
		Name:        req.Payload.Name,
		Description: req.Payload.Description,
		UserID:      userID,
		TenantID:    tenantID,
	}

	if req.Payload.ConfluenceConfig != nil {
		p.ConfluenceEnabled = req.Payload.ConfluenceConfig.Enabled
		p.ConfluenceSpaceKey = req.Payload.ConfluenceConfig.SpaceKey
		p.ConfluenceParentPageID = req.Payload.ConfluenceConfig.ParentPageId
	}

	result, err := c.uc.Create(ctx, p)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProjectReply{
		Result:  &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Project created successfully"},
		Payload: convertToProto(result),
	}, nil
}

func (c *ProjectController) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.ProjectReply, error) {
	if req.Payload == nil || req.Payload.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Project ID is required")
	}

	p := &postgres.Project{
		ID:          req.Payload.Id,
		Name:        req.Payload.Name,
		Description: req.Payload.Description,
		Status:      req.Payload.Status,
	}
	if req.Payload.ConfluenceConfig != nil {
		p.ConfluenceEnabled = req.Payload.ConfluenceConfig.Enabled
		p.ConfluenceSpaceKey = req.Payload.ConfluenceConfig.SpaceKey
		p.ConfluenceParentPageID = req.Payload.ConfluenceConfig.ParentPageId
	}

	result, err := c.uc.Update(ctx, p)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProjectReply{
		Result:  &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Project updated successfully"},
		Payload: convertToProto(result),
	}, nil
}

func (c *ProjectController) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.ProjectReply, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Project ID is required")
	}

	result, err := c.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Project not found")
	}

	return &pb.ProjectReply{
		Result:  &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Project found"},
		Payload: convertToProto(result),
	}, nil
}

func (c *ProjectController) ListProjects(ctx context.Context, req *pb.ListProjectsRequest) (*pb.ListProjectsReply, error) {
	var page, limit int32 = 1, 10
	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = req.Pagination.Page
		}
		if req.Pagination.Limit > 0 {
			limit = req.Pagination.Limit
		}
	}

	statusFilter := ""
	searchFilter := ""
	if req.Payload != nil {
		statusFilter = req.Payload.Status
		searchFilter = req.Payload.Search
	}

	projects, total, err := c.uc.List(ctx, statusFilter, searchFilter, page, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoProjects []*pb.Project
	for _, p := range projects {
		protoProjects = append(protoProjects, convertToProto(p))
	}

	return &pb.ListProjectsReply{
		Result: &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Projects retrieved successfully"},
		Pagination: &pb.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
		Payload: protoProjects,
	}, nil
}

func (c *ProjectController) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteReply, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Project ID is required")
	}

	if err := c.uc.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteReply{
		Result: &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "Project deleted successfully"},
	}, nil
}

// Helper to convert DB model to Proto model
func convertToProto(p *postgres.Project) *pb.Project {
	return &pb.Project{
		Id:          p.ID,
		UserId:      p.UserID,
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status,
		ConfluenceConfig: &pb.ConfluenceConfig{
			Enabled:      p.ConfluenceEnabled,
			SpaceKey:     p.ConfluenceSpaceKey,
			ParentPageId: p.ConfluenceParentPageID,
		},
		FeatureCount: p.FeatureCount,
		Progress:     p.Progress,
		CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
