package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vtno/zypher/internal/server/handlers"
	"github.com/vtno/zypher/internal/store"
)

// Server is a struct that represents a server
type Server struct {
	srv   *http.Server
	store store.Store
}

type ServerOption func(*Server)

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.srv.Addr = fmt.Sprintf(":%d", port)
	}
}

const defaultDbPath = "zypher.db"

// NewServer returns a new Server
func NewServer(opts ...ServerOption) (*Server, error) {
	mux := http.NewServeMux()
	bbStore, err := store.NewBBoltStore(defaultDbPath)
	if err != nil {
		return nil, fmt.Errorf("error creating bbStore: %v", err)
	}
	kh := handlers.NewKeyHandler(bbStore)
	mux.HandleFunc("/key", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			kh.Get(w, r)
		case "POST":
			kh.Post(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte("up"))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	httpSrv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	srv := &Server{
		srv:   &httpSrv,
		store: bbStore,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv, nil
}

// Start starts the server
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	if err := s.store.Close(); err != nil {
		return fmt.Errorf("error closing store: %v", err)
	}

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %v", err)
	}

	return nil
}
