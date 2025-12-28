package requests

import (
	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/google/uuid"
)

type UpdateDilemma struct {
	Topic     string `form:"topic" json:"topic" validate:"required,min=1,max=256" example:"Tests in software"`
	RootName  string `form:"name" json:"name" validate:"required,min=1" example:"Tests"`
	RootValue string `form:"value" json:"value" validate:"required,min=1" example:"What should be if Ivan don't test his program?'"`
}

func (req *UpdateDilemma) ToModel(did uuid.UUID, file *dilemma_dto.FileDto) *dilemma_dto.UpdateDilemmaDto {
	return &dilemma_dto.UpdateDilemmaDto{
		ID:        did,
		Topic:     req.Topic,
		RootName:  req.RootName,
		RootValue: req.RootValue,
		RootImage: file,
	}
}
