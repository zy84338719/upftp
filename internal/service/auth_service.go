package service

import (
	"github.com/zy84338719/upftp/internal/auth"
	"github.com/zy84338719/upftp/internal/logger"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type LoginResult struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Expires  int64  `json:"expires"`
}

func (s *Service) Authenticate(username, password string) bool {
	return username == s.cfg.HTTPAuth.Username && password == s.cfg.HTTPAuth.Password
}

func (s *Service) Login(input LoginInput) (*LoginResult, error) {
	session, err := s.sessions.CreateSession(input.Username)
	if err != nil {
		return nil, err
	}
	logger.Info("User logged in: %s", input.Username)
	return &LoginResult{
		Token:    session.Token,
		Username: session.Username,
		Expires:  session.ExpiresAt.Unix(),
	}, nil
}

func (s *Service) Logout(token string) {
	if token != "" {
		s.sessions.DeleteSession(token)
	}
	logger.Info("User logged out")
}

func (s *Service) ValidateSession(token string) (*auth.Session, bool) {
	return s.sessions.ValidateSession(token)
}

func (s *Service) ValidateBasicAuth(user, pass string) bool {
	return user == s.cfg.HTTPAuth.Username && pass == s.cfg.HTTPAuth.Password
}

func (s *Service) IsHTTPAuthEnabled() bool {
	return s.cfg.HTTPAuth.Enabled
}
