package httpclient

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

// Client defines the interface for HTTP requests with retry logic.
type Client interface {
	// Do executes an HTTP request with retry logic.
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// RetryableClient implements Client with exponential backoff retry.
type RetryableClient struct {
	// client is the underlying HTTP client
	client *http.Client

	// maxRetries is the maximum number of retry attempts
	maxRetries int

	// retryDelay is the base delay for exponential backoff
	retryDelay time.Duration
}

// NewRetryableClient creates a new client with retry logic.
func NewRetryableClient(timeout time.Duration, maxRetries int, retryDelay time.Duration) *RetryableClient {
	return &RetryableClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
				MaxConnsPerHost:    5,
			},
		},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// Do executes an HTTP request with retry logic.
func (c *RetryableClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 2^attempt * baseDelay
			delay := c.retryDelay * time.Duration(math.Pow(2, float64(attempt-1)))
			time.Sleep(delay)
		}

		// Clone the request body for retry attempts
		reqClone := req.Clone(ctx)

		resp, lastErr = c.client.Do(reqClone)
		if lastErr == nil && resp.StatusCode < 500 {
			// Success or client error (4xx) - don't retry
			return resp, nil
		}

		// Close response body if we got one but will retry
		if resp != nil {
			resp.Body.Close()
		}

		// Don't retry context cancellation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("all %d retry attempts failed: %w", c.maxRetries, lastErr)
}
