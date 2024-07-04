package fast

import (
	"log/slog"
	"net/http"

	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

// Handler is the interface that links the endpoint to the router
type Handler interface {
	// Register registers the endpoint to the given router
	Register(router fiber.Router, validator validator.Validator)
}

// endpoint is a struct that represents an endpoint
type endpoint[In, Out any] struct {
	path    string
	method  string
	handler func(Context, In) (Out, error)
}

var _ Handler = (*endpoint[any, any])(nil)

func Endpoint[In, Out any]() *endpoint[In, Out] {
	return &endpoint[In, Out]{}
}

// Path sets the path of the endpoint
func (e *endpoint[In, Out]) Path(path string) *endpoint[In, Out] {
	e.path = path
	return e
}

// Method sets the method of the endpoint
func (e *endpoint[In, Out]) Method(method string) *endpoint[In, Out] {
	e.method = method
	return e
}

// Handle sets the handler of the endpoint
func (e *endpoint[In, Out]) Handle(fn func(Context, In) (Out, error)) Handler {
	e.handler = fn
	return e
}

// Register registers the endpoint to the given router
func (e *endpoint[In, Out]) Register(r fiber.Router, v validator.Validator) {
	r.Add(
		e.method,
		e.path,
		func(c *fiber.Ctx) error {
			var input In

			parser := c.QueryParser

			if e.method == http.MethodPost {
				parser = c.BodyParser
			}

			if err := parser(&input); err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			if err := v.ValidateStruct(&input); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(validator.ValidationErrorSerializer{
					Errors: v.Translate(err),
				})
			}

			context := newContext(c)

			output, err := e.handler(context, input)
			if err != nil {
				slog.Error("error in handler %s %s: %w", e.method, e.path, err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			if err := v.ValidateStruct(&output); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(validator.ValidationErrorSerializer{
					Errors: v.Translate(err),
				})
			}

			return c.Status(fiber.StatusOK).JSON(output)
		},
	)
}
