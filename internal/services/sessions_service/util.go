package sessions_service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/dto/security_dto"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/security_errors"
)

const loginExpirationTime = 24 * time.Hour

// checkTelegramAuthorization check data from Telegram Login Widget.
func checkTelegramAuthorization(req *security_dto.TgLoginDto, botToken string) error {
	const op = "sessions - checkTelegramAuthorization"
	// Секретный ключ — SHA-256 хеш токена бота, бинарное значение
	sha := sha256.Sum256([]byte(botToken))
	secretKey := sha[:] // []byte

	// Вычисляем HMAC-SHA-256
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(req.CheckString))
	expectedMAC := mac.Sum(nil)
	expectedHash := hex.EncodeToString(expectedMAC)

	// Сравниваем с переданным hash
	if !hmac.Equal([]byte(expectedHash), []byte(req.Hash)) {
		return berrors.FromErr(op, security_errors.ErrDataNotFromLoginSource)
	}

	// Check data not expired
	if time.Since(req.AuthDate) > loginExpirationTime {
		return berrors.FromErr(op, security_errors.ErrExternalLoginExpired)
	}

	return nil
}
