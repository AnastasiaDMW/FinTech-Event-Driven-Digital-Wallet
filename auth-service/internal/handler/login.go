package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AnastasiaDMW/auth-service/internal/auth"
	"github.com/AnastasiaDMW/auth-service/internal/dto"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.UserRepo.FindByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	if err := user.ComparePassword(req.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	tokens, err := auth.GenerateTokenPair(user.ID, user.Email, "user")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	h.Logger.Debug("USER ID", "id", user.ID)

	if err := h.TokenStore.DeleteAllUserTokens(user.ID); err != nil {
		h.Logger.Debug("DELETE TOKENS ERROR", "error", err)
	}

	err = h.TokenStore.SaveRefreshToken(
		tokens.RefreshToken,
		user.ID,
		7*24*time.Hour,
	)
	if err != nil {
		log.Printf("REDIS SAVE ERROR: %v", err)
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
