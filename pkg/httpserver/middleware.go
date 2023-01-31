package httpserver

import (
	"net/http"
	"net/http/httputil"
	"time"

	libhttputil "github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"github.com/nexmoinc/gosrvlib/pkg/traceid"
	"github.com/nexmoinc/gosrvlib/pkg/uidc"
	"go.uber.org/zap"
)

// MiddlewareArgs contains extra optional arguments to be passed to the middleware handler function MiddlewareFn.
type MiddlewareArgs struct {
	// Method is the HTTP method (e.g.: GET, POST, PUT, DELETE, ...).
	Method string

	// Path is the URL path.
	Path string

	// Description is the description of the route or a general description for the handler.
	Description string

	// TraceIDHeaderName is the Trace ID header name.
	TraceIDHeaderName string

	// RedactFunc is the function used to redact HTTP request and response dumps in the logs.
	RedactFunc RedactFn

	// Logger is the logger.
	Logger *zap.Logger
}

// MiddlewareFn is a function that wraps an http.Handler.
type MiddlewareFn func(args MiddlewareArgs, next http.Handler) http.Handler

// RequestInjectHandler wraps all incoming requests and injects a logger in the request scoped context.
func RequestInjectHandler(logger *zap.Logger, traceIDHeaderName string, redactFn RedactFn, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqTime := time.Now().UTC()
		reqID := traceid.FromHTTPRequestHeader(r, traceIDHeaderName, uidc.NewID128())

		l := logger.With(
			zap.String(traceid.DefaultLogKey, reqID),
			zap.Time("request_time", reqTime),
			zap.String("request_method", r.Method),
			zap.String("request_path", logging.Sanitize(r.URL.Path)),
			zap.String("request_query", logging.Sanitize(r.URL.RawQuery)),
			zap.String("request_remote_address", r.RemoteAddr),
			zap.String("request_uri", r.RequestURI),
			zap.String("request_user_agent", logging.Sanitize(r.UserAgent())),
			zap.String("request_x_forwarded_for", logging.Sanitize(r.Header.Get("X-Forwarded-For"))),
		)

		if l.Check(zap.DebugLevel, "debug") != nil {
			reqDump, _ := httputil.DumpRequest(r, true)
			l = l.With(zap.String("request", redactFn(string(reqDump))))
		}

		ctx := r.Context()
		ctx = libhttputil.WithRequestTime(ctx, reqTime)
		ctx = traceid.NewContext(ctx, reqID)
		ctx = logging.WithLogger(ctx, l)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// LoggerMiddlewareFn returns the middleware handler function to handle logs.
func LoggerMiddlewareFn(args MiddlewareArgs, next http.Handler) http.Handler {
	return RequestInjectHandler(args.Logger, args.TraceIDHeaderName, args.RedactFunc, next)
}

// ApplyMiddleware returns an http Handler with all middleware handler functions applied.
func ApplyMiddleware(arg MiddlewareArgs, next http.Handler, middleware ...MiddlewareFn) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		next = middleware[i](arg, next)
	}

	return next
}
