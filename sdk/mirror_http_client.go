package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// The policy half of the prototype: everything the transport deliberately does not do.
// See mirror_http_transport.go for the boundary and .claude/docs/go-mirror-http-prototype.md
// for how this maps onto sdk-collaboration-hub#283.

var errMirrorHttpRetriesExhausted = errors.New("retryable status persisted after every attempt")

// mirrorHttpRetryOptions is the mirror HTTP retry and timeout policy, as one value.
//
// It is deliberately separate from the gRPC knobs on Client. perAttemptTimeout is its own
// field rather than a reinterpretation of Client.requestTimeout, which is documented as the
// total budget for a whole operation — reusing that name for a per-attempt bound would
// silently change what SetRequestTimeout means for existing callers.
//
// retryableStatusCodes is data rather than a rule in prose, so "which 5xx" has one answer
// instead of one per reader.
//
// The connect timeout is deliberately absent: it is fixed when the transport is built, not per
// call, so a per-call field for it would silently do nothing. It lives on the Client alongside
// the transport instead — see client_mirror_http.go.
type mirrorHttpRetryOptions struct {
	maxAttempts          int
	perAttemptTimeout    time.Duration
	totalDeadline        time.Duration
	minBackoff           time.Duration
	maxBackoff           time.Duration
	retryableStatusCodes []int
}

func defaultMirrorHttpRetryOptions() mirrorHttpRetryOptions {
	return mirrorHttpRetryOptions{
		maxAttempts:       mirrorNodeDefaultMaxAttempts,
		perAttemptTimeout: mirrorNodeDefaultTimeout,
		// Off by default. Nothing bounded total wall clock before this layer, so switching one
		// on is a behaviour change that belongs to the cross-SDK policy decision, not to the
		// migration of a call site. The machinery and its tests are here for when it lands.
		totalDeadline: 0,
		minBackoff:    250 * time.Millisecond,
		maxBackoff:    8 * time.Second,
		// 501, 505, 507, 508 and 511 are absent on purpose: they describe a server that will
		// answer the same way next time.
		retryableStatusCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}
}

// mirrorHttpClient binds a mirror node base URL to a transport and a policy. It owns neither
// the transport's lifetime nor a Client — one transport can back several of these, which is
// what lets a single connection pool serve several mirror nodes.
type mirrorHttpClient struct {
	baseURL   string
	transport httpTransport
	options   mirrorHttpRetryOptions
}

func newMirrorHttpClient(baseURL string, transport httpTransport, options mirrorHttpRetryOptions) *mirrorHttpClient {
	return &mirrorHttpClient{baseURL: baseURL, transport: transport, options: options}
}

// newDefaultMirrorHttpClient builds a client over a transport it creates itself, for callers
// that do not want to inject one. The caller still owns closing it.
func newDefaultMirrorHttpClient(baseURL string) *mirrorHttpClient {
	return newMirrorHttpClient(baseURL, newDefaultHttpTransport(0), defaultMirrorHttpRetryOptions())
}

func (c *mirrorHttpClient) get(ctx context.Context, path mirrorRestPath) (httpResponse, error) {
	return c.do(ctx, httpRequest{
		method: httpMethodGet,
		url:    resolveMirrorPath(c.baseURL, path),
	})
}

// post carries a body. Every mirror REST endpoint in scope is read-only, so a POST is retried
// like a GET; the method is not treated as a reason to suppress retries.
func (c *mirrorHttpClient) post(ctx context.Context, path mirrorRestPath, contentType string, body []byte) (httpResponse, error) {
	return c.do(ctx, httpRequest{
		method:      httpMethodPost,
		url:         resolveMirrorPath(c.baseURL, path),
		body:        body,
		contentType: contentType,
	})
}

func (c *mirrorHttpClient) close(closeTimeout time.Duration) error {
	return c.transport.close(closeTimeout)
}

// do runs the attempt loop. A terminal transport failure ends it immediately rather than
// consuming the budget; a retryable status that survives every attempt is returned alongside
// an error, so a caller can still read the body to build its own message.
func (c *mirrorHttpClient) do(ctx context.Context, req httpRequest) (httpResponse, error) {
	if c.options.maxAttempts < 1 {
		return httpResponse{}, fmt.Errorf("maxAttempts must be at least 1, got %d", c.options.maxAttempts)
	}

	if c.options.totalDeadline > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.totalDeadline)
		defer cancel()
	}

	var lastResp httpResponse
	var lastErr error

	for attempt := 0; attempt < c.options.maxAttempts; attempt++ {
		resp, err := c.attempt(ctx, req)
		switch {
		case err != nil:
			// The caller's cancellation or the total deadline outranks the attempt budget.
			if ctxErr := ctx.Err(); ctxErr != nil {
				return httpResponse{}, fmt.Errorf("mirror node HTTP request stopped after %d attempt(s): %w", attempt+1, ctxErr)
			}
			if kind, _ := httpErrorKindOf(err); kind != httpTransient {
				return httpResponse{}, err
			}
			lastResp, lastErr = httpResponse{}, err
		case !c.shouldRetryStatus(resp.statusCode):
			return resp, nil
		default:
			lastResp, lastErr = resp, nil
		}

		if attempt == c.options.maxAttempts-1 {
			break
		}
		if err := c.waitBeforeRetry(ctx, attempt, lastResp); err != nil {
			return lastResp, err
		}
	}

	return lastResp, c.exhaustedError(lastResp, lastErr)
}

// attempt bounds one exchange. The deadline travels on the context rather than on the HTTP
// client, so a caller can cancel a single mirror read.
func (c *mirrorHttpClient) attempt(ctx context.Context, req httpRequest) (httpResponse, error) {
	if c.options.perAttemptTimeout <= 0 {
		return c.transport.roundTrip(ctx, req)
	}

	attemptCtx, cancel := context.WithTimeout(ctx, c.options.perAttemptTimeout)
	defer cancel()

	return c.transport.roundTrip(attemptCtx, req)
}

func (c *mirrorHttpClient) shouldRetryStatus(status int) bool {
	if status >= http.StatusOK && status < http.StatusMultipleChoices {
		return false
	}

	return slices.Contains(c.options.retryableStatusCodes, status)
}

// waitBeforeRetry sleeps, preferring a server-supplied Retry-After over computed backoff when
// it is no longer than the cap we would have waited anyway.
func (c *mirrorHttpClient) waitBeforeRetry(ctx context.Context, attempt int, resp httpResponse) error {
	delay := c.backoffDelay(attempt)
	if after, ok := retryAfterDelay(resp.headers); ok && after <= c.options.maxBackoff {
		delay = after
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return nil
	}
}

// backoffDelay is exponential backoff with full jitter: a uniform draw from
// [0, min(maxBackoff, minBackoff · 2^attempt)). Jitter matters because every SDK retrying a
// throttled mirror node on the same curve re-converges on it.
func (c *mirrorHttpClient) backoffDelay(attempt int) time.Duration {
	ceiling := min(c.options.minBackoff<<min(attempt, 20), c.options.maxBackoff)

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

func (c *mirrorHttpClient) exhaustedError(resp httpResponse, lastErr error) error {
	if lastErr != nil {
		return fmt.Errorf("mirror node HTTP request failed after %d attempt(s): %w", c.options.maxAttempts, lastErr)
	}

	return &httpError{
		op:   "retry",
		kind: httpTransient,
		err:  fmt.Errorf("%w after %d attempt(s), last status %d", errMirrorHttpRetriesExhausted, c.options.maxAttempts, resp.statusCode),
	}
}

// retryAfterDelay reads a Retry-After header in either permitted form: delta-seconds, or an
// HTTP date. A date already in the past means "now".
func retryAfterDelay(headers http.Header) (time.Duration, bool) {
	if headers == nil {
		return 0, false
	}

	raw := strings.TrimSpace(headers.Get(httpHeaderRetryAfter))
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

// mirrorNodeStatusError formats a non-200 the way mirrorNodeReadBody does, so a call site
// migrated onto this layer reports the same message it always has.
func mirrorNodeStatusError(resp httpResponse) error {
	return fmt.Errorf("received non-200 response from mirror node: %d, details: %s", resp.statusCode, resp.body)
}
