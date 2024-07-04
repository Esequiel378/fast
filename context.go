package fast

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Context struct {
	ctx *fiber.Ctx
}

func newContext(ctx *fiber.Ctx) Context {
	return Context{
		ctx: ctx,
	}
}

// GetReqHeaders returns the HTTP request headers.
func (c Context) Headers() http.Header {
	return http.Header(c.ctx.GetReqHeaders())
}
