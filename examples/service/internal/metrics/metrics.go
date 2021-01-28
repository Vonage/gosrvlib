// Package metrics defines the instrumentation metrics for this program.
package metrics

import (
	"github.com/nexmoinc/gosrvlib/pkg/metrics"
	prom "github.com/nexmoinc/gosrvlib/pkg/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// NameExample is the name of an example custom collector.
	NameExample = "example_collector"

	labelCode = "code"
)

// Metrics groups the custom collectors to be shared with other packages.
type Metrics struct {
	// CollectorExample is an example collector.
	CollectorExample *prometheus.CounterVec
}

// New creates a new Metrics instance.
func New() *Metrics {
	return &Metrics{
		CollectorExample: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: NameExample,
				Help: "Example of custom collector.",
			},
			[]string{labelCode},
		),
	}
}

// CreateMetricsClientFunc returns the metrics Client.
func (m *Metrics) CreateMetricsClientFunc() (metrics.Client, error) {
	coll := prom.DefaultCollectorOptions
	coll = append(coll, prom.WithCollector(m.CollectorExample))
	return prom.New(coll...)
}

// IncExampleCounter is an example function to increment a counter.
func (m *Metrics) IncExampleCounter(code string) {
	m.CollectorExample.With(prometheus.Labels{labelCode: code}).Inc()
}
