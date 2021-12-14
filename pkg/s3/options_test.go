package s3

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/require"
)

func Test_WithEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		immutable bool
		want      *awsConfig
	}{
		{
			name:      "Immutable URL",
			url:       "test",
			immutable: true,
			want: &awsConfig{awsOpts: []func(*config.LoadOptions) error{
				config.WithEndpointResolver(endpointResolver{url: "test", isImmutable: true})},
			},
		},
		{
			name:      "Mutable URL",
			url:       "test",
			immutable: false,
			want: &awsConfig{awsOpts: []func(*config.LoadOptions) error{
				config.WithEndpointResolver(endpointResolver{url: "test", isImmutable: false})},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &awsConfig{}
			gotFn := WithEndpoint(tt.url, tt.immutable)

			gotFn(cfg)

			require.Equal(t, len(tt.want.awsOpts), len(cfg.awsOpts))

			for i, opt := range tt.want.awsOpts {
				reflect.DeepEqual(opt, cfg.awsOpts[i])
			}
		})
	}
}
