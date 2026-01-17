package pages_router

import (
	"github.com/Woland-prj/dilemator/internal/view/data"
	"github.com/Woland-prj/dilemator/internal/view/pages"
	"github.com/Woland-prj/dilemator/internal/view/ui"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Register - register user routes for fiber app router.
func Register(
	router fiber.Router,
	l logger.Interface,
	botName string,
) error {
	c := &PagesController{
		l:       l,
		botName: botName,
	}

	landingGroup := router.Group("")
	{
		landingGroup.Get("", c.landing)
	}

	platformGroup := router.Group("/platform")
	{
		platformGroup.Get("", c.platform)
	}

	loginGroup := router.Group("/login")
	{
		loginGroup.Get("", c.login)
	}

	editorGroup := router.Group("/editor")
	{
		editorGroup.Get("", c.editor)
	}

	viewerGroup := router.Group("/viewer")
	{
		viewerGroup.Get("", c.viewer)
	}

	return nil
}

type PagesController struct {
	l       logger.Interface
	botName string
}

func (c PagesController) landing(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	return pages.Landing(ui.MenuProps{
		Links: data.LandingMenuLinks(),
	}).Render(ctx.Context(), ctx)
}

func (c PagesController) platform(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	return pages.Platform(ui.MenuProps{
		Links: data.PlatformMenuLinks(),
	}).Render(ctx.Context(), ctx)
}

func (c PagesController) login(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	return pages.Login(c.botName).Render(ctx.Context(), ctx)
}

func (c PagesController) editor(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	didStr := ctx.Query("did")

	if didStr == "" {
		return pages.Editor(ui.MenuProps{
			Links: data.EditorMenuLinks(),
		}, "").Render(ctx.Context(), ctx)
	}

	_, err := uuid.Parse(didStr)
	if err != nil {
		return pages.Editor(ui.MenuProps{
			Links: data.EditorMenuLinks(),
		}, "").Render(ctx.Context(), ctx)
	}

	return pages.Editor(ui.MenuProps{
		Links: data.EditorMenuLinks(),
	}, didStr).Render(ctx.Context(), ctx)
}

func (c PagesController) viewer(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	didStr := ctx.Query("did")

	if didStr == "" {
		return pages.Viewer(ui.MenuProps{
			Links: data.EditorMenuLinks(),
		}, "").Render(ctx.Context(), ctx)
	}

	_, err := uuid.Parse(didStr)
	if err != nil {
		return pages.Viewer(ui.MenuProps{
			Links: data.EditorMenuLinks(),
		}, "").Render(ctx.Context(), ctx)
	}

	return pages.Viewer(ui.MenuProps{
		Links: data.EditorMenuLinks(),
	}, didStr).Render(ctx.Context(), ctx)
}
