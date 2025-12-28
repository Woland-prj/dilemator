package responses

import (
	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

type ProfileResponse struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
	Avatar  *string `json:"avatar"`
}

type UserResponse struct {
	ID      uuid.UUID        `json:"id"`
	Email   *string          `json:"email"`
	TgID    *int64           `json:"tgId"`
	Profile *ProfileResponse `json:"profile"`
}

func NewUserResponse(user *user_entity.User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		TgID:  user.TgID,
		Profile: &ProfileResponse{
			Name:    user.Profile.Name,
			Surname: user.Profile.Surname,
			Avatar:  user.Profile.Avatar,
		},
	}
}
