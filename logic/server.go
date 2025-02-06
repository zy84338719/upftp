package logic

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templates embed.FS

func GinServer(ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	router := gin.New()
	htmlFS, _ := fs.Sub(templates, "templates")
	router.SetHTMLTemplate(template.Must(template.ParseFS(htmlFS, "*.html")))
	router.StaticFS("/files", gin.Dir(Root, true))

	router.GET("/", func(c *gin.Context) {
		files, _ := ioutil.ReadDir(Root)
		fileList := []map[string]interface{}{}
		for _, file := range files {
			fileType := getFileType(file.Name())
			fileMap := map[string]interface{}{
				"Name":       file.Name(),
				"Size":       formatFileSize(file.Size()),
				"ModTime":    file.ModTime().Format("2006-01-02 15:04:05"),
				"CanPreview": canPreviewFile(file.Name()),
				"FileType":   fileType,
			}
			fileList = append(fileList, fileMap)
		}
		c.HTML(http.StatusOK, "index.html", fileList)
	})

	router.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.FileAttachment(path.Join(Root, filename), filename)
	})

	router.GET("/preview/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		if !isPathSafe(filename) {
			c.String(http.StatusForbidden, "Access denied")
			return
		}
		filePath := path.Join(Root, filename)
		c.File(filePath)
	})

	if err := router.Run(Port); err != nil {
		panic(fmt.Errorf("Server start error = %s", err))
	}
}
