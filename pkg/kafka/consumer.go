package kafka

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	cfg    *config
	client *kafka.Consumer
}

// NewConsumer creates a new instance of Consumer.
func NewConsumer(urls, topics []string, groupId string, opts ...Option) (*Consumer, error) {
	cfg := defaultConfig()

	for _, applyOpt := range opts {
		applyOpt(cfg)
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    strings.Join(urls, ","),
		"group.id":             groupId,
		"auto.offset.reset":    cfg.autoOffsetResetPolicy,
		"session.timeoutMs.ms": cfg.timeoutMs,
	})
	if err != nil {
		return nil, err
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		cfg:    cfg,
		client: consumer,
	}, nil
}

// ReadMessage reads one message from the Kafka; is blocked if no messages in the queue.
func (c *Consumer) ReadMessage() ([]byte, error) {
	msg, err := c.client.ReadMessage(-1)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}
