package server

import (
	"context"
	"log"
	"net/http"

	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/handlers"
)

type HTTPServer struct {
	server *http.Server
	info   *handlers.ServerInfo
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
		log.Println("Stopping HTTP server...")
		s.server.Shutdown(context.Background())
	}()

	log.Printf("HTTP server starting on %s%s", ip, config.AppConfig.HTTPAddr())
	log.Printf("Web interface: http://%s%s", ip, config.AppConfig.HTTPAddr())

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
