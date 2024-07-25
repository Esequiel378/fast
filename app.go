package fast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	validator validator.Validator
	server    *fiber.App
}

// WithFiberApp sets the fiber app to use.
// This is useful to pre-configure the fiber app
func WithFiberApp(app *fiber.App) func(*App) {
	return func(a *App) {
		a.server = app
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
func (a App) MustRegister(prefix string, handler any) {
	mustValidateAndRegisterHandler(handler, a.server.Group(prefix), a.validator)
}

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

// Group creates a new group of routes
func (a App) Group(prefix string) Group {
	return Group{
		router:    a.server.Group(prefix),
		validator: a.validator,
	}
}

func mustValidateAndRegisterHandler(handler any, router fiber.Router, validator validator.Validator) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Struct {
		panic("handler is not a struct")
	}

	handlerValue := reflect.ValueOf(handler)

	for i := 0; i < handlerType.NumMethod(); i++ {
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

		handler.Register(router, validator)
	}
}
