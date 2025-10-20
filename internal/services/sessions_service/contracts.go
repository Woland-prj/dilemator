package sessions_service

import (
	"context"

	"github.com/Woland-prj/dilemator/internal/domain/dto/security_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock_sessions.go -package security . SessionService

type SessionService interface {
	Get(ctx context.Context, id uuid.UUID) (*security_entity.Session, error)
	Login(ctx context.Context, req security_dto.LoginDto, userAgent, ip string) (*security_entity.Session, error)
	TgLogin(ctx context.Context, req *security_dto.TgLoginDto, userAgent, ip string) (*security_entity.Session, error)
	Verify(ctx context.Context, sessionToken string) (*security_entity.UserDetails, error)
	Logout(ctx context.Context, sessionToken string) error
}
