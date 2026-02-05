package controllers

import (
	"context"

	"github.com/blcvn/backend/services/prompt-service/entities"
	"github.com/blcvn/backend/services/prompt-service/usecases"
	"github.com/blcvn/kratos-proto/go/prompt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PromptController struct {
	prompt.UnimplementedPromptServiceServer
	usecase *usecases.PromptUsecase
}

func NewPromptController(usecase *usecases.PromptUsecase) *PromptController {
	return &PromptController{usecase: usecase}
}

func (c *PromptController) CreateTemplate(ctx context.Context, req *prompt.CreateTemplateRequest) (*prompt.CreateTemplateResponse, error) {
	if req.Payload == nil {
		return nil, nil // Return proper error
	}

	// Convert proto variables to entity variables
	var vars []entities.Variable
	for _, v := range req.Payload.Variables {
		vars = append(vars, entities.Variable{
			Name:         v.Name,
			Type:         v.Type,
			Description:  v.Description,
			DefaultValue: v.DefaultValue,
			Required:     v.Required,
		})
	}

	payload := &entities.CreateTemplatePayload{
		Name:        req.Payload.Name,
		Version:     req.Payload.Version,
		Content:     req.Payload.Template, // Mapping 'Template' field to 'Content'
		Environment: int32(req.Payload.Environment),
		Variables:   vars,
		Metadata:    req.Payload.Metadata,
	}

	tmpl, err := c.usecase.CreateTemplate(ctx, payload)
	if err != nil {
		return &prompt.CreateTemplateResponse{
			Result: &prompt.Result{
				Code:    prompt.ResultCode_INTERNAL, // Simplified mapping
				Message: err.Error(),
			},
		}, nil
	}

	return &prompt.CreateTemplateResponse{
		Result: &prompt.Result{
			Code:    prompt.ResultCode_SUCCESS,
			Message: "Success",
		},
		Template: &prompt.PromptTemplate{
			Id:        tmpl.ID.String(),
			Name:      tmpl.Name,
			CreatedAt: timestamppb.New(tmpl.CreatedAt),
			UpdatedAt: timestamppb.New(tmpl.UpdatedAt),
		},
	}, nil
}

func (c *PromptController) GetTemplate(ctx context.Context, req *prompt.GetTemplateRequest) (*prompt.GetTemplateResponse, error) {
	tmpl, err := c.usecase.RetrieveTemplate(ctx, req.Id) // assuming ID lookup
	// If lookup is by Name, usecase needs refactoring or we handle it here.
	// But GetTemplateRequest only has Id now based on prompt.yaml?
	// Wait, proto reconstruction had `string id = 3`.

	if err != nil {
		return &prompt.GetTemplateResponse{
			Result: &prompt.Result{Code: prompt.ResultCode_NOT_FOUND, Message: err.Error()},
		}, nil
	}

	return &prompt.GetTemplateResponse{
		Result: &prompt.Result{Code: prompt.ResultCode_SUCCESS},
		Template: &prompt.PromptTemplate{
			Id:        tmpl.ID.String(),
			Name:      tmpl.Name,
			CreatedAt: timestamppb.New(tmpl.CreatedAt),
			UpdatedAt: timestamppb.New(tmpl.UpdatedAt),
		},
	}, nil
}

func (c *PromptController) RenderTemplate(ctx context.Context, req *prompt.RenderTemplateRequest) (*prompt.RenderTemplateResponse, error) {
	if req.Payload == nil {
		return nil, nil
	}
	rendered, versionUsed, err := c.usecase.RenderPrompt(ctx, req.Payload.TemplateId, "", req.Payload.Variables)
	if err != nil {
		return &prompt.RenderTemplateResponse{
			Result: &prompt.Result{Code: prompt.ResultCode_INTERNAL, Message: err.Error()},
		}, nil
	}

	return &prompt.RenderTemplateResponse{
		Result: &prompt.Result{Code: prompt.ResultCode_SUCCESS},
		Rendered: &prompt.RenderedPrompt{
			TemplateId:    req.Payload.TemplateId,
			RenderedText:  rendered,
			VariablesUsed: req.Payload.Variables,
			VersionUsed:   versionUsed,
		},
	}, nil
}

func (c *PromptController) ListTemplates(ctx context.Context, req *prompt.ListTemplatesRequest) (*prompt.ListTemplatesResponse, error) {
	return nil, nil // Implementation deferred
}

func (c *PromptController) UpdateTemplate(ctx context.Context, req *prompt.UpdateTemplateRequest) (*prompt.UpdateTemplateResponse, error) {
	return nil, nil // Implementation deferred
}

func (c *PromptController) DeleteTemplate(ctx context.Context, req *prompt.DeleteTemplateRequest) (*prompt.ResponseEmpty, error) {
	return nil, nil // Implementation deferred
}

func (c *PromptController) CreateExperiment(ctx context.Context, req *prompt.CreateExperimentRequest) (*prompt.CreateExperimentResponse, error) {
	return nil, nil // Implementation deferred
}

func (c *PromptController) GetExperiment(ctx context.Context, req *prompt.GetExperimentRequest) (*prompt.GetExperimentResponse, error) {
	return nil, nil // Implementation deferred
}

func (c *PromptController) CompleteExperiment(ctx context.Context, req *prompt.CompleteExperimentRequest) (*prompt.CompleteExperimentResponse, error) {
	return nil, nil // Implementation deferred
}
