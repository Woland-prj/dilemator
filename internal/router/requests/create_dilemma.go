package requests

import (
	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/google/uuid"
)

type CreateDilemma struct {
	Topic     string `json:"topic" validate:"required,min=1,max=256" example:"Tests in software"`
	RootName  string `json:"name" validate:"required,min=1" example:"Tests"`
	RootValue string `json:"value" validate:"required,min=1" example:"What should be if Ivan don't test his program?'"`
}

func (req *CreateDilemma) ToModel(ownerID uuid.UUID) *dilemma_dto.CreateDilemmaDto {
	return &dilemma_dto.CreateDilemmaDto{
		OwnerID:   ownerID,
		Topic:     req.Topic,
		RootName:  req.RootName,
		RootValue: req.RootValue,
	}
}
