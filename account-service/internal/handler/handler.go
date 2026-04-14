package handler

import (
	"log/slog"

	"github.com/AnastasiaDMW/account-service/internal/dto"
	"github.com/AnastasiaDMW/account-service/internal/store"
)

type Handler struct {
	repo   *store.AccountRepository
	logger *slog.Logger
}

func New(repo *store.AccountRepository, logger *slog.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

func (h *Handler) HandleUserChanged(e dto.ChangedUserEvent) error {
	if err := e.Validate(); err != nil {
		h.logger.Debug("Failed to create user profile", "error", err.Error())
		return err
	}

	if err := h.repo.CreateUserProfile(e); err != nil {
		h.logger.Debug("Failed to create user profile", "error", err)
		return err
	}

	if err := h.repo.CreateAccount(int64(e.UserID)); err != nil {
		h.logger.Debug("Failed to create account", "error", err)
		return err
	}

	h.logger.Debug("User profile and account created", "user_id", e.UserID)
	return nil
}
