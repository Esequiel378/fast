package main

import (
	"log"
	"net/http"

	"github.com/esequiel378/fast"
)

func main() {
	app, err := fast.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Group("/admin").
		MustRegister("/users", UserHandler{}).
		MustRegister("/accounts", AccountsHandler{})

	log.Fatal(app.Run(":3000"))
}

type User struct {
	Name string `json:"name"`
}

type UserHandler struct{}

func (h UserHandler) HandleList() fast.Handler {
	type Out struct {
		Users []User `json:"users"`
	}

	return fast.
		Endpoint[fast.In, Out]().
		Method(http.MethodGet).
		Path("/").
		Handle(func(*fast.Context, fast.In) (Out, error) {
			output := Out{
				Users: []User{
					{Name: "Alice"},
					{Name: "Bob"},
				},
			}

			return output, nil
		})
}

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AccountsHandler struct{}

func (h AccountsHandler) HandleList() fast.Handler {
	type In struct {
		Limit int `json:"limit" validate:"omitempty,gte=10,lte=100"`
		Page  int `json:"page" validate:"omitempty,gte=0"`
	}

	type Out struct {
		Accounts []Account `json:"accounts"`
	}

	return fast.
		Endpoint[In, Out]().
		Method(http.MethodGet).
		Path("/").
		Handle(func(_ *fast.Context, input In) (Out, error) {
			output := Out{
				Accounts: []Account{
					{ID: 1, Name: "Alice"},
					{ID: 2, Name: "Bob"},
				},
			}

			return output, nil
		})
}
