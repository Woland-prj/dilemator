package middleware

import (
	"errors"
	"net/http"

	"github.com/Woland-prj/dilemator/internal/router/managers"
	"github.com/Woland-prj/dilemator/internal/router/responses"
	"github.com/gofiber/fiber/v2"
)

const AuthContextKey = "requester"

type authHandler struct {
	cookieManager *managers.CookieManager
}

func (h *authHandler) Handle(c *fiber.Ctx) error {
	userDetails, err := h.cookieManager.VerifyCookie(c)
	if err != nil {
		if errors.Is(err, managers.ErrNoCookiePresent) || errors.Is(err, managers.ErrInvalidTokenFormat) {
			c.Set("HX-Redirect", "/login")

			return responses.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		}

		if errors.Is(err, managers.ErrInvalidSession) {
			c.Set("HX-Redirect", "/login")

			return responses.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		}

		return responses.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal Server Error")
	}

	c.Locals(AuthContextKey, userDetails)

	return c.Next()
}

// WithAuth creates middleware to auth check.
func WithAuth(cookieManager *managers.CookieManager) fiber.Handler {
	handler := &authHandler{cookieManager: cookieManager}

	return func(c *fiber.Ctx) error {
		return handler.Handle(c)
	}
}
