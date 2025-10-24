package users_service

import (
	"context"

	"github.com/Woland-prj/dilemator/internal/domain/dto/users_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock_service.go -package users_service . UserService

type UserService interface {
	Register(ctx context.Context, req *users_dto.RegisterDto) (*user_entity.User, error)
	TgRegister(ctx context.Context, req *users_dto.TgRegisterDto) (*user_entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user_entity.User, error)
	GetByTgID(ctx context.Context, id int64) (*user_entity.User, error)
	GetByEmail(ctx context.Context, email string) (*user_entity.User, error)
}
