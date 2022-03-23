package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Config contains the AWS configuration options.
type Config struct {
	awsOpts []func(*config.LoadOptions) error
}

// loadConfig loads the AWS configuration with the specified options.
func loadConfig(ctx context.Context, opts ...Option) (aws.Config, error) {
	cfg := &Config{}

	for _, apply := range opts {
		apply(cfg)
	}

	return config.LoadDefaultConfig(ctx, cfg.awsOpts...) // nolint: wrapcheck
}
