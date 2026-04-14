package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AnastasiaDMW/auth-service/internal/auth"
	"github.com/AnastasiaDMW/auth-service/internal/dto"
	"github.com/AnastasiaDMW/auth-service/internal/kafka"
	"github.com/AnastasiaDMW/auth-service/internal/store"
)

func (h *Handler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	var req dto.ChangeEmailRequest

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

	claims, ok := auth.GetClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user, err := h.UserRepo.FindById(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if err := user.ComparePassword(req.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	err = h.UserRepo.UpdateEmail(userID, req.NewEmail)
	if errors.Is(err, store.ErrEmailExists) {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	event := kafka.UserChangedEvent{
		ID:    userID,
		Email: req.NewEmail,
	}

	payload, _ := json.Marshal(event)

	if err := h.Producer.Send(kafka.UserChangedTopic, strconv.Itoa(userID), payload); err != nil {
		h.Logger.Debug("Failed to send kafka event in ChangeEmail", "error", err)
	}
	h.Logger.Debug("Send messages")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email updated",
	})
}
