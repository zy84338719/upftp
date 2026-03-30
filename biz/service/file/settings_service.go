package file

import (
	"github.com/zy84338719/upftp/pkg/conf"
	"github.com/zy84338719/upftp/pkg/errors"
)

func (s *Service) GetSettings() map[string]interface{} {
	return map[string]interface{}{
		"language":     s.cfg.GetLanguage(),
		"httpAuth":     s.cfg.HTTPAuth.Enabled,
		"httpAuthUser": s.cfg.HTTPAuth.Username,
		"httpAuthSet":  s.cfg.HTTPAuth.Password != "",
		"ftpUser":      s.cfg.Username,
		"ftpSet":       s.cfg.Password != "",
		"ftpEnabled":   s.cfg.EnableFTP,
		"mcpEnabled":   s.cfg.EnableMCP,
		"httpPort":     s.cfg.GetHTTPPort(),
		"ftpPort":      s.cfg.GetFTPPort(),
	}
}

func (s *Service) SetLanguage(language string) error {
	if language != "en" && language != "zh" {
		return errors.ErrInvalidLanguage
	}
	s.cfg.Language = language
	return conf.SaveConfig()
}

func (s *Service) SetHTTPAuth(enabled bool, username, password string) error {
	s.cfg.HTTPAuth.Enabled = enabled
	if username != "" {
		s.cfg.HTTPAuth.Username = username
	}
	if password != "" {
		s.cfg.HTTPAuth.Password = password
	}
	return conf.SaveConfig()
}

func (s *Service) SetFTP(username, password string) error {
	if username != "" {
		s.cfg.Username = username
	}
	if password != "" {
		s.cfg.Password = password
	}
	return conf.SaveConfig()
}
