package index

import (
	"embed"
	"io/fs"
)

//go:embed templates/*
var templates embed.FS

func getTemplatesFS() fs.FS {
	sub, _ := fs.Sub(templates, "templates")
	return sub
}
