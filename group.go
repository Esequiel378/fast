package fast

import (
	"path"

	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

// Group is a group of routes
type Group struct {
	router      fiber.Router
	validator   validator.Validator
	path        string
	apiSchema   *OpenAPIGenerator
	middlewares []Middleware
}

// MustRegister registers a handler to the app
// The handler must be a struct with methods that return a Handler
func (g Group) MustRegister(prefix string, handler any, middlewares ...Middleware) Group {
	mustValidateAndRegisterHandler(
		path.Join(g.path, prefix),
		handler,
		g.router.Group(prefix),
		g.validator,
		g.apiSchema,
		append(g.middlewares, middlewares...)...,
	)
	return g
}
