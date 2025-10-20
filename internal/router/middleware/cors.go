package middleware

import (
	"github.com/Woland-prj/dilemator/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(cfg *config.Config) func(c *fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.HTTP.CorsAllowedOrigins,
		AllowHeaders:     cfg.HTTP.CorsAllowedHeaders,
		AllowCredentials: cfg.HTTP.CorsAllowCredentials,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	})
}
