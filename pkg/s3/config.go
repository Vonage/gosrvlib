package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type awsConfig struct {
	endpointResolver aws.EndpointResolverFunc
}

// nolint: gochecknoglobals
var awsLoadDefaultConfigFn = config.LoadDefaultConfig

func loadConfig(ctx context.Context, opts ...Option) (aws.Config, error) {
	cfg := awsConfig{}

	for _, apply := range opts {
		apply(&cfg)
	}

	var awsOpts []func(*config.LoadOptions) error

	if cfg.endpointResolver != nil {
		awsOpts = append(awsOpts, config.WithEndpointResolver(cfg.endpointResolver))
	}

	return awsLoadDefaultConfigFn(ctx, awsOpts...)
}
