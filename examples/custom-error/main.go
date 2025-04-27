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

	app.MustRegister("/test-error", Handler{})

	log.Fatal(app.Listen(":3003"))
}

type Handler struct{}

func (h Handler) HandleGet() fast.Handler {
	return fast.
		Endpoint[fast.In, fast.Out]().
		Handle(func(*fast.Context, fast.In) (fast.Out, error) {
			return "", fast.NewHTTPError(http.StatusNotFound, "Resource not found")
		})
}
