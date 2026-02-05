package controllers

import (
	"context"

	"github.com/blcvn/backend/services/prompt-service/common/errors"
	"github.com/blcvn/backend/services/prompt-service/entities"
	"github.com/blcvn/backend/services/prompt-service/helper"
	pb "github.com/blcvn/kratos-proto/go/prompt"
)

type iPromptUsecase interface {
	CreateTemplate(ctx context.Context, payload *entities.CreateTemplatePayload) (*entities.PromptTemplate, errors.BaseError)
	GetTemplate(ctx context.Context, id string) (*entities.PromptTemplate, errors.BaseError)
	ListTemplates(ctx context.Context, filter *entities.TemplateFilter) ([]*entities.PromptTemplate, int64, errors.BaseError)
	UpdateTemplate(ctx context.Context, payload *entities.UpdateTemplatePayload) (*entities.PromptTemplate, errors.BaseError)
	DeleteTemplate(ctx context.Context, id string) errors.BaseError
	RenderTemplate(ctx context.Context, name string, variables map[string]string) (*entities.RenderedPrompt, errors.BaseError)
}

type promptController struct {
	pb.UnimplementedPromptServiceServer
	usecase   iPromptUsecase
	transform *helper.Transform
}

func NewPromptController(usecase iPromptUsecase) *promptController {
	return &promptController{
		usecase:   usecase,
		transform: helper.NewTransform(),
	}
}

func (c *promptController) CreateTemplate(ctx context.Context, req *pb.CreateTemplateRequest) (*pb.CreateTemplateResponse, error) {
	// Extract description and tags from metadata if possible, or leave empty
	description := req.Payload.Metadata["description"]
	// Tags handling could be added here if passed in metadata

	payload := &entities.CreateTemplatePayload{
		Name:        req.Payload.Name,
		Description: description,
		Content:     req.Payload.Template,
		Variables:   c.transform.Pb2Variable(req.Payload.Variables),
		// Tags:        req.Payload.Tags,
	}

	template, err := c.usecase.CreateTemplate(ctx, payload)
	if err != nil {
		return &pb.CreateTemplateResponse{
			Metadata: req.Metadata,
			Result:   &pb.Result{Code: pb.ResultCode(err.GetCode()), Message: err.Error()},
		}, nil
	}

	return &pb.CreateTemplateResponse{
		Metadata: req.Metadata,
		Result:   &pb.Result{Code: pb.ResultCode_SUCCESS, Message: "created successfully"},
		Template: c.transform.Template2Pb(template),
	}, nil
}

func (c *promptController) GetTemplate(ctx context.Context, req *pb.GetTemplateRequest) (*pb.GetTemplateResponse, error) {
	template, err := c.usecase.GetTemplate(ctx, req.Id)
	if err != nil {
		return &pb.GetTemplateResponse{
			Metadata: req.Metadata,
			Result:   &pb.Result{Code: pb.ResultCode(err.GetCode()), Message: err.Error()},
		}, nil
	}
	return &pb.GetTemplateResponse{
		Metadata: req.Metadata,
		Result:   &pb.Result{Code: pb.ResultCode_SUCCESS},
		Template: c.transform.Template2Pb(template),
	}, nil
}

func (c *promptController) ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error) {
	filter := &entities.TemplateFilter{
		Status: entities.TemplateStatus(req.Status.String()),
		// Tags:     req.Tags,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	templates, total, err := c.usecase.ListTemplates(ctx, filter)
	if err != nil {
		return &pb.ListTemplatesResponse{
			Metadata: req.Metadata,
			Result:   &pb.Result{Code: pb.ResultCode(err.GetCode()), Message: err.Error()},
		}, nil
	}

	pbTemplates := make([]*pb.PromptTemplate, len(templates))
	for i, t := range templates {
		pbTemplates[i] = c.transform.Template2Pb(t)
	}

	return &pb.ListTemplatesResponse{
		Metadata:  req.Metadata,
		Result:    &pb.Result{Code: pb.ResultCode_SUCCESS},
		Templates: pbTemplates,
		Total:     int32(total),
	}, nil
}

func (c *promptController) UpdateTemplate(ctx context.Context, req *pb.UpdateTemplateRequest) (*pb.UpdateTemplateResponse, error) {
	payload := &entities.UpdateTemplatePayload{
		ID:        req.Payload.Id,
		Content:   req.Payload.Template,
		Variables: c.transform.Pb2Variable(req.Payload.Variables),
		Status:    entities.TemplateStatus(req.Payload.Status.String()),
	}

	template, err := c.usecase.UpdateTemplate(ctx, payload)
	if err != nil {
		return &pb.UpdateTemplateResponse{
			Metadata: req.Metadata,
			Result:   &pb.Result{Code: pb.ResultCode(err.GetCode()), Message: err.Error()},
		}, nil
	}

	return &pb.UpdateTemplateResponse{
		Metadata: req.Metadata,
		Result:   &pb.Result{Code: pb.ResultCode_SUCCESS},
		Model:    c.transform.Template2Pb(template),
	}, nil
}

func (c *promptController) DeleteTemplate(ctx context.Context, req *pb.DeleteTemplateRequest) (*pb.ResponseEmpty, error) {
	_ = c.usecase.DeleteTemplate(ctx, req.Id)
	return &pb.ResponseEmpty{
		Result: &pb.Result{Code: pb.ResultCode_SUCCESS},
	}, nil
}

func (c *promptController) RenderTemplate(ctx context.Context, req *pb.RenderTemplateRequest) (*pb.RenderTemplateResponse, error) {
	rendered, err := c.usecase.RenderTemplate(ctx, req.Payload.TemplateId, req.Payload.Variables)
	if err != nil {
		return &pb.RenderTemplateResponse{
			Metadata: req.Metadata,
			Result:   &pb.Result{Code: pb.ResultCode(err.GetCode()), Message: err.Error()},
		}, nil
	}
	return &pb.RenderTemplateResponse{
		Metadata: req.Metadata,
		Result:   &pb.Result{Code: pb.ResultCode_SUCCESS},
		Rendered: &pb.RenderedPrompt{
			TemplateId:    req.Payload.TemplateId,
			RenderedText:  rendered.Content,
			VariablesUsed: rendered.Variables,
		},
	}, nil
}
