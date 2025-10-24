package users_router

import (
	"errors"
	"net/http"

	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/security_errors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/user_errors"
	"github.com/Woland-prj/dilemator/internal/router/managers"
	"github.com/Woland-prj/dilemator/internal/router/middleware"
	"github.com/Woland-prj/dilemator/internal/router/requests"
	"github.com/Woland-prj/dilemator/internal/router/responses"
	"github.com/Woland-prj/dilemator/internal/services/factory"
	"github.com/Woland-prj/dilemator/internal/services/users_service"
	"github.com/Woland-prj/dilemator/internal/view/ui"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Register - register user routes for fiber app router.
func Register(
	apiRouter fiber.Router,
	componentsRouter fiber.Router,
	f *factory.ServiceFactory,
	cm *managers.CookieManager,
	l logger.Interface,
) error {
	const op = "users - Register"

	s, err := factory.InstantiateService[users_service.UserService](f)
	if err != nil {
		return berrors.FromErr(op, err)
	}

	c := &UserController{
		s: s,
		l: l,
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	userAPIGroup := apiRouter.Group("/user")
	{
		userAPIGroup.Post("/register", c.register)
	}

	userComponentsGroup := componentsRouter.Group("/user")
	{
		userComponentsGroup.Get("/profilebadge", middleware.WithAuth(cm), c.profileBadge)
	}

	return nil
}

type UserController struct {
	s users_service.UserService
	l logger.Interface
	v *validator.Validate
}

func (c *UserController) register(ctx *fiber.Ctx) error {
	const op = "http - UserController - register"

	var body requests.Register

	if err := ctx.BodyParser(&body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	if err := c.v.Struct(body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

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

	_, err := c.s.Register(ctx.Context(), body.ToModel())
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		if errors.Is(err, user_errors.ErrUserAlreadyExists) {
			return responses.ErrorResponse(
				ctx,
				http.StatusConflict,
				"CONFLICT",
				err.Error(),
			)
		}

		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
		)
	}

	return ctx.SendStatus(http.StatusCreated)
}

func (c *UserController) profileBadge(ctx *fiber.Ctx) error {
	requester, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return ui.ErrorBlock(security_errors.ErrSession).Render(ctx.Context(), ctx)
	}

	user, err := c.s.GetByID(ctx.Context(), requester.GetID())
	if err != nil {
		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
		)
	}

	return ui.ProfileBadge(*user).Render(ctx.Context(), ctx)
}
