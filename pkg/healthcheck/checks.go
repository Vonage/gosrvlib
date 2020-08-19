package healthcheck

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPCheckMethod represents the list of allowed HTTP verbs to perform an healthcheck
type HTTPCheckMethod string

const (
	// MethodGet represents is equivalent to HTTP GET
	MethodGet HTTPCheckMethod = http.MethodGet

	// MethodHead represents is equivalent to HTTP HEAD
	MethodHead = http.MethodHead
)

// CheckHTTPStatus checks if the given HTTP request responds with the expected status code
func CheckHTTPStatus(method HTTPCheckMethod, url string, wantStatusCode int, timeout time.Duration) error {
	httpClient := http.Client{
		Timeout: timeout,
	}

	var resp *http.Response
	var err error
	switch method {
	case MethodHead:
		resp, err = httpClient.Head(url)
	case MethodGet:
		resp, err = httpClient.Get(url)
	default:
		return fmt.Errorf("unsupported http healthcheck method %v", method)
	}

	if err == nil {
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != wantStatusCode {
			return fmt.Errorf("unexpected http healthcheck status code: %d", resp.StatusCode)
		}
	}
	return err
}
