package requests

import (
	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/google/uuid"
)

type CreateNode struct {
	Name  string `json:"name" validate:"required,min=1" example:"Tests"`
	Value string `json:"value" validate:"required,min=1" example:"What should be if Ivan don't test his program?'"`
}

func (req *CreateNode) ToModel(pid uuid.UUID) *dilemma_dto.CreateDilemmaNodeDto {
	return &dilemma_dto.CreateDilemmaNodeDto{
		ParentID: pid,
		Name:     req.Name,
		Value:    req.Value,
	}
}
