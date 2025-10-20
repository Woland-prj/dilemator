package pentity

import (
	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

const ProfileTableName = "user_profile"

type ProfileEntity struct {
	UserID     uuid.UUID `gorm:"primaryKey;column:user_id"`
	Name       *string   `gorm:"column:name"`
	Surname    *string   `gorm:"column:surname"`
	Patronymic *string   `gorm:"column:patronymic"`
	Avatar     *string   `gorm:"column:avatar"`
}

func (*ProfileEntity) TableName() string {
	return "user_profile"
}

func (u *ProfileEntity) ToModel() *user_entity.Profile {
	return &user_entity.Profile{
		Name:    u.Name,
		Surname: u.Surname,
		Avatar:  u.Avatar,
	}
}

func ProfileEntityFromModel(uid uuid.UUID, profile *user_entity.Profile) ProfileEntity {
	return ProfileEntity{
		UserID:  uid,
		Name:    profile.Name,
		Surname: profile.Surname,
		Avatar:  profile.Avatar,
	}
}
