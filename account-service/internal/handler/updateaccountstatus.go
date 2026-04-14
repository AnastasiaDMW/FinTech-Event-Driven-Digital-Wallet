package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AnastasiaDMW/account-service/internal/dto"
)

func (h *Handler) UpdateAccountStatus(w http.ResponseWriter, r *http.Request) {
	var req dto.AccountStatusRequest

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

	err := h.repo.UpdateAccountStatus(req.AccountNumber, req.Status)
	if err != nil {
		h.logger.Debug("Failed to update account status", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Debug("Account status updated", "account_number", req.AccountNumber, "status", req.Status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "updated",
	})
}
