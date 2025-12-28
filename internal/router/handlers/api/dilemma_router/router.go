package dilemma_router

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/Woland-prj/dilemator/internal/domain/entity/security_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/dilemma_errors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/security_errors"
	"github.com/Woland-prj/dilemator/internal/router/managers"
	"github.com/Woland-prj/dilemator/internal/router/middleware"
	"github.com/Woland-prj/dilemator/internal/router/requests"
	"github.com/Woland-prj/dilemator/internal/router/responses"
	"github.com/Woland-prj/dilemator/internal/services/dilemma_service"
	"github.com/Woland-prj/dilemator/internal/services/factory"
	"github.com/Woland-prj/dilemator/internal/view/ui"
	"github.com/Woland-prj/dilemator/internal/view/ui/nodeeditor"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		dilemmaAPIGroup.Get("", middleware.WithAuth(cm), c.getDilemmas)
		dilemmaAPIGroup.Post("", middleware.WithAuth(cm), c.createDilemma)
		dilemmaAPIGroup.Put("", middleware.WithAuth(cm), c.updateDilemma)
		dilemmaAPIGroup.Post("/node", middleware.WithAuth(cm), c.createDilemmaNode)
		dilemmaAPIGroup.Put("/node", middleware.WithAuth(cm), c.updateDilemmaNode)
	}

	dilemmaComponentsGroup := componentsRouter.Group("/dilemma")
	{
		dilemmaComponentsGroup.Get("/dashboard", middleware.WithAuth(cm), c.dashboard)
		dilemmaComponentsGroup.Get("/editor", middleware.WithAuth(cm), c.editor)
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

	requester, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
	}

	ctx.Request().WriteTo(os.Stdout)

	var body requests.CreateDilemma

	if err := ctx.BodyParser(&body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	if err := c.validate(body, ctx, op); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())
		return err
	}

	img, err := readFile(ctx, "image")
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())
		return err
	}

	c.l.Debug("dilema to save:", slog.Any("dilemma", body.ToModel(requester.GetID(), img)))

	dilemma, err := c.s.CreateDilemma(ctx.Context(), body.ToModel(requester.GetID(), img))
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		if errors.Is(err, dilemma_errors.ErrDilemmaAlreadyExists) {
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

	return nodeeditor.EditorContainer(*dilemma, *dilemma.RootNode, false).Render(ctx.Context(), ctx)
}

func (c *DilemmaController) updateDilemma(ctx *fiber.Ctx) error {
	const op = "http - DilemmaController - updateDilemma"

	_, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
	}

	var body requests.UpdateDilemma

	if err := ctx.BodyParser(&body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	if err := c.validate(body, ctx, op); err != nil {
		return err
	}

	img, err := readFile(ctx, "image")
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())
		return err
	}

	didStr := ctx.Query("did")

	if didStr == "" {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "did not set")
	}

	did, err := uuid.Parse(didStr)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	dilemma, err := c.s.UpdateDilemma(ctx.Context(), body.ToModel(did, img))
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		if errors.Is(err, dilemma_errors.ErrDilemmaNotFound) {
			return responses.ErrorResponse(
				ctx,
				http.StatusNotFound,
				"NOT_FOUND",
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

	return nodeeditor.EditorContainer(*dilemma, *dilemma.RootNode, false).Render(ctx.Context(), ctx)
}

func (c *DilemmaController) createDilemmaNode(ctx *fiber.Ctx) error {
	const op = "http - DilemmaController - createDilemmaNode"

	_, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
	}

	didStr := ctx.Query("did")
	pidStr := ctx.Query("pid")

	if didStr == "" || pidStr == "" {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "empty params")
	}

	did, err := uuid.Parse(didStr)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	pid, err := uuid.Parse(pidStr)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	var body requests.CreateNode

	if err := ctx.BodyParser(&body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	if err := c.validate(body, ctx, op); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())
		return err
	}

	img, err := readFile(ctx, "image")
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())
		return err
	}

	dilemma, err := c.s.GetByID(ctx.Context(), did)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	node, err := c.s.CreateDilemmaNode(ctx.Context(), body.ToModel(pid, img))
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		if errors.Is(err, dilemma_errors.ErrNodeAlreadyExists) {
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

	return nodeeditor.EditorContainer(*dilemma, *node, false).Render(ctx.Context(), ctx)
}

func (c *DilemmaController) updateDilemmaNode(ctx *fiber.Ctx) error {
	const op = "http - DilemmaController - updateDilemmaNode"

	_, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
	}

	didStr := ctx.Query("did")
	nidStr := ctx.Query("nid")

	if didStr == "" {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "empty params")
	}

	if nidStr == "" {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "empty params")
	}

	did, err := uuid.Parse(didStr)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	nid, err := uuid.Parse(nidStr)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	var body requests.UpdateNode

	if err := ctx.BodyParser(&body); err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
	}

	if err := c.validate(body, ctx, op); err != nil {
		return err
	}

	img, err := readFile(ctx, "image")
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())
		return err
	}

	dilemma, err := c.s.GetByID(ctx.Context(), did)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	node, err := c.s.UpdateDilemmaNode(ctx.Context(), *body.ToModel(nid, img))
	if err != nil {
		c.l.Debug(berrors.FromErr(op, err).Error())

		return responses.ErrorResponse(
			ctx,
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
		)
	}

	return nodeeditor.EditorContainer(*dilemma, *node, false).Render(ctx.Context(), ctx)
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

func (c *DilemmaController) getDilemmas(ctx *fiber.Ctx) error {
	requester, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return responses.ErrorResponse(ctx, http.StatusInternalServerError, "FORBIDDEN", "unknown requester")
	}

	page, size, err := parsePagination(ctx)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	}

	dilemmas, err := c.s.GetByOwner(ctx.Context(), requester.GetID(), page, size)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	return ctx.Status(http.StatusOK).JSON(responses.NewDilemmaResponseList(dilemmas))
}

func (c *DilemmaController) editor(ctx *fiber.Ctx) error {
	_, ok := ctx.Locals(middleware.AuthContextKey).(*security_entity.UserDetails)
	if !ok {
		return ui.ErrorBlock(security_errors.ErrSession).Render(ctx.Context(), ctx)
	}

	didStr := ctx.Query("did")
	nidStr := ctx.Query("nid")
	pidStr := ctx.Query("pid")

	if didStr == "" {
		return nodeeditor.EditorContainer(*dilemma_entity.NewEmptyDilemma(), *dilemma_entity.NewEmptyNode(uuid.Nil), true).Render(ctx.Context(), ctx)
	}

	did, err := uuid.Parse(didStr)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	dilemma, err := c.s.GetByID(ctx.Context(), did)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	if pidStr != "" {
		pid, err := uuid.Parse(pidStr)
		if err != nil {
			return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
		}

		c.l.Debug("empty node is", slog.Any("node", *dilemma_entity.NewEmptyNode(pid)))

		return nodeeditor.EditorContainer(*dilemma, *dilemma_entity.NewEmptyNode(pid), true).Render(ctx.Context(), ctx)
	}

	if nidStr == "" {
		return nodeeditor.EditorContainer(*dilemma, *dilemma.RootNode, true).Render(ctx.Context(), ctx)
	}

	nid, err := uuid.Parse(nidStr)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	node, err := c.s.GetNodeByID(ctx.Context(), nid)
	if err != nil {
		return ui.ErrorBlock(err).Render(ctx.Context(), ctx)
	}

	return nodeeditor.EditorContainer(*dilemma, *node, false).Render(ctx.Context(), ctx)
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

func readFile(c *fiber.Ctx, fieldname string) (*dilemma_dto.FileDto, error) {
	file, err := c.FormFile(fieldname)
	if err != nil {
		return nil, responses.ErrorResponse(
			c, http.StatusBadRequest,
			"BAD_REQUEST", "no file present")
	}

	openedFile, err := file.Open()
	if err != nil {
		return nil, responses.ErrorResponse(
			c, http.StatusUnprocessableEntity,
			"UNPROCESSABLE_ENTITY", "error opening file")
	}
	defer openedFile.Close()

	data, err := io.ReadAll(openedFile)
	if err != nil {
		return nil, responses.ErrorResponse(
			c, http.StatusUnprocessableEntity,
			"UNPROCESSABLE_ENTITY", "error reading file")
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	return &dilemma_dto.FileDto{
		Data:        data,
		ContentType: contentType,
	}, nil
}
