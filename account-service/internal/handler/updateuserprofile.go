package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AnastasiaDMW/account-service/internal/auth"
	"github.com/AnastasiaDMW/account-service/internal/dto"
	"github.com/AnastasiaDMW/account-service/internal/model"
)

func (h *Handler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserProfileRequest

	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value(auth.UserIdKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var birthDate *time.Time

	if req.BirthDate != nil {
		t, err := dto.ParseBirthDate(*req.BirthDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		birthDate = &t
	}

	err := h.repo.UpdateUserProfile(userId, model.UserProfile {
		BirthDate: birthDate,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		h.logger.Debug("Failed to update user profile", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Debug("User profile updated", "user_id", userId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "updated",
	})
}