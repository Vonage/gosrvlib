package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricInFlightRequest is a default metric label.
	MetricInFlightRequest = "in_flight_requests"

	// MetricAPIRequests is a default metric label.
	MetricAPIRequests = "api_requests_total"

	// MetricRequestDuration is a default metric label.
	MetricRequestDuration = "request_duration_seconds"

	// MetricResponseSize is a default metric label.
	MetricResponseSize = "response_size_bytes"

	// MetricErrorLevel is a default metric label.
	MetricErrorLevel = "error_level_total"

	// MetricGoRuntime is a default metric label.
	MetricGoRuntime = "go_runtime"

	// MetricGoProcess is a default metric label.
	MetricGoProcess = "go_process"
)

var (
	// DefaultCollectors contains the list of default collectors
	DefaultCollectors = []Option{
		WithCollector(
			MetricGoRuntime,
			prometheus.NewGoCollector(),
		),
		WithCollector(
			MetricGoProcess,
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		),
		WithCollectorGauge(
			MetricInFlightRequest,
			prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: MetricInFlightRequest,
					Help: "A gauge of requests being served by the wrapped handler.",
				},
			),
		),
		WithCollectorCounterVec(
			MetricAPIRequests,
			prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: MetricAPIRequests,
					Help: "A counter for requests to the wrapped handler.",
				},
				[]string{"code", "method"},
			),
		),
		WithCollectorHistogramVec(
			MetricRequestDuration,
			prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    MetricRequestDuration,
					Help:    "A histogram of latencies for requests.",
					Buckets: prometheus.ExponentialBuckets(0.001, 10, 6),
				},
				[]string{"handler", "method"},
			),
		),
		WithCollectorHistogramVec(
			MetricResponseSize,
			prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    MetricResponseSize,
					Help:    "A histogram of response sizes for requests.",
					Buckets: prometheus.ExponentialBuckets(100, 2, 6),
				},
				[]string{},
			),
		),
		WithCollectorCounterVec(
			MetricErrorLevel,
			prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: MetricErrorLevel,
					Help: "Number of errors by levels.",
				},
				[]string{"level"},
			),
		),
	}
)
