package awscli

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// Option is a type to allow setting custom client options.
type Option func(*awsConfig)

// WithEndpoint overrides the AWS endpoint for the service.
func WithEndpoint(url string) Option {
	return func(cfg *awsConfig) {
		cfg.endpointResolver = func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: url,
				// HostnameImmutable: true,
			}, nil
		}
	}
}

// WithEndpointFromEnv overrides the AWS endpoint for the service reading it from the AWS_ENDPOINT environment variable.
func WithEndpointFromEnv() Option {
	return func(cfg *awsConfig) {
		cfg.endpointResolver = func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: os.Getenv("AWS_ENDPOINT"),
				// HostnameImmutable: true,
			}, nil
		}
	}
}
