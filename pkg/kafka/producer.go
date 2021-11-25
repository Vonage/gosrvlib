package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer represents a wrapper around kafka.Producer.
type Producer struct {
	cfg    *config
	client *kafka.Producer
}

// NewProducer creates a new instance of Producer.
func NewProducer(urls []string, opts ...Option) (*Producer, error) {
	cfg := defaultConfig()

	for _, applyOpt := range opts {
		applyOpt(cfg)
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       strings.Join(urls, ","),
		"delivery.timeout.ms":     int(cfg.timeout.Milliseconds()),
		"go.produce.channel.size": cfg.produceChannelSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new kafka producer: %w", err)
	}

	return &Producer{
		cfg:    cfg,
		client: producer,
	}, nil
}

// Close cleans up Producer's internal resources.
func (p *Producer) Close() {
	p.client.Close()
}

// ProduceMessage sends a message to Kafka topic.
func (p *Producer) ProduceMessage(topic string, msg []byte) error {
	err := p.client.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: msg,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to send kafka message: %w", err)
	}

	return nil
}
