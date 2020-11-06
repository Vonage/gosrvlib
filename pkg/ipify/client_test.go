package ipify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		opts        []ClientOption
		wantTimeout time.Duration
		wantAPIURL  string
		wantErrorIP string
		wantErr     bool
	}{
		{
			name:        "succeeds with defaults",
			wantTimeout: defaultTimeout,
			wantAPIURL:  defaultAPIURL,
			wantErrorIP: defaultErrorIP,
			wantErr:     false,
		},
		{
			name: "succeeds with custom values",
			opts: []ClientOption{
				WithTimeout(2 * time.Second),
				WithURL("http://test.ipify.invalid"),
				WithErrorIP("0.0.0.0"),
			},
			wantTimeout: 2 * time.Second,
			wantAPIURL:  "http://test.ipify.invalid",
			wantErrorIP: "0.0.0.0",
			wantErr:     false,
		},
		{
			name:    "fails with invalid character in URL",
			opts:    []ClientOption{WithURL("http://invalid-url.domain.invalid\u007F")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := NewClient(tt.opts...)
			if tt.wantErr {
				require.Nil(t, c, "NewClient() returned client should be nil")
				require.Error(t, err, "NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotNil(t, c, "NewClient() returned client should not be nil")
			require.NoError(t, err, "NewClient() unexpected error = %v", err)
			require.Equal(t, tt.wantTimeout, c.timeout, "NewClient() unexpected timeout = %d got %d", tt.wantTimeout, c.timeout)
			require.Equal(t, tt.wantAPIURL, c.apiURL, "NewClient() unexpected apiURL = %d got %d", tt.wantAPIURL, c.apiURL)
			require.Equal(t, tt.wantErrorIP, c.errorIP, "NewClient() unexpected errorIP = %d got %d", tt.wantErrorIP, c.errorIP)
		})
	}
}
