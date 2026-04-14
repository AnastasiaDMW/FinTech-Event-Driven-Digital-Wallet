package auth

import "errors"

var (
	ErrInvalidToken     = errors.New("Invalid token")
	ErrInvalidTokenType = errors.New("Invalid token type")
	ErrInvalidAlg       = errors.New("Invalid alg type in token header")
)
