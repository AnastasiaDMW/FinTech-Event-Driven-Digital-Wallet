package redisstore

import "time"

type TokenStore interface {
	SaveRefreshToken(token string, userID string, ttl time.Duration) error
	GetRefreshToken(token string) (string, error)
	DeleteRefreshToken(token string, userID string) error
	DeleteAllUserTokens(userID string) error
	HasUserTokens(userID string) (bool, error)
}
