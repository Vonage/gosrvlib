// Package httpretrier allow to retry HTTP requests based on HTTP status code boolean conditions.
package httpretrier

import (
	"math/rand"
	"net/http"
	"time"
)

// Boolean operations.
const (
	OpNotEqual = iota + 1
	OpLessThan
	OpLessOrEqualThan
	OpEqual
	OpGreaterOrEqualThan
	OpGreaterThan
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

// ORGroup maps different boolean operations to a list of values to be checked against a reference value.
type ORGroup map[int][]int

// HTTPRetrier represents an instance of the HTTP retrier.
type HTTPRetrier struct {
	conditions  []ORGroup
	delayFactor float64
	delay       int64
	jitter      int64
	attempts    uint
}

func defaultHTTPRetrier() *HTTPRetrier {
	return &HTTPRetrier{
		conditions:  []ORGroup{{OpGreaterOrEqualThan: []int{500}}},
		attempts:    DefaultAttempts,
		delay:       DefaultDelay,
		delayFactor: DefaultDelayFactor,
		jitter:      DefaultJitter,
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

// HTTPDoFn is thr function signature of an http.Do function.
type HTTPDoFn func(req *http.Request) (*http.Response, error)

// Retry execute the Do function and retry in case of error.
func (c *HTTPRetrier) Retry(do HTTPDoFn, req *http.Request) (*http.Response, error) {
	delay := float64(c.delay)

	for i := c.attempts; i > 1; i-- {
		resp, err := do(req)
		if !c.check(resp.StatusCode) {
			return resp, err
		}

		delay *= c.delayFactor

		time.Sleep(time.Duration(int64(delay)+rand.Int63n(c.jitter)) * time.Millisecond) // nolint:gosec
	}

	return do(req)
}

func (c *HTTPRetrier) check(ref int) bool {
	for _, group := range c.conditions {
		if !checkORGroup(group, ref) {
			return false
		}
	}

	return true
}

func checkORGroup(group ORGroup, ref int) bool {
	for op, vals := range group {
		for _, val := range vals {
			if checkCondition(op, ref, val) {
				return true
			}
		}
	}

	return false
}

func checkCondition(op, ref, val int) bool {
	switch op {
	case OpNotEqual:
		return ref != val
	case OpLessThan:
		return ref < val
	case OpLessOrEqualThan:
		return ref <= val
	case OpEqual:
		return ref == val
	case OpGreaterOrEqualThan:
		return ref >= val
	case OpGreaterThan:
		return ref > val
	default:
		return false
	}
}
