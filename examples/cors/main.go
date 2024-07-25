package main

import (
	"log"
	"net/http"

	"github.com/esequiel378/fast"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	a := fiber.New()
	a.Use(cors.New())

	app, err := fast.New(fast.WithFiberApp(a))
	if err != nil {
		log.Fatal(err)
	}

	app.MustRegister("/", GreetingHandler{})

	log.Fatal(app.Listen(":3000"))
}

type GreetingHandler struct{}

func (h GreetingHandler) HandleGet() fast.Handler {
	return fast.
		Endpoint[fast.In, fast.Out]().
		Method(http.MethodGet).
		Path("/greeting").
		Handle(func(*fast.Context, fast.In) (fast.Out, error) {
			return "Hello, World!", nil
		})
}
