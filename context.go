package fast

import (
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	*fiber.Ctx
}

func newContext(ctx *fiber.Ctx) *Context {
	return &Context{
		Ctx: ctx,
	}
}
