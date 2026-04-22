package handler

import (
	"crypto/rsa"
	"log/slog"
	"os"

	"github.com/AnastasiaDMW/auth-service/internal/event"
	"github.com/AnastasiaDMW/auth-service/internal/store/postgresstore"
	"github.com/AnastasiaDMW/auth-service/internal/store/redisstore"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	UserRepo   *postgresstore.UserRepository
	Logger     *slog.Logger
	TokenStore redisstore.TokenStore
	Producer event.Producer
}

var publicKey *rsa.PublicKey

func InitGatewayKey() error {
	pub := os.Getenv("PUBLIC_KEY")

	keyPublic, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pub))
	if err != nil {
		return err
	}

	publicKey = keyPublic
	return nil
}