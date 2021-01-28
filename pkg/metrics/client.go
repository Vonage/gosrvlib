package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Client represents the state type of this client.
type Client struct {
	Registry              *prometheus.Registry
	handlerOpts           promhttp.HandlerOpts
	Collector             map[string]prometheus.Collector
	CollectorGauge        map[string]prometheus.Gauge
	CollectorCounter      map[string]prometheus.Counter
	CollectorSummary      map[string]prometheus.Summary
	CollectorHistogram    map[string]prometheus.Histogram
	CollectorGaugeVec     map[string]*prometheus.GaugeVec
	CollectorCounterVec   map[string]*prometheus.CounterVec
	CollectorSummaryVec   map[string]*prometheus.SummaryVec
	CollectorHistogramVec map[string]*prometheus.HistogramVec
}

// New creates a new metrics instance.
func New(opts ...Option) (*Client, error) {
	c := initClient()
	err := c.Configure(opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Configure allows to specify more options.
func (c *Client) Configure(opts ...Option) error {
	for _, applyOpt := range opts {
		if err := applyOpt(c); err != nil {
			return err
		}
	}
	return nil
}

func initClient() *Client {
	return &Client{
		Registry:              prometheus.NewRegistry(),
		handlerOpts:           promhttp.HandlerOpts{},
		Collector:             make(map[string]prometheus.Collector),
		CollectorGauge:        make(map[string]prometheus.Gauge),
		CollectorCounter:      make(map[string]prometheus.Counter),
		CollectorSummary:      make(map[string]prometheus.Summary),
		CollectorHistogram:    make(map[string]prometheus.Histogram),
		CollectorGaugeVec:     make(map[string]*prometheus.GaugeVec),
		CollectorCounterVec:   make(map[string]*prometheus.CounterVec),
		CollectorSummaryVec:   make(map[string]*prometheus.SummaryVec),
		CollectorHistogramVec: make(map[string]*prometheus.HistogramVec),
	}
}

// InstrumentHandler wraps an http.Handler to collect Prometheus metrics.
func (c *Client) InstrumentHandler(path string, handler http.HandlerFunc) http.Handler {
	var h http.Handler
	h = handler
	collectorRequestSize, ok := c.CollectorHistogramVec[RequestSize]
	if ok {
		h = promhttp.InstrumentHandlerRequestSize(collectorRequestSize, h)
	}
	collectorResponseSize, ok := c.CollectorHistogramVec[ResponseSize]
	if ok {
		h = promhttp.InstrumentHandlerResponseSize(collectorResponseSize, h)
	}
	collectorAPIRequests, ok := c.CollectorCounterVec[APIRequests]
	if ok {
		h = promhttp.InstrumentHandlerCounter(collectorAPIRequests, h)
	}
	collectorRequestDuration, ok := c.CollectorHistogramVec[RequestDuration]
	if ok {
		h = promhttp.InstrumentHandlerDuration(collectorRequestDuration.MustCurryWith(prometheus.Labels{labelHandler: path}), h)
	}
	collectorInFlightRequests, ok := c.CollectorGauge[InFlightRequests]
	if ok {
		h = promhttp.InstrumentHandlerInFlight(collectorInFlightRequests, h)
	}
	return h
}

// MetricsHandlerFunc returns an http handler function to serve the metrics endpoint.
func (c *Client) MetricsHandlerFunc() http.HandlerFunc {
	h := promhttp.HandlerFor(c.Registry, c.handlerOpts)
	return promhttp.InstrumentMetricHandler(c.Registry, h).ServeHTTP
}

// DefaultMetricsHandlerFunc returns a default http handler function to serve the metrics endpoint.
func DefaultMetricsHandlerFunc() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}

// IncLogLevelCounter counts the number of errors for each log severity level.
func (c *Client) IncLogLevelCounter(level string) {
	m, ok := c.CollectorCounterVec[ErrorLevel]
	if ok {
		m.With(prometheus.Labels{labelLevel: level}).Inc()
	}
}

// InstrumentRoundTripper is a middleware that wraps the provided http.RoundTripper to observe the request result with default metrics.
func (c *Client) InstrumentRoundTripper(next http.RoundTripper) http.RoundTripper {
	collectorOutboundRequests, ok := c.CollectorCounterVec[OutboundRequests]
	if ok {
		next = promhttp.InstrumentRoundTripperCounter(collectorOutboundRequests, next)
	}
	collectorOutboundRequestsDuration, ok := c.CollectorHistogramVec[OutboundRequestsDuration]
	if ok {
		next = promhttp.InstrumentRoundTripperDuration(collectorOutboundRequestsDuration, next)
	}
	collectorOutboundInFlightRequests, ok := c.CollectorGauge[OutboundInFlightRequests]
	if ok {
		next = promhttp.InstrumentRoundTripperInFlight(collectorOutboundInFlightRequests, next)
	}
	return next
}
