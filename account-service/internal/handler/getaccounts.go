package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AnastasiaDMW/account-service/internal/auth"
)

func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := r.Context().Value(auth.UserIdKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accounts, err := h.repo.GetAccountsByUserID(userId)
	if err != nil {
		h.logger.Debug("Failed to get accounts", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}
