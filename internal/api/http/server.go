package http

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"
	"testAnalyticService/internal"
)

type httpServer struct {
	host           string
	port           int
	analysticsRepo internal.AnalyticsRepository
	logger         *zap.Logger
}

func NewHTTPServer(host string, port int, analysticsRepo internal.AnalyticsRepository, logger *zap.Logger) *httpServer {
	return &httpServer{
		host:           host,
		port:           port,
		analysticsRepo: analysticsRepo,
		logger:         logger,
	}
}

func (s *httpServer) Start(ctx context.Context) error {
	s.registerRoutes()
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *httpServer) registerRoutes() {
	handlers := NewHandlers(s.analysticsRepo, s.logger)
	http.HandleFunc("/analitycs", handlers.analitycsHandler)
}
