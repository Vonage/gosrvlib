package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Client is the public metrics client interface.
type Client interface {
	InstrumentHandler(path string, handler http.HandlerFunc) http.Handler
	InstrumentRoundTripper(next http.RoundTripper) http.RoundTripper
	MetricsHandlerFunc() http.HandlerFunc
	IncLogLevelCounter(level string)
	IncErrorCounter(task, operation, code string)
	PromRegistry() *prometheus.Registry
	HandlerOpts() promhttp.HandlerOpts
	SetHandlerOpts(promhttp.HandlerOpts)
}

// client represents the state type of this client.
type client struct {
	registry    *prometheus.Registry
	handlerOpts promhttp.HandlerOpts
}

// New creates a new metrics instance.
func New(opts ...Option) (Client, error) {
	c := &client{
		registry:    prometheus.NewRegistry(),
		handlerOpts: promhttp.HandlerOpts{},
	}
	for _, applyOpt := range opts {
		if err := applyOpt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// InstrumentHandler wraps an http.Handler to collect Prometheus metrics.
func (c *client) InstrumentHandler(path string, handler http.HandlerFunc) http.Handler {
	var h http.Handler
	h = handler
	h = promhttp.InstrumentHandlerRequestSize(collectorRequestSize, h)
	h = promhttp.InstrumentHandlerResponseSize(collectorResponseSize, h)
	h = promhttp.InstrumentHandlerCounter(collectorAPIRequests, h)
	h = promhttp.InstrumentHandlerDuration(collectorRequestDuration.MustCurryWith(prometheus.Labels{labelHandler: path}), h)
	h = promhttp.InstrumentHandlerInFlight(collectorInFlightRequests, h)
	return h
}

// InstrumentRoundTripper is a middleware that wraps the provided http.RoundTripper to observe the request result with default metrics.
func (c *client) InstrumentRoundTripper(next http.RoundTripper) http.RoundTripper {
	next = promhttp.InstrumentRoundTripperCounter(collectorOutboundRequests, next)
	next = promhttp.InstrumentRoundTripperDuration(collectorOutboundRequestsDuration, next)
	next = promhttp.InstrumentRoundTripperInFlight(collectorOutboundInFlightRequests, next)
	return next
}

// MetricsHandlerFunc returns an http handler function to serve the metrics endpoint.
func (c *client) MetricsHandlerFunc() http.HandlerFunc {
	h := promhttp.HandlerFor(c.registry, c.handlerOpts)
	return promhttp.InstrumentMetricHandler(c.registry, h).ServeHTTP
}

// IncLogLevelCounter counts the number of errors for each log severity level.
func (c *client) IncLogLevelCounter(level string) {
	collectorErrorLevel.With(prometheus.Labels{labelLevel: level}).Inc()
}

// IncErrorCounter increments the number of errors by task, operation and error code.
func (c *client) IncErrorCounter(task, operation, code string) {
	collectorErrorCode.With(prometheus.Labels{labelTask: task, labelOperation: operation, labelCode: code}).Inc()
}

// PromRegistry retrieves the internal Prometheus registry.
func (c *client) PromRegistry() *prometheus.Registry {
	return c.registry
}

// HandlerOpts exposes the internal handlerOpts property.
func (c *client) HandlerOpts() promhttp.HandlerOpts {
	return c.handlerOpts
}

// SetHandlerOpts populates the internal handlerOpts property.
func (c *client) SetHandlerOpts(opts promhttp.HandlerOpts) {
	c.handlerOpts = opts
}
