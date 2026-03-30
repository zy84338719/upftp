package server

import (
	"context"
	"crypto/tls"

	"github.com/cloudwego/hertz/pkg/app/server"
	hertzConfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/handlers"
	"github.com/zy84338719/upftp/internal/logger"
)

type HTTPServer struct {
	h *server.Hertz
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start(ctx context.Context, ip string, httpPort, ftpPort int, root string) error {
	handlers.SetServerInfo(ip, httpPort, ftpPort, root)

	opts := []hertzConfig.Option{
		server.WithHostPorts(config.AppConfig.HTTPAddr()),
		server.WithDisablePrintRoute(true),
	}

	if config.AppConfig.HTTPS.Enabled {
		cert, err := tls.LoadX509KeyPair(config.AppConfig.HTTPS.CertFile, config.AppConfig.HTTPS.KeyFile)
		if err == nil {
			opts = append(opts, server.WithTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
			}))
		}
	}

	s.h = server.Default(opts...)

	handlers.RegisterRoutes(s.h.Engine)

	go func() {
		<-ctx.Done()
		logger.Info("Stopping HTTP server...")
		s.h.Shutdown(ctx)
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

	err := s.h.Run()
	if err != nil && err != context.Canceled {
		return err
	}

	return nil
}
