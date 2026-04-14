package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ChangedUserEvent struct {
	UserID int    `json:"id"`
	Email  string `json:"email"`
}

func (r *ChangedUserEvent) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.UserID, validation.Required, validation.Min(int64(1))),
		validation.Field(&r.Email, validation.Required, is.Email),
	)
}
