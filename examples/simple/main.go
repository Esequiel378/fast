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

	app.MustRegister("/", GreetingHandler{})

	log.Fatal(app.Run(":3000"))
}

type GreetingHandler struct{}

func (h GreetingHandler) HandleGet() fast.Handler {
	return fast.
		Endpoint[fast.In, fast.Out]().
		Method(http.MethodGet).
		Path("/greeting").
		Handle(func(fast.Context, fast.In) (fast.Out, error) {
			return "Hello, World!", nil
		})
}
