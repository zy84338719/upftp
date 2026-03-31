package file

import (
	"io"
	"os"

	"github.com/zy84338719/upftp/pkg/conf"
	"github.com/zy84338719/upftp/pkg/file/dal"
	"github.com/zy84338719/upftp/pkg/file/model"
)

type Service struct {
	cfg      *conf.Config
	fileSvc  *FileService
	serverIP string
	httpPort int
	ftpPort  int
	root     string
}

func NewService(cfg *conf.Config) *Service {
	store := dal.NewLocalFileStore()
	path := dal.NewPathResolver(cfg.Root)
	fileSvc := NewFileService(store, path, cfg)
	return &Service{
		cfg:     cfg,
		fileSvc: fileSvc,
	}
}

func (s *Service) SetServerInfo(ip string, httpPort, ftpPort int, root string) {
	s.serverIP = ip
	s.httpPort = httpPort
	s.ftpPort = ftpPort
	s.root = root
}

func (s *Service) ServerIP() string     { return s.serverIP }
func (s *Service) HTTPPort() int        { return s.httpPort }
func (s *Service) FTPPort() int         { return s.ftpPort }
func (s *Service) Root() string         { return s.root }
func (s *Service) Config() *conf.Config { return s.cfg }

// --- File Operations ---

func (s *Service) ListFiles(path string) ([]model.FileInfo, error) {
	return s.fileSvc.ListFiles(path)
}

func (s *Service) ListFilesLegacy(path string) ([]model.LegacyFileInfo, error) {
	return s.fileSvc.ListFilesLegacy(path)
}

func (s *Service) GetFile(path string) (io.ReadCloser, int64, string, error) {
	return s.fileSvc.GetFile(path)
}

func (s *Service) Upload(path string, content io.Reader, size int64) error {
	return s.fileSvc.Upload(path, content, size)
}

func (s *Service) Delete(path string) error {
	return s.fileSvc.Delete(path)
}

func (s *Service) CreateFolder(path string) error {
	return s.fileSvc.CreateFolder(path)
}

func (s *Service) Rename(oldPath, newName string) error {
	return s.fileSvc.Rename(oldPath, newName)
}

func (s *Service) BuildTree(rootPath string, depth int) (*model.TreeNode, error) {
	return s.fileSvc.BuildTree(rootPath, depth)
}

func (s *Service) UploadEnabled() bool {
	return s.cfg.Upload.Enabled
}

func (s *Service) GetServerInfo() map[string]interface{} {
	return map[string]interface{}{
		"ip":          s.serverIP,
		"httpPort":    s.httpPort,
		"ftpPort":     s.ftpPort,
		"root":        s.root,
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
		"HTTPAuthPass":  s.cfg.HTTPAuth.Password,
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

// Stat 获取文件信息
func (s *Service) Stat(path string) (os.FileInfo, error) {
	return s.fileSvc.Stat(path)
}
