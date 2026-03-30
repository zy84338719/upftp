package handler

import (
	"context"
	"encoding/json"
	"html/template"
	"io/fs"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/model"
)

func HandleIndexPage(ctx context.Context, c *app.RequestContext) {
	legacyFiles, err := svc.ListFilesLegacy("/")
	if err != nil {
		logger.Error("Failed to list root files: %v", err)
		legacyFiles = []model.LegacyFileInfo{}
	}

	data := svc.GetIndexPageData(legacyFiles)

	funcMap := template.FuncMap{
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}

	tmpl, err := template.New("index").Funcs(funcMap).ParseFS(templates, "templates/index.html")
	if err != nil {
		c.String(consts.StatusInternalServerError, "Template parse error: "+err.Error())
		logger.Error("Failed to parse index template: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(c.Response.BodyWriter(), data); err != nil {
		logger.Error("Failed to execute index template: %v", err)
	}
}

func getTemplatesFS() fs.FS {
	sub, _ := fs.Sub(templates, "templates")
	return sub
}
