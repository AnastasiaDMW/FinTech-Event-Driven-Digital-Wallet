package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AnastasiaDMW/account-service/internal/auth"
)

func (h *Handler) GetUserProfiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := r.Context().Value(auth.UserIdKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userProfiles, err := h.repo.GetUserProfileByUserID(userId)
	if err != nil {
		h.logger.Debug("Failed to get user profiles", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfiles)
}
