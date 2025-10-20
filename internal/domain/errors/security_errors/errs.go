package security_errors

import "errors"

var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrSessionNotFound         = errors.New("session not found")
	ErrSessionExpiredOrRevoked = errors.New("session is expired or revoked")
	ErrDataNotFromLoginSource  = errors.New("data not from login source")
	ErrExternalLoginExpired    = errors.New("external login is expired")
	ErrSessionAlreadyExists    = errors.New("session already exists")
)
