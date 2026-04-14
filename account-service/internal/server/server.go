package server

import (
	"context"
	"net/http"

	"github.com/AnastasiaDMW/account-service/internal/auth"
	"github.com/AnastasiaDMW/account-service/internal/handler"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func New(bindAddr string, h *handler.Handler, authMw *auth.AuthMiddleware) *Server {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(api chi.Router) {

		api.Use(authMw.Middleware)

		api.Route("/accounts", func(accounts chi.Router) {
			accounts.Get("/", h.GetAccounts)
			accounts.Post("/", h.AddAccount)
			accounts.Patch("/status", h.UpdateAccountStatus)
			accounts.Post("/valid", h.Valid)
			accounts.Get("/balance/{accountNumber}", h.GetBalance)
		})

		api.Route("/user-profiles", func(up chi.Router) {
			up.Get("/", h.GetUserProfiles)
			up.Patch("/", h.UpdateUserProfile)
		})
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    bindAddr,
			Handler: r,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}