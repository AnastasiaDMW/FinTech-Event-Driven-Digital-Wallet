package server

import (
	"context"
	"net/http"

	"github.com/AnastasiaDMW/auth-service/internal/auth"
	"github.com/AnastasiaDMW/auth-service/internal/handler"
)

type Server struct {
	httpServer *http.Server
}

func New(bindAddr string, h *handler.Handler) *Server {
	rootMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/login", h.Login)
	apiMux.HandleFunc("/signup", h.SignUp)
	apiMux.HandleFunc("/refresh", h.Refresh)
	apiMux.HandleFunc("/logout", h.Logout)

	apiMux.Handle("/profiles/email", auth.JWTMiddleware(h.TokenStore)(http.HandlerFunc(h.ChangeEmail)))
	apiMux.Handle("/profiles/verify-email", auth.JWTMiddleware(h.TokenStore)(http.HandlerFunc(h.VerifyEmail)))

	rootMux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	return &Server{
		httpServer: &http.Server{
			Addr:    bindAddr,
			Handler: rootMux,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}