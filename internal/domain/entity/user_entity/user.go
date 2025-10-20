package user_entity

import "github.com/google/uuid"

// User -.
type User struct {
	ID       uuid.UUID
	Email    *string
	Password *string
	TgID     *int64
	Profile  *Profile
}

func NewUser(id uuid.UUID, email, password *string, tgID *int64, profile *Profile) *User {
	return &User{
		ID:       id,
		Email:    email,
		Password: password,
		TgID:     tgID,
		Profile:  profile,
	}
}
