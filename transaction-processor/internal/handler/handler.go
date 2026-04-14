package handler

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/AnastasiaDMW/transaction-processor/internal/dto"
	"github.com/AnastasiaDMW/transaction-processor/internal/event"
	"github.com/AnastasiaDMW/transaction-processor/internal/store"
)

type Handler struct {
	repo     *store.TransactionRepository
	logger   *slog.Logger
	Producer event.Producer
}

func New(repo *store.TransactionRepository, logger *slog.Logger, producer event.Producer) *Handler {
	return &Handler{
		repo:     repo,
		logger:   logger,
		Producer: producer,
	}
}

func (h *Handler) HandleTransaction(e dto.Transaction, topic string) error {

	if err := e.Validate(); err != nil {
		h.logger.Debug("Invalid transaction", "error", err.Error())
		return err
	}

	if !e.IsProcessable() {
		return nil
	}

	exist, err := h.repo.CheckIdempotent(e.Idempotent)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	hasTo := e.AccountTo != ""
	hasFrom := e.AccountFrom != ""

	add := func(account string, amount string) (int64, error) {
		return h.repo.AddLedger(dto.Ledger{
			AccountNumber: account,
			Amount:        amount,
			Idempotent:    e.Idempotent,
		})
	}

	if !hasFrom && hasTo {
		id, err := add(e.AccountTo, e.Amount)
		if err != nil {
			h.logger.Debug("Failed deposit ledger", "error", err)
			return err
		}

		h.logger.Debug("Deposit ledger created", "id", id)
		return nil
	}

	if hasFrom && !hasTo {
		id, err := add(e.AccountFrom, "-"+e.Amount)
		if err != nil {
			h.logger.Debug("Failed withdraw ledger", "error", err)
			return err
		}

		h.logger.Debug("Withdraw ledger created", "id", id)
		return nil
	}

	if hasFrom && hasTo {
		fromID, err := add(e.AccountFrom, "-"+e.Amount)
		if err != nil {
			h.logger.Debug("Failed transfer (from)", "error", err)
			return err
		}

		toID, err := add(e.AccountTo, e.Amount)
		if err != nil {
			h.logger.Debug("Failed transfer (to)", "error", err)
			return err
		}

		h.logger.Debug("Transfer ledger created",
			"from_id", fromID,
			"to_id", toID,
		)

		return nil
	}

	err = h.sendMessage(e, topic)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) sendMessage(e dto.Transaction, topic string) error {
	event := dto.Transaction{
		Id:          e.Id,
		AccountFrom: e.AccountFrom,
		AccountTo:   e.AccountTo,
		Amount:      e.Amount,
		Idempotent:  e.Idempotent,
		Type:        e.Type,
		EventType:   dto.EventProcessed,
	}

	payload, _ := json.Marshal(event)
	if err := h.Producer.Send(topic, strconv.FormatInt(e.Id, 10), payload); err != nil {
		h.logger.Debug("Failed to send kafka event in Transaction Processor", "error", err)
		return err
	}
	h.logger.Debug("Send messages")

	return nil
}
