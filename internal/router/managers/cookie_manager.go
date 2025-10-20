package managers

import (
	"errors"
	"time"

	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/Woland-prj/dilemator/internal/services/sessions_service"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrNoCookiePresent    = errors.New("no cookie present")
	ErrInvalidSession     = errors.New("invalid session")
)

type CookieManager struct {
	sessionService sessions_service.SessionService
	log            logger.Interface
	cfg            *CookieConfig
}

type CookieConfig struct {
	CookieName string
	MaxAgeDays int
	Secure     bool
	SameSite   string
}

func NewCookieManager(
	sessionService sessions_service.SessionService,
	log logger.Interface,
	cfg *CookieConfig,
) *CookieManager {
	return &CookieManager{
		sessionService: sessionService,
		log:            log,
		cfg:            cfg,
	}
}

func (cm *CookieManager) GetCookie(sessionToken string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     cm.cfg.CookieName,
		Value:    sessionToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   cm.cfg.Secure,
		MaxAge:   int(time.Duration(cm.cfg.MaxAgeDays) * 24 * time.Hour / time.Second),
		SameSite: cm.cfg.SameSite,
	}
}

func (cm *CookieManager) VerifyCookie(c *fiber.Ctx) (*security_entity.UserDetails, error) {
	sessionToken := c.Cookies(cm.cfg.CookieName)
	if sessionToken == "" {
		return nil, ErrNoCookiePresent
	}

	cm.log.Debug("Session token is: ", sessionToken)

	details, err := cm.sessionService.Verify(c.Context(), sessionToken)
	if err != nil {
		return nil, ErrInvalidSession
	}

	return details, nil
}

func (cm *CookieManager) GetClearCookie() *fiber.Cookie {
	cookie := cm.GetCookie("")
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.MaxAge = -1

	return cookie
}
