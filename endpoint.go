package fast

import (
	"net/http"
)

type (
	// In is the default input type for an endpoint
	In struct{}
	// Out is the default output type for an endpoint
	Out string
)

// EndpointBuilder is the builder for creating endpoints
type EndpointBuilder[I, O any] struct {
	path        string
	method      string
	middlewares []func(*Context) error
}

// Endpoint creates a new endpoint builder
func Endpoint[I, O any]() *EndpointBuilder[I, O] {
	return &EndpointBuilder[I, O]{
		path:   "/",
		method: http.MethodGet,
	}
}

// Path sets the path of the endpoint
func (b *EndpointBuilder[I, O]) Path(path string) *EndpointBuilder[I, O] {
	b.path = path
	return b
}

// Method sets the method of the endpoint
func (b *EndpointBuilder[I, O]) Method(method string) *EndpointBuilder[I, O] {
	b.method = method
	return b
}

// Middlewares sets the middlewares of the endpoint
func (b *EndpointBuilder[I, O]) Middlewares(middlewares ...func(*Context) error) *EndpointBuilder[I, O] {
	b.middlewares = middlewares
	return b
}

// Handle finalizes the builder and returns a Handler that can be registered
func (b *EndpointBuilder[I, O]) Handle(fn func(*Context, I) (O, error)) Handler {
	var (
		input  I
		output O
	)

	return &endpointHandler[I, O]{
		path:        b.path,
		method:      b.method,
		handler:     fn,
		middlewares: b.middlewares,
		input:       input,
		output:      output,
	}
}
