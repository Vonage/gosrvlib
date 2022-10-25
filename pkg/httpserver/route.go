package httpserver

import (
	"net/http"
)

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

	// Middlewares is a list of middlewares to apply to this route.
	Middlewares []Middleware `json:"-"`
}

// Index contains the list of routes attached to the current service.
type Index struct {
	// Routes is the list of routes attached to the current service.
	Routes []Route `json:"routes"`
}

// MiddlewareInfo contains extra information to be passed to the middleware.
type MiddlewareInfo {
	// Method is the HTTP method (e.g.: GET, POST, PUT, DELETE, ...).
	Method string

	// Path is the URL path.
	Path string
}

// MiddlewareFn is a function that wraps an http.Handler.
type MiddlewareFn func(info MiddlewareInfo, next http.Handler) http.Handler
