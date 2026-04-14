package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100)),
	)
}

type TokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (r *TokenRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.RefreshToken, validation.Required),
	)
}

type ChangeEmailRequest struct {
	NewEmail string `json:"newEmail"`
	Password string `json:"password"`
}

func (r *ChangeEmailRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.NewEmail, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required),
	)
}
