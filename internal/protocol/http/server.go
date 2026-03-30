package http

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	hertzConfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zy84338719/upftp/internal/conf"
	"github.com/zy84338719/upftp/internal/handler"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/service"
)

type HTTPServer struct {
	h   *server.Hertz
	svc *service.Service
}

func NewHTTPServer(svc *service.Service) *HTTPServer {
	return &HTTPServer{svc: svc}
}

func (s *HTTPServer) GetEngine() *route.Engine {
	if s.h == nil {
		s.h = server.Default()
	}
	return s.h.Engine
}

func (s *HTTPServer) Start(ctx context.Context) error {
	ip := s.svc.ServerIP()
	httpPort := s.svc.HTTPPort()

	handler.SetService(s.svc)

	if s.h == nil {
		opts := []hertzConfig.Option{
			server.WithHostPorts(conf.AppConfig.HTTPAddr()),
			server.WithDisablePrintRoute(true),
		}

		if conf.AppConfig.HTTPS.Enabled {
			cert, err := tls.LoadX509KeyPair(conf.AppConfig.HTTPS.CertFile, conf.AppConfig.HTTPS.KeyFile)
			if err == nil {
				opts = append(opts, server.WithTLS(&tls.Config{
					Certificates: []tls.Certificate{cert},
				}))
			}
		}

		s.h = server.Default(opts...)
	}

	handler.RegisterRoutes(s.h.Engine)

	go func() {
		<-ctx.Done()
		logger.Info("Stopping HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.h.Shutdown(shutdownCtx)
	}()

	protocol := "HTTP"
	if conf.AppConfig.HTTPS.Enabled {
		protocol = "HTTPS"
	}

	logger.Info("%s server starting on %s:%d", protocol, ip, httpPort)
	if conf.AppConfig.HTTPS.Enabled {
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
