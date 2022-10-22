package httpserver

import (
	"net/http"
	"net/http/httputil"

	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"github.com/nexmoinc/gosrvlib/pkg/traceid"
	"github.com/nexmoinc/gosrvlib/pkg/uidc"
	"go.uber.org/zap"
)

// Middleware is an HTTP middleware.
type Middleware func(http.Handler) http.Handler

// ApplyMiddleware applies a middleware chain to the provided HTTP handler.
func ApplyMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}

	wrapped := h

	// apply middlewares in reverse to preserve their execution order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

// defaultLogMiddleware implements a Middleware.
type defaultLogMiddleware struct {
	Logger            *zap.Logger
	TraceIDHeaderName string
	RedactFn          RedactFn
}

// requestInjectHandler wraps all incoming requests and injects a logger in the request scoped context.
func (m *defaultLogMiddleware) requestInjectHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqID := traceid.FromHTTPRequestHeader(r, m.TraceIDHeaderName, uidc.NewID128())

		l := m.Logger.With(
			zap.String("traceid", reqID),
			zap.String("request_method", r.Method),
			zap.String("request_path", r.URL.Path),
			zap.String("request_query", r.URL.RawQuery),
			zap.String("request_uri", r.RequestURI),
			zap.String("request_user_agent", r.UserAgent()),
			zap.String("request_remote_address", r.RemoteAddr),
			zap.String("request_x_forwarded_for", r.Header.Get("X-Forwarded-For")),
		)

		if l.Check(zap.DebugLevel, "debug") != nil {
			reqDump, _ := httputil.DumpRequest(r, true)
			l = l.With(zap.String("request", m.RedactFn(string(reqDump))))
		}

		ctx = logging.WithLogger(ctx, l)
		ctx = traceid.NewContext(ctx, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// NewDefaultLogMiddleware builds a new default log middleware.
func NewDefaultLogMiddleware(l *zap.Logger, traceIDHeaderName string, redactFn RedactFn) Middleware {
	m := &defaultLogMiddleware{
		Logger:            l,
		TraceIDHeaderName: traceIDHeaderName,
		RedactFn:          redactFn,
	}

	return m.requestInjectHandler
}
