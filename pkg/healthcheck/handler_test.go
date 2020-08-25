// +build unit

package healthcheck

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name           string
		checks         []HealthCheck
		opts           []HandlerOption
		wantStatus     int
		wantBody       string
		wantMaxElapsed time.Duration
	}{
		{
			name: "success multiple OK",
			checks: []HealthCheck{
				New("test_01", &testHealthChecker{delay: 100 * time.Millisecond, err: nil}),
				New("test_02", &testHealthChecker{delay: 100 * time.Millisecond, err: nil}),
			},
			wantStatus:     http.StatusOK,
			wantBody:       `{"test_01":"OK","test_02":"OK"}`,
			wantMaxElapsed: 200 * time.Millisecond,
		},
		{
			name: "success multiple OK with custom response writer",
			checks: []HealthCheck{
				New("test_11", &testHealthChecker{delay: 100 * time.Millisecond, err: nil}),
				New("test_12", &testHealthChecker{delay: 100 * time.Millisecond, err: nil}),
			},
			opts: []HandlerOption{
				WithResultWriter(func(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
					type wrapper struct {
						Data interface{} `json:"data"`
					}
					httputil.SendJSON(ctx, w, statusCode, &wrapper{
						Data: data,
					})
				}),
			},
			wantStatus:     http.StatusOK,
			wantBody:       `{"data":{"test_11":"OK","test_12":"OK"}}`,
			wantMaxElapsed: 200 * time.Millisecond,
		},
		{
			name: "success mixed results",
			checks: []HealthCheck{
				New("test_31", &testHealthChecker{delay: 100 * time.Millisecond, err: nil}),
				New("test_32", &testHealthChecker{delay: 200 * time.Millisecond, err: fmt.Errorf("check error")}),
			},
			wantStatus:     http.StatusServiceUnavailable,
			wantBody:       `{"test_31":"OK","test_32":"check error"}`,
			wantMaxElapsed: 300 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(testutil.Context(), http.MethodGet, "/", nil)
			require.NoError(t, err, "no error expected reading body data")

			h := NewHandler(tt.checks, tt.opts...)

			st := time.Now()
			h.ServeHTTP(rr, req)
			el := time.Since(st)

			resp := rr.Result()
			payloadData, _ := ioutil.ReadAll(resp.Body)
			payload := string(payloadData)

			require.Equal(t, tt.wantStatus, resp.StatusCode)
			require.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
			require.Equal(t, tt.wantBody+"\n", payload)

			// ensure we are running concurrently
			require.True(t, el < tt.wantMaxElapsed, "check time = %s, want < %s", el, tt.wantMaxElapsed)
		})
	}
}

func Test_runCheckWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		checker HealthChecker
		timeout time.Duration
		wantErr error
	}{
		{
			name:    "check fails with timeout",
			checker: &testHealthChecker{delay: 100 * time.Millisecond},
			timeout: 10 * time.Millisecond,
			wantErr: context.DeadlineExceeded,
		},
		{
			name:    "check succeed with OK result",
			checker: &testHealthChecker{delay: 100 * time.Millisecond, err: nil},
			timeout: 500 * time.Millisecond,
			wantErr: nil,
		},
		{
			name:    "check succeed with ERR result",
			checker: &testHealthChecker{delay: 100 * time.Millisecond, err: fmt.Errorf("check failed")},
			timeout: 500 * time.Millisecond,
			wantErr: fmt.Errorf("check failed"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := testutil.Context()
			got := runCheckWithTimeout(ctx, tt.checker, tt.timeout)
			if !reflect.DeepEqual(got, tt.wantErr) {
				t.Errorf("runCheck() = %#v, want = %#v", got, tt.wantErr)
			}
		})
	}
}
