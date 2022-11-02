package httpreverseproxy

import (
	"net/http"
	"net/http/httptest"
	stdhttputil "net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		serviceAddr string
		opts        []Option
		wantErr     bool
	}{
		{
			name:        "fails with invalid character in URL",
			serviceAddr: "http://invalid-url.domain.invalid\u007F",
			wantErr:     true,
		},
		{
			name:        "succeeds with defaults",
			serviceAddr: "http://service.domain.invalid:1234/",
			wantErr:     false,
		},
		{
			name:        "succeeds with custom http client",
			serviceAddr: "http://service.domain.invalid:1235/",
			opts:        []Option{WithHTTPClient(&testHTTPClient{})},
			wantErr:     false,
		},
		{
			name:        "succeeds with custom reverse proxy",
			serviceAddr: "http://service.domain.invalid:1235/",
			opts:        []Option{WithReverseProxy(&stdhttputil.ReverseProxy{})},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := New(tt.serviceAddr, tt.opts...)
			if tt.wantErr {
				require.Nil(t, c, "New() returned client should be nil")
				require.Error(t, err, "New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.NotNil(t, c, "New() returned client should not be nil")
			require.NoError(t, err, "New() unexpected error = %v", err)
		})
	}
}

func TestClient_ForwardRequest(t *testing.T) {
	t.Parallel()

	doneCh := make(chan struct{})

	var proxyServerURL string

	// setup target test server
	targetMux := http.NewServeMux()
	targetMux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			httputil.SendStatus(r.Context(), w, http.StatusOK)
			close(doneCh)
		}()

		rd, err := stdhttputil.DumpRequest(r, false)
		require.NoError(t, err)
		t.Logf("%s", string(rd))

		proxyTestURL, err := url.Parse(proxyServerURL)
		require.NoError(t, err)

		require.Equal(t, r.Host, proxyTestURL.Host)
		require.Equal(t, r.Header.Get("X-Forwarded-For"), "127.0.0.1")
	})

	targetServer := httptest.NewServer(targetMux)
	defer targetServer.Close()

	// setup proxy test server
	c, err := New(targetServer.URL)
	require.NoError(t, err)

	proxyMux := testutil.RouterWithHandler(http.MethodGet, "/v1/*path", c.ForwardRequest)

	proxyServer := httptest.NewServer(proxyMux)
	defer proxyServer.Close()

	// perform test!
	proxyServerURL = proxyServer.URL

	req, _ := http.NewRequestWithContext(testutil.Context(), http.MethodGet, proxyServerURL+"/v1/test", nil)

	hc := &http.Client{Timeout: 1 * time.Second}
	resp, err := hc.Do(req)
	require.NotNil(t, resp)

	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	}()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	<-doneCh
}
