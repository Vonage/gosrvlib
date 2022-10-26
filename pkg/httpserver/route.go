package httpserver

import (
	"net/http"

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

	// RootLogger is the logger.
	RootLogger *zap.Logger
}

// MiddlewareFn is a function that wraps an http.Handler.
type MiddlewareFn func(args MiddlewareArgs, next http.Handler) http.Handler

// Route contains the HTTP route description.
type Route struct {
	// Method is the HTTP method (e.g.: GET, POST, PUT, DELETE, ...).
	Method string `json:"method"`

	// Path is the URL path.
	Path string `json:"path"`

	// Description is the description of this route that is displayed by the /index endpoint.
	Description string `json:"description"`

	// Handler is the handler function.
	Handler http.HandlerFunc `json:"-"`

	// Middleware is a set of middleware to apply to this route.
	Middleware []MiddlewareFn `json:"-"`
}

// Index contains the list of routes attached to the current service.
type Index struct {
	// Routes is the list of routes attached to the current service.
	Routes []Route `json:"routes"`
}
