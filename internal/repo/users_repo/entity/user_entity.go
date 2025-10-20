package pentity

import (
	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

const (
	UserTableName             = "users"
	UserRoleRelationTableName = "user_roles"
)

// UserEntity -.
type UserEntity struct {
	ID       uuid.UUID     `gorm:"primaryKey;column:id"`
	Email    *string       `gorm:"column:email;index:,unique"`
	Password *string       `gorm:"column:password"`
	TgID     *int64        `gorm:"column:tg_id"`
	Profile  ProfileEntity `gorm:"foreignKey:UserID;references:ID"`
}

func (*UserEntity) TableName() string {
	return "users"
}

func (u *UserEntity) ToModel() *user_entity.User {
	return &user_entity.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		TgID:     u.TgID,
		Profile:  u.Profile.ToModel(),
	}
}

func UserEntityFromModel(user *user_entity.User) *UserEntity {
	return &UserEntity{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		TgID:     user.TgID,
		Profile:  ProfileEntityFromModel(user.ID, user.Profile),
	}
}
