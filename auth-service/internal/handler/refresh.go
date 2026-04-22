package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AnastasiaDMW/auth-service/internal/auth"
	"github.com/AnastasiaDMW/auth-service/internal/dto"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.TokenRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.TokenStore.GetRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(publicKey, req.RefreshToken, "refresh")
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if claims.UserID != userID {
		http.Error(w, "Token mismatch", http.StatusUnauthorized)
		return
	}

	tokens, err := auth.GenerateTokenPair(userID, claims.Email, claims.Role)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	_ = h.TokenStore.DeleteRefreshToken(req.RefreshToken, userID)

	_ = h.TokenStore.SaveRefreshToken(
		tokens.RefreshToken,
		userID,
		7*24*time.Hour,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
