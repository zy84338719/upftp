package handlers

import (
	"context"
	"encoding/json"
	"html/template"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/auth"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/logger"
)

var sessionManager = auth.NewSessionManager()

func HandleLogin(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.SetStatusCode(consts.StatusBadRequest)
		json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}

	if err := json.Unmarshal(body, &creds); err != nil {
		c.Header("Content-Type", "application/json")
		c.SetStatusCode(consts.StatusBadRequest)
		json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	if creds.Username != config.AppConfig.HTTPAuth.Username ||
		creds.Password != config.AppConfig.HTTPAuth.Password {
		c.Header("Content-Type", "application/json")
		c.SetStatusCode(consts.StatusUnauthorized)
		json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]string{
			"error":   "Unauthorized",
			"message": "用户名或密码错误",
		})
		logger.Warn("Failed login attempt for user: %s from %s", creds.Username, c.ClientIP())
		return
	}

	session, err := sessionManager.CreateSession(creds.Username)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.SetStatusCode(consts.StatusInternalServerError)
		json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]string{
			"error": "Failed to create session",
		})
		logger.Error("Failed to create session: %v", err)
		return
	}

	cookieMaxAge := 0
	if creds.Remember {
		cookieMaxAge = 86400 * 30
	}

	c.SetCookie("auth_token", session.Token, cookieMaxAge, "/", "",
		protocol.CookieSameSiteStrictMode, false, true)

	if creds.Remember {
		c.SetCookie("auth_username", creds.Username, cookieMaxAge, "/", "",
			protocol.CookieSameSiteStrictMode, false, false)
	}

	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]interface{}{
		"success":  true,
		"token":    session.Token,
		"username": session.Username,
		"expires":  session.ExpiresAt.Unix(),
	})

	logger.Info("User logged in: %s from %s", creds.Username, c.ClientIP())
}

func HandleLogout(ctx context.Context, c *app.RequestContext) {
	if cookie := string(c.Cookie("auth_token")); cookie != "" {
		sessionManager.DeleteSession(cookie)
	}

	c.SetCookie("auth_token", "", -1, "/", "",
		protocol.CookieSameSiteStrictMode, false, true)
	c.SetCookie("auth_username", "", -1, "/", "",
		protocol.CookieSameSiteStrictMode, false, false)

	accept := string(c.GetHeader("Accept"))
	contentType := string(c.GetHeader("Content-Type"))
	if accept == "application/json" || contentType == "application/json" {
		c.Header("Content-Type", "application/json")
		json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{
			"success": true,
		})
		return
	}

	c.Redirect(consts.StatusSeeOther, []byte("/login"))
	logger.Info("User logged out")
}

func HandleLoginPage(ctx context.Context, c *app.RequestContext) {
	if !config.AppConfig.HTTPAuth.Enabled {
		c.Redirect(consts.StatusSeeOther, []byte("/"))
		return
	}

	if cookie := string(c.Cookie("auth_token")); cookie != "" {
		if _, valid := sessionManager.ValidateSession(cookie); valid {
			c.Redirect(consts.StatusSeeOther, []byte("/"))
			return
		}
	}

	tmpl, err := template.ParseFS(templates, "templates/login.html")
	if err != nil {
		c.String(consts.StatusInternalServerError, "Internal server error")
		logger.Error("Failed to parse login template: %v", err)
		return
	}

	data := struct {
		Version    string
		ProjectURL string
	}{
		Version:    config.AppConfig.Version,
		ProjectURL: config.AppConfig.ProjectURL,
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(c.Response.BodyWriter(), data); err != nil {
		logger.Error("Failed to execute login template: %v", err)
	}
}

func GetSessionManager() *auth.SessionManager {
	return sessionManager
}
