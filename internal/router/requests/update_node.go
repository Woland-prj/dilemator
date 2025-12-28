package requests

import (
	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/google/uuid"
)

type UpdateNode struct {
	Name  string `from:"name" json:"name" validate:"required,min=1" example:"Tests"`
	Value string `form:"value" json:"value" validate:"required,min=1" example:"What should be if Ivan don't test his program?'"`
	Image []byte `form:"-" json:"image" validate:"omitempty,min=1" example:"image/base64,oiuFhdDhjjvc..."`
}

func (req *UpdateNode) ToModel(nid uuid.UUID, img *dilemma_dto.FileDto) *dilemma_dto.UpdateDilemmaNodeDto {
	return &dilemma_dto.UpdateDilemmaNodeDto{
		ID:    nid,
		Name:  req.Name,
		Value: req.Value,
		Image: img,
	}
}
