package httpretrier

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "succeeds with defaults",
			wantErr: false,
		},
		{
			name: "succeeds with custom values",
			opts: []Option{
				WithRetryIfFn(func(statusCode int, err error) bool { return true }),
				WithAttempts(5),
				WithDelay(601),
				WithDelayFactor(1.3),
				WithJitter(109),
			},
			wantErr: false,
		},
		{
			name:    "fails with invalid option",
			opts:    []Option{WithJitter(0)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := New(tt.opts...)

			if tt.wantErr {
				require.Nil(t, c, "New() returned value should be nil")
				require.Error(t, err, "New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.NotNil(t, c, "New() returned value should not be nil")
			require.NoError(t, err, "New() unexpected error = %v", err)
		})
	}
}
