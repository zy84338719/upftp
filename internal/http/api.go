package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zy84338719/upftp/internal/core"
)

// --- 通用响应 ---

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

// --- 健康检查 ---

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok"})
}

// --- 会话控制 ---

func (s *Server) listSessions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"sessions": s.mgr.List(),
		"host":     s.host,
		"ftp_port": s.ftpPort,
	})
}

func (s *Server) createSession(w http.ResponseWriter, r *http.Request) {
	var spec core.SessionSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		writeErr(w, 400, "invalid body: "+err.Error())
		return
	}
	if spec.Dir == "" {
		writeErr(w, 400, "dir is required")
		return
	}
	// 未指定认证字段时,继承默认:匿名
	if spec.User == "" && !spec.Anonymous {
		spec.Anonymous = true
	}
	sess, err := s.mgr.Create(r.Context(), spec)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, sess)
}

func (s *Server) stopSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErr(w, 400, "missing id")
		return
	}
	if !s.mgr.Stop(id) {
		writeErr(w, 404, "session not found")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "stopped", "id": id})
}

// --- 文件操作 ---

type listFilesResp struct {
	Host    string           `json:"host"`
	FTPPort int              `json:"ftp_port"`
	Files   []core.FileEntry `json:"files"`
	Path    string           `json:"path"`
}

func (s *Server) listFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}
	entries, err := s.files.List(path)
	if err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	writeJSON(w, 200, listFilesResp{
		Host: s.host, FTPPort: s.ftpPort, Files: entries, Path: path,
	})
}

func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeErr(w, 400, "path is required")
		return
	}
	f, info, err := s.files.Open(path)
	if err != nil {
		writeErr(w, 404, err.Error())
		return
	}
	defer f.Close()
	w.Header().Set("Content-Disposition", "attachment; filename=\""+baseName(info.Name())+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	_, _ = io.Copy(w, f)
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	// 支持两种调用:
	//   1) multipart:字段 file=文件, path=目标相对路径(浏览器用)
	//   2) raw body:path 查询参数指定目标,请求体即文件内容(API 用)
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			writeErr(w, 400, "parse multipart: "+err.Error())
			return
		}
		file, hdr, err := r.FormFile("file")
		if err != nil {
			writeErr(w, 400, "missing file field")
			return
		}
		defer file.Close()
		target := r.FormValue("path")
		if target == "" {
			target = "/" + hdr.Filename
		}
		if err := s.files.Save(target, file); err != nil {
			writeErr(w, 500, err.Error())
			return
		}
		writeJSON(w, 201, map[string]string{"status": "uploaded", "path": target})
		return
	}

	target := r.URL.Query().Get("path")
	if target == "" {
		writeErr(w, 400, "path query is required for raw upload")
		return
	}
	if err := s.files.Save(target, r.Body); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]string{"status": "uploaded", "path": target})
}

func (s *Server) mkdir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeErr(w, 400, "path is required")
		return
	}
	if err := s.files.Mkdir(path); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]string{"status": "created", "path": path})
}

func (s *Server) deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeErr(w, 400, "path is required")
		return
	}
	if err := s.files.Remove(path); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted", "path": path})
}

func baseName(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[i+1:]
		}
	}
	return p
}
