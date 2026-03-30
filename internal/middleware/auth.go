package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/auth"
	"github.com/zy84338719/upftp/internal/conf"
	"github.com/zy84338719/upftp/internal/logger"
)

func AuthMiddleware(cfg *conf.Config, sm *auth.SessionManager) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !cfg.HTTPAuth.Enabled {
			c.Next(ctx)
			return
		}

		if cookie := string(c.Cookie("auth_token")); cookie != "" {
			if session, valid := sm.ValidateSession(cookie); valid {
				c.Set("username", session.Username)
				c.Next(ctx)
				return
			}
		}

		if user, pass, ok := c.Request.BasicAuth(); ok {
			if user == cfg.HTTPAuth.Username && pass == cfg.HTTPAuth.Password {
				c.Next(ctx)
				return
			}
		}

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
