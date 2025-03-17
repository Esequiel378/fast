package fast

import (
	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

// Group is a group of routes
type Group struct {
	router    fiber.Router
	validator validator.Validator
}

// MustRegister registers a handler to the app
// The handler must be a struct with methods that return a Handler
func (g Group) MustRegister(prefix string, handler any) Group {
	mustValidateAndRegisterHandler(handler, g.router.Group(prefix), g.validator)
	return g
}
