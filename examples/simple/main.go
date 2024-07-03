package main

import (
	"log"
	"net/http"

	"github.com/esequiel378/fast"
)

func main() {
	app, _ := fast.New()

	app.MustRegister("/", UserHandler{})

	app.Group("/admin").
		MustRegister("/users", UserHandler{}).
		MustRegister("/accounts", UserHandler{})

	log.Fatal(app.Run(":3000"))
}

type User struct {
	Name string
}

type UserHandler struct{}

func (h UserHandler) HandleList() fast.Handler {
	type In struct {
		Limit int `json:"limit" validate:"omitempty,gte=10,lte=100"`
		Page  int `json:"page" validate:"omitempty,gte=0"`
	}

	type Out struct {
		Users []User
	}

	return fast.
		NewEndpoint[In, Out]().
		Path("/").
		Method(http.MethodGet).
		Handle(func(c fast.Context, input In) (Out, error) {
			// Perfom database query
			output := Out{
				Users: []User{
					{Name: "Alice"},
					{Name: "Bob"},
				},
			}

			return output, nil
		})
}
