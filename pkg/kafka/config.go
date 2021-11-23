package kafka

import (
	"time"
)

const (
	defaultTimeout               = 6 * time.Second // default timeout
	defaultAutoOffsetResetPolicy = OffsetEarliest
)

type Offset string

const (
	OffsetLatest   Offset = "latest"
	OffsetEarliest Offset = "earliest"
	OffsetNone     Offset = "none"
)

func defaultConfig() *config {
	return &config{
		timeout:               defaultTimeout,
		autoOffsetResetPolicy: defaultAutoOffsetResetPolicy,
	}
}

type config struct {
	timeout               time.Duration
	autoOffsetResetPolicy Offset
}
