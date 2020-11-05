// Package ipify allows to get the public IP address
// Ref. https://www.ipify.org/
package ipify

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// APIURL is the ipify API URL
	APIURL = "https://api64.ipify.org"

	// DefaultTimeout is the default timeout value
	DefaultTimeout = 2 * time.Second

	// ErrorIP is the string to return in case of error
	ErrorIP = ""
)

// GetPublicIP retrieves the public IP of this service instance via ipify.com API
func GetPublicIP(ctx context.Context) (string, error) {
	httpClient := http.Client{
		Timeout: DefaultTimeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIURL, nil)
	if err != nil {
		return ErrorIP, fmt.Errorf("build request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ErrorIP, fmt.Errorf("failed performing ipify request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return ErrorIP, fmt.Errorf("unexpected ipify status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrorIP, fmt.Errorf("failed reading response body: %w", err)
	}

	return string(body), nil
}
