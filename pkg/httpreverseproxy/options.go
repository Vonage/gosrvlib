package httpreverseproxy

import (
	"log"
	"net/http/httputil"
)

// Option is the interface that allows to set client options.
type Option func(c *Client)

// WithReverseProxy overrides the default httputil.ReverseProxy.
// Leave the Director and Transport entries nil to be automatically set.
// If the Director entry is specified, then the addr argument of the New function is ignored.
// If the Transport entry is specified, then the HTTP client specified with WithHTTPClient is ignored.
func WithReverseProxy(p *httputil.ReverseProxy) Option {
	return func(c *Client) {
		c.proxy = p
	}
}

// WithHTTPClient overrides the default HTTP client used to forward the requests.
// The HTTP client can contain extra logic for logging.
func WithHTTPClient(h HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = h
	}
}

// WithLogger specifies an optional logger for errors that occur when attempting to proxy the request.
// If nil, logging is done via the log package's standard logger.
func WithLogger(l *log.Logger) Option {
	return func(c *Client) {
		c.logger = l
	}
}
