package fast

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	validator validator.Validator
	server    *fiber.App
	path      string
	apiSchema *OpenAPIGenerator
}

// WithFiberApp sets the fiber app to use.
// This is useful to pre-configure the fiber app
func WithFiberApp(app *fiber.App) func(*App) {
	return func(a *App) {
		a.server = app
	}
}

// WithExperimentalOpenAPISchema enables the OpenAPI schema generator
// and serves the schema at /swagger.json and Swagger UI at /swagger
// WARN: This is experimental and not recommended for production use
func WithExperimentalOpenAPISchema() func(*App) {
	return func(a *App) {
		a.apiSchema = NewOpenAPIGenerator(OpenAPIInfo{
			Title:       "Fast",
			Description: "Fast",
			Version:     "0.0.1",
		})

		// Endpoint to serve the OpenAPI JSON
		a.server.Get("/swagger.json", func(c *fiber.Ctx) error {
			schema, err := a.apiSchema.GenerateJSON()
			if err != nil {
				return c.Status(500).JSON(map[string]string{
					"error": "Failed to generate OpenAPI schema",
				})
			}

			// Generate ETag based on content
			hash := sha256.Sum256([]byte(schema))
			etag := fmt.Sprintf(`"%x"`, hash[:])

			// Check if client sent If-None-Match header
			if c.Get("If-None-Match") == etag {
				return c.Status(304).Send(nil) // Not Modified
			}

			// Send full response with ETag
			c.Set("ETag", etag)
			// TODO: Add cache control based on environment
			// c.Set("Cache-Control", "max-age=3600") // Cache for 1 hour
			c.Set("Content-Type", "application/json")

			return c.SendString(schema)
		})

		// Serve Swagger UI
		a.server.Get("/swagger", func(c *fiber.Ctx) error {
			c.Set("Content-Type", "text/html")
			return c.SendString(swaggerUIHTML)
		})
	}
}

func New(opts ...func(*App)) (App, error) {
	v, err := validator.NewValidatorV10()
	if err != nil {
		return App{}, err
	}

	server := fiber.New()

	instance := App{
		validator: v,
		server:    server,
		path:      "",
	}

	for _, opt := range opts {
		opt(&instance)
	}

	return instance, nil
}

// Listen serves HTTP requests from the given addr.
//
//	app.Listen(":8080")
//	app.Listen("127.0.0.1:8080")
func (a App) Listen(addr string) error {
	// TODO: Improve endpoints listing on startup
	data, _ := json.MarshalIndent(a.server.Stack(), "", "  ")
	fmt.Print(string(data))
	return a.server.Listen(addr)
}

var handlerReturnType = reflect.TypeOf((*Handler)(nil)).Elem()

// MustRegister registers a handler to the app
// The handler must be a struct with Hanlder methods.
func (a App) MustRegister(prefix string, handler any, middlewares ...Middleware) {
	mustValidateAndRegisterHandler(
		path.Join(a.path, prefix),
		handler,
		a.server.Group(prefix),
		a.validator,
		a.apiSchema,
		middlewares...,
	)
}

// Group creates a new group of routes
func (a App) Group(prefix string, middlewares ...Middleware) Group {
	return Group{
		router:      a.server.Group(prefix),
		validator:   a.validator,
		path:        path.Join(a.path, prefix),
		apiSchema:   a.apiSchema,
		middlewares: middlewares,
	}
}

func mustValidateAndRegisterHandler(
	path string,
	handler any,
	router fiber.Router,
	validator validator.Validator,
	apiSchema *OpenAPIGenerator,
	middlewares ...Middleware,
) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Struct {
		panic("handler is not a struct")
	}

	handlerValue := reflect.ValueOf(handler)

	for i := range handlerType.NumMethod() {
		method := handlerType.Method(i)

		hasCorrectName := strings.HasPrefix(method.Name, "Handle")
		hasCorrectReturnType := method.Type.NumOut() == 1 && method.Type.Out(0).Implements(handlerReturnType)

		if !hasCorrectName || !hasCorrectReturnType {
			continue
		}

		handler, ok := method.Func.Call([]reflect.Value{handlerValue})[0].Interface().(Handler)
		if !ok {
			panic("methods starting with `Handle` must return fast.Handler")
		}

		handler.Register(router, validator, middlewares...)
		fmt.Printf("Registered %s %s\n", handler.Method(), path)
		if apiSchema != nil {
			apiSchema.RegisterHandler(path, handler)
		}
	}
}
