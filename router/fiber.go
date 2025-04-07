package router

import (
	"github.com/gofiber/fiber/v2"
)

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{Ctx: c}
}

func (c FiberCtx) Bind(i any) error {
	return c.Ctx.BodyParser(i)
}

func (c FiberCtx) TodoID() string {
	return c.Ctx.Params("id")
}

func (c FiberCtx) TransactionID() string {
	return string(c.Ctx.Request().Header.Peek("TransactionID"))
}

func (c FiberCtx) Audience() string {
	return c.Ctx.Get("aud")
}

func (c FiberCtx) Status(code int) {
	c.Ctx.Status(code)
}

func (c FiberCtx) JSON(code int, i any) {
	c.Ctx.Status(code).JSON(i)
}
