package logic

import (
	"archive/zip"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templates embed.FS

type FileInfo struct {
	Name       string
	Size       string
	ModTime    string
	IsDir      bool
	CanPreview bool
	FileType   string
	Path       string
}

func GinServer(ctx context.Context) {
	mux := http.NewServeMux()

	// 静态文件服务
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(Root))))

	// 主页和文件夹浏览
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if !strings.HasPrefix(urlPath, "/") {
			urlPath = "/" + urlPath
		}

		fsPath := path.Join(Root, urlPath)
		fileInfo, err := os.Stat(fsPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if !fileInfo.IsDir() {
			http.ServeFile(w, r, fsPath)
			return
		}

		files, _ := ioutil.ReadDir(fsPath)
		fileList := []FileInfo{}

		// 添加返回上级目录的链接
		if urlPath != "/" {
			fileList = append(fileList, FileInfo{
				Name:  "..",
				IsDir: true,
				Path:  path.Dir(urlPath),
			})
		}

		for _, file := range files {
			filePath := path.Join(urlPath, file.Name())
			fileList = append(fileList, FileInfo{
				Name:       file.Name(),
				Size:       formatFileSize(file.Size()),
				ModTime:    file.ModTime().Format("2006-01-02 15:04:05"),
				IsDir:      file.IsDir(),
				CanPreview: !file.IsDir() && canPreviewFile(file.Name()),
				FileType:   getFileType(file.Name()),
				Path:       filePath,
			})
		}

		tmpl, _ := template.ParseFS(templates, "templates/index.html")
		tmpl.Execute(w, fileList)
	})

	// 下载（支持文件夹打包下载）
	mux.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/download/")
		if !isPathSafe(filename) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		filePath := path.Join(Root, filename)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if !fileInfo.IsDir() {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))
			http.ServeFile(w, r, filePath)
			return
		}

		// 处理文件夹下载
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(filename)))

		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 获取相对路径
			relPath, err := filepath.Rel(filePath, path)
			if err != nil {
				return err
			}

			// 跳过根目录
			if relPath == "." {
				return nil
			}

			// 创建zip文件头
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = relPath

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(writer, file)
			}
			return err
		})

		if err != nil {
			http.Error(w, "Error creating zip file", http.StatusInternalServerError)
			return
		}
	})

	// 预览
	mux.HandleFunc("/preview/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/preview/")
		if !isPathSafe(filename) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		filePath := path.Join(Root, filename)
		http.ServeFile(w, r, filePath)
	})

	server := &http.Server{
		Addr:    Port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("Server start error = %s", err))
	}
}
