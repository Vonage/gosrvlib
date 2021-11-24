package kafka

import (
	"time"
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

	defaultTimeout               = 6 * time.Second
	defaultAutoOffsetResetPolicy = OffsetEarliest
	defaultProduceChannelSize    = 10_000
)

type config struct {
	timeout               time.Duration
	autoOffsetResetPolicy Offset
	produceChannelSize    int
}

func defaultConfig() *config {
	return &config{
		timeout:               defaultTimeout,
		autoOffsetResetPolicy: defaultAutoOffsetResetPolicy,
		produceChannelSize:    defaultProduceChannelSize,
	}
}
