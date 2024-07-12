package fast

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

type (
	// In is the default input type for an endpoint
	In struct{}
	// Out is the default output type for an endpoint
	Out string
)

// Handler is the interface that links the endpoint to the router
type Handler interface {
	// Register registers the endpoint to the given router
	Register(router fiber.Router, validator validator.Validator)
}

// endpoint is a struct that represents an endpoint
type endpoint[In, Out any] struct {
	path        string
	method      string
	handler     func(Context, In) (Out, error)
	middlewares []func(Context) error
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

// Middlewares sets the middlewares of the endpoint
func (e *endpoint[In, Out]) Middlewares(middlewares ...func(Context) error) *endpoint[In, Out] {
	e.middlewares = middlewares
	return e
}

// Handle sets the handler of the endpoint
func (e *endpoint[In, Out]) Handle(fn func(Context, In) (Out, error)) Handler {
	e.handler = fn
	return e
}

// Register registers the endpoint to the given router
func (e *endpoint[In, Out]) Register(r fiber.Router, v validator.Validator) {
	handlers := make([]fiber.Handler, len(e.middlewares))

	for idx, middleware := range e.middlewares {
		handlers[idx] = func(c *fiber.Ctx) error {
			err := middleware(newContext(c))

			var httpErr httpError
			if errors.As(err, &httpErr) {
				return c.Status(httpErr.status).SendString(httpErr.message)
			}

			return c.Next()
		}
	}

	var out Out
	shouldValidateOutput := reflect.TypeOf(out).Kind() == reflect.Struct

	handlers = append(handlers, func(c *fiber.Ctx) error {
		var input In

		parser := c.QueryParser

		if e.method == http.MethodPost {
			parser = c.BodyParser
		}

		if err := parser(&input); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := v.ValidateStruct(&input); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(validator.ValidationErrorSerializer{
				Errors: v.Translate(err),
			})
		}

		output, err := e.handler(newContext(c), input)
		if err != nil {
			slog.Error("error in handler %s %s: %w", e.method, e.path, err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if shouldValidateOutput {
			if err := v.ValidateStruct(&output); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(validator.ValidationErrorSerializer{
					Errors: v.Translate(err),
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(output)
	})

	r.Add(e.method, e.path, handlers...)
}
