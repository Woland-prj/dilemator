package pages_router

import (
	"github.com/Woland-prj/dilemator/internal/view/data"
	"github.com/Woland-prj/dilemator/internal/view/pages"
	"github.com/Woland-prj/dilemator/internal/view/ui"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/gofiber/fiber/v2"
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
