package handler

import (
	"log/slog"

	"github.com/AnastasiaDMW/notification-service/internal/model"
	"github.com/AnastasiaDMW/notification-service/internal/notifier"
	"github.com/AnastasiaDMW/notification-service/internal/store"
)

type Handler struct {
	repo     *store.NotificationRepository
	notifier *notifier.Notifier
	logger   *slog.Logger
}

func New(repo *store.NotificationRepository, notifier *notifier.Notifier, logger *slog.Logger) *Handler {
	return &Handler{
		repo:     repo,
		notifier: notifier,
		logger:   logger,
	}
}

func (h *Handler) HandleUserChanged(e model.ChangedUserEvent) error {
	err := h.repo.AddUser(e.UserID, e.Email)
	if err != nil {
		h.logger.Debug("Failed to upsert user",
			"error", err,
			"user_id", e.UserID,
			"email", e.Email,
		)
		return err
	}
	h.logger.Debug("User saved/updated", "user_id", e.UserID, "email", e.Email)
	return nil
}

func (h *Handler) HandleTransaction(e model.TransactionEvent) error {
	email, err := h.repo.GetEmailByUserID(e.UserID)
	if err != nil {
		h.logger.Debug("user email not found", "error", err, "user_id", e.UserID)
		return err
	}

	err = h.notifier.Send(email, e)
	if err != nil {
		h.logger.Debug("failed to send notification", "error", err, "user_id", e.UserID, "email", email)
		return err
	}

	h.logger.Debug("Notification sent",
		"user_id", e.UserID,
		"email", email,
		"status", e.Status,
		"amount", e.Amount)
	return nil
}
