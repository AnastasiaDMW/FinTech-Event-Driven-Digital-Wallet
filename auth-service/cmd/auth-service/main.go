package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnastasiaDMW/auth-service/internal/auth"
	"github.com/AnastasiaDMW/auth-service/internal/handler"
	"github.com/AnastasiaDMW/auth-service/internal/kafka"
	"github.com/AnastasiaDMW/auth-service/internal/logger"
	"github.com/AnastasiaDMW/auth-service/internal/server"
	"github.com/AnastasiaDMW/auth-service/internal/store/postgresstore"
	"github.com/AnastasiaDMW/auth-service/internal/store/redisstore"
)

const tomlPath = "./config/authservice.toml"

var kafkaAddress = []string{"localhost:19092", "localhost:29092", "localhost:39092"}

func main() {
	cfg, err := server.LoadConfig(tomlPath)
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New(cfg.LogLevel, cfg.LogFormat)

	err = auth.InitKeys()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgresstore.New(cfg.Store.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgresstore.NewUserRepository(db)

	redisStore, err := redisstore.NewRedisTokenStore(
		cfg.RedisAddr,
		cfg.RedisUsername,
		cfg.RedisPassword,
	)
	if err != nil {
		logg.Debug(err.Error())
		log.Fatal(err)
	}
	logg.Debug("Redis connected successfully")

	producer, err := kafka.NewProducer(kafkaAddress)
	if err != nil {
		logg.Debug("Error created producer")
		log.Fatal(err)
	}

	h := &handler.Handler{
		UserRepo:   userRepo,
		Logger:     logg,
		TokenStore: redisStore,
		Producer:   producer,
	}

	srv := server.New(cfg.BindAddr, h)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logg.Debug("Server starting", "addr", cfg.BindAddr)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logg.Error("Server failed", "error", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logg.Debug("Shutting down...")
	srv.Stop(shutdownCtx)
	producer.Close()
	db.DB.Close()
}
