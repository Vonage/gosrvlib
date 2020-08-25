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

	"github.com/nexmoinc/gosrvlib/pkg/httputil/jsendx"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
)

type testHealthChecker struct {
	delay time.Duration
	err   error
}

func (th *testHealthChecker) HealthCheck(ctx context.Context) error {
	if th.delay != 0 {
		time.Sleep(th.delay)
	}
	return th.err
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		appInfo        *jsendx.AppInfo
		checkers       []HealthChecker
		wantStatus     int
		wantBody       string
		wantMaxElapsed time.Duration
	}{
		{
			name: "success multiple OK",
			checkers: []HealthChecker{
				&testHealthChecker{delay: 100 * time.Millisecond, err: nil},
				&testHealthChecker{delay: 100 * time.Millisecond, err: nil},
			},
			wantStatus:     http.StatusOK,
			wantBody:       `{"0":"OK","1":"OK"}`,
			wantMaxElapsed: 200 * time.Millisecond,
		},
		{
			name: "success multiple OK (JSendX)",
			appInfo: &jsendx.AppInfo{
				ProgramName:    "Test",
				ProgramVersion: "0.0.0",
				ProgramRelease: "test",
			},
			checkers: []HealthChecker{
				&testHealthChecker{delay: 100 * time.Millisecond, err: nil},
				&testHealthChecker{delay: 100 * time.Millisecond, err: nil},
			},
			wantStatus:     http.StatusOK,
			wantBody:       `{"program":"Test","version":"0.0.0","release":"test","url":"","datetime":"<DT>","timestamp":<TS>,"status":"success","code":200,"message":"OK","data":{"0":"OK","1":"OK"}}`,
			wantMaxElapsed: 200 * time.Millisecond,
		},
		{
			name: "success mixed results",
			checkers: []HealthChecker{
				&testHealthChecker{delay: 100 * time.Millisecond, err: nil},
				&testHealthChecker{
					delay: 200 * time.Millisecond,
					err:   fmt.Errorf("check error"),
				},
			},
			wantStatus:     http.StatusServiceUnavailable,
			wantBody:       `{"0":"OK","1":"check error"}`,
			wantMaxElapsed: 300 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			checkers := make(map[string]HealthChecker)
			for i, v := range tt.checkers {
				checkers[fmt.Sprintf("%d", i)] = v
			}

			rr := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(testutil.Context(), http.MethodGet, "/", nil)
			require.NoError(t, err, "no error expected reading body data")

			handler := Handler(checkers, tt.appInfo)

			st := time.Now()
			handler(rr, req)
			el := time.Since(st)

			resp := rr.Result()
			payloadData, _ := ioutil.ReadAll(resp.Body)
			payload := string(payloadData)

			if tt.appInfo != nil {
				payload = testutil.ReplaceDateTime(payload, "<DT>")
				payload = testutil.ReplaceUnixTimestamp(payload, "<TS>")
			}

			require.Equal(t, tt.wantStatus, resp.StatusCode)
			require.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
			require.Equal(t, tt.wantBody+"\n", payload)

			// ensure we are running concurrently
			t.Logf("check time = %s, want = <%s", el, tt.wantMaxElapsed)
			require.True(t, el < tt.wantMaxElapsed, "the check took longer than %v", tt.wantMaxElapsed)
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
