# Fast

Fast is conviniet wrapper for the [Fiber framework](https://gofiber.io/). It provides a set of tools to simplify the development of web applications.


### Example

```go
func main() {
  app, _ := fast.New()

  app.MustRegister("/", GreetingHandler{})

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
```

```shell
$ curl http://localhost:3003
Hello, World!
```

The `In` and `Out` types are used to define the input and output of the endpoint.
Fast will perform [validations](https://github.com/go-playground/validator) under the hood and will automatically serialize the output to JSON.
