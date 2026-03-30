package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func CORSMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if string(c.Request.Method()) == "OPTIONS" {
			c.SetStatusCode(204)
			return
		}

		c.Next(ctx)
	}
}
