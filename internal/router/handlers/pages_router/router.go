package pages_router

import (
	"github.com/Woland-prj/dilemator/internal/view/pages"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// Register - register user routes for fiber app router.
func Register(
	router fiber.Router,
	l logger.Interface,
) error {
	const op = "users - Register"

	c := &PagesController{
		l: l,
	}

	landingGroup := router.Group("")
	{
		landingGroup.Get("", c.landing)
	}

	platformGroup := router.Group("/platform")
	{
		platformGroup.Get("", c.platform)
	}

	return nil
}

type PagesController struct {
	l logger.Interface
}

func (c PagesController) landing(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return pages.Landing().Render(ctx.Context(), ctx)
}

func (c PagesController) platform(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return pages.Platform().Render(ctx.Context(), ctx)
}
