package sessions_repo

import (
	"context"
	"errors"
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/security_errors"
	pentity "github.com/Woland-prj/dilemator/internal/repo/sessions_repo/entity"
	"github.com/Woland-prj/dilemator/internal/services/sessions_service"
	"github.com/Woland-prj/dilemator/pkg/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepositoryAdapter struct {
	*postgres.Postgres
}

var _ sessions_service.SessionRepositoryPort = (*SessionRepositoryAdapter)(nil)

func NewSessionRepositoryAdapter(pg *postgres.Postgres) *SessionRepositoryAdapter {
	return &SessionRepositoryAdapter{pg}
}

func (s *SessionRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*security_entity.Session, error) {
	const op = "repo - SessionRepositoryAdapter - FindByID"

	var sessionEntity pentity.SessionEntity
	if err := s.DB.WithContext(ctx).First(&sessionEntity, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, security_errors.ErrSessionNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return sessionEntity.ToModel(), nil
}

func (s *SessionRepositoryAdapter) FindByUserID(ctx context.Context, userID uuid.UUID) (*security_entity.Session, error) {
	const op = "repo - SessionRepositoryAdapter - FindByUserID"

	var sessionEntity pentity.SessionEntity
	if err := s.DB.WithContext(ctx).First(&sessionEntity, "user_id = ?", userID.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, security_errors.ErrSessionNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return sessionEntity.ToModel(), nil
}

func (s *SessionRepositoryAdapter) Save(ctx context.Context, session *security_entity.Session) error {
	const op = "repo - SessionRepositoryAdapter - Save"

	sesEn := pentity.SessionEntityFromModel(session)

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(sesEn).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return berrors.FromErr(op, security_errors.ErrSessionAlreadyExists)
			}

			return berrors.InternalFromErr(op, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionRepositoryAdapter) FindByTokenHash(ctx context.Context, tokenHash string) (*security_entity.Session, error) {
	const op = "repo - SessionRepositoryAdapter - FindByTokenHash"

	var sessionEntity pentity.SessionEntity

	err := s.DB.WithContext(ctx).
		Where("session_token = ? AND is_revoked = false", tokenHash).
		Where("expires_at > ?", time.Now().UTC()).
		First(&sessionEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, security_errors.ErrSessionNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return sessionEntity.ToModel(), nil
}

func (s *SessionRepositoryAdapter) UpdateLastUsedAt(ctx context.Context, id uuid.UUID, lastUsedAt time.Time) error {
	const op = "repo - SessionRepositoryAdapter - UpdateLastUsedAt"

	result := s.DB.WithContext(ctx).
		Model(&pentity.SessionEntity{}).
		Where("id = ?", id).
		Update("last_used_at", lastUsedAt)

	if result.Error != nil {
		return berrors.InternalFromErr(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return berrors.FromErr(op, security_errors.ErrSessionNotFound)
	}

	return nil
}

func (s *SessionRepositoryAdapter) UpdateExpiresAt(ctx context.Context, id uuid.UUID, expiresAt time.Time) error {
	const op = "repo - SessionRepositoryAdapter - UpdateExpiresAt"

	result := s.DB.WithContext(ctx).
		Model(&pentity.SessionEntity{}).
		Where("id = ?", id).
		Update("expires_at", expiresAt)

	if result.Error != nil {
		return berrors.InternalFromErr(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return berrors.FromErr(op, security_errors.ErrSessionNotFound)
	}

	return nil
}

func (s *SessionRepositoryAdapter) UpdateRevokeStatus(ctx context.Context, id uuid.UUID, revoked bool) error {
	const op = "repo - SessionRepositoryAdapter - UpdateRevokeStatus"

	result := s.DB.WithContext(ctx).
		Model(&pentity.SessionEntity{}).
		Where("id = ?", id).
		Update("is_revoked", revoked)

	if result.Error != nil {
		return berrors.InternalFromErr(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return berrors.FromErr(op, security_errors.ErrSessionNotFound)
	}

	return nil
}

func (s *SessionRepositoryAdapter) FindActiveByUserAgentAndIP(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*security_entity.Session, error) {
	const op = "repo - SessionRepositoryAdapter - FindActiveByUserAgentAndIP"

	var sessionEntity pentity.SessionEntity

	err := s.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("user_agent = ?", userAgent).
		Where("ip_address = ?", ipAddress).
		Where("is_revoked = false").
		Where("expires_at > ?", time.Now().UTC()).
		First(&sessionEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, security_errors.ErrSessionNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return sessionEntity.ToModel(), nil
}

func (s *SessionRepositoryAdapter) Update(ctx context.Context, session *security_entity.Session) error {
	const op = "repo - SessionRepositoryAdapter - Update"

	ent := pentity.SessionEntityFromModel(session)

	result := s.DB.WithContext(ctx).
		Model(&ent).
		Select("session_token", "expires_at", "last_used_at").
		Updates(ent)

	if result.Error != nil {
		return berrors.InternalFromErr(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return berrors.FromErr(op, security_errors.ErrSessionNotFound)
	}

	return nil
}
