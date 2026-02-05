package helper

import (
	"github.com/blcvn/backend/services/prompt-service/entities"
	pb "github.com/blcvn/kratos-proto/go/prompt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Transform struct{}

func NewTransform() *Transform { return &Transform{} }

func (t *Transform) Template2Pb(entity *entities.PromptTemplate) *pb.PromptTemplate {
	vars := make([]*pb.Variable, len(entity.Variables))
	for i, v := range entity.Variables {
		vars[i] = &pb.Variable{
			Name:         v.Name,
			Description:  v.Description,
			Type:         v.Type,
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
		}
	}

	return &pb.PromptTemplate{
		Id:   entity.ID,
		Name: entity.Name,
		// Description: entity.Description, // Not in proto
		Version:   entity.Version,
		Template:  entity.Content, // Mapped to Content
		Variables: vars,
		// Tags:        entity.Tags, // Not in proto
		Metadata:  map[string]string{"description": entity.Description}, // Use metadata for desc
		Status:    pb.TemplateStatus(pb.TemplateStatus_value[string(entity.Status)]),
		CreatedAt: timestamppb.New(entity.CreatedAt),
		UpdatedAt: timestamppb.New(entity.UpdatedAt),
	}
}

func (t *Transform) Pb2Variable(pbVars []*pb.Variable) []entities.Variable {
	vars := make([]entities.Variable, len(pbVars))
	for i, v := range pbVars {
		vars[i] = entities.Variable{
			Name:         v.Name,
			Description:  v.Description,
			Type:         v.Type,
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
		}
	}
	return vars
}
