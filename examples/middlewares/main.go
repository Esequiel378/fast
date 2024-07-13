package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/esequiel378/fast"
)

func main() {
	app, err := fast.New()
	if err != nil {
		log.Fatal(err)
	}

	app.MustRegister("/", NewGreetingHandler(AuthMiddleware{}))

	log.Fatal(app.Run(":3000"))
}

type GreetingHandler struct {
	apiKeyValidator APIKeyValidator
}

func NewGreetingHandler(apiKeyValidator APIKeyValidator) GreetingHandler {
	return GreetingHandler{
		apiKeyValidator: apiKeyValidator,
	}
}

func (h GreetingHandler) HandleGreeting() fast.Handler {
	type In struct {
		Name string `json:"name" validate:"required"`
	}

	type Out struct {
		Message string `json:"message" validate:"required"`
	}

	return fast.
		Endpoint[In, Out]().
		Method(http.MethodGet).
		Path("/greeting").
		Middlewares(h.apiKeyValidator.HandleValidateAPIKey()).
		Handle(func(_ *fast.Context, in In) (Out, error) {
			output := Out{
				Message: fmt.Sprintf("Hello, %s!", in.Name),
			}

			return output, nil
		})
}

// APIKeyValidator is an interface that defines the method to validate an API key in a request
type APIKeyValidator interface {
	HandleValidateAPIKey() fast.Middleware
}

type AuthMiddleware struct {
	// Normally, here you would have a dependency to a service that validates the API key or a database connection
}

var _ APIKeyValidator = (*AuthMiddleware)(nil)

func (m AuthMiddleware) HandleValidateAPIKey() fast.Middleware {
	return func(c *fast.Context) error {
		apiKey := c.Get("API-Key")

		if apiKey != "fast-is-awesome" {
			return fast.UnauthorizedError("invalid API key")
		}

		return nil
	}
}
