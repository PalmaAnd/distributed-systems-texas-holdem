package api

import (
	"log"
	"net/http"
	"os"
)

type Server struct {
	mux          *http.ServeMux
	allowedOrigin string
}

func New() *Server {
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://34.58.122.79"
	}
	s := &Server{
		mux:          http.NewServeMux(),
		allowedOrigin: allowedOrigin,
	}
	s.mux.HandleFunc("/api/v1/evaluate", s.handleEvaluate)
	s.mux.HandleFunc("/api/v1/compare", s.handleCompare)
	s.mux.HandleFunc("/api/v1/probability", s.handleProbability)
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// Explicitly allow the configured frontend origin
	if origin == s.allowedOrigin {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}

	// Handle OPTIONS preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 204
		return
	}

	s.mux.ServeHTTP(w, r)
}

func (s *Server) Listen(addr string) error {
	log.Printf("Server listening on %s", addr)
	return http.ListenAndServe(addr, s)
}
