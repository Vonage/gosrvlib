package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	// DefaultWaitTimeSeconds is the default duration (in seconds) for which the call waits for a message to arrive in the queue before returning.
	DefaultWaitTimeSeconds = 60

	// DefaultVisibilityTimeout is the default duration (in seconds) that the received messages are hidden from subsequent retrieve requests after being retrieved by a ReceiveMessage request.
	DefaultVisibilityTimeout = 300
)

type cfg struct {
	awsOpts           []func(*config.LoadOptions) error
	awsConfig         aws.Config
	waitTimeSeconds   int32
	visibilityTimeout int32
}

func loadConfig(ctx context.Context, opts ...Option) (*cfg, error) {
	c := &cfg{
		waitTimeSeconds:   DefaultWaitTimeSeconds,
		visibilityTimeout: DefaultVisibilityTimeout,
	}

	for _, apply := range opts {
		apply(c)
	}

	if c.waitTimeSeconds < 0 {
		return nil, fmt.Errorf("waitTimeSeconds must be greater or equal zero")
	}

	if c.visibilityTimeout < 0 || c.visibilityTimeout > 43200 {
		return nil, fmt.Errorf("visibilityTimeout must be between 0 and 43200")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, c.awsOpts...)

	if err == nil {
		c.awsConfig = awsConfig
	}

	return c, err // nolint: wrapcheck
}
