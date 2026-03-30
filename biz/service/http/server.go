package http

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	hertzConfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zy84338719/upftp/biz/router"
	"github.com/zy84338719/upftp/biz/service/file"
	"github.com/zy84338719/upftp/pkg/conf"
	"github.com/zy84338719/upftp/pkg/logger"
)

type HTTPServer struct {
	h   *server.Hertz
	svc *file.Service
}

func NewHTTPServer(svc *file.Service) *HTTPServer {
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

	// 设置静态文件服务
	s.h.Static("/static", "./biz/handler/index/templates")

	// 使用新生成的路由逻辑
	router.GeneratedRegister(s.h)

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
