package main

import (
	"log"

	"github.com/esequiel378/fast"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	a := fiber.New()
	a.Use(
		cors.New(
			cors.Config{
				AllowMethods: "GET",
			},
		),
	)

	app, err := fast.New(fast.WithFiberApp(a))
	if err != nil {
		log.Fatal(err)
	}

	app.MustRegister("/greeting", GreetingHandler{})

	log.Fatal(app.Listen(":3003"))
}

type GreetingHandler struct{}

func (h GreetingHandler) HandleGet() fast.Handler {
	return fast.
		Endpoint[fast.In, fast.Out]().
		Handle(func(*fast.Context, fast.In) (fast.Out, error) {
			return "Hello, World!", nil
		})
}
