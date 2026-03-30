package service

import (
	"io"
	"os"

	"github.com/zy84338719/upftp/internal/auth"
	"github.com/zy84338719/upftp/internal/biz"
	"github.com/zy84338719/upftp/internal/conf"
	"github.com/zy84338719/upftp/internal/model"
)

type Service struct {
	fileSvc  *biz.FileService
	cfg      *conf.Config
	sessions *auth.SessionManager
	serverIP string
	httpPort int
	ftpPort  int
	root     string
}

func New(fileSvc *biz.FileService, cfg *conf.Config, sessions *auth.SessionManager) *Service {
	return &Service{
		fileSvc:  fileSvc,
		cfg:      cfg,
		sessions: sessions,
	}
}

func (s *Service) SetServerInfo(ip string, httpPort, ftpPort int, root string) {
	s.serverIP = ip
	s.httpPort = httpPort
	s.ftpPort = ftpPort
	s.root = root
}

func (s *Service) ServerIP() string               { return s.serverIP }
func (s *Service) HTTPPort() int                  { return s.httpPort }
func (s *Service) FTPPort() int                   { return s.ftpPort }
func (s *Service) Root() string                   { return s.root }
func (s *Service) Config() *conf.Config           { return s.cfg }
func (s *Service) Sessions() *auth.SessionManager { return s.sessions }
func (s *Service) FileService() *biz.FileService  { return s.fileSvc }

// --- File Operations ---

func (s *Service) ListFiles(path string) ([]model.FileInfo, error) {
	return s.fileSvc.ListFiles(path)
}

func (s *Service) ListFilesLegacy(path string) ([]model.LegacyFileInfo, error) {
	files, err := s.fileSvc.ListFiles(path)
	if err != nil {
		return nil, err
	}
	legacy := make([]model.LegacyFileInfo, 0, len(files))
	for _, f := range files {
		legacy = append(legacy, model.NewLegacyFileInfo(f))
	}
	return legacy, nil
}

func (s *Service) Upload(path string, reader io.Reader, size int64) error {
	return s.fileSvc.Upload(path, reader, size)
}

func (s *Service) Download(path string) (string, error) {
	return s.fileSvc.Download(path)
}

func (s *Service) Delete(path string) error {
	return s.fileSvc.Delete(path)
}

func (s *Service) Rename(path string, newName string) error {
	return s.fileSvc.Rename(path, newName)
}

func (s *Service) Move(src string, dst string) error {
	return s.fileSvc.Move(src, dst)
}

func (s *Service) Copy(src string, dst string) error {
	return s.fileSvc.Copy(src, dst)
}

func (s *Service) CreateFolder(path string) error {
	return s.fileSvc.CreateFolder(path)
}

func (s *Service) Stat(path string) (os.FileInfo, error) {
	return s.fileSvc.Stat(path)
}

func (s *Service) SafePath(relativePath string) (string, error) {
	return s.fileSvc.SafePath(relativePath)
}

func (s *Service) OpenFile(path string) (*os.File, error) {
	return s.fileSvc.OpenFile(path)
}

func (s *Service) CreateFileForWrite(path string) (*os.File, error) {
	return s.fileSvc.CreateFileForWrite(path)
}

func (s *Service) ReadFileContent(path string) ([]byte, error) {
	return s.fileSvc.ReadFileContent(path)
}

func (s *Service) WriteFileContent(path string, content []byte) error {
	return s.fileSvc.WriteFileContent(path, content)
}

func (s *Service) ListDir(path string) ([]os.DirEntry, error) {
	return s.fileSvc.ListDir(path)
}

func (s *Service) BuildTree(path string, depth int) (*model.TreeNode, error) {
	return s.fileSvc.BuildTree(path, depth)
}

func (s *Service) SearchFiles(pattern string, path string) ([]string, error) {
	return s.fileSvc.SearchFiles(pattern, path)
}

// --- Upload Config ---

func (s *Service) UploadEnabled() bool {
	return s.cfg.Upload.Enabled
}

// --- Server Info ---

func (s *Service) GetServerInfo() map[string]interface{} {
	return map[string]interface{}{
		"version":     s.cfg.Version,
		"buildDate":   s.cfg.BuildDate,
		"goVersion":   s.cfg.GoVersion,
		"platform":    s.cfg.Platform,
		"projectURL":  s.cfg.ProjectURL,
		"projectName": s.cfg.ProjectName,
		"lastCommit":  s.cfg.LastCommit,
		"language":    s.cfg.GetLanguage(),
	}
}

func (s *Service) GetIndexPageData(files []model.LegacyFileInfo) map[string]interface{} {
	return map[string]interface{}{
		"Files":         files,
		"ServerInfo":    s.serverInfoStruct(),
		"CurrentPath":   "/",
		"UploadEnabled": s.cfg.Upload.Enabled,
		"Version":       s.cfg.Version,
		"BuildDate":     s.cfg.BuildDate,
		"GoVersion":     s.cfg.GoVersion,
		"Platform":      s.cfg.Platform,
		"ProjectURL":    s.cfg.ProjectURL,
		"ProjectName":   s.cfg.ProjectName,
		"LastCommit":    s.cfg.LastCommit,
		"Language":      s.cfg.GetLanguage(),
		"HTTPAuthOn":    s.cfg.HTTPAuth.Enabled,
		"HTTPAuthUser":  s.cfg.HTTPAuth.Username,
	}
}

func (s *Service) serverInfoStruct() *ServerInfo {
	return &ServerInfo{
		IP:       s.serverIP,
		HTTPPort: s.httpPort,
		FTPPort:  s.ftpPort,
		Root:     s.root,
	}
}

type ServerInfo struct {
	IP       string
	HTTPPort int
	FTPPort  int
	Root     string
}
