package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type config struct {
	*kafka.ConfigMap
}

func defaultConfig() *config {
	return &config{
		&kafka.ConfigMap{},
	}
}
