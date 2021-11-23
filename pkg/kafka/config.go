package kafka

import (
	"time"
)

const (
	defaultTimeout               = 5 * time.Second // default timeout
	defaultAutoOffsetResetPolicy = "earliest"
)

func defaultConfig() *config {
	return &config{
		timeout:               defaultTimeout,
		autoOffsetResetPolicy: defaultAutoOffsetResetPolicy,
	}
}

type config struct {
	timeout               time.Duration
	autoOffsetResetPolicy string
}
