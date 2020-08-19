package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/stretchr/testify/require"
)

// nolint:gocognit
func TestCheckHttpStatus(t *testing.T) {
	tests := []struct {
		name              string
		handlerMethod     string
		handlerDelay      time.Duration
		handlerStatusCode int
		checkMethod       HTTPCheckMethod
		checkTimeout      time.Duration
		checkWantStatus   int
		wantErr           bool
	}{
		{
			name:        "fails with unsupported healthcheck method",
			checkMethod: HTTPCheckMethod("INVALID"),
			wantErr:     true,
		},
		{
			name:              "fails with wrong check method HEAD",
			checkMethod:       MethodHead,
			handlerMethod:     http.MethodGet,
			handlerStatusCode: http.StatusOK,
			wantErr:           true,
		},
		{
			name:              "fails with wrong check method GET",
			checkMethod:       MethodGet,
			handlerMethod:     http.MethodHead,
			handlerStatusCode: http.StatusOK,
			wantErr:           true,
		},
		{
			name:              "fails with handler timeout",
			checkMethod:       MethodGet,
			checkTimeout:      1 * time.Second,
			handlerMethod:     http.MethodGet,
			handlerStatusCode: http.StatusOK,
			handlerDelay:      2 * time.Second,
			wantErr:           true,
		},
		{
			name:              "succeed HEAD with 200 response",
			checkMethod:       MethodHead,
			checkTimeout:      1 * time.Second,
			checkWantStatus:   http.StatusOK,
			handlerMethod:     http.MethodHead,
			handlerStatusCode: http.StatusOK,
			wantErr:           false,
		},
		{
			name:              "succeed GET with 200 response",
			checkMethod:       MethodGet,
			checkTimeout:      1 * time.Second,
			checkWantStatus:   http.StatusOK,
			handlerMethod:     http.MethodGet,
			handlerStatusCode: http.StatusOK,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if tt.handlerMethod != r.Method {
					httputil.SendStatus(r.Context(), w, http.StatusMethodNotAllowed)
					return
				}
				if tt.handlerMethod == r.Method {
					if tt.handlerDelay != 0 {
						time.Sleep(tt.handlerDelay)
					}
					httputil.SendStatus(r.Context(), w, tt.handlerStatusCode)
				}
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			err := CheckHTTPStatus(tt.checkMethod, ts.URL, tt.checkWantStatus, tt.checkTimeout)
			t.Logf("check error: %v", err)
			if tt.wantErr {
				require.Error(t, err, "CheckHTTPStatus() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.Nil(t, err, "CheckHTTPStatus() unexpected error = %v", err)
			}
		})
	}
}
