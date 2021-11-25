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

	_ = cfg.ConfigMap.SetKey("bootstrap.servers", strings.Join(urls, ","))

	producer, err := kafka.NewProducer(cfg.ConfigMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kafka producer: %w", err)
	}

	return &Producer{cfg: cfg, client: producer}, nil
}

// Close cleans up Producer's internal resources.
func (p *Producer) Close() {
	p.client.Close()
}

// ProduceMessage sends a message to Kafka topic.
func (p *Producer) ProduceMessage(topic string, msg []byte) error {
	return p.client.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: msg,
		},
		nil,
	)
}
