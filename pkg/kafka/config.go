package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Offset points to where Kafka should start to read messages from.
type Offset string

const (
	// OffsetLatest automatically reset the offset to the latest offset.
	OffsetLatest Offset = "latest"
	// OffsetEarliest automatically reset the offset to the earliest offset.
	OffsetEarliest Offset = "earliest"
	// OffsetNone throw an error to the consumer if no previous offset is found for the consumer's group.
	OffsetNone Offset = "none"
)

type config struct {
	*kafka.ConfigMap
}

func defaultConfig() *config {
	return &config{
		&kafka.ConfigMap{
			"auto.offset.reset":  string(OffsetEarliest),
			"session.timeout.ms": int((6 * time.Second).Milliseconds()),
		},
	}
}
