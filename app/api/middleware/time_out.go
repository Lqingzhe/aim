package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func SetTimeOut(timeOut time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx, cancel := context.WithTimeout(ctx, timeOut)
		defer cancel()
		c.Set("ctx", ctx)
		c.Next()
	}
}
