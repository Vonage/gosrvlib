package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// collector names

	// GoRuntime is the name of the collector which exports metrics about the current go process.
	GoRuntime = "go_runtime"

	// GoProcess is the name of the collector which exports the current state of process metrics
	// including cpu, memory and file descriptor usage as well as the process start time for
	// the given process id under the given namespace.
	GoProcess = "go_process"

	// APIRequests is the name of the collector that counts the total inbound http requests.
	APIRequests = "api_requests_total"

	// ErrorLevel is the name of the collector that counts the number of errors for each log severity level.
	ErrorLevel = "error_level_total"

	// InFlightRequests is the name of the collector that counts in-flight inbound http requests.
	InFlightRequests = "in_flight_requests"

	// RequestDuration is the name of the collector that measures the inbound http request duration in seconds.
	RequestDuration = "request_duration_seconds"

	// RequestSize is the name of the collector that measures the http request size in bytes.
	RequestSize = "requeste_size_bytes"

	// ResponseSize is the name of the collector that measures the http response size in bytes.
	ResponseSize = "response_size_bytes"

	// OutboundRequests is the name of the collector that measures the number of outbound requests.
	OutboundRequests = "outbound_requests_total"

	// OutboundRequestsDuration is the name of the collector that measures the outbound requests duration in seconds.
	OutboundRequestsDuration = "outbound_request_duration_seconds"

	// OutboundInFlightRequests is the name of the collector that counts in-flight outbound http requests.
	OutboundInFlightRequests = "outbound_in_flight_requests"

	// labels

	labelStatusCode = "code"
	labelHandler    = "handler"
	labelLevel      = "level"
	labelMethod     = "method"
)

var (

	// DefaultSizeBuckets default prometheus buckets for size in bytes.
	DefaultSizeBuckets = prometheus.ExponentialBuckets(100, 10, 6)

	// DefaultDurationBuckets default prometheus buckets for duration in seconds.
	DefaultDurationBuckets = prometheus.ExponentialBuckets(0.001, 10, 6)

	// DefaultCollectors contains the list of default collectors
	DefaultCollectors = []Option{
		WithCollector(
			GoRuntime,
			prometheus.NewGoCollector(),
		),
		WithCollector(
			GoProcess,
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		),
		WithCollectorGauge(
			InFlightRequests,
			prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: InFlightRequests,
					Help: "Number of In-flight http requests.",
				},
			),
		),
		WithCollectorCounterVec(
			APIRequests,
			prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: APIRequests,
					Help: "Total number of http requests.",
				},
				[]string{labelStatusCode, labelMethod},
			),
		),
		WithCollectorHistogramVec(
			RequestDuration,
			prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    RequestDuration,
					Help:    "Requests duration in seconds.",
					Buckets: DefaultDurationBuckets,
				},
				[]string{labelHandler, labelMethod},
			),
		),
		WithCollectorHistogramVec(
			ResponseSize,
			prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    ResponseSize,
					Help:    "Response size in bytes.",
					Buckets: DefaultSizeBuckets,
				},
				[]string{},
			),
		),
		WithCollectorHistogramVec(
			RequestSize,
			prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    RequestSize,
					Help:    "Requests size in bytes.",
					Buckets: DefaultSizeBuckets,
				},
				[]string{},
			),
		),
		WithCollectorCounterVec(
			ErrorLevel,
			prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: ErrorLevel,
					Help: "Number of errors by severity level.",
				},
				[]string{labelLevel},
			),
		),
		WithCollectorCounterVec(
			OutboundRequests,
			prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: OutboundRequests,
					Help: "Total number of outbound http requests.",
				},
				[]string{labelStatusCode, labelMethod},
			),
		),
		WithCollectorHistogramVec(
			OutboundRequestsDuration,
			prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    OutboundRequestsDuration,
					Help:    "Outbound requests duration in seconds.",
					Buckets: DefaultDurationBuckets,
				},
				[]string{labelStatusCode, labelMethod},
			),
		),
		WithCollectorGauge(
			OutboundInFlightRequests,
			prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: OutboundInFlightRequests,
					Help: "Number of outbound In-flight http requests.",
				},
			),
		),
	}
)
