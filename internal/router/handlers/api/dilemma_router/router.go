package dilemma_router

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
	"github.com/Woland-prj/dilemator/internal/services/dilemma_service"
	"github.com/Woland-prj/dilemator/internal/services/factory"
	"github.com/Woland-prj/dilemator/internal/view/ui"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	defaultPageSize = 10
	defaultPage     = 1
)

var (
	errInvalidPage = errors.New("page must be >= 1")
	errInvalidSize = errors.New("size must be between 1 and 50")
)

// Register - register dilemma routes for fiber app router.
func Register(
	apiRouter fiber.Router,
	componentsRouter fiber.Router,
	f *factory.ServiceFactory,
	cm *managers.CookieManager,
	l logger.Interface,
) error {
	const op = "dilemma_router - Register"

	s, err := factory.InstantiateService[dilemma_service.DilemmaService](f)
	if err != nil {
		return berrors.FromErr(op, err)
	}

	c := &DilemmaController{
		s: s,
		l: l,
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	dilemmaAPIGroup := apiRouter.Group("/dilemma")
	{
		dilemmaAPIGroup.Post("", middleware.WithAuth(cm), c.createDilemma)
	}

	dilemmaComponentsGroup := componentsRouter.Group("/dilemma")
	{
		dilemmaComponentsGroup.Get("/dashboard", middleware.WithAuth(cm), c.dashboard)
	}

	return nil
}

type DilemmaController struct {
	s dilemma_service.DilemmaService
	l logger.Interface
	v *validator.Validate
}

func (c *DilemmaController) createDilemma(ctx *fiber.Ctx) error {
	const op = "http - DilemmaController - createDilemma"

	var body requests.CreateDilemma

	if err := ctx.BodyParser(&body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	if err := c.validate(body, ctx, op); err != nil {
		return err
	}

	requester, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
	}

	_, err := c.s.CreateDilemma(ctx.Context(), body.ToModel(requester.GetID()))
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

func (c *DilemmaController) dashboard(ctx *fiber.Ctx) error {
	requester, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return ui.ErrorBlock(security_errors.ErrSession).Render(ctx.Context(), ctx)
	}

	page, size, err := parsePagination(ctx)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	dilemmas, err := c.s.GetByOwner(ctx.Context(), requester.GetID(), page, size)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	return ui.Dashboard(dilemmas).Render(ctx.Context(), ctx)
}

func parsePagination(ctx *fiber.Ctx) (page, size int, err error) {
	page = ctx.QueryInt("page", defaultPage)
	size = ctx.QueryInt("size", defaultPageSize)

	if page < 1 {
		return 0, 0, errInvalidPage
	}

	if size < 1 || size > 50 {
		return 0, 0, errInvalidSize
	}

	return page, size, err
}

func (c *DilemmaController) validate(s interface{}, ctx *fiber.Ctx, op string) error {
	if err := c.v.Struct(s); err != nil {
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

	return nil
}
