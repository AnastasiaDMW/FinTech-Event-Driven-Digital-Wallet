package handler

import (
	"log/slog"

	"github.com/AnastasiaDMW/auth-service/internal/event"
	"github.com/AnastasiaDMW/auth-service/internal/store/postgresstore"
	"github.com/AnastasiaDMW/auth-service/internal/store/redisstore"
)

type Handler struct {
	UserRepo   *postgresstore.UserRepository
	Logger     *slog.Logger
	TokenStore redisstore.TokenStore
	Producer event.Producer
}
