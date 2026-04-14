package store

import "errors"

var (
	ErrEmailExists     = errors.New("Email already exists")
	ErrUserNotFound    = errors.New("User not found")
	ErrRedisConnection = errors.New("Redis connection failed")
)
