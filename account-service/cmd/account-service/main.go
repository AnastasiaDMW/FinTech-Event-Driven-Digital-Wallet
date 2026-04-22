package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnastasiaDMW/account-service/internal/auth"
	"github.com/AnastasiaDMW/account-service/internal/dto"
	"github.com/AnastasiaDMW/account-service/internal/handler"
	"github.com/AnastasiaDMW/account-service/internal/kafka"
	"github.com/AnastasiaDMW/account-service/internal/logger"
	"github.com/AnastasiaDMW/account-service/internal/server"
	"github.com/AnastasiaDMW/account-service/internal/store"
	"github.com/joho/godotenv"
)

const tomlPath = "./config/accountservice.toml"

var kafkaAddress = []string{"kafka1:9092", "kafka2:9092", "kafka3:9092"}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not found, using system env")
	}

	cfg, err := server.LoadConfig(tomlPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logg := logger.New(cfg.LogLevel, cfg.LogFormat)

	logg.Debug("Starting service")

	publicKey, err := auth.LoadPublicKey()
	if err != nil {
		log.Fatal(err)
	}

	authMw := &auth.AuthMiddleware{
		PublicKey: publicKey,
	}

	db, err := store.New(cfg.Store.DatabaseURL)
	if err != nil {
		logg.Debug("Failed to connect database", "error", err)
		log.Fatal(err)
	}

	accountRepo := store.NewAccountRepository(db, logg)
	h := handler.New(accountRepo, logg)

	srv := server.New(cfg.BindAddr, h, authMw)

	consumer, err := kafka.NewConsumer(kafkaAddress, cfg.KafkaGroupId, logg)
	if err != nil {
		logg.Debug("Failed to create kafka consumer", "error", err)
		log.Fatal(err)
	}

	if err := consumer.Subscribe(cfg.KafkaUserChangedTopic); err != nil {
		logg.Debug("Failed to subscribe kafka topic", "error", err)
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		consumer.Start(ctx, func(topic string, raw []byte) error {
			var e dto.ChangedUserEvent

			if err := json.Unmarshal(raw, &e); err != nil {
				logg.Debug("Invalid kafka message", "error", err)
				return err
			}

			if err := h.HandleUserChanged(e); err != nil {
				logg.Debug("Failed to handle event", "error", err)
				return err
			}

			return nil
		})
	}()

	go func() {
		logg.Info("http server started", "addr", cfg.BindAddr)

		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logg.Debug("http server failed", "error", err)
		}
	}()

	<-ctx.Done()

	logg.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := consumer.Close(); err != nil {
		logg.Error("kafka shutdown error", "error", err)
	}

	if err := srv.Stop(shutdownCtx); err != nil {
		logg.Error("http shutdown error", "error", err)
	}

	if err := db.DB.Close(); err != nil {
		logg.Error("db close error", "error", err)
	}
}
