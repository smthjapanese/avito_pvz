package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server представляет HTTP-сервер
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

// NewServer создает новый HTTP-сервер
func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
	}
}

// Run запускает HTTP-сервер
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown выполняет корректное завершение работы HTTP-сервера
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
