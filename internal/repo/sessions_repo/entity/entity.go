package entity

import (
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/google/uuid"
)

const SessionTableName = "sessions"

type SessionEntity struct {
	ID           uuid.UUID `gorm:"primaryKey;column:id;type:uuid"`
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	SessionToken string    `gorm:"column:session_token;type:text;not null;uniqueIndex"`
	UserAgent    *string   `gorm:"column:user_agent;type:text"`
	IPAddress    *string   `gorm:"column:ip_address;type:inet"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null;index"`
	LastUsedAt   time.Time `gorm:"column:last_used_at;not null"`
	IsRevoked    bool      `gorm:"column:is_revoked;not null;default:false;index"`
	DeviceName   *string   `gorm:"column:device_name;type:text"`
}

func (*SessionEntity) TableName() string {
	return "sessions"
}

func (e *SessionEntity) ToModel() *security_entity.Session {
	return &security_entity.Session{
		ID:           e.ID,
		UserID:       e.UserID,
		SessionToken: e.SessionToken,
		UserAgent:    e.UserAgent,
		IPAddress:    e.IPAddress,
		CreatedAt:    e.CreatedAt,
		ExpiresAt:    e.ExpiresAt,
		LastUsedAt:   e.LastUsedAt,
		IsRevoked:    e.IsRevoked,
		DeviceName:   e.DeviceName,
	}
}

func SessionEntityFromModel(session *security_entity.Session) *SessionEntity {
	return &SessionEntity{
		ID:           session.ID,
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		UserAgent:    session.UserAgent,
		IPAddress:    session.IPAddress,
		CreatedAt:    session.CreatedAt,
		ExpiresAt:    session.ExpiresAt,
		LastUsedAt:   session.LastUsedAt,
		IsRevoked:    session.IsRevoked,
		DeviceName:   session.DeviceName,
	}
}
