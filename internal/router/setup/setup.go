package router_setup

import (
	"net/http"

	"github.com/Woland-prj/dilemator/config"
	// swagger docs import.
	_ "github.com/Woland-prj/dilemator/docs"
	"github.com/Woland-prj/dilemator/internal/router/middleware"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Giashka Go API
// @description API for dilemator.wwoland.ru
// @version     1.0
// @host        dilemator.wwoland.ru
// @BasePath    /.
func NewRouter(app *fiber.App, cfg *config.Config, l logger.Interface) (apiGroup, componentsGroup fiber.Router) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(middleware.MultipartFormData())
	app.Use(middleware.Cors(cfg))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("giashka-api")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Static assets
	app.Use(cfg.HTTP.AssetsPrefix, filesystem.New(filesystem.Config{
		Root: http.Dir(cfg.HTTP.AssetsDir),
	}))

	apiGroup = app.Group("/api")
	componentsGroup = app.Group("/components")

	return apiGroup, componentsGroup
}
