package users_repo

import (
	"context"
	"errors"

	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/user_errors"
	pentity "github.com/Woland-prj/dilemator/internal/repo/users_repo/entity"
	"github.com/Woland-prj/dilemator/internal/services/users_service"
	"github.com/Woland-prj/dilemator/pkg/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryAdapter struct {
	*postgres.Postgres
}

var _ users_service.UserRepositoryPort = (*UserRepositoryAdapter)(nil)

func NewUserRepositoryAdapter(pg *postgres.Postgres) *UserRepositoryAdapter {
	return &UserRepositoryAdapter{
		Postgres: pg,
	}
}

func (u *UserRepositoryAdapter) Save(ctx context.Context, user *user_entity.User) error {
	const op = "repo - persistent - UserRepositoryAdapter - Save"

	uEn := pentity.UserEntityFromModel(user)

	err := u.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(uEn).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return berrors.FromErr(op, user_errors.ErrUserAlreadyExists)
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

func (u *UserRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*user_entity.User, error) {
	const op = "repository - UserRepositoryAdapter - FindByID"

	var userEntity pentity.UserEntity
	if err := u.DB.WithContext(ctx).
		Model(&pentity.UserEntity{}).
		Preload("Profile").
		First(&userEntity, "id = ?", id.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, user_errors.ErrUserNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return userEntity.ToModel(), nil
}

func (u *UserRepositoryAdapter) FindByEmail(ctx context.Context, email string) (*user_entity.User, error) {
	const op = "repository - UserRepositoryAdapter - FindByEmail"

	var userEntity pentity.UserEntity
	if err := u.DB.WithContext(ctx).
		Model(&pentity.UserEntity{}).
		Preload("Profile").
		First(&userEntity, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, user_errors.ErrUserNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return userEntity.ToModel(), nil
}

func (u *UserRepositoryAdapter) FindByTgID(ctx context.Context, tgID int64) (*user_entity.User, error) {
	const op = "repository - UserRepositoryAdapter - FindByTgID"

	var userEntity pentity.UserEntity
	if err := u.DB.WithContext(ctx).
		Model(&pentity.UserEntity{}).
		Preload("Profile").
		First(&userEntity, "tg_id = ?", tgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, user_errors.ErrUserNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return userEntity.ToModel(), nil
}
