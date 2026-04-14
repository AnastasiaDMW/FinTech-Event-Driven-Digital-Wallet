package kafka

import (
	"context"
	"log/slog"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const sessionTimeout = 7000 //ms

type Consumer struct {
	consumer *kafka.Consumer
	logger   *slog.Logger
	stopCh   chan struct{}
}

func NewConsumer(address []string, groupID string, logger *slog.Logger) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 groupID,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		logger:   logger,
		stopCh:   make(chan struct{}),
	}, nil
}

func (c *Consumer) Subscribe(topic string) error {
	return c.consumer.SubscribeTopics([]string{topic}, nil)
}

func (c *Consumer) Start(
	ctx context.Context,
	handler func(topic string, msg []byte) error,
) {
	c.logger.Info("kafka consumer started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("kafka consumer stopping")
			_ = c.consumer.Close()
			return
		default:
		}

		event := c.consumer.Poll(100)
		if event == nil {
			continue
		}

		switch e := event.(type) {

		case *kafka.Message:
			topic := *e.TopicPartition.Topic

			if err := handler(topic, e.Value); err != nil {
				c.logger.Error("handler error",
					slog.String("topic", topic),
					slog.Any("error", err),
				)
			}
		}
	}
}

func (c *Consumer) Close() error {
	close(c.stopCh)
	return c.consumer.Close()
}