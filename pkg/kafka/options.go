package kafka

import (
	"time"
)

// Option is a type alias for a function that configures Kafka client.
type Option func(*config)

// WithTimeout sets the timeout.
func WithTimeout(t time.Duration) Option {
	return func(cfg *config) {
		cfg.timeout = t
	}
}

// WithAutoOffsetResetPolicy sets respective parameter.
func WithAutoOffsetResetPolicy(p string) Option {
	return func(cfg *config) {
		cfg.autoOffsetResetPolicy = p
	}
}
