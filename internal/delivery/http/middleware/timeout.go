package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		finishedRequest := make(chan struct{})
		panickedRequest := make(chan any)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panickedRequest <- p
				}
			}()

			c.Request = c.Request.WithContext(timeoutCtx)
			c.Next()
			close(finishedRequest)
		}()

		select {
		case <-ctx.Done():
			c.Error(errors.New("request timed out"))
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "request timed out",
			})
		case p := <-panickedRequest:
			panic(p)
		case <-finishedRequest:
		}
	}
}
