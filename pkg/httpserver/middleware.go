package httpserver

import (
	"net/http"
	"net/http/httputil"

	"github.com/nexmoinc/gosrvlib/pkg/httpserver/route"
	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"github.com/nexmoinc/gosrvlib/pkg/traceid"
	"github.com/nexmoinc/gosrvlib/pkg/uidc"
	"go.uber.org/zap"
)

func loggerMiddleware(rootLogger *zap.Logger, traceIDHeaderName string, redactFn RedactFn) MiddlewareFn {
	return func(next http.Handler) http.Handler {
		return RequestInjectHandler(rootLogger, traceIDHeaderName, redactFn, next)
	}
}

func applyMiddleware(r Route, middleware ...MiddlewareFn) http.Handler {
	info := MiddlewareInfo {
		Method: r.Method,
		Path: r.Path,
	}

	for i := len(middleware) - 1; i >= 0; i-- {
		next = middleware[i](info, next)
	}

	return next
}

// RequestInjectHandler wraps all incoming requests and injects a logger in the request scoped context.
func RequestInjectHandler(rootLogger *zap.Logger, traceIDHeaderName string, redactFn RedactFn, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqID := traceid.FromHTTPRequestHeader(r, traceIDHeaderName, uidc.NewID128())

		l := rootLogger.With(
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
			l = l.With(zap.String("request", redactFn(string(reqDump))))
		}

		ctx = logging.WithLogger(ctx, l)
		ctx = traceid.NewContext(ctx, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
