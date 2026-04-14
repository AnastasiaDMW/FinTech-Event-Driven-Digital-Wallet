package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AnastasiaDMW/auth-service/internal/dto"
	"github.com/AnastasiaDMW/auth-service/internal/kafka"
	"github.com/AnastasiaDMW/auth-service/internal/model"
	"github.com/AnastasiaDMW/auth-service/internal/store"
)

type SignUpResponse struct {
	Message string `json:"message"`
	UserId  int    `json:"user_id"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest

	if r.Method != http.MethodPost {
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

	user := &model.User{
		Email:             req.Email,
		EncryptedPassword: "",
		Password:          req.Password,
	}

	if err := user.BeforeCreate(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, err := h.UserRepo.CreateUser(user)
	if errors.Is(err, store.ErrEmailExists) {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	event := kafka.UserChangedEvent{
		ID:    id,
		Email: user.Email,
	}

	payload, _ := json.Marshal(event)

	if err := h.Producer.Send(kafka.UserChangedTopic, strconv.Itoa(id), payload); err != nil {
		h.Logger.Debug("Failed to send kafka event in SignUp", "error", err)
	}
	h.Logger.Debug("Send messages")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignUpResponse{
		Message: "Пользователь зарегистрирован!",
		UserId:  id,
	})
}
