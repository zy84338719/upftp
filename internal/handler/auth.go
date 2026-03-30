package handler

import (
	"context"
	"encoding/json"
	"html/template"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/response"
	"github.com/zy84338719/upftp/internal/service"
)

func HandleLogin(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	var input service.LoginInput
	if err := json.Unmarshal(body, &input); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	if !svc.Authenticate(input.Username, input.Password) {
		response.Unauthorized(c, "Invalid username or password")
		logger.Warn("Failed login attempt for user: %s from %s", input.Username, c.ClientIP())
		return
	}

	result, err := svc.Login(input)
	if err != nil {
		response.InternalError(c, "Failed to create session")
		logger.Error("Failed to create session: %v", err)
		return
	}

	cookieMaxAge := 0
	if input.Remember {
		cookieMaxAge = 86400 * 30
	}

	c.SetCookie("auth_token", result.Token, cookieMaxAge, "/", "",
		protocol.CookieSameSiteStrictMode, false, true)

	if input.Remember {
		c.SetCookie("auth_username", input.Username, cookieMaxAge, "/", "",
			protocol.CookieSameSiteStrictMode, false, false)
	}

	response.Success(c, map[string]interface{}{
		"success":  true,
		"token":    result.Token,
		"username": result.Username,
		"expires":  result.Expires,
	})
}

func HandleLogout(ctx context.Context, c *app.RequestContext) {
	svc.Logout(string(c.Cookie("auth_token")))

	c.SetCookie("auth_token", "", -1, "/", "",
		protocol.CookieSameSiteStrictMode, false, true)
	c.SetCookie("auth_username", "", -1, "/", "",
		protocol.CookieSameSiteStrictMode, false, false)

	accept := string(c.GetHeader("Accept"))
	contentType := string(c.GetHeader("Content-Type"))
	if accept == "application/json" || contentType == "application/json" {
		response.Success(c, map[string]bool{"success": true})
		return
	}

	c.Redirect(consts.StatusSeeOther, []byte("/login"))
}

func HandleLoginPage(ctx context.Context, c *app.RequestContext) {
	if !svc.IsHTTPAuthEnabled() {
		c.Redirect(consts.StatusSeeOther, []byte("/"))
		return
	}

	if cookie := string(c.Cookie("auth_token")); cookie != "" {
		if _, valid := svc.ValidateSession(cookie); valid {
			c.Redirect(consts.StatusSeeOther, []byte("/"))
			return
		}
	}

	tmpl, err := template.ParseFS(getTemplatesFS(), "login.html")
	if err != nil {
		c.String(consts.StatusInternalServerError, "Internal server error")
		logger.Error("Failed to parse login template: %v", err)
		return
	}

	data := struct {
		Version    string
		ProjectURL string
	}{
		Version:    svc.Config().Version,
		ProjectURL: svc.Config().ProjectURL,
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(c.Response.BodyWriter(), data); err != nil {
		logger.Error("Failed to execute login template: %v", err)
	}
}
