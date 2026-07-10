package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	// Attempts before giving up on transient transport failures.
	mirrorNodeDefaultMaxAttempts = 3
	// Per-request timeout for a mirror node call.
	mirrorNodeDefaultTimeout = 30 * time.Second
)

// mirrorNodePostWithRetry POSTs to a mirror node REST/JSON-RPC endpoint, retrying transport
// failures and 5xx/429 responses with exponential backoff. Genuine 4xx responses are the
// intended result of the call and are not retried.
//
// It returns the raw result of the final attempt so callers can format their own error:
//   - success: non-nil response, StatusCode 200, Body open;
//   - non-retryable/exhausted non-200: response with Body open and a nil error;
//   - transport failure: nil response and the transport error.
//
// The caller must close a non-nil Body. A timeout of 0 disables the per-request timeout.
func mirrorNodePostWithRetry(client *Client, url, contentType string, body []byte, maxAttempts uint64, timeout time.Duration) (*http.Response, error) {
	if maxAttempts == 0 {
		return nil, errors.New("maxAttempts must be at least 1")
	}

	httpClient := &http.Client{Timeout: timeout}

	var resp *http.Response
	var err error

	for attempt := uint64(0); attempt < maxAttempts; attempt++ {
		resp, err = httpClient.Post(url, contentType, bytes.NewBuffer(body)) // #nosec

		// Success or a non-retryable outcome is terminal; return the raw result.
		if (err == nil && resp != nil && resp.StatusCode == http.StatusOK) || !mirrorNodeShouldRetry(err, resp) {
			return resp, err
		}

		// Retryable, but no point backing off after the last attempt.
		if attempt == maxAttempts-1 {
			return resp, err
		}

		// Discard the retryable response body before the next attempt.
		if resp != nil {
			resp.Body.Close()
		}

		// Exponential backoff capped at 8s; exp is clamped to avoid shift overflow.
		exp := min(attempt, uint64(5))
		delayMs := 250.0 * float64(uint64(1)<<exp)
		if delayMs > 8000 {
			delayMs = 8000
		}

		select {
		case <-client.networkUpdateContext.Done():
			return nil, client.networkUpdateContext.Err()
		case <-time.After(time.Duration(delayMs) * time.Millisecond):
		}
	}

	// Unreachable: the final iteration always returns.
	return resp, err
}

// mirrorNodeShouldRetry reports whether a failed POST should be retried: transport errors
// and 5xx/429 are transient; 4xx responses are genuine results the caller expects.
func mirrorNodeShouldRetry(err error, resp *http.Response) bool {
	if err == nil && resp != nil {
		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			return true
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return false
		}
	}

	if err == nil {
		return false
	}

	return true
}
