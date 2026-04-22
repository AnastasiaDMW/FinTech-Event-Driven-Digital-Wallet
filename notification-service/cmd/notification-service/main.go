package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnastasiaDMW/notification-service/internal/config"
	"github.com/AnastasiaDMW/notification-service/internal/handler"
	"github.com/AnastasiaDMW/notification-service/internal/kafka"
	"github.com/AnastasiaDMW/notification-service/internal/logger"
	"github.com/AnastasiaDMW/notification-service/internal/model"
	"github.com/AnastasiaDMW/notification-service/internal/notifier"
	"github.com/AnastasiaDMW/notification-service/internal/store"
	"github.com/joho/godotenv"
)

const tomlPath = "./config/notificationservice.toml"

var kafkaAddress = []string{"kafka1:9092", "kafka2:9092", "kafka3:9092"}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadConfig(tomlPath)
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New(cfg.LogLevel, cfg.LogFormat)

	db, err := store.New(cfg.Store.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	notifRepo := store.NewNotificationRepository(db)

	notif := notifier.New(logg, notifier.SMTPConfig{
		Host:     mustEnv("SMTP_HOST"),
		Port:     mustEnv("SMTP_PORT"),
		Username: mustEnv("SMTP_USERNAME"),
		Password: mustEnv("SMTP_PASSWORD"),
		From:     mustEnv("SMTP_FROM"),
	})

	h := handler.New(notifRepo, notif, logg)

	consumer, err := kafka.NewConsumer(kafkaAddress, cfg.KafkaGroupId, logg)
	if err != nil {
		log.Fatal(err)
	}

	consumer.SubscribeTopics([]string{
		cfg.KafkaUserChangedTopic,
		cfg.KafkaSendNotifyTopic,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		consumer.Start(ctx, func(topic string, raw []byte) error {
			switch topic {

			case cfg.KafkaUserChangedTopic:
				var e model.ChangedUserEvent
				if err := json.Unmarshal(raw, &e); err != nil {
					return err
				}
				return h.HandleUserChanged(e)

			case cfg.KafkaSendNotifyTopic:
				var e model.TransactionEvent
				if err := json.Unmarshal(raw, &e); err != nil {
					return err
				}
				return h.HandleTransaction(e)
			}

			return nil
		})
	}()

	<-ctx.Done()

	if err := consumer.Close(); err != nil {
		logg.Debug("Failed to close kafka consumer", "error", err)
	}

	if err := db.DB.Close(); err != nil {
		logg.Debug("Failed to close db", "error", err)
	}
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env %s is required", key)
	}
	return val
}
