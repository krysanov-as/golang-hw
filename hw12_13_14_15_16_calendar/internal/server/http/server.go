package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

type Logger interface {
	Info(string)
	Error(string)
}

type Application interface{}

func NewServer(logger Logger, app Application, host string, port int) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return &Server{
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", host, port),
			Handler:           loggingMiddleware(mux),
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
