package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Woland-prj/dilemator/internal/domain/dto/users_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/user_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/user_errors"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/google/uuid"
)

type usersService struct {
	log      logger.Interface
	userRepo UserRepositoryPort
	hash     HashProviderPort
}

var _ UserService = (*usersService)(nil)

func NewUserService(
	log logger.Interface,
	userRepo UserRepositoryPort,
	hash HashProviderPort,
) UserService {
	return &usersService{
		log:      log,
		userRepo: userRepo,
		hash:     hash,
	}
}

func (s *usersService) Register(
	ctx context.Context,
	req *users_dto.RegisterDto,
) (*user_entity.User, error) {
	return s.basicRegister(ctx, req)
}

func (s *usersService) basicRegister(
	ctx context.Context,
	req *users_dto.RegisterDto,
) (*user_entity.User, error) {
	const op = "users - usersService - basicRegister"

	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, berrors.FromErr(op, user_errors.ErrUserAlreadyExists)
	}

	if !errors.Is(err, user_errors.ErrUserNotFound) {
		return nil, berrors.InternalFromErr(op, err)
	}

	hash, err := s.hash.HashPassword(req.Password)
	if err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	profile := user_entity.NewProfile(req.Name, req.Surname, nil)

	user := user_entity.NewUser(uuid.New(), &req.Email, &hash, nil, profile)

	err = s.userRepo.Save(ctx, user)
	if err != nil {
		if errors.Is(err, user_errors.ErrUserAlreadyExists) {
			s.log.Debug(fmt.Sprintf("%s: %s", op, err))

			return nil, berrors.Wrap(
				op,
				fmt.Sprintf("user %s already exists", *user.Email),
				err,
			)
		}

		s.log.Error(fmt.Sprintf("%s: %s", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}

	return user, nil
}

func (s *usersService) TgRegister(
	ctx context.Context,
	req *users_dto.TgRegisterDto,
) (*user_entity.User, error) {
	const op = "users - usersService - TgRegister"

	_, err := s.userRepo.FindByTgID(ctx, req.TgID)
	if err == nil {
		return nil, berrors.Wrap(
			op,
			fmt.Sprintf("User with telegram account %d already exists", req.TgID),
			user_errors.ErrUserNotFound,
		)
	}

	if !errors.Is(err, user_errors.ErrUserNotFound) {
		return nil, berrors.InternalFromErr(op, err)
	}

	profile := user_entity.NewProfile(&req.Name, &req.Surname, &req.Avatar)

	user := user_entity.NewUser(uuid.New(), nil, nil, &req.TgID, profile)

	err = s.userRepo.Save(ctx, user)
	if err != nil {
		if errors.Is(err, user_errors.ErrUserAlreadyExists) {
			s.log.Debug(fmt.Sprintf("%s: %s", op, err))

			return nil, berrors.Wrap(
				op,
				fmt.Sprintf("user %s already exists", *user.Email),
				err,
			)
		}

		s.log.Error(fmt.Sprintf("%s: %s", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}

	return user, nil
}

func (s *usersService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*user_entity.User, error) {
	const op = "users - usersService - GetByID"

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, user_errors.ErrUserNotFound) {
			s.log.Debug(fmt.Sprintf("%s: %s", op, err))

			return nil, berrors.Wrap(
				op,
				fmt.Sprintf("User with id %s is not found", id),
				err,
			)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return user, nil
}

func (s *usersService) GetByTgID(
	ctx context.Context,
	id int64,
) (*user_entity.User, error) {
	const op = "users - usersService - GetByTgID"

	user, err := s.userRepo.FindByTgID(ctx, id)
	if err != nil {
		if errors.Is(err, user_errors.ErrUserNotFound) {
			s.log.Debug(fmt.Sprintf("%s: %s", op, err))

			return nil, berrors.Wrap(
				op,
				fmt.Sprintf("User with tg id %d is not found", id),
				err,
			)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return user, nil
}

func (s *usersService) GetByEmail(
	ctx context.Context,
	email string,
) (*user_entity.User, error) {
	const op = "users - usersService - GetByEmail"

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user_errors.ErrUserNotFound) {
			s.log.Debug(fmt.Sprintf("%s: %s", op, err))

			return nil, berrors.Wrap(
				op,
				fmt.Sprintf("User with email %s is not found", email),
				err,
			)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return user, nil
}
