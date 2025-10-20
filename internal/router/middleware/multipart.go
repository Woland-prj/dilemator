package middleware

import (
	"mime"

	"github.com/Woland-prj/dilemator/internal/router/responses"
	"github.com/gofiber/fiber/v2"
)

func MultipartFormData() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type", "")
		if contentType == "" {
			return c.Next()
		}

		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil || mediaType != "multipart/form-data" {
			return c.Next()
		}

		form, err := c.MultipartForm()
		if err != nil {
			return responses.ErrorResponse(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid multipart form data")
		}

		c.Locals("multipartForm", form)

		return c.Next()
	}
}
