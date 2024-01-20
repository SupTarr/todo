package todos

import "github.com/gin-gonic/gin"

type MyContext struct {
	*gin.Context
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

func NewGinHandler(handler func(Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(&MyContext{Context: c})
	}
}
