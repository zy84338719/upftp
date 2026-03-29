package server

import (
	"context"
	"net/http"

	"github.com/zy84338719/upftp/v2/internal/config"
	"github.com/zy84338719/upftp/v2/internal/handlers"
	"github.com/zy84338719/upftp/v2/internal/logger"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start(ctx context.Context, ip string, httpPort, ftpPort int, root string) error {
	handlers.SetServerInfo(ip, httpPort, ftpPort, root)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	s.server = &http.Server{
		Addr:    config.AppConfig.HTTPAddr(),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		logger.Info("Stopping HTTP server...")
		s.server.Shutdown(context.Background())
	}()

	protocol := "HTTP"
	if config.AppConfig.HTTPS.Enabled {
		protocol = "HTTPS"
	}

	logger.Info("%s server starting on %s:%d", protocol, ip, httpPort)
	if config.AppConfig.HTTPS.Enabled {
		logger.Info("Web interface: https://%s:%d", ip, httpPort)
	} else {
		logger.Info("Web interface: http://%s:%d", ip, httpPort)
	}

	var err error
	if config.AppConfig.HTTPS.Enabled {
		if config.AppConfig.HTTPS.CertFile == "" || config.AppConfig.HTTPS.KeyFile == "" {
			logger.Error("HTTPS enabled but cert_file or key_file not specified")
			return nil
		}
		err = s.server.ListenAndServeTLS(config.AppConfig.HTTPS.CertFile, config.AppConfig.HTTPS.KeyFile)
	} else {
		err = s.server.ListenAndServe()
	}

	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
