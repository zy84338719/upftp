package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/response"
	"github.com/zy84338719/upftp/internal/service"
)

func HandleGetSettings(ctx context.Context, c *app.RequestContext) {
	response.Success(c, svc.GetSettings())
}

func HandleSetLanguage(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	var req struct {
		Language string `json:"language"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	if err := svc.SetLanguage(req.Language); err != nil {
		if err == service.ErrInvalidLanguage {
			response.BadRequest(c, "Invalid language")
			return
		}
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	response.Success(c, map[string]bool{"success": true})
}

func HandleSetHTTPAuth(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	var req struct {
		Enabled  bool   `json:"enabled"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	if err := svc.SetHTTPAuth(req.Enabled, req.Username, req.Password); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	response.Success(c, map[string]bool{"success": true})
}

func HandleSetFTP(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	if err := svc.SetFTP(req.Username, req.Password); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	response.Success(c, map[string]bool{"success": true})
}
