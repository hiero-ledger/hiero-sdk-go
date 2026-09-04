//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testMirrorAccountPath is the stand-in endpoint these tests exercise.
const testMirrorAccountPath = "/accounts/0.0.1"

// fakeHttpTurn is one scripted transport outcome. The last turn repeats once the script runs
// out, so a test can say "always 503" with a single entry.
type fakeHttpTurn struct {
	resp  httpResponse
	err   error
	delay time.Duration
}

type fakeHttpTransport struct {
	mu         sync.Mutex
	turns      []fakeHttpTurn
	calls      int
	requests   []httpRequest
	closeCalls int
}

func (f *fakeHttpTransport) roundTrip(ctx context.Context, req httpRequest) (httpResponse, error) {
	f.mu.Lock()
	turn := f.turns[min(f.calls, len(f.turns)-1)]
	f.calls++
	f.requests = append(f.requests, req)
	f.mu.Unlock()

	if turn.delay > 0 {
		select {
		case <-ctx.Done():
			return httpResponse{}, &httpError{op: "send", kind: httpTransient, err: ctx.Err()}
		case <-time.After(turn.delay):
		}
	}

	return turn.resp, turn.err
}

func (f *fakeHttpTransport) close(_ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closeCalls++

	return nil
}

func (f *fakeHttpTransport) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.calls
}

func (f *fakeHttpTransport) lastRequest() httpRequest {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.requests[len(f.requests)-1]
}

// fastRetryOptions removes the real waits so retry behaviour can be asserted without sleeping.
func fastRetryOptions(maxAttempts int) mirrorHttpRetryOptions {
	options := defaultMirrorHttpRetryOptions()
	options.maxAttempts = maxAttempts
	options.minBackoff = time.Millisecond
	options.maxBackoff = 2 * time.Millisecond

	return options
}

func statusResponse(status int) httpResponse {
	return httpResponse{statusCode: status, headers: http.Header{}}
}

func testMirrorPath(t *testing.T, raw string) mirrorRestPath {
	t.Helper()

	path, err := newMirrorRestPath(raw)
	require.NoError(t, err)

	return path
}

func TestUnitMirrorHttpClientRetriesTransientStatusThenSucceeds(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusServiceUnavailable)},
		{resp: httpResponse{statusCode: http.StatusOK, body: []byte(`{"balance":1}`), headers: http.Header{}}},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(3))

	resp, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, `{"balance":1}`, string(resp.body))
	assert.Equal(t, 2, transport.callCount())
}

func TestUnitMirrorHttpClientDoesNotRetryTerminalStatus(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusBadRequest, http.StatusNotFound, http.StatusNotImplemented} {
		transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(status)}}}
		client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(3))

		resp, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
		require.NoError(t, err, "a terminal status is a result, not a transport failure")
		assert.Equal(t, status, resp.statusCode)
		assert.Equal(t, 1, transport.callCount(), "status %d must not be retried", status)
	}
}

func TestUnitMirrorHttpClientRetriesRequestTimeoutStatus(t *testing.T) {
	t.Parallel()

	// 408 is retryable here; the shipping helper treats every 4xx as terminal, which is the
	// one classification axis Go and Java currently disagree on.
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusRequestTimeout)},
		{resp: statusResponse(http.StatusOK)},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(3))

	resp, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, 2, transport.callCount())
}

func TestUnitMirrorHttpClientReturnsResponseAndErrorWhenRetriesExhausted(t *testing.T) {
	t.Parallel()

	const maxAttempts = 3
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: httpResponse{statusCode: http.StatusServiceUnavailable, body: []byte("still down"), headers: http.Header{}}},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(maxAttempts))

	resp, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.Error(t, err, "exhausting the budget is a failure, not a silent non-200")
	require.ErrorIs(t, err, errMirrorHttpRetriesExhausted)
	assert.Equal(t, http.StatusServiceUnavailable, resp.statusCode, "the last response is still available for the message")
	assert.Equal(t, "still down", string(resp.body))
	assert.Equal(t, maxAttempts, transport.callCount())
}

func TestUnitMirrorHttpClientStopsImmediatelyOnTerminalTransportError(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{err: &httpError{op: "send", kind: httpTerminal, err: errors.New("no such host")}},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(5))

	_, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.Error(t, err)
	assert.Equal(t, 1, transport.callCount(), "a terminal failure must not consume the attempt budget")
}

func TestUnitMirrorHttpClientStopsImmediatelyWhenTransportClosed(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{err: &httpError{op: "send", kind: httpClosed, err: errHttpTransportClosed}},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(5))

	_, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.ErrorIs(t, err, errHttpTransportClosed)
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorHttpClientRetriesTransientTransportError(t *testing.T) {
	t.Parallel()

	const maxAttempts = 3
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{err: &httpError{op: "send", kind: httpTransient, err: errors.New("connection reset")}},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(maxAttempts))

	_, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection reset")
	assert.Equal(t, maxAttempts, transport.callCount())
}

func TestUnitMirrorHttpClientPrefersRetryAfterOverBackoff(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Set(httpHeaderRetryAfter, "0")

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: httpResponse{statusCode: http.StatusTooManyRequests, headers: headers}},
		{resp: statusResponse(http.StatusOK)},
	}}

	options := defaultMirrorHttpRetryOptions()
	options.maxAttempts = 2
	// Backoff that would dominate the test if Retry-After were ignored.
	options.minBackoff = 5 * time.Second
	options.maxBackoff = 5 * time.Second
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, options)

	start := time.Now()
	resp, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Less(t, elapsed, time.Second, "Retry-After: 0 should override the computed backoff")
}

func TestUnitMirrorHttpClientIgnoresRetryAfterAboveCap(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Set(httpHeaderRetryAfter, "3600")

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: httpResponse{statusCode: http.StatusTooManyRequests, headers: headers}},
		{resp: statusResponse(http.StatusOK)},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(2))

	start := time.Now()
	_, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Less(t, elapsed, time.Second, "an hour-long Retry-After must not be honoured over the cap")
}

func TestUnitMirrorHttpClientBoundsTotalWallClock(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusServiceUnavailable), delay: 200 * time.Millisecond},
	}}

	options := fastRetryOptions(20)
	options.perAttemptTimeout = time.Second
	options.totalDeadline = 300 * time.Millisecond
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, options)

	start := time.Now()
	_, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Less(t, elapsed, 2*time.Second, "the total deadline bounds the loop regardless of maxAttempts")
	assert.Less(t, transport.callCount(), 20, "the deadline should cut the attempt budget short")
}

func TestUnitMirrorHttpClientAppliesPerAttemptTimeout(t *testing.T) {
	t.Parallel()

	const maxAttempts = 2
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusOK), delay: 500 * time.Millisecond},
	}}

	options := fastRetryOptions(maxAttempts)
	options.perAttemptTimeout = 50 * time.Millisecond
	options.totalDeadline = 5 * time.Second
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, options)

	_, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.Error(t, err)
	assert.Equal(t, maxAttempts, transport.callCount(), "a per-attempt timeout is transient and retried")
}

func TestUnitMirrorHttpClientHonoursCallerCancellation(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusOK), delay: 5 * time.Second},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(3))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := client.get(ctx, testMirrorPath(t, testMirrorAccountPath))

	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(start), 2*time.Second)
	assert.Equal(t, 1, transport.callCount(), "cancellation is not a reason to try again")
}

func TestUnitMirrorHttpClientRejectsZeroMaxAttempts(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(0))

	_, err := client.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maxAttempts must be at least 1")
	assert.Equal(t, 0, transport.callCount())
}

func TestUnitMirrorHttpClientRetriesPostLikeGet(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusServiceUnavailable)},
		{resp: statusResponse(http.StatusOK)},
	}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(3))

	resp, err := client.post(context.Background(), testMirrorPath(t, "/contracts/call"), "application/json", []byte(`{}`))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, 2, transport.callCount(), "mirror REST is read-only, so POST is retried too")
	assert.Equal(t, httpMethodPost, transport.lastRequest().method)
	assert.Equal(t, []byte(`{}`), transport.lastRequest().body, "the body must survive a retry")
}

func TestUnitMirrorHttpClientResolvesPathAgainstBaseURL(t *testing.T) {
	t.Parallel()

	for _, base := range []string{"https://mirror.example.com/api/v1", "https://mirror.example.com/api/v1/"} {
		transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
		client := newMirrorHttpClient(base, transport, fastRetryOptions(1))

		_, err := client.get(context.Background(), testMirrorPath(t, "/accounts/0.0.1?limit=25"))
		require.NoError(t, err)
		assert.Equal(t,
			"https://mirror.example.com/api/v1/accounts/0.0.1?limit=25",
			transport.lastRequest().url,
			"both spellings of the base URL must resolve identically")
	}
}

func TestUnitMirrorHttpClientCloseDelegatesToTransport(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", transport, fastRetryOptions(1))

	require.NoError(t, client.close(time.Second))
	assert.Equal(t, 1, transport.closeCalls)
}

func TestUnitMirrorHttpShouldRetryStatus(t *testing.T) {
	t.Parallel()

	client := newMirrorHttpClient("https://mirror.example.com/api/v1", &fakeHttpTransport{}, defaultMirrorHttpRetryOptions())

	retryable := []int{
		http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
	for _, status := range retryable {
		assert.True(t, client.shouldRetryStatus(status), "status %d should be retryable", status)
	}

	// The pinned list is the point: "5xx" would sweep these in, and each of them describes a
	// server that will answer the same way next time.
	terminal := []int{
		http.StatusOK,
		http.StatusNoContent,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusNotImplemented,
		http.StatusHTTPVersionNotSupported,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNetworkAuthenticationRequired,
	}
	for _, status := range terminal {
		assert.False(t, client.shouldRetryStatus(status), "status %d should not be retryable", status)
	}
}

func TestUnitMirrorHttpBackoffIsJitteredAndCapped(t *testing.T) {
	t.Parallel()

	options := defaultMirrorHttpRetryOptions()
	client := newMirrorHttpClient("https://mirror.example.com/api/v1", &fakeHttpTransport{}, options)

	distinct := make(map[time.Duration]struct{})
	for attempt := range 6 {
		for range 40 {
			delay := client.backoffDelay(attempt)
			require.GreaterOrEqual(t, delay, time.Duration(0))
			require.LessOrEqual(t, delay, options.maxBackoff, "backoff must never exceed the cap")
			distinct[delay] = struct{}{}
		}
	}

	assert.Greater(t, len(distinct), 1, "full jitter must not collapse to a fixed curve")
	assert.LessOrEqual(t, client.backoffDelay(0), options.minBackoff, "the first retry draws from [0, minBackoff)")
}

func TestUnitMirrorHttpRetryAfterDelay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  time.Duration
		ok    bool
	}{
		{name: "absent", value: "", ok: false},
		{name: "delta seconds", value: "2", want: 2 * time.Second, ok: true},
		{name: "zero seconds", value: "0", want: 0, ok: true},
		{name: "padded", value: "  5  ", want: 5 * time.Second, ok: true},
		{name: "negative", value: "-1", ok: false},
		{name: "garbage", value: "soon", ok: false},
		{name: "past http date", value: "Mon, 02 Jan 2006 15:04:05 GMT", want: 0, ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			headers := http.Header{}
			if tt.value != "" {
				headers.Set(httpHeaderRetryAfter, tt.value)
			}

			got, ok := retryAfterDelay(headers)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.Equal(t, tt.want, got)
			}
		})
	}

	t.Run("nil headers", func(t *testing.T) {
		t.Parallel()

		_, ok := retryAfterDelay(nil)
		assert.False(t, ok)
	})

	t.Run("future http date", func(t *testing.T) {
		t.Parallel()

		headers := http.Header{}
		headers.Set(httpHeaderRetryAfter, time.Now().Add(30*time.Second).UTC().Format(http.TimeFormat))

		got, ok := retryAfterDelay(headers)
		require.True(t, ok)
		assert.Greater(t, got, 20*time.Second)
		assert.LessOrEqual(t, got, 30*time.Second)
	})
}

// End-to-end through the real transport rather than the fake, so the two halves are known to
// fit together: policy above, dumb pipe below, one seam.
func TestUnitMirrorHttpDefaultClientRetriesEndToEnd(t *testing.T) {
	t.Parallel()

	var attempts int32
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.RequestURI()
		if atomic.AddInt32(&attempts, 1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"balances":[]}`))
	}))
	defer server.Close()

	client := newDefaultMirrorHttpClient(server.URL + mirrorHttpAPIVersionPrefix)
	defer func() { require.NoError(t, client.close(time.Second)) }()

	resp, err := client.get(context.Background(), testMirrorPath(t, "/balances?account.id=0.0.1"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, `{"balances":[]}`, string(resp.body))
	assert.Equal(t, int32(2), atomic.LoadInt32(&attempts))
	assert.Equal(t, "/api/v1/balances?account.id=0.0.1", gotPath, "the path resolves under the configured base")
}
