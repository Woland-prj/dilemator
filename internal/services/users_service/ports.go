package users_service

import (
	"context"

	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks_test.go -package users_service_test . UserRepositoryPort,HashProviderPort

type UserRepositoryPort interface {
	Save(ctx context.Context, user *user_entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*user_entity.User, error)
	FindByEmail(ctx context.Context, email string) (*user_entity.User, error)
	FindByTgID(ctx context.Context, id int64) (*user_entity.User, error)
}

type HashProviderPort interface {
	HashPassword(string) (string, error)
	VerifyPassword(string, string) (bool, error)
}
