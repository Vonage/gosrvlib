package bootstrap

import (
	"context"
	"github.com/nexmoinc/gosrvlib/pkg/metrics"

	"go.uber.org/zap"
)

// Option is a type alias for a function that configures the application logger.
type Option func(*config)

// WithContext overrides the application context (useful for testing).
func WithContext(ctx context.Context) Option {
	return func(cfg *config) {
		cfg.context = ctx
	}
}

// WithLogger overrides the default application logger.
func WithLogger(l *zap.Logger) Option {
	return func(cfg *config) {
		cfg.createLoggerFunc = func() (*zap.Logger, error) {
			return l, nil
		}
	}
}

// WithCreateLoggerFunc overrides the root logger creation function.
func WithCreateLoggerFunc(fn CreateLoggerFunc) Option {
	return func(cfg *config) {
		cfg.createLoggerFunc = fn
	}
}

// WithCreateMetricsClientFunc overrides the default metrics client register.
func WithCreateMetricsClientFunc(fn CreateMetricsClientFunc) Option {
	return func(cfg *config) {
		cfg.createMetricsClientFunc = fn
	}
}

// WithMetricsOptions overrides the default metrics client register.
func WithMetricsOptions(opts ...metrics.Option) Option {
	return func(cfg *config) {
		cfg.createMetricsClientFunc = func() (metrics.Client, error) {
			return metrics.New(opts...)
		}
	}
}
