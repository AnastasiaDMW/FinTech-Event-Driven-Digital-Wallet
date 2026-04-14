package store

import "errors"

var (
	ErrAccountLimit = errors.New("Account limit exceeded")
	ErrGenerateUniqueAccNum = errors.New("Failed to generate unique account number after retries")
	ErrAccountNotFound = errors.New("Account not found")
)
