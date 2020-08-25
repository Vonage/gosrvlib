package healthcheck

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/nexmoinc/gosrvlib/pkg/httputil/jsendx"
)

const (
	checkTimeout = 1 * time.Second
)

// HealthCheckerMap is a type alias for a map of healthchecker instances
type HealthCheckerMap map[string]HealthChecker

// Handler returns an HTTP handler function performing the healthcheck
// This is a basic fanout implementation, it could be smarter and run a background collection process
// independent from how many times we call the status endpoint
func Handler(checks HealthCheckerMap, appInfo *jsendx.AppInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkCount := len(checks)

		var wg sync.WaitGroup
		wg.Add(checkCount)

		type resultWrap struct {
			id  string
			err error
		}
		resCh := make(chan resultWrap, checkCount)
		defer close(resCh)

		for id, hc := range checks {
			go func(id string, hc HealthChecker) {
				defer wg.Done()
				resCh <- resultWrap{
					id:  id,
					err: runCheckWithTimeout(r.Context(), hc, checkTimeout),
				}
			}(id, hc)
		}

		wg.Wait()

		status := http.StatusOK
		data := make(map[string]string, checkCount)
		for i := 0; i < checkCount; i++ {
			r := <-resCh
			if r.err != nil {
				status = http.StatusServiceUnavailable
				data[r.id] = r.err.Error()
				continue
			}
			data[r.id] = "OK"
		}

		if appInfo != nil {
			jsendx.Send(r.Context(), w, status, appInfo, data)
			return
		}
		httputil.SendJSON(r.Context(), w, status, data)
	}
}

func runCheckWithTimeout(ctx context.Context, c HealthChecker, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resCh := make(chan error)
	go func() {
		resCh <- c.HealthCheck(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case r := <-resCh:
		return r
	}
}
