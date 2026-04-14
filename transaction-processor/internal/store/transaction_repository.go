package store

import (
	"log/slog"

	"github.com/AnastasiaDMW/transaction-processor/internal/dto"
)

type TransactionRepository struct {
	store  *PostgresDB
	logger *slog.Logger
}

func NewTransactionRepository(s *PostgresDB, logger *slog.Logger) *TransactionRepository {
	return &TransactionRepository{store: s, logger: logger}
}

func (r *TransactionRepository) AddLedger(e dto.Ledger) (int64, error) {
	var ledgerId int64
	err := r.store.DB.QueryRow(
		"INSERT INTO ledgers(account_number, amount, idempotent) VALUES($1, $2, $3) RETURNING id",
		e.AccountNumber,
		e.Amount,
		e.Idempotent,
	).Scan(&ledgerId)
	if err != nil {
		r.logger.Debug("failed to insert ledger", "error", err)
		return 0, err
	}

	return ledgerId, nil
}

func (r *TransactionRepository) CheckIdempotent(idempotent string) (bool, error) {
	var exists bool

	err := r.store.DB.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM ledgers WHERE idempotent = $1)",
		idempotent,
	).Scan(&exists)

	if err != nil {
		r.logger.Error("Failed to check idempotent", "error", err)
		return false, err
	}

	return exists, nil
}
