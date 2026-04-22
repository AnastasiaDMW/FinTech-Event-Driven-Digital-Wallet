package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AnastasiaDMW/account-service/internal/auth"
	"github.com/AnastasiaDMW/account-service/internal/dto"
)

func (h *Handler) Valid(w http.ResponseWriter, r *http.Request) {
	var req dto.ValidRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AllIsEmpty() {
		http.Error(w, "Empty request", http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value(auth.UserIdKey).(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	hasTo := req.AccountTo != ""
	hasFrom := req.AccountFrom != ""

	var resp dto.ValidRespone

	if hasFrom && !hasTo {
		isValidFrom, err := h.repo.IsValidAccountNumberFrom(userId, req.AccountFrom)
		if err != nil {
			http.Error(w, "Failed to validate account_from", http.StatusInternalServerError)
			return
		}
		resp.IsValidAccountFrom = isValidFrom
		resp.IsValidAccountTo = true
	}

	if hasTo && !hasFrom {
		isValidTo, err := h.repo.IsValidAccountNumberTo(req.AccountTo)
		if err != nil {
			http.Error(w, "Failed to validate account_to", http.StatusInternalServerError)
			return
		}
		resp.IsValidAccountTo = isValidTo
		resp.IsValidAccountFrom = true
	}

	if hasTo && hasFrom {
		isValidFrom, err := h.repo.IsValidAccountNumberFrom(userId, req.AccountFrom)
		if err != nil {
			http.Error(w, "Failed to validate account_from", http.StatusInternalServerError)
			return
		}

		isValidTo, err := h.repo.IsValidAccountNumberTo(req.AccountTo)
		if err != nil {
			http.Error(w, "Failed to validate account_to", http.StatusInternalServerError)
			return
		}

		resp.IsValidAccountFrom = isValidFrom
		resp.IsValidAccountTo = isValidTo
	}

	h.logger.Debug("Validate successfull")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
