package jsendx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/nexmoinc/gosrvlib/pkg/internal/mocks"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	t.Parallel()

	params := &AppInfo{
		ProgramName:    "test",
		ProgramVersion: "1.2.3",
		ProgramRelease: "12345",
		ServerAddress:  "/api/test",
	}

	rr := httptest.NewRecorder()
	Send(testutil.Context(), rr, http.StatusOK, params, "hello test")

	resp := rr.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	var okResp Response
	_ = json.Unmarshal(body, &okResp)

	require.Equal(t, "test", okResp.Program, "uncexpected response: %s", body)
	require.Equal(t, "1.2.3", okResp.Version, "uncexpected response: %s", body)
	require.Equal(t, "12345", okResp.Release, "uncexpected response: %s", body)
	require.Equal(t, "/api/test", okResp.URL, "uncexpected response: %s", body)
	require.Equal(t, "OK", okResp.Message, "uncexpected response: %s", body)
	require.Equal(t, "hello test", okResp.Data, "uncexpected response: %s", body)

	// add coverage for error handling
	mockWriter := mocks.NewMockTestHTTPResponseWriter(gomock.NewController(t))
	mockWriter.EXPECT().Header().AnyTimes().Return(http.Header{})
	mockWriter.EXPECT().WriteHeader(http.StatusOK)
	mockWriter.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("io error"))
	Send(testutil.Context(), mockWriter, http.StatusOK, params, "message")
}

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		setupRouter func(*httprouter.Router)
		wantStatus  int
	}{
		{
			name:       "should handle 404",
			method:     http.MethodGet,
			path:       "/not/found",
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "should handle 405",
			method: http.MethodPost,
			setupRouter: func(r *httprouter.Router) {
				fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
				})
				r.Handler(http.MethodGet, "/not/allowed", fn)
			},
			path:       "/not/allowed",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "should handle panic in handler",
			method: http.MethodGet,
			setupRouter: func(r *httprouter.Router) {
				fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("panicking!")
				})
				r.Handler(http.MethodGet, "/panic", fn)
			},
			path:       "/panic",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			params := &AppInfo{
				ProgramName:    "test",
				ProgramVersion: "1.2.3",
				ProgramRelease: "12345",
				ServerAddress:  "/api/test",
			}
			r := NewRouter(params)

			if tt.setupRouter != nil {
				tt.setupRouter(r)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, httptest.NewRequest(tt.method, tt.path, nil))

			resp := rr.Result()
			require.Equal(t, tt.wantStatus, resp.StatusCode, "status code got = %d, want = %d", resp.StatusCode, tt.wantStatus)
		})
	}
}
