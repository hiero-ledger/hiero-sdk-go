package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

const httpOpCall = "call"

// mirrorNodeHttpClient sends the requests of one call to one mirror node, retrying per policy.
// It does not own transport.
type mirrorNodeHttpClient struct {
	baseURL   string
	transport HttpTransport
	policy    MirrorNodeHttpRetryPolicy
	deadline  time.Time
	headers   map[string]string
	// clientClosed is closed when the owning Client closes; nil if there is none.
	clientClosed <-chan struct{}
}

func newMirrorNodeHttpClient(baseURL string, transport HttpTransport, policy MirrorNodeHttpRetryPolicy) *mirrorNodeHttpClient {
	client := &mirrorNodeHttpClient{baseURL: baseURL, transport: transport, policy: policy}
	if total := policy.GetTotalDeadline(); total > 0 {
		client.deadline = time.Now().Add(total)
	}

	return client
}

// get sends a GET request for path.
func (c *mirrorNodeHttpClient) get(path mirrorNodeRestPath, cancellation Cancellation) (HttpResponse, error) {
	return c.do(HttpRequest{
		method:  HttpMethodGet,
		url:     resolveMirrorPath(c.baseURL, path),
		headers: c.headers,
	}, cancellation)
}

// post sends a POST request for path. Mirror node POST endpoints are read-only, so POST is retried like GET.
func (c *mirrorNodeHttpClient) post(path mirrorNodeRestPath, contentType string, body []byte, cancellation Cancellation) (HttpResponse, error) {
	return c.do(HttpRequest{
		method:      HttpMethodPost,
		url:         resolveMirrorPath(c.baseURL, path),
		body:        body,
		contentType: contentType,
		headers:     c.headers,
	}, cancellation)
}

// do sends req, retrying per policy. When retries run out on a retryable status, the last
// response is returned with the error.
func (c *mirrorNodeHttpClient) do(req HttpRequest, cancellation Cancellation) (HttpResponse, error) {
	var lastResp HttpResponse
	var lastErr error

	maxAttempts := int(c.policy.GetMaxAttempts())
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := c.stopReason(attempt, lastResp, cancellation); err != nil {
			return lastResp, err
		}

		req.deadline = c.attemptDeadline()
		resp, err := c.transport.RoundTrip(req, cancellation)
		switch {
		case err != nil:
			// Client close, cancellation and the total deadline take precedence over the transport error.
			if stopErr := c.stopReason(attempt+1, HttpResponse{}, cancellation); stopErr != nil {
				return HttpResponse{}, stopErr
			}
			if kind, _ := httpErrorKindOf(err); kind != httpTransient {
				return HttpResponse{}, err
			}
			lastResp, lastErr = HttpResponse{}, err
		case !c.shouldRetryStatus(resp.statusCode):
			return resp, nil
		default:
			lastResp, lastErr = resp, nil
		}

		if attempt == maxAttempts-1 {
			break
		}
		if err := c.waitBeforeRetry(attempt, lastResp, cancellation); err != nil {
			return lastResp, err
		}
	}

	return lastResp, c.exhaustedError(lastResp, lastErr)
}

// attemptDeadline returns the smaller of the per-attempt timeout and the time left in the call, or 0 if neither is set.
func (c *mirrorNodeHttpClient) attemptDeadline() time.Duration {
	bound := c.policy.GetPerAttemptTimeout()
	if !c.deadline.IsZero() {
		remaining := max(time.Until(c.deadline), time.Nanosecond)
		if bound <= 0 || remaining < bound {
			bound = remaining
		}
	}

	return bound
}

// stopReason returns an error if the Client closed, the caller cancelled, or the total deadline passed.
func (c *mirrorNodeHttpClient) stopReason(attempts int, last HttpResponse, cancellation Cancellation) error {
	select {
	case <-c.clientClosed:
		return clientClosedError()
	default:
	}
	if cancellation.IsCancelled() {
		return cancelledError(cancellation)
	}
	if !c.deadline.IsZero() && !time.Now().Before(c.deadline) {
		return deadlineExceededError(attempts, last, context.DeadlineExceeded)
	}

	return nil
}

func (c *mirrorNodeHttpClient) shouldRetryStatus(status int) bool {
	if status >= http.StatusOK && status < http.StatusMultipleChoices {
		return false
	}

	return slices.Contains(c.policy.retryableStatuses(), uint16(status))
}

// waitBeforeRetry sleeps before the next attempt, using Retry-After on a 429 or 5xx. It fails at
// once if the wait would pass the total deadline.
func (c *mirrorNodeHttpClient) waitBeforeRetry(attempt int, resp HttpResponse, cancellation Cancellation) error {
	delay := c.backoffDelay(attempt)
	if resp.statusCode == http.StatusTooManyRequests || resp.statusCode >= http.StatusInternalServerError {
		if after, ok := retryAfterDelay(resp.headers); ok {
			delay = after
		}
	}
	if !c.deadline.IsZero() && delay >= time.Until(c.deadline) {
		return deadlineExceededError(attempt+1, resp, fmt.Errorf("a wait of %s is longer than the time left", delay))
	}

	select {
	case <-cancellation.GetContext().Done():
		return cancelledError(cancellation)
	case <-c.clientClosed:
		return clientClosedError()
	case <-time.After(delay):
		return nil
	}
}

func clientClosedError() error {
	return &HttpTransportError{op: httpOpCall, ID: HttpTransportClientClosedError, Err: errMirrorHttpClientClosed}
}

func cancelledError(cancellation Cancellation) error {
	return &HttpTransportError{op: httpOpCall, ID: HttpTransportCancelledError, Err: cancellation.GetContext().Err()}
}

func deadlineExceededError(attempts int, last HttpResponse, cause error) error {
	if last.statusCode != 0 {
		return fmt.Errorf("%w after %d attempt(s), last status %d: %w", errMirrorHttpDeadlineExceeded, attempts, last.statusCode, cause)
	}

	return fmt.Errorf("%w after %d attempt(s): %w", errMirrorHttpDeadlineExceeded, attempts, cause)
}

// backoffDelay returns a random delay in [0, min(maxBackoff, initialBackoff*2^attempt)).
func (c *mirrorNodeHttpClient) backoffDelay(attempt int) time.Duration {
	ceiling := min(c.policy.GetInitialBackoff()<<min(attempt, 20), c.policy.GetMaxBackoff())

	return mirrorHttpJitter(ceiling)
}

func mirrorHttpJitter(ceiling time.Duration) time.Duration {
	if ceiling <= 0 {
		return 0
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(ceiling)))
	if err != nil {
		// Waiting the whole ceiling is the safe direction to fail.
		return ceiling
	}

	return time.Duration(n.Int64())
}

func (c *mirrorNodeHttpClient) exhaustedError(resp HttpResponse, lastErr error) error {
	if lastErr != nil {
		return fmt.Errorf("mirror node HTTP request failed after %d attempt(s): %w", c.policy.GetMaxAttempts(), lastErr)
	}

	return fmt.Errorf("%w after %d attempt(s), last status %d", errMirrorHttpRetriesExhausted, c.policy.GetMaxAttempts(), resp.statusCode)
}

// retryAfterDelay reads a Retry-After header in either permitted form: delta-seconds, or an
// HTTP date. A date already in the past means "now".
func retryAfterDelay(headers map[string][]string) (time.Duration, bool) {
	values := headers[httpHeaderRetryAfter]
	if len(values) == 0 {
		return 0, false
	}

	raw := strings.TrimSpace(values[0])
	if raw == "" {
		return 0, false
	}

	if seconds, err := strconv.Atoi(raw); err == nil {
		if seconds < 0 {
			return 0, false
		}
		return time.Duration(seconds) * time.Second, true
	}

	if when, err := http.ParseTime(raw); err == nil {
		return max(time.Until(when), 0), true
	}

	return 0, false
}

// mirrorNodeStatusError returns the error for a non-200 mirror node response.
func mirrorNodeStatusError(resp HttpResponse) error {
	return fmt.Errorf("received non-200 response from mirror node: %d, details: %s", resp.statusCode, resp.body)
}
