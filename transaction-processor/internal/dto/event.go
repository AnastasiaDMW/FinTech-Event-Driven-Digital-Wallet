package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Transaction struct {
	Id          int64           `json:"id"`
	AccountFrom string          `json:"accountFrom"`
	AccountTo   string          `json:"accountTo"`
	Amount      string          `json:"amount"`
	Idempotent  string          `json:"idempotent"`
	Type        TransactionType `json:"type"`
	EventType   EventType       `json:"eventType"`
}

type EventType string

const (
	EventCreated   EventType = "CREATED"
	EventChanged   EventType = "CHANGED"
	EventFailed    EventType = "FAILED"
	EventProcessed EventType = "PROCESSED"
	EventReserved  EventType = "RESERVED"
)

type TransactionType string

const (
	TransactionTransfer TransactionType = "TRANSFER"
	TransactionWithdraw TransactionType = "WITHDRAW"
	TransactionDeposit  TransactionType = "DEPOSIT"
)

func (t *Transaction) Validate() error {
	return validation.ValidateStruct(
		t,

		validation.Field(&t.AccountFrom,
			validation.Required,
			is.UUIDv4.Error("accountFrom must be UUID"),
		),

		validation.Field(&t.AccountTo,
			validation.Required,
			is.UUIDv4.Error("accountTo must be UUID"),
		),

		validation.Field(&t.Idempotent,
			validation.Required,
			is.UUIDv4.Error("idempotent must be UUID"),
		),

		validation.Field(&t.Amount,
			validation.Required,
		),

		validation.Field(&t.Type,
			validation.Required,
		),

		validation.Field(&t.EventType,
			validation.When(t.EventType != EventReserved,
				validation.By(func(value interface{}) error {
					return nil
				}),
			),
		),
	)
}

func (t *Transaction) IsProcessable() bool {
	return t.EventType == EventReserved
}
