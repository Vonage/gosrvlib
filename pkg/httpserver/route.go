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

	// Middleware is a set of middleware to apply to this route.
	Middleware []MiddlewareFn `json:"-"`
}

// Index contains the list of routes attached to the current service.
type Index struct {
	// Routes is the list of routes attached to the current service.
	Routes []Route `json:"routes"`
}
