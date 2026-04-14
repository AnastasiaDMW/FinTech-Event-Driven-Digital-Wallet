package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnastasiaDMW/transaction-processor/internal/config"
	"github.com/AnastasiaDMW/transaction-processor/internal/dto"
	"github.com/AnastasiaDMW/transaction-processor/internal/handler"
	"github.com/AnastasiaDMW/transaction-processor/internal/kafka"
	"github.com/AnastasiaDMW/transaction-processor/internal/logger"
	"github.com/AnastasiaDMW/transaction-processor/internal/store"
)

const tomlPath = "./config/transactionprocessor.toml"

var kafkaAddress = []string{"localhost:19092", "localhost:29092", "localhost:39092"}

func main() {
	cfg, err := config.LoadConfig(tomlPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logg := logger.New(cfg.LogLevel, cfg.LogFormat)

	logg.Debug("Starting service")

	db, err := store.New(cfg.Store.DatabaseURL)
	if err != nil {
		logg.Debug("Failed to connect database", "error", err)
		log.Fatal(err)
	}

	producer, err := kafka.NewProducer(kafkaAddress)
	if err != nil {
		logg.Debug("Error created producer")
		log.Fatal(err)
	}

	transactionRepo := store.NewTransactionRepository(db, logg)
	h := handler.New(transactionRepo, logg, producer)

	consumer, err := kafka.NewConsumer(kafkaAddress, cfg.KafkaGroupId, logg)
	if err != nil {
		logg.Debug("Failed to create kafka consumer", "error", err)
		log.Fatal(err)
	}

	if err := consumer.Subscribe(cfg.KafkaTransactionTopic); err != nil {
		logg.Debug("Failed to subscribe kafka topic", "error", err)
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		consumer.Start(ctx, func(topic string, raw []byte) error {
			var e dto.Transaction

			if err := json.Unmarshal(raw, &e); err != nil {
				logg.Debug("Invalid kafka message", "error", err)
				return err
			}

			if err := h.HandleTransaction(e, cfg.KafkaTransactionTopic); err != nil {
				logg.Debug("Failed to handle transaction", "error", err)
				return err
			}

			return nil
		})
	}()

	<-ctx.Done()

	logg.Debug("shutdown signal received")

	if err := consumer.Close(); err != nil {
		logg.Error("kafka shutdown error", "error", err)
	}

	if err := db.DB.Close(); err != nil {
		logg.Error("db close error", "error", err)
	}

	logg.Debug("service stopped")
}
