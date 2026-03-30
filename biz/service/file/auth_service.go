package file

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

func (s *Service) ValidateBasicAuth(user, pass string) bool {
	return user == s.cfg.HTTPAuth.Username && pass == s.cfg.HTTPAuth.Password
}

func (s *Service) IsHTTPAuthEnabled() bool {
	return s.cfg.HTTPAuth.Enabled
}
