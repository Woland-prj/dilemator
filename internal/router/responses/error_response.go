package responses

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Timestamp string                 `json:"timestamp,omitempty" example:"2025-07-24T09:51:43"`
	Status    int                    `json:"status,omitempty" example:"404"`
	Error     string                 `json:"error,omitempty" example:"NOT_FOUND"`
	Message   string                 `json:"message,omitempty" example:"Something Not Found"`
	Path      string                 `json:"path,omitempty" example:"/api/target/resource"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func ErrorResponse(ctx *fiber.Ctx, code int, err, msg string) error {
	return ctx.Status(code).JSON(Error{
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    code,
		Error:     err,
		Message:   msg,
		Path:      ctx.Path(),
		Details:   nil,
	})
}

func ErrorResponseWithDetails(ctx *fiber.Ctx, code int, err, msg string, details map[string]interface{}) error {
	return ctx.Status(code).JSON(Error{
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    code,
		Error:     err,
		Message:   msg,
		Path:      ctx.Path(),
		Details:   details,
	})
}

func ValidationErrorsToDetails(validationErrors validator.ValidationErrors) map[string]interface{} {
	details := make(map[string]interface{})

	for _, err := range validationErrors {
		field := err.Field()

		if existing, exists := details[field]; exists {
			if errors, ok := existing.([]string); ok {
				details[field] = append(errors, formatValidationError(err))
			}
		} else {
			details[field] = []string{formatValidationError(err)}
		}
	}

	return details
}

func formatValidationError(err validator.FieldError) string {
	if fn, exists := validationMessages[err.Tag()]; exists {
		return fn(err)
	}

	return fmt.Sprintf("failed validation: %s=%s", err.Tag(), err.Param())
}

// validationMessages contains user-friendly error messages for common validator tags.
// It is initialized once at startup and is safe for concurrent use.
// We disable gochecknoglobals because this is a benign global used for static configuration.
//
//nolint:gochecknoglobals // Static lookup table for validation error messages, safe and efficient
var validationMessages = map[string]func(validator.FieldError) string{
	"required": func(err validator.FieldError) string {
		return "field is required"
	},
	"min": func(err validator.FieldError) string {
		if err.Kind() == reflect.String {
			return fmt.Sprintf("must be at least %s characters long", err.Param())
		}

		return fmt.Sprintf("must be at least %s", err.Param())
	},
	"max": func(err validator.FieldError) string {
		if err.Kind() == reflect.String {
			return fmt.Sprintf("must be at most %s characters long", err.Param())
		}

		return fmt.Sprintf("must be at most %s", err.Param())
	},
	"email": func(err validator.FieldError) string {
		return "must be a valid email address"
	},
	"len": func(err validator.FieldError) string {
		return fmt.Sprintf("must be exactly %s characters long", err.Param())
	},
	"oneof": func(err validator.FieldError) string {
		return fmt.Sprintf("must be one of: %s", strings.ReplaceAll(err.Param(), " ", ", "))
	},
	"numeric": func(err validator.FieldError) string {
		return "must be a numeric value"
	},
	"alphanum": func(err validator.FieldError) string {
		return "must contain only alphanumeric characters"
	},
	"eqfield": func(err validator.FieldError) string {
		return fmt.Sprintf("must be equal to %s field", err.Param())
	},
	"nefield": func(err validator.FieldError) string {
		return fmt.Sprintf("must not be equal to %s field", err.Param())
	},
	"gt": func(err validator.FieldError) string {
		return fmt.Sprintf("must be greater than %s", err.Param())
	},
	"lt": func(err validator.FieldError) string {
		return fmt.Sprintf("must be less than %s", err.Param())
	},
	"gte": func(err validator.FieldError) string {
		return fmt.Sprintf("must be greater than or equal to %s", err.Param())
	},
	"lte": func(err validator.FieldError) string {
		return fmt.Sprintf("must be less than or equal to %s", err.Param())
	},
}
