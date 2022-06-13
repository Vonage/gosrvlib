package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type consumerClient interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

// Consumer represents a wrapper around kafka.Consumer.
type Consumer struct {
	cfg    *config
	client consumerClient
}

// NewConsumer creates a new instance of Consumer.
func NewConsumer(urls []string, topic, groupID string, opts ...Option) (*Consumer, error) {
	cfg := defaultConfig()

	for _, applyOpt := range opts {
		applyOpt(cfg)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: urls,
		Topic:   topic,
		GroupID: groupID,
		MaxWait: 1,
	})

	return &Consumer{cfg: cfg, client: r}, nil
}

// Close cleans up Consumer's internal resources.
func (c *Consumer) Close() error {
	return c.client.Close() // nolint: wrapcheck
}

// ReadMessage reads one message from the Kafka; is blocked if no messages in the queue.
func (c *Consumer) ReadMessage(ctx context.Context) ([]byte, error) {
	msg, err := c.client.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read kafka message: %w", err)
	}

	return msg.Value, nil
}
