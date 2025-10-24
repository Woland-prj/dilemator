package sessions_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/dto/security_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/security_errors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/user_errors"
	"github.com/Woland-prj/dilemator/internal/services/users_service"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/google/uuid"
)

type sessionService struct {
	log                 logger.Interface
	tgBotToken          string
	sessionRepo         SessionRepositoryPort
	hash                HashProviderPort
	userService         users_service.UserService
	sessionLifetimeDays int
}

var _ SessionService = (*sessionService)(nil)

func NewSessionService(
	log logger.Interface,
	sessionRepo SessionRepositoryPort,
	hash HashProviderPort,
	userService users_service.UserService,
	sessionLifetimeDays int,
	tgBotToken string,
) SessionService {
	return &sessionService{
		log:                 log,
		sessionRepo:         sessionRepo,
		hash:                hash,
		userService:         userService,
		sessionLifetimeDays: sessionLifetimeDays,
		tgBotToken:          tgBotToken,
	}
}

func (s *sessionService) Get(ctx context.Context, id uuid.UUID) (*security_entity.Session, error) {
	const op = "sessions - sessionService - Get"

	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, berrors.Wrap(
			op,
			fmt.Sprintf("Session with id %s is not found", id),
			err,
		)
	}

	return session, nil
}

func (s *sessionService) Login(ctx context.Context, req security_dto.LoginDto, userAgent, ip string) (*security_entity.Session, error) {
	const op = "sessions - sessionService - Login"

	// Authenticate
	usr, err := s.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user_errors.ErrUserNotFound) {
			return nil, berrors.FromErr(op, security_errors.ErrInvalidCredentials)
		}

		return nil, berrors.FromErr(op, berrors.ErrInternal)
	}

	// Check password
	ok, err := s.hash.VerifyPassword(*usr.Password, req.Password)
	if err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	if !ok {
		return nil, berrors.FromErr(op, security_errors.ErrInvalidCredentials)
	}

	// Get session
	session, err := s.createOrUpdateSession(ctx, usr.ID, userAgent, ip)
	if err != nil {
		return nil, berrors.FromErr(op, err)
	}

	return session, nil
}

func (s *sessionService) TgLogin(ctx context.Context, req *security_dto.TgLoginDto, userAgent, ip string) (*security_entity.Session, error) {
	const op = "sessions - sessionService - TgLogin"

	s.log.Info(fmt.Sprintf("checking telegram authorization %+v", req))

	// Verify tg data
	err := checkTelegramAuthorization(req, s.tgBotToken)
	if err != nil {
		s.log.Error(fmt.Sprintf("checking failed: %s", err.Error()))

		return nil, berrors.FromErr(op, err)
	}

	// Authenticate
	usr, err := s.userService.GetByTgID(ctx, req.TgID)
	if err != nil {
		if !errors.Is(err, user_errors.ErrUserNotFound) {
			return nil, berrors.InternalFromErr(op, err)
		}

		usr, err = s.userService.TgRegister(ctx, req.ToRegRequest())
		if err != nil {
			return nil, berrors.FromErr(op, err)
		}
	}

	// Get session
	session, err := s.createOrUpdateSession(ctx, usr.ID, userAgent, ip)
	if err != nil {
		return nil, berrors.FromErr(op, err)
	}

	return session, nil
}

func (s *sessionService) createOrUpdateSession(ctx context.Context, userID uuid.UUID, userAgent, ip string) (*security_entity.Session, error) {
	const op = "sessions - sessionService - createOrUpdateSession"

	existingSession, err := s.sessionRepo.FindActiveByUserAgentAndIP(ctx, userID, userAgent, ip)
	if err != nil {
		if !errors.Is(err, security_errors.ErrSessionNotFound) {
			return nil, berrors.InternalFromErr(op, err)
		}
	}

	if existingSession != nil && existingSession.IsActive() {
		newToken := s.hash.GenerateRandomToken()
		newTokenHash := s.hash.Hash(newToken)

		existingSession.SessionToken = newTokenHash
		existingSession.ExpiresAt = time.Now().UTC().Add(7 * 24 * time.Hour)
		existingSession.LastUsedAt = time.Now().UTC()

		if err = s.sessionRepo.Update(ctx, existingSession); err != nil {
			return nil, berrors.FromErr(op, err)
		}

		existingSession.SessionToken = newToken

		return existingSession, nil
	}

	newToken := s.hash.GenerateRandomToken()
	newTokenHash := s.hash.Hash(newToken)

	session := security_entity.NewSession(
		uuid.New(),
		userID,
		newTokenHash,
		time.Now().UTC().Add(time.Duration(s.sessionLifetimeDays)*24*time.Hour),
		&userAgent,
		&ip,
		nil,
	)

	if err = s.sessionRepo.Save(ctx, session); err != nil {
		return nil, berrors.FromErr(op, err)
	}

	session.SessionToken = newToken

	return session, nil
}

func (s *sessionService) Verify(ctx context.Context, sessionToken string) (*security_entity.UserDetails, error) {
	const op = "sessions - sessionService - Verify"

	tokenHash := s.hash.Hash(sessionToken)

	session, err := s.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, berrors.FromErr(op, security_errors.ErrSessionNotFound)
	}

	if !session.IsActive() {
		return nil, berrors.FromErr(op, security_errors.ErrSessionExpiredOrRevoked)
	}

	session.Touch()

	if err = s.sessionRepo.UpdateLastUsedAt(ctx, session.ID, session.LastUsedAt); err != nil {
		s.log.Warn("failed to update last_used_at: ", err)
	}

	timeLeft := time.Until(session.ExpiresAt)
	if timeLeft < 72*time.Hour {
		session.Extend(time.Now().UTC().Add(time.Duration(s.sessionLifetimeDays) * 24 * time.Hour))
	}

	if err = s.sessionRepo.UpdateExpiresAt(ctx, session.ID, session.ExpiresAt); err != nil {
		s.log.Warn("failed to update expires_at: ", err)
	}

	usr, err := s.userService.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, berrors.FromErr(op, err)
	}

	return security_entity.NewUserDetails(usr), nil
}

func (s *sessionService) Logout(ctx context.Context, sessionToken string) error {
	const op = "sessions - sessionService - Logout"

	tokenHash := s.hash.Hash(sessionToken)

	session, err := s.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil
	}

	session.Revoke()

	if err = s.sessionRepo.UpdateRevokeStatus(ctx, session.ID, true); err != nil {
		return berrors.InternalFromErr(op, err)
	}

	return nil
}
