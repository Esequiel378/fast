package fast

type Middleware = func(ctx *Context) error
