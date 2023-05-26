package server

import (
	"context"
	"net/http"
)

// Server is a struct that represents a server
type Server struct {
	srv http.Server
}

// NewServer returns a new Server
func NewServer() *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte("up"))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &Server{
		srv: srv,
	}
}

// Start starts the server
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
