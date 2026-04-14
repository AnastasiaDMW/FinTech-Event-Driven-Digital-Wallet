package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/AnastasiaDMW/account-service/internal/auth"
	"github.com/go-chi/chi"
)

const (
	urlBalanceService = "http://localhost:8085/api/v1/balance/"
	waitTimeout       = 4
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(auth.UserIdKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accountNumber := chi.URLParam(r, "accountNumber")
	if accountNumber == "" {
		http.Error(w, "AccountNumber is required", http.StatusBadRequest)
		return
	}

	isValid, err := h.repo.IsValidAccountNumberFrom(userId, accountNumber)
	if err != nil {
		h.logger.Debug("Failed to validate account", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isValid {
		http.Error(w, "Account not found or not yours", http.StatusForbidden)
		return
	}

	url := urlBalanceService + accountNumber

	client := &http.Client{
		Timeout: waitTimeout * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		h.logger.Debug("Failed to create request", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	resp, err := client.Do(req)
	if err != nil {
		h.logger.Debug("Failed to call balance service", "error", err)
		http.Error(w, "Balance service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get balance", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}
