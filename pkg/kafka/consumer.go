package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer represents a wrapper around kafka.Consumer.
type Consumer struct {
	cfg    *config
	client *kafka.Consumer
}

// NewConsumer creates a new instance of Consumer.
func NewConsumer(urls, topics []string, groupID string, opts ...Option) (*Consumer, error) {
	cfg := defaultConfig()

	for _, applyOpt := range opts {
		applyOpt(cfg)
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(urls, ","),
		"group.id":           groupID,
		"auto.offset.reset":  string(cfg.autoOffsetResetPolicy),
		"session.timeout.ms": int(cfg.timeout.Milliseconds()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new kafka consumer: %w", err)
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe kafka topic: %w", err)
	}

	return &Consumer{
		cfg:    cfg,
		client: consumer,
	}, nil
}

// Close cleans up Consumer's internal resources.
func (c *Consumer) Close() error {
	return c.client.Close() // nolint: wrapcheck
}

// ReadMessage reads one message from the Kafka; is blocked if no messages in the queue.
func (c *Consumer) ReadMessage() ([]byte, error) {
	msg, err := c.client.ReadMessage(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read kafka message: %w", err)
	}

	return msg.Value, nil
}
