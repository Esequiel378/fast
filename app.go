package fast

import (
	"reflect"
	"strings"

	"github.com/esequiel378/fast/internal/validator"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	validator validator.Validator
	server    *fiber.App
}

func New() (App, error) {
	v, err := validator.NewValidatorV10()
	if err != nil {
		return App{}, err
	}

	server := fiber.New()

	instance := App{
		validator: v,
		server:    server,
	}

	return instance, nil
}

func (a App) Run(addr string) error {
	return a.server.Listen(addr)
}

var handlerReturnType = reflect.TypeOf((*Handler)(nil)).Elem()

// MustRegister registers a handler to the app
// The handler must be a struct with methods that return a Handler
func (a App) MustRegister(prefix string, handler any) {
	router := a.server.Group(prefix)

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

		handler.Register(router, a.validator)
	}
}
