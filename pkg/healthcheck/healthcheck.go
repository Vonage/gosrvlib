package healthcheck

import (
	"context"
)

// HealthChecker is the interface that wraps the HealthCheck method.
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}
