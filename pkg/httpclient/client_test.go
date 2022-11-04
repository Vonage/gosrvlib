package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Parallel()

	timeout := 17 * time.Second
	traceid := "test-header-123"
	component := "test-component"
	logPrefix := "prefixtest_"
	fn := func(next http.RoundTripper) http.RoundTripper { return next }
	opts := []Option{
		WithTimeout(timeout),
		WithRoundTripper(fn),
		WithTraceIDHeaderName(traceid),
		WithComponent(component),
		WithLogPrefix(logPrefix),
	}
	got := New(opts...)
	require.NotNil(t, got, "New() returned client should not be nil")
	require.Equal(t, traceid, got.traceIDHeaderName)
	require.Equal(t, component, got.component)
	require.Equal(t, timeout, got.client.Timeout)
	require.Equal(t, fn(http.DefaultTransport), got.client.Transport)
}

func TestClient_Do(t *testing.T) {
	t.Parallel()

	bodyStr := `TEST BODY OK`
	body := make([]byte, 0)

	for i := 0; i < 100; i++ {
		body = append(body, []byte(bodyStr+`\n`)...)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer server.Close()

	client := New()

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/error", nil)
	require.NoError(t, err, "failed creating http request: %v", err)

	resp, err := client.Do(req)
	require.Nil(t, resp)
	require.Error(t, err, "client.Do with invalid URL: an error was expected")

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	require.NoError(t, err, "failed creating http request: %v", err)

	resp, err = client.Do(req)
	require.NoError(t, err, "client.Do(): unexpected error = %v", err)

	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "error closing resp.Body")
	}()

	require.NotNil(t, resp, "returned response should not be nil")

	// Create a sink instance, and register it with zap for the "memory" protocol.
	sink := &MemorySink{new(bytes.Buffer)}
	err = zap.RegisterSink("memdiff", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})
	require.NoError(t, err)

	l, err := logging.NewLogger(
		logging.WithFields(
			zap.String("program", "test_log"),
			zap.String("version", "1.2.3"),
			zap.String("release", "4"),
		),
		logging.WithFormatStr("json"),
		logging.WithLevelStr("debug"),
		logging.WithOutputPaths([]string{"memdiff://"}),      // Redirect all messages to the MemorySink.
		logging.WithErrorOutputPaths([]string{"memdiff://"}), // Redirect all errors to the MemorySink.
	)
	require.NoError(t, err)
	require.NotNil(t, l)

	ctx = logging.WithLogger(ctx, l)
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)

	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	}()

	require.NotNil(t, resp)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, body, responseBody)

	// check logs
	out := sink.String()
	require.NotEmpty(t, out, "captured log output")
	require.Contains(t, out, `"request":"GET / HTTP/1.1`)
	require.Contains(t, out, `"response":"HTTP/1.1 200 OK`)
	require.Contains(t, out, bodyStr)
}

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

// Implement Close and Sync as no-ops to satisfy the interface. The Write
// method is provided by the embedded buffer.

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }
