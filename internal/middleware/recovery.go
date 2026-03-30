package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zy84338719/upftp/internal/logger"
)

func RecoveryMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered: %v", r)
				c.Header("Content-Type", "application/json")
				c.String(500, fmt.Sprintf(`{"error":"Internal server error"}`))
			}
		}()
		c.Next(ctx)
	}
}
