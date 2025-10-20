package requests

import (
	"github.com/Woland-prj/dilemator/internal/domain/dto/users_dto"
)

type Register struct {
	Email    string  `json:"email"      validate:"required,email"         example:"example@mail.ru"`
	Password string  `json:"password"   validate:"required,min=1,max=50"  example:"12345"`
	Name     *string `json:"name"       validate:"omitempty,min=1,max=80" example:"Ivan"`
	Surname  *string `json:"surname"    validate:"omitempty,min=1,max=80" example:"Berezin"`
}

func (req *Register) ToModel() *users_dto.RegisterDto {
	return &users_dto.RegisterDto{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Surname:  req.Surname,
	}
}
