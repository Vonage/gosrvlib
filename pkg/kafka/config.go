package kafka

import (
	"time"
)

// Offset points to where Kafka should start to read messages from.
type Offset string

const (
	OffsetLatest   Offset = "latest"
	OffsetEarliest Offset = "earliest"
	OffsetNone     Offset = "none"

	defaultTimeout               = 6 * time.Second
	defaultAutoOffsetResetPolicy = OffsetEarliest
)

type config struct {
	timeout               time.Duration
	autoOffsetResetPolicy Offset
}

func defaultConfig() *config {
	return &config{
		timeout:               defaultTimeout,
		autoOffsetResetPolicy: defaultAutoOffsetResetPolicy,
	}
}
