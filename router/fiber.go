package router

import (
	"github.com/SupTarr/todo/todos"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

type FiberRouter struct {
	*fiber.App
}

func NewGinContext(c *gin.Context) *MyContext {
	return &MyContext{Context: c}
}

func (c MyContext) Bind(i interface{}) error {
	return c.ShouldBindJSON(i)
}

func (c MyContext) TodoID() string {
	return c.Param("id")
}

func (c MyContext) TransactionID() string {
	return c.Request.Header.Get("TransactionID")
}

func (c MyContext) Audience() string {
	if aud, ok := c.Get("aud"); ok {
		if s, ok := aud.(string); ok {
			return s
		}
	}

	return ""
}

func (c MyContext) Status(code int) {
	c.Context.Status(code)
}

func (c MyContext) JSON(code int, i interface{}) {
	c.Context.JSON(code, i)
}

func NewGinHandler(handler func(todos.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(&MyContext{Context: c})
	}
}
