package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AnastasiaDMW/account-service/internal/auth"
	"github.com/AnastasiaDMW/account-service/internal/store"
)

func (h *Handler) AddAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	userId, ok := r.Context().Value(auth.UserIdKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.repo.CreateAccount(userId)
	if err != nil {
		if errors.Is(err, store.ErrAccountLimit) {
			http.Error(w, "Account limit reached", http.StatusBadRequest)
			return
		}

		h.logger.Debug("Failed to create account", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Debug("Account created", "user_id", userId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
