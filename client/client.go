// Package client provides a production-grade HTTP client for the PixiGPT API.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the main PixiGPT API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	retryMax   int
}

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpClient = c
	}
}

// WithRetryMax sets maximum retry attempts for failed requests.
func WithRetryMax(max int) Option {
	return func(client *Client) {
		client.retryMax = max
	}
}

// New creates a new PixiGPT client with production-grade defaults.
//
// Default configuration:
//   - Connection pooling: 100 max idle connections, 10 per host
//   - Timeouts: 30s client, 10s dial, 5s TLS handshake
//   - Keep-alive: enabled
//   - Retries: 3 max attempts with exponential backoff
func New(apiKey, baseURL string, opts ...Option) *Client {
	// Production-grade HTTP transport for high volume
	transport := &http.Transport{
		MaxIdleConns:        100,             // Total connection pool
		MaxIdleConnsPerHost: 10,              // Per-host connection reuse
		IdleConnTimeout:     90 * time.Second, // Keep connections alive
		DisableCompression:  false,            // Enable gzip
		DisableKeepAlives:   false,            // Enable keep-alive for connection reuse

		// Timeouts to prevent hanging
		DialContext: (&dialContext{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	c := &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second, // Overall request timeout
		},
		retryMax: 3, // Retry up to 3 times
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest executes an HTTP request with retries and proper error handling.
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, result interface{}) error {
	url := c.baseURL + path
	var lastErr error

	// Retry loop with exponential backoff
	for attempt := 0; attempt <= c.retryMax; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 100ms, 200ms, 400ms, 800ms...
			backoff := time.Duration(100*(1<<uint(attempt-1))) * time.Millisecond
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed (attempt %d/%d): %w", attempt+1, c.retryMax+1, err)
			continue // Retry on network errors
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		// Handle HTTP errors
		if resp.StatusCode >= 400 {
			var apiErr APIError
			if err := json.Unmarshal(respBody, &apiErr); err != nil {
				// Not a JSON error, return raw
				lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
			} else {
				lastErr = &apiErr
			}

			// Retry on 5xx errors (server-side issues)
			if resp.StatusCode >= 500 {
				continue
			}

			// Don't retry 4xx errors (client errors)
			return lastErr
		}

		// Success - parse response
		if result != nil {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}
		}

		return nil
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}
