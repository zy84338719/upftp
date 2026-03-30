package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/config"
)

func HandleGetSettings(ctx context.Context, c *app.RequestContext) {
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]interface{}{
		"language":     config.AppConfig.GetLanguage(),
		"httpAuth":     config.AppConfig.HTTPAuth.Enabled,
		"httpAuthUser": config.AppConfig.HTTPAuth.Username,
		"httpAuthPass": config.AppConfig.HTTPAuth.Password,
		"ftpUser":      config.AppConfig.Username,
		"ftpPass":      config.AppConfig.Password,
		"ftpEnabled":   config.AppConfig.EnableFTP,
		"mcpEnabled":   config.AppConfig.EnableMCP,
		"httpPort":     config.AppConfig.GetHTTPPort(),
		"ftpPort":      config.AppConfig.GetFTPPort(),
	})
}

func HandleSetLanguage(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	var req struct {
		Language string `json:"language"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	if req.Language != "en" && req.Language != "zh" {
		c.String(consts.StatusBadRequest, `{"error":"Invalid language"}`)
		return
	}
	config.AppConfig.Language = req.Language
	if err := config.SaveConfig(); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}

func HandleSetHTTPAuth(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	var req struct {
		Enabled  bool   `json:"enabled"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	config.AppConfig.HTTPAuth.Enabled = req.Enabled
	if req.Username != "" {
		config.AppConfig.HTTPAuth.Username = req.Username
	}
	if req.Password != "" {
		config.AppConfig.HTTPAuth.Password = req.Password
	}
	if err := config.SaveConfig(); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}

func HandleSetFTP(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	if req.Username != "" {
		config.AppConfig.Username = req.Username
	}
	if req.Password != "" {
		config.AppConfig.Password = req.Password
	}
	if err := config.SaveConfig(); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}
