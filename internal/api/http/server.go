package http

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"
	"testAnalyticService/internal/worker"
)

type httpServer struct {
	host   string
	port   int
	worker worker.Worker
	logger *zap.Logger
}

func NewHTTPServer(host string, port int, worker worker.Worker, logger *zap.Logger) *httpServer {
	return &httpServer{
		host:   host,
		port:   port,
		worker: worker,
		logger: logger,
	}
}

func (s *httpServer) Start(ctx context.Context) error {
	s.registerRoutes()
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.logger.Info(fmt.Sprintf("Starting HTTP-server on %s", addr))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *httpServer) registerRoutes() {
	handlers := NewHandlers(s.worker, s.logger)
	http.HandleFunc("/analitycs", handlers.analitycsHandler)
}
