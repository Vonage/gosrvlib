// +build unit

package bootstrap_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/bootstrap"
	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// nolint:gocognit
func TestBootstrap(t *testing.T) {
	tests := []struct {
		opts             []bootstrap.Option
		name             string
		bindFunc         bootstrap.BindFunc
		createLoggerFunc bootstrap.CreateLoggerFunc
		stopAfter        time.Duration
		checkLogs        bool
		wantErr          bool
	}{
		{
			name: "should fail due to create logger function",
			createLoggerFunc: func() (*zap.Logger, error) {
				return nil, fmt.Errorf("log error")
			},
			wantErr: true,
		},
		{
			name: "should fail due to bind function",
			bindFunc: func(context.Context, *zap.Logger) error {
				return fmt.Errorf("bind error")
			},
			wantErr: true,
		},
		{
			name: "should succeed",
			bindFunc: func(context.Context, *zap.Logger) error {
				return nil
			},
			stopAfter: 500 * time.Millisecond,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			ctx, logs := testutil.ContextWithLogObserver(zap.DebugLevel)

			if tt.stopAfter != 0 {
				stopCtx, stop := context.WithTimeout(ctx, tt.stopAfter)
				defer stop()

				ctx = stopCtx
			}

			opts := []bootstrap.Option{
				bootstrap.WithContext(ctx),
			}
			opts = append(opts, tt.opts...)

			if tt.createLoggerFunc != nil {
				opts = append(opts, bootstrap.WithCreateLoggerFunc(tt.createLoggerFunc))
			} else {
				fn := func() (*zap.Logger, error) {
					return logging.FromContext(ctx), nil
				}
				opts = append(opts, bootstrap.WithCreateLoggerFunc(fn))
			}

			if err := bootstrap.Bootstrap(tt.bindFunc, opts...); (err != nil) != tt.wantErr {
				t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkLogs {
				entries := logs.All()
				require.Equal(t, "application started", entries[0].Message)
				require.Equal(t, "application stopped", entries[1].Message)
			}
		})
	}
}
