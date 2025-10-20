package security_entity

import (
	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

type UserDetails struct {
	usr *user_entity.User
}

func NewUserDetails(usr *user_entity.User) *UserDetails {
	return &UserDetails{usr: usr}
}

func (d *UserDetails) GetUsername() *string {
	return d.usr.Email
}

func (d *UserDetails) GetID() uuid.UUID {
	return d.usr.ID
}
