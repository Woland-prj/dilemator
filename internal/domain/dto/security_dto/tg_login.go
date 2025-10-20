package security_dto

import (
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/dto/users_dto"
)

type TgLoginDto struct {
	TgID        int64
	Name        string
	Surname     string
	Username    string
	Avatar      string
	AuthDate    time.Time
	Hash        string
	CheckString string
}

func (req *TgLoginDto) ToRegRequest() *users_dto.TgRegisterDto {
	return &users_dto.TgRegisterDto{
		TgID:    req.TgID,
		Name:    req.Name,
		Surname: req.Surname,
		Avatar:  req.Avatar,
	}
}
