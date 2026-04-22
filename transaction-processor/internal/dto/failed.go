package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Failed struct {
	Id           int64           `json:"id"`
	AccountFrom  string          `json:"accountFrom"`
	AccountTo    string          `json:"accountTo"`
	Amount       int64         `json:"amount"`
	Idempotent   string          `json:"idempotent"`
	MessageError string          `json:"messageError"`
	Type         TransactionType `json:"type"`
}

func (t *Failed) Validate() error {
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

		validation.Field(&t.MessageError,
			validation.Required,
		),
	)
}
