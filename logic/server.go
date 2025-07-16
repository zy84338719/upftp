package logic

import (
	"archive/zip"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zy84338719/upftp/config"
	"github.com/zy84338719/upftp/filehandler"
)

//go:embed templates/*
var templates embed.FS

type ServerInfo struct {
	IP       string
	HTTPPort int
	FTPPort  int
	Root     string
}

var serverInfo *ServerInfo

func SetServerInfo(ip string, httpPort, ftpPort int, root string) {
	serverInfo = &ServerInfo{
		IP:       ip,
		HTTPPort: httpPort,
		FTPPort:  ftpPort,
		Root:     root,
	}
}

func StartHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()

	// API接口
	mux.HandleFunc("/api/info", handleServerInfo)

	// 静态文件服务
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(config.AppConfig.Root))))

	// 主页和文件夹浏览
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if !strings.HasPrefix(urlPath, "/") {
			urlPath = "/" + urlPath
		}

		fsPath := path.Join(config.AppConfig.Root, urlPath)
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
		fileList := []filehandler.FileInfo{}

		// 添加返回上级目录的链接
		if urlPath != "/" {
			fileList = append(fileList, filehandler.FileInfo{
				Name:  "..",
				IsDir: true,
				Path:  path.Dir(urlPath),
				Icon:  "📁",
			})
		}

		for _, file := range files {
			filePath := path.Join(urlPath, file.Name())
			fileType := filehandler.GetFileType(file.Name())
			
			fileInfo := filehandler.FileInfo{
				Name:        file.Name(),
				Size:        filehandler.FormatFileSize(file.Size()),
				ModTime:     file.ModTime().Format("2006-01-02 15:04:05"),
				IsDir:       file.IsDir(),
				CanPreview:  !file.IsDir() && filehandler.CanPreviewFile(fileType),
				FileType:    fileType,
				FileTypeStr: filehandler.GetFileTypeString(fileType),
				Path:        filePath,
				Icon:        getFileIcon(file.IsDir(), fileType),
				MimeType:    filehandler.GetMimeType(fileType),
			}
			fileList = append(fileList, fileInfo)
		}

		data := struct {
			Files      []filehandler.FileInfo
			ServerInfo *ServerInfo
			CurrentPath string
		}{
			Files:       fileList,
			ServerInfo:  serverInfo,
			CurrentPath: urlPath,
		}

		tmpl, _ := template.ParseFS(templates, "templates/index.html")
		tmpl.Execute(w, data)
	})

	// 下载（支持文件夹打包下载）
	mux.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/download/")
		if !filehandler.IsPathSafe(filename) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		filePath := path.Join(config.AppConfig.Root, filename)
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
		if !filehandler.IsPathSafe(filename) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		filePath := path.Join(config.AppConfig.Root, filename)
		
		// 获取文件类型并设置适当的Content-Type
		fileType := filehandler.GetFileType(filename)
		mimeType := filehandler.GetMimeType(fileType)
		w.Header().Set("Content-Type", mimeType)
		
		http.ServeFile(w, r, filePath)
	})

	server := &http.Server{
		Addr:    config.AppConfig.Port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Println("Stopping HTTP server...")
		server.Shutdown(context.Background())
	}()

	log.Printf("HTTP server starting on %s%s", serverInfo.IP, config.AppConfig.Port)
	log.Printf("Web interface: http://%s%s", serverInfo.IP, config.AppConfig.Port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %v", err)
	}

	return nil
}

func handleServerInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := fmt.Sprintf(`{
		"version": "%s",
		"lastCommit": "%s",
		"httpPort": %d,
		"ftpPort": %d,
		"ftpEnabled": %t,
		"rootPath": "%s"
	}`, 
		config.AppConfig.Version,
		config.AppConfig.LastCommit,
		serverInfo.HTTPPort,
		serverInfo.FTPPort,
		config.AppConfig.EnableFTP,
		config.AppConfig.Root,
	)
	
	w.Write([]byte(response))
}

func getFileIcon(isDir bool, fileType filehandler.FileType) string {
	if isDir {
		return "📁"
	}
	return filehandler.GetFileIcon(fileType)
}
