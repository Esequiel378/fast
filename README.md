# Fast

Fast is conviniet wrapper for the [Fiber framework](https://gofiber.io/). It provides a set of tools to simplify the development of web applications.


### Example

```go
func main() {
	app, _ := fast.New()

	app.MustRegister("/users", UserHandler{})

	log.Fatal(app.Run(":3000"))
}

type User struct {
	Name string
}

type UserHandler struct{}

func (h UserHandler) HandleList() fast.Handler {
    type In struct {}
    type Out struct {}

	return fast.
		Endpoint[In, Out]().
		Method(http.MethodGet).
		Path("/").
		Handle(func(c fast.Context, input In) (Out, error) {
			return Out{}, nil
		})
}
```

The `In` and `Out` types are used to define the input and output of the endpoint.
Fast will perform [validations](https://github.com/go-playground/validator) under the hood and will automatically serialize the output to JSON.
