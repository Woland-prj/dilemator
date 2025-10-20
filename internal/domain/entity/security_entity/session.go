package security_entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	SessionToken string
	UserAgent    *string
	IPAddress    *string
	CreatedAt    time.Time
	ExpiresAt    time.Time
	LastUsedAt   time.Time
	IsRevoked    bool
	DeviceName   *string
}

func NewSession(
	id, userID uuid.UUID,
	sessionToken string,
	expiresAt time.Time,
	userAgent, ipAddress, deviceName *string,
) *Session {
	now := time.Now().UTC()

	return &Session{
		ID:           id,
		UserID:       userID,
		SessionToken: sessionToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		CreatedAt:    now,
		ExpiresAt:    expiresAt,
		LastUsedAt:   now,
		IsRevoked:    false,
		DeviceName:   deviceName,
	}
}

// IsActive returns true if the session is not withdrawn and has not expired.
func (s *Session) IsActive() bool {
	return !s.IsRevoked && time.Now().UTC().Before(s.ExpiresAt)
}

// Touch updates the time of the latest use of the session.
func (s *Session) Touch() {
	s.LastUsedAt = time.Now().UTC()
}

// Revoke отзывает сессию.
func (s *Session) Revoke() {
	s.IsRevoked = true
}

// Extend extends the validity of the session.
func (s *Session) Extend(newExpiresAt time.Time) {
	s.ExpiresAt = newExpiresAt
}
