package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Option is a type alias for a function that configures Kafka client.
type Option func(*config)

// WithConfigParameter extends kafka.ConfigMap with additional parameters.
// Full list of parameters: https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
func WithConfigParameter(key string, val kafka.ConfigValue) Option {
	return func(c *config) {
		_ = c.ConfigMap.SetKey(key, val) // it never returns an error
	}
}

// WithTimeout sets the timeout; must be >= 6 sec.
func WithTimeout(t time.Duration) Option {
	return WithConfigParameter("session.timeout.ms", int(t.Milliseconds()))
}

// WithAutoOffsetResetPolicy sets respective parameter.
func WithAutoOffsetResetPolicy(p Offset) Option {
	return WithConfigParameter("auto.offset.reset", string(p))
}

// WithProduceChannelSize sets respective parameter.
func WithProduceChannelSize(size int) Option {
	return WithConfigParameter("go.produce.channel.size", size)
}
