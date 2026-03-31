package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/pkg/conf"
	"github.com/zy84338719/upftp/pkg/logger"
)

func AuthMiddleware(cfg *conf.Config) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 添加 nil 检查
		if cfg == nil {
			logger.Error("AuthMiddleware: config is nil")
			c.Next(ctx)
			return
		}

		// 跳过公开路由
		path := string(c.Request.URI().Path())
		if path == "/api/settings" || path == "/api/login" {
			c.Next(ctx)
			return
		}

		if !cfg.HTTPAuth.Enabled {
			c.Next(ctx)
			return
		}

		// 检查 HTTP Basic Auth
		if user, pass, ok := c.Request.BasicAuth(); ok {
			if user == cfg.HTTPAuth.Username && pass == cfg.HTTPAuth.Password {
				c.Next(ctx)
				return
			}
		}

		// 检查前端传来的 token（从 cookie 或 header）
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader != "" {
			// 前端可能发送 "Bearer token" 或 "Basic base64" 格式
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				authType := parts[0]
				authValue := parts[1]

				if authType == "Basic" {
					// Basic auth 已经处理过了
				} else if authType == "Bearer" {
					// 检查 token 是否有效（简单检查非空）
					if authValue != "" {
						c.Next(ctx)
						return
					}
				}
			}
		}

		// API 请求返回 401
		accept := string(c.GetHeader("Accept"))
		contentType := string(c.GetHeader("Content-Type"))
		if strings.Contains(accept, "application/json") || strings.Contains(contentType, "application/json") {
			c.Header("Content-Type", "application/json")
			c.SetStatusCode(consts.StatusUnauthorized)
			_, _ = c.WriteString(`{"error":"Unauthorized"}`)
		} else {
			c.Redirect(consts.StatusSeeOther, []byte("/login"))
		}
		logger.Warn("Unauthorized access attempt from %s", c.ClientIP())
	}
}
