package httpreverseproxy

import (
	"net/http/httputil"
)

// Option is the interface that allows to set client options.
type Option func(c *Client)

// WithReverseProxy overrides the default HTTP client used to forward the requests.
func WithReverseProxy(p *httputil.ReverseProxy) Option {
	return func(c *Client) {
		c.proxy = p
	}
}

// WithHTTPClient overrides the default HTTP client used to forward the requests.
func WithHTTPClient(hc HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}
