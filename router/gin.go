package router

import (
	"github.com/SupTarr/todo/my_context"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	*gin.Context
}

func NewGinContext(c *gin.Context) *GinContext {
	return &GinContext{Context: c}
}

func (c GinContext) Bind(i any) error {
	return c.ShouldBindJSON(i)
}

func (c GinContext) TodoID() string {
	return c.Param("id")
}

func (c GinContext) TransactionID() string {
	return c.Request.Header.Get("TransactionID")
}

func (c GinContext) Audience() string {
	if aud, ok := c.Get("aud"); ok {
		if s, ok := aud.(string); ok {
			return s
		}
	}

	return ""
}

func (c GinContext) Status(code int) {
	c.Context.Status(code)
}

func (c GinContext) JSON(code int, i any) {
	c.Context.JSON(code, i)
}

func NewGinHandler(handler func(my_context.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(&GinContext{Context: c})
	}
}
