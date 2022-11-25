package cli

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/gosrvlibexampleowner/gosrvlibexample/internal/metrics"
	"github.com/nexmoinc/gosrvlib/pkg/bootstrap"
	"github.com/nexmoinc/gosrvlib/pkg/httputil/jsendx"
	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
)

//nolint:gocognit,paralleltest
func Test_bind(t *testing.T) {
	tests := []struct {
		name           string
		fcfg           func(cfg appConfig) appConfig
		preBindAddr    string
		pingAddr       string
		wantErr        bool
		wantTimeoutErr bool
	}{
		{
			name: "fails with monitor server port already bound",
			fcfg: func(cfg appConfig) appConfig {
				cfg.Enabled = false
				cfg.Servers.Monitoring.Address = ":30044"
				cfg.Servers.Public.Address = ":30045"
				return cfg
			},
			preBindAddr:    ":30044",
			wantErr:        true,
			wantTimeoutErr: false,
		},
		{
			name: "fails with public server port already bound",
			fcfg: func(cfg appConfig) appConfig {
				cfg.Enabled = false
				cfg.Servers.Monitoring.Address = ":30046"
				cfg.Servers.Public.Address = ":30047"
				return cfg
			},
			preBindAddr:    ":30047",
			wantErr:        true,
			wantTimeoutErr: false,
		},
		{
			name: "fails with same server ports",
			fcfg: func(cfg appConfig) appConfig {
				cfg.Enabled = false
				cfg.Servers.Monitoring.Address = ":30043"
				cfg.Servers.Public.Address = ":30043"
				return cfg
			},
			wantErr: true,
		},
		{
			name: "fails with bad ipify client address",
			fcfg: func(cfg appConfig) appConfig {
				cfg.Clients.Ipify.Address = "test.ipify.url.invalid\u007F"
				return cfg
			},
			wantErr:        true,
			wantTimeoutErr: false,
		},
		{
			name: "succeed with separate server ports",
			fcfg: func(cfg appConfig) appConfig {
				cfg.Enabled = false
				cfg.Servers.Monitoring.Address = ":30041"
				cfg.Servers.Public.Address = ":30042"
				return cfg
			},
			wantErr: false,
		},
		{
			name: "success with all features enabled",
			fcfg: func(cfg appConfig) appConfig {
				return cfg
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preBindAddr != "" {
				l, err := net.Listen("tcp", tt.preBindAddr)
				require.NoError(t, err)
				defer func() { _ = l.Close() }()
			}

			cfg := tt.fcfg(getValidTestConfig())

			mtr := metrics.New()

			testBindFn := bind(
				&cfg,
				&jsendx.AppInfo{
					ProgramName:    "test",
					ProgramVersion: "0.0.0",
					ProgramRelease: "0",
				},
				mtr,
			)

			testCtx, cancel := context.WithTimeout(testutil.Context(), 1*time.Second)
			defer cancel()

			testBootstrapOpts := []bootstrap.Option{
				bootstrap.WithContext(testCtx),
				bootstrap.WithLogger(logging.FromContext(testCtx)),
			}
			err := bootstrap.Bootstrap(testBindFn, testBootstrapOpts...)
			if tt.wantErr {
				require.Error(t, err, "bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantTimeoutErr {
				require.True(t, errors.Is(err, context.DeadlineExceeded),
					"bind() error = %v, wantErr %v", err, context.DeadlineExceeded)
			} else {
				require.False(t, errors.Is(err, context.DeadlineExceeded), "bind() unexpected timeout error")
			}
		})
	}
}
