package fast

import (
	"errors"
	"log"
	"reflect"

	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

// Handler is the interface that links the endpoint to the router
type Handler interface {
	// Register registers the endpoint to the given router
	Register(router fiber.Router, validator validator.Validator, middlewares ...Middleware)
	// Path returns the endpoint path
	Path() string
	// Method returns the HTTP method
	Method() string
	// Middlewares returns the middleware functions
	Middlewares() []func(*Context) error
	// InputSerializer returns the input serializer
	InputSerializer() any
	// OutputSerializer returns the output serializer
	OutputSerializer() any
}

// endpointHandler implements the Handler interface
type endpointHandler[I, O any] struct {
	path        string
	method      string
	handler     func(*Context, I) (O, error)
	middlewares []func(*Context) error
	input       I
	output      O
}

// Path returns the endpoint path
func (h *endpointHandler[I, O]) Path() string {
	return h.path
}

// Method returns the HTTP method
func (h *endpointHandler[I, O]) Method() string {
	return h.method
}

// Middlewares returns the middleware functions
func (h *endpointHandler[I, O]) Middlewares() []func(*Context) error {
	return h.middlewares
}

// Register registers the endpoint to the given router
func (h *endpointHandler[I, O]) Register(r fiber.Router, v validator.Validator, middlewares ...Middleware) {
	handlers := make([]fiber.Handler, len(h.middlewares)+len(middlewares))

	for idx, middleware := range append(middlewares, h.middlewares...) {
		handlers[idx] = func(c *fiber.Ctx) error {
			err := middleware(newContext(c))
			var httpErr httpError
			if errors.As(err, &httpErr) {
				return c.Status(httpErr.status).SendString(httpErr.message)
			}
			return err
		}
	}

	var out O
	shouldValidateOutput := reflect.TypeOf(out).Kind() == reflect.Struct

	handlers = append(handlers, func(c *fiber.Ctx) error {
		var input I
		parser := c.QueryParser

		if len(c.BodyRaw()) > 0 {
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

		output, err := h.handler(newContext(c), input)
		if err != nil {
			log.Printf("error in handler %s %s: %s", h.method, h.path, err)
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

	r.Add(h.method, h.path, handlers...)
}

func (h *endpointHandler[I, O]) InputSerializer() any {
	return h.input
}

func (h *endpointHandler[I, O]) OutputSerializer() any {
	return h.output
}
