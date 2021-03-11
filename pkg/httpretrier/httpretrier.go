// Package httpretrier allow to retry HTTP requests based on HTTP status code boolean conditions.
package httpretrier

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	// DefaultAttempts is the default maximum number of retry attempts.
	DefaultAttempts = 3

	// DefaultDelay is the default base delay amount in milliseconds.
	DefaultDelay = 500

	// DefaultDelayFactor is the default multiplication factor to get the successive delay value.
	DefaultDelayFactor = 2

	// DefaultJitter is the maximum random Jitter between retries in milliseconds.
	DefaultJitter = 100
)

// RetryIfFn is the signature of the function used to decide when retry.
type RetryIfFn func(statusCode int, err error) bool

// HTTPDoFn is the signature of the http.Do function to be retried.
type HTTPDoFn func(req *http.Request) (*http.Response, error)

// HTTPRetrier represents an instance of the HTTP retrier.
type HTTPRetrier struct {
	delayFactor float64
	delay       int64
	jitter      int64
	attempts    uint
	retryIfFn   RetryIfFn
}

func defaultHTTPRetrier() *HTTPRetrier {
	return &HTTPRetrier{
		attempts:    DefaultAttempts,
		delay:       DefaultDelay,
		delayFactor: DefaultDelayFactor,
		jitter:      DefaultJitter,
		retryIfFn:   defaultRetryIfFn,
	}
}

// New creates a new instance.
func New(opts ...Option) (*HTTPRetrier, error) {
	c := defaultHTTPRetrier()

	for _, applyOpt := range opts {
		if err := applyOpt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Retry execute the Do function and retry in case of error.
func (c *HTTPRetrier) Retry(do HTTPDoFn, req *http.Request) (*http.Response, error) {
	delay := float64(c.delay)

	for i := c.attempts; i > 1; i-- {
		resp, err := do(req)
		if !c.retryIfFn(resp.StatusCode, err) {
			return resp, err
		}

		delay *= c.delayFactor

		time.Sleep(time.Duration(int64(delay)+rand.Int63n(c.jitter)) * time.Millisecond) // nolint:gosec
	}

	return do(req)
}

func defaultRetryIfFn(statusCode int, err error) bool {
	if err != nil {
		return false
	}

	switch statusCode {
	case http.StatusNotFound, http.StatusRequestTimeout, http.StatusConflict, http.StatusMisdirectedRequest, http.StatusTooEarly, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusInsufficientStorage:
		return true
	}

	return false
}
