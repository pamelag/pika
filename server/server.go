package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pamelag/pika/insight"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Insights insight.Service
	router   chi.Router
}

// New returns a new HTTP server.
func New(is insight.Service) *Server {
	s := &Server{
		Insights: is,
	}

	r := chi.NewRouter()
	r.Use(secureHeaders)
	// set recoverer middleware
	r.Use(middleware.Recoverer)

	r.Route("/insights", func(r chi.Router) {
		ins := insightHandler{s.Insights}
		r.Mount("/v1", ins.router())
	})

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Set secure headers to the response
func secureHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Add("X-Frame-Options", "DENY")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("Content-Security-Policy", "frame-ancestors 'none'")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Encode error for response
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
