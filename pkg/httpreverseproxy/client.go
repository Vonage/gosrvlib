package httpreverseproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	libhttputil "github.com/Vonage/gosrvlib/pkg/httputil"
	"github.com/Vonage/gosrvlib/pkg/logging"
	"github.com/Vonage/gosrvlib/pkg/traceid"
	"go.uber.org/zap"
)

// HTTPClient contains the function to perform the HTTP request to the proxied service.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client implements the Reverse Proxy.
type Client struct {
	proxy      *httputil.ReverseProxy
	httpClient HTTPClient
	logger     *zap.Logger
}

type errHandler = func(w http.ResponseWriter, r *http.Request, err error)

// New returns a new instance of the Client.
//
//nolint:gocognit
func New(addr string, opts ...Option) (*Client, error) {
	c := &Client{}

	for _, applyOpt := range opts {
		applyOpt(c)
	}

	if c.proxy == nil {
		c.proxy = &httputil.ReverseProxy{}
	}

	if c.proxy.Director == nil {
		addr = strings.TrimRight(addr, "/")

		proxyURL, err := url.Parse(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid service address: %s", addr)
		}

		c.proxy.Director = func(r *http.Request) {
			r.URL.Scheme = proxyURL.Scheme
			r.URL.Host = proxyURL.Host
			r.URL.Path = "/" + libhttputil.PathParam(r, "path")
			r.Host = proxyURL.Host
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		}
	}

	if c.proxy.Transport == nil {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}

		c.proxy.Transport = &httpWrapper{client: c.httpClient}
	}

	if c.logger == nil {
		c.logger, _ = logging.NewLogger(
			logging.WithFormatStr("json"),
			logging.WithLevelStr("error"),
		)
	}

	// Override the default logger to write to the zap one.
	el, err := zap.NewStdLogAt(c.logger, zap.ErrorLevel)
	if err == nil {
		c.proxy.ErrorLog = el
	}

	if c.proxy.ErrorHandler == nil {
		c.proxy.ErrorHandler = defaultErrorHandler(c.logger)
	}

	return c, nil
}

// ForwardRequest forwards a request to the proxied service.
func (c *Client) ForwardRequest(w http.ResponseWriter, r *http.Request) {
	c.proxy.ServeHTTP(w, r)
}

type httpWrapper struct {
	client HTTPClient
}

// RoundTrip implements the RoundTripper interface.
func (c *httpWrapper) RoundTrip(r *http.Request) (*http.Response, error) {
	// Request.RequestURI can't be set in client requests.
	// Ref.: https://github.com/golang/go/blob/f3c39a83a3076eb560c7f687cbb35eef9b506e7d/src/net/http/client.go#L219
	r.RequestURI = ""

	return c.client.Do(r) //nolint:wrapcheck
}

func defaultErrorHandler(logger *zap.Logger) errHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		resTime := time.Now().UTC()
		ctx := r.Context()

		reqTime, ok := libhttputil.GetRequestTimeFromContext(ctx)
		if !ok {
			reqTime = resTime
		}

		logger.With(
			zap.String(traceid.DefaultLogKey, traceid.FromContext(ctx, "")),
			zap.String("request_method", r.Method),
			zap.String("request_path", r.URL.Path),
			zap.String("request_query", r.URL.RawQuery),
			zap.String("request_uri", r.RequestURI),
			zap.Int("response_code", http.StatusBadGateway),
			zap.String("response_message", http.StatusText(http.StatusBadGateway)),
			zap.Any("response_status", libhttputil.Status(http.StatusBadGateway)),
			zap.Time("request_time", reqTime),
			zap.Time("response_time", resTime),
			zap.Duration("response_duration", resTime.Sub(reqTime)),
		).Error("proxy_error", zap.Error(err))

		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
}
