package sessions_service

import (
	"context"
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks_test.go -package sessions_service_test . SessionRepositoryPort,HashProviderPort

type SessionRepositoryPort interface {
	Save(ctx context.Context, session *security_entity.Session) error
	FindByID(ctx context.Context, id uuid.UUID) (*security_entity.Session, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*security_entity.Session, error)
	UpdateLastUsedAt(ctx context.Context, id uuid.UUID, lastUsedAt time.Time) error
	UpdateExpiresAt(ctx context.Context, id uuid.UUID, expiresAt time.Time) error
	UpdateRevokeStatus(ctx context.Context, id uuid.UUID, revoked bool) error
	FindActiveByUserAgentAndIP(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*security_entity.Session, error)
	Update(ctx context.Context, session *security_entity.Session) error
}

type HashProviderPort interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) (bool, error)
	GenerateRandomToken() string
	Hash(data string) string
}
