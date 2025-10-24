package sessions_router

import (
	"errors"
	"net/http"

	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/security_errors"
	"github.com/Woland-prj/dilemator/internal/router/managers"
	"github.com/Woland-prj/dilemator/internal/router/requests"
	"github.com/Woland-prj/dilemator/internal/router/responses"
	"github.com/Woland-prj/dilemator/internal/services/factory"
	"github.com/Woland-prj/dilemator/internal/services/sessions_service"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Register - register security routes for fiber app router.
func Register(
	router fiber.Router,
	f *factory.ServiceFactory,
	cm *managers.CookieManager,
	l logger.Interface,
) error {
	const op = "auth - Register"

	s, err := factory.InstantiateService[sessions_service.SessionService](f)
	if err != nil {
		return berrors.FromErr(op, err)
	}

	c := &SessionController{
		s:  s,
		l:  l,
		v:  validator.New(validator.WithRequiredStructEnabled()),
		cm: cm,
	}

	authGroup := router.Group("/auth")
	{
		authGroup.Post("/login/tg", c.tgLogin)
		authGroup.Post("/logout", c.logout)
	}

	return nil
}

type SessionController struct {
	s  sessions_service.SessionService
	l  logger.Interface
	v  *validator.Validate
	cm *managers.CookieManager
}

// @Summary     User login via telegram widget API
// @Description Authenticates user via telegram widget
// @ID          loginByTelegram
// @Tags        auth-controller
// @Accept      json
// @Produce     json
// @Param       credentials body requests.TgLogin true "Tg auth data"
// @Success     204 "Successfully authenticated, session cookie set"
// @Failure     400 {object} responses.Error "Invalid request body or validation failed"
// @Failure     401 {object} responses.Error "User with this id not found"
// @Failure     403 {object} responses.Error "Data not from telegram or tg login expired"
// @Failure     500 {object} responses.Error "Internal server error"
// @Router      /auth/login/tg [post].
func (c *SessionController) tgLogin(ctx *fiber.Ctx) error {
	const op = "http - UserController - register"

	var body requests.TgLogin

	if err := ctx.BodyParser(&body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	if err := validate(body, c.v, ctx); err != nil {
		return err
	}

	ip := ctx.IP()
	userAgent := ctx.Get("User-Agent")

	session, err := c.s.TgLogin(ctx.Context(), body.ToModel(), userAgent, ip)
	if err != nil {
		return c.handleSessionErrors(ctx, op, err)
	}

	cookie := c.cm.GetCookie(session.SessionToken)
	ctx.Cookie(cookie)

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     User logout
// @Description Revokes current session
// @ID          logoutUser
// @Tags        auth-controller
// @Produce     json
// @Success     204 "Successfully logged out, session cookie cleared"
// @Failure     500 {object} responses.Error "Internal server error"
// @Router      /auth/logout [post].
func (c *SessionController) logout(ctx *fiber.Ctx) error {
	const op = "http - SessionController - logout"

	token := ctx.Cookies("sessionToken")
	if token == "" {
		return ctx.SendStatus(http.StatusNoContent)
	}

	if err := c.s.Logout(ctx.Context(), token); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to revoke session")
	}

	cookie := c.cm.GetClearCookie()
	ctx.Cookie(cookie)

	return ctx.SendStatus(http.StatusNoContent)
}

func (c *SessionController) handleSessionErrors(ctx *fiber.Ctx, op string, err error) error {
	c.l.Debug(berrors.FromErr(op, err).Error())

	if errors.Is(err, security_errors.ErrInvalidCredentials) {
		return responses.ErrorResponse(ctx, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	}

	if errors.Is(err, security_errors.ErrSessionNotFound) {
		return responses.ErrorResponse(ctx, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	}

	if errors.Is(err, security_errors.ErrSessionExpiredOrRevoked) {
		return responses.ErrorResponse(ctx, http.StatusForbidden, "UNAUTHORIZED", err.Error())
	}

	if errors.Is(err, security_errors.ErrSessionAlreadyExists) {
		return responses.ErrorResponse(ctx, http.StatusConflict, "CONFLICT", err.Error())
	}

	if errors.Is(err, security_errors.ErrDataNotFromLoginSource) {
		return responses.ErrorResponse(ctx, http.StatusForbidden, "FORBIDDEN", err.Error())
	}

	if errors.Is(err, security_errors.ErrExternalLoginExpired) {
		return responses.ErrorResponse(ctx, http.StatusConflict, "FORBIDDEN", err.Error())
	}

	return responses.ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
}

func validate(s interface{}, v *validator.Validate, ctx *fiber.Ctx) error {
	if err := v.Struct(s); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return responses.ErrorResponseWithDetails(
				ctx, http.StatusBadRequest,
				"BAD_REQUEST",
				"Validation failed",
				responses.ValidationErrorsToDetails(validationErrors),
			)
		}

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	return nil
}
