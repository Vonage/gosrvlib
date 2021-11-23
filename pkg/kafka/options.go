package kafka

import (
	"time"
)

// Option is a type alias for a function that configures Kafka client.
type Option func(*config)

// WithTimeout sets the timeoutMs.
func WithTimeout(t time.Duration) Option {
	return func(cfg *config) {
		cfg.timeoutMs = t.Milliseconds()
	}
}

// WithAutoOffsetResetPolicy sets respective parameter.
func WithAutoOffsetResetPolicy(p string) Option {
	return func(cfg *config) {
		cfg.autoOffsetResetPolicy = p
	}
}
