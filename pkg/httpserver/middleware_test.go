package httpserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"github.com/nexmoinc/gosrvlib/pkg/redact"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/nexmoinc/gosrvlib/pkg/traceid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestApplyMiddleware(t *testing.T) {
	t.Parallel()

	// Custom middleware chain with logger.
	ctx, logs := testutil.ContextWithLogObserver(zapcore.DebugLevel)

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logging.FromContext(ctx)
			require.NotNil(t, l, "logger not found")

			l = l.With(
				zap.Int64("content_length", r.ContentLength),
			)

			ctx = logging.WithLogger(ctx, l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.FromContext(r.Context()).Info("called")
	})

	handler := ApplyMiddleware(nextHandler, middleware)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("123")))
	handler.ServeHTTP(nil, req)

	logEntries := logs.All()
	require.Len(t, logEntries, 1, "expected only 1 log message")

	logEntry := logEntries[0]
	logContextMap := logEntry.ContextMap()

	// content_length
	rqValue, rqExists := logContextMap["content_length"]
	require.True(t, rqExists, "content_length field missing")
	require.Equal(t, int64(3), rqValue)

	// message
	require.Equal(t, "called", logEntry.Message)

	// Empty middleware chain.
	ctx, logs = testutil.ContextWithLogObserver(zapcore.DebugLevel)
	handler = ApplyMiddleware(nextHandler)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(nil, req)

	logEntries = logs.All()
	require.Len(t, logEntries, 0, "expected no message")
}

func Test_requestInjectHandler(t *testing.T) {
	t.Parallel()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logging.FromContext(r.Context())
		require.NotNil(t, l, "logger not found")

		l.Info("injected")
	})

	ctx, logs := testutil.ContextWithLogObserver(zapcore.DebugLevel)
	defaultMiddleware := DefaultMiddleware(ctx, traceid.DefaultHeader, redact.HTTPData)
	handler := ApplyMiddleware(nextHandler, defaultMiddleware...)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(nil, req)

	logEntries := logs.All()
	require.Len(t, logEntries, 1, "expected only 1 log message")

	logEntry := logEntries[0]
	logContextMap := logEntry.ContextMap()

	// check request id
	idValue, idExists := logContextMap["traceid"]
	require.True(t, idExists, "traceid field missing")
	require.NotEmpty(t, idValue, "expected requestId value not found")

	// check method
	mValue, mExists := logContextMap["request_method"]
	require.True(t, mExists, "request_method field missing")
	require.Equal(t, http.MethodGet, mValue)

	// check path
	pValue, pExists := logContextMap["request_path"]
	require.True(t, pExists, "request_path field missing")
	require.Equal(t, "/", pValue)

	// request_query
	rqValue, rqExists := logContextMap["request_query"]
	require.True(t, rqExists, "request_query field missing")
	require.Equal(t, "", rqValue)

	// request_uri
	ruValue, ruExists := logContextMap["request_uri"]
	require.True(t, ruExists, "request_uri field missing")
	require.Equal(t, "/", ruValue)

	// request_user_agent
	uaValue, uaExists := logContextMap["request_user_agent"]
	require.True(t, uaExists, "request_user_agent field missing")
	require.Equal(t, "", uaValue)

	// request_user_agent
	_, ipExists := logContextMap["request_remote_address"]
	require.True(t, ipExists, "request_remote_address field missing")

	// message
	require.Equal(t, "injected", logEntry.Message)
}
