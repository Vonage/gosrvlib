package kafka

import (
	"time"
)

type config struct {
	sessionTimeout time.Duration
}

func defaultConfig() *config {
	return &config{
		sessionTimeout: time.Second * 10,
	}
}
