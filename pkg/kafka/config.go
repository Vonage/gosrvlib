package kafka

import (
	"time"
)

const (
	defaultTimeout               = 5 * time.Second // default timeoutMs
	defaultAutoOffsetResetPolicy = "earliest"
)

func defaultConfig() *config {
	return &config{
		timeoutMs:             defaultTimeout.Milliseconds(),
		autoOffsetResetPolicy: defaultAutoOffsetResetPolicy,
	}
}

type config struct {
	timeoutMs             int64
	autoOffsetResetPolicy string
}
