package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const (
	// Attempts before giving up on transient transport failures.
	mirrorNodeDefaultMaxAttempts = 3
	// Per-request timeout for a mirror node call.
	mirrorNodeDefaultTimeout = 30 * time.Second

	// The only URL schemes a mirror node REST call can be made over.
	mirrorNodeSchemeHTTP  = "http"
	mirrorNodeSchemeHTTPS = "https"
)

// mirrorNodeRestBaseURL returns the client's mirror node REST API base URL, erroring when no
// mirror network is configured.
//
// The canonical value comes from _MirrorNode.getBaseRestUrl, which maps a localhost mirror network
// to port 38081. Endpoints that a local node serves elsewhere override it afterwards: the fee
// estimate and registered-node endpoints on 8084, the contract-call endpoint on 8545 (the JSON-RPC
// relay). Callers that need no override use the returned URL as-is.
func mirrorNodeRestBaseURL(client *Client) (string, error) {
	if client == nil {
		return "", errNoClientProvided
	}
	// The nil check is not redundant: GetMirrorNetwork dereferences the field.
	if client.mirrorNetwork == nil || len(client.GetMirrorNetwork()) == 0 {
		return "", errors.New("mirror node is not set")
	}

	return client.GetMirrorRestApiBaseUrl()
}

// mirrorNodeValidateURL rejects a URL no HTTP request could ever satisfy, so the retry loop does
// not spend its whole budget and backoff on a permanent failure.
func mirrorNodeValidateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid mirror node URL %q: %w", rawURL, err)
	}
	if parsed.Scheme != mirrorNodeSchemeHTTP && parsed.Scheme != mirrorNodeSchemeHTTPS {
		return fmt.Errorf("unsupported mirror node URL scheme %q in %q", parsed.Scheme, rawURL)
	}
	if parsed.Host == "" {
		return fmt.Errorf("mirror node URL %q has no host", rawURL)
	}

	return nil
}

// mirrorNodeRequestWithRetry runs send with exponential backoff, retrying transport failures and
// 5xx/429; a 4xx is the intended result of the call and is not retried. send runs once per attempt
// so it can rebuild the request body.
//
// The final attempt is returned raw so callers can format their own error: a 200 response, a
// non-200 response with a nil error, or a nil response and the transport error. The caller must
// close a non-nil Body (mirrorNodeReadBody does). A timeout of 0 disables the per-request timeout.
func mirrorNodeRequestWithRetry(client *Client, maxAttempts uint64, timeout time.Duration, send func(*http.Client) (*http.Response, error)) (*http.Response, error) {
	if maxAttempts == 0 {
		return nil, errors.New("maxAttempts must be at least 1")
	}

	httpClient := &http.Client{Timeout: timeout}

	var resp *http.Response
	var err error

	for attempt := uint64(0); attempt < maxAttempts; attempt++ {
		resp, err = send(httpClient)

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

// mirrorNodePostWithRetry POSTs; see mirrorNodeRequestWithRetry for the retry contract.
func mirrorNodePostWithRetry(client *Client, url, contentType string, body []byte, maxAttempts uint64, timeout time.Duration) (*http.Response, error) {
	if err := mirrorNodeValidateURL(url); err != nil {
		return nil, err
	}

	return mirrorNodeRequestWithRetry(client, maxAttempts, timeout, func(httpClient *http.Client) (*http.Response, error) {
		return httpClient.Post(url, contentType, bytes.NewBuffer(body)) // #nosec
	})
}

// mirrorNodeGetWithRetry GETs; see mirrorNodeRequestWithRetry for the retry contract.
func mirrorNodeGetWithRetry(client *Client, url string, maxAttempts uint64, timeout time.Duration) (*http.Response, error) {
	if err := mirrorNodeValidateURL(url); err != nil {
		return nil, err
	}

	return mirrorNodeRequestWithRetry(client, maxAttempts, timeout, func(httpClient *http.Client) (*http.Response, error) {
		return httpClient.Get(url) // #nosec
	})
}

// mirrorNodeShouldRetry reports whether a failed request should be retried: transport errors
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

// mirrorNodeReadBody closes the response body and turns a nil response or a non-200 status into an
// error carrying the response details.
func mirrorNodeReadBody(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("received nil response from mirror node")
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// Report the status even when the body could not be read; the status is the useful part.
		if readErr != nil {
			return nil, fmt.Errorf("received non-200 response from mirror node: %d, and its body could not be read: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("received non-200 response from mirror node: %d, details: %s", resp.StatusCode, body)
	}
	if readErr != nil {
		return nil, fmt.Errorf("failed to read mirror node response: %w", readErr)
	}

	return body, nil
}

// mirrorNodeWalkPages walks a paginated endpoint from pageURL, passing each fetched body to
// handlePage, which returns the raw links.next. Stops when next is absent, or at maxPages so a
// misbehaving mirror node cannot loop forever.
//
// A links.next is a root-relative path, so it is resolved against pageURL itself rather than
// against a separately supplied base -- one argument cannot then disagree with the other.
func mirrorNodeWalkPages(
	pageURL string,
	maxPages int,
	fetchPage func(pageURL string) ([]byte, error),
	handlePage func(body []byte) (*string, error),
) error {
	base, err := url.Parse(pageURL)
	if err != nil {
		return fmt.Errorf("invalid mirror node page URL %q: %w", pageURL, err)
	}

	for range maxPages {
		body, err := fetchPage(pageURL)
		if err != nil {
			return err
		}

		next, err := handlePage(body)
		if err != nil {
			return err
		}

		if next == nil || *next == "" {
			return nil
		}

		resolved, err := resolveNextURL(base, *next)
		if err != nil {
			return fmt.Errorf("invalid pagination next link %q: %w", *next, err)
		}
		pageURL = resolved
	}

	return fmt.Errorf("exceeded pagination cap of %d pages", maxPages)
}

// resolveNextURL resolves a pagination next link against the mirror base URL.
func resolveNextURL(base *url.URL, next string) (string, error) {
	parsed, err := url.Parse(next)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(parsed).String(), nil
}

// linksJSON is the pagination block on every paginated mirror node response.
type linksJSON struct {
	Next *string `json:"next"`
}
