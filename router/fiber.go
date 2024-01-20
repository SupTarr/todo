package router

import (
	"github.com/SupTarr/todo/todos"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberRouter struct {
	*fiber.App
}

func NewFiberRouter() *FiberRouter {
	r := fiber.New()
	r.Use(cors.New())
	r.Use(logger.New())
	return &FiberRouter{r}
}

func (r *FiberRouter) POST(path string, handler func(todos.Context)) {
	r.App.Post(path, func(c *fiber.Ctx) error {
		handler(NewFiberCtx(c))
		return nil
	})
}

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{Ctx: c}
}

func (c FiberCtx) Bind(i interface{}) error {
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

func (c FiberCtx) JSON(code int, i interface{}) {
	c.Ctx.Status(code).JSON(i)
}
