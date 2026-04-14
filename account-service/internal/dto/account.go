package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ValidRequest struct {
	AccountTo   string `json:"accountTo"`
	AccountFrom string `json:"accountFrom"`
}

type ValidRespone struct {
	IsValidAccountTo   bool `json:"isValidAccountTo"`
	IsValidAccountFrom bool `json:"isValidAccountFrom"`
}

func (v *ValidRequest) AllIsEmpty() bool {
	return v.AccountFrom == "" && v.AccountTo == ""
}

type AccountStatusRequest struct {
	AccountNumber string `json:"accountNumber"`
	Status        string `json:"status"`
}

func (r *AccountStatusRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.AccountNumber, validation.Required),
		validation.Field(&r.Status, validation.Required),
	)
}
