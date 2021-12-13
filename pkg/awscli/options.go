package awscli

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

// Option is a type to allow setting custom client options.
type Option func(*awsConfig)

// WithEndpoint overrides the AWS endpoint for the service.
func WithEndpoint(url string, isImmutable bool) Option {
	return func(cfg *awsConfig) {
		cfg.endpointResolver = func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               url,
				HostnameImmutable: isImmutable,
			}, nil
		}
	}
}
