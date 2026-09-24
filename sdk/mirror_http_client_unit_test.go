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

const testMirrorAccountPath = "/accounts/0.0.1"

// fakeHttpTurn is one scripted transport outcome. The last turn repeats once the script runs
// out, so a test can say "always 503" with a single entry.
type fakeHttpTurn struct {
	resp  HttpResponse
	err   error
	delay time.Duration
}

type fakeHttpTransport struct {
	mu            sync.Mutex
	turns         []fakeHttpTurn
	calls         int
	requests      []HttpRequest
	cancellations []Cancellation
	closeCalls    int
}

// RoundTrip honours the request deadline and the cancellation the way the default transport
// does: its own deadline is a timeout, the cancellation firing is the caller's.
func (f *fakeHttpTransport) RoundTrip(req HttpRequest, cancellation Cancellation) (HttpResponse, error) {
	f.mu.Lock()
	turn := f.turns[min(f.calls, len(f.turns)-1)]
	f.calls++
	f.requests = append(f.requests, req)
	f.cancellations = append(f.cancellations, cancellation)
	f.mu.Unlock()

	if turn.delay > 0 {
		ctx := cancellation.GetContext()
		if req.deadline > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, req.deadline)
			defer cancel()
		}
		select {
		case <-ctx.Done():
			id := HttpTransportTimeoutError
			if cancellation.IsCancelled() {
				id = HttpTransportCancelledError
			}
			return HttpResponse{}, &HttpTransportError{op: httpOpSend, ID: id, Err: ctx.Err()}
		case <-time.After(turn.delay):
		}
	}

	return turn.resp, turn.err
}

func (f *fakeHttpTransport) Close(_ time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closeCalls++
}

func (f *fakeHttpTransport) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.calls
}

func (f *fakeHttpTransport) lastRequest() HttpRequest {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.requests[len(f.requests)-1]
}

// fastRetryPolicy removes the real waits so retry behaviour can be asserted without sleeping.
func fastRetryPolicy(maxAttempts uint16) MirrorNodeHttpRetryPolicy {
	return DefaultMirrorNodeHttpRetryPolicy().
		WithMaxAttempts(maxAttempts).
		WithInitialBackoff(time.Millisecond).
		WithMaxBackoff(2 * time.Millisecond)
}

func statusResponse(status int) HttpResponse {
	return HttpResponse{statusCode: status, headers: map[string][]string{}}
}

func testMirrorPath(t *testing.T, raw string) mirrorNodeRestPath {
	t.Helper()

	path, err := newMirrorNodeRestPath(raw)
	require.NoError(t, err)

	return path
}

func TestUnitMirrorHttpClientRetriesTransientStatusThenSucceeds(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusServiceUnavailable)},
		{resp: HttpResponse{statusCode: http.StatusOK, body: []byte(`{"balance":1}`), headers: map[string][]string{}}},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(3))

	resp, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, `{"balance":1}`, string(resp.body))
	assert.Equal(t, 2, transport.callCount())
}

func TestUnitMirrorHttpClientDoesNotRetryTerminalStatus(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusBadRequest, http.StatusNotFound, http.StatusNotImplemented} {
		transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(status)}}}
		client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(3))

		resp, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
		require.NoError(t, err)
		assert.Equal(t, status, resp.statusCode)
		assert.Equal(t, 1, transport.callCount(), "status %d must not be retried", status)
	}
}

func TestUnitMirrorHttpClientRetriesRequestTimeoutStatus(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusRequestTimeout)},
		{resp: statusResponse(http.StatusOK)},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(3))

	resp, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, 2, transport.callCount())
}

func TestUnitMirrorHttpClientReturnsResponseAndErrorWhenRetriesExhausted(t *testing.T) {
	t.Parallel()

	const maxAttempts = 3
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: HttpResponse{statusCode: http.StatusServiceUnavailable, body: []byte("still down"), headers: map[string][]string{}}},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(maxAttempts))

	resp, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.Error(t, err)
	require.ErrorIs(t, err, errMirrorHttpRetriesExhausted)
	assert.Equal(t, http.StatusServiceUnavailable, resp.statusCode)
	assert.Equal(t, "still down", string(resp.body))
	assert.Equal(t, maxAttempts, transport.callCount())
}

func TestUnitMirrorHttpClientStopsImmediatelyOnTerminalTransportError(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{err: &HttpTransportError{op: httpOpSend, ID: HttpTransportUnknownHostError, Err: errors.New("no such host")}},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(5))

	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.Error(t, err)
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorHttpClientStopsImmediatelyWhenTransportClosed(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{err: &HttpTransportError{op: httpOpSend, ID: HttpTransportClientClosedError, Err: errHttpTransportClosed}},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(5))

	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.ErrorIs(t, err, errHttpTransportClosed)
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorHttpClientTreatsUnrecognisedTransportErrorAsTerminal(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{err: errors.New("corporate proxy said no")},
		{resp: statusResponse(http.StatusOK)},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(5))

	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.Error(t, err)
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorHttpClientRetriesTransientTransportError(t *testing.T) {
	t.Parallel()

	const maxAttempts = 3
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{err: &HttpTransportError{op: httpOpSend, ID: HttpTransportConnectionError, Err: errors.New("connection reset")}},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(maxAttempts))

	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection reset")
	assert.Equal(t, maxAttempts, transport.callCount())
}

func TestUnitMirrorHttpClientPrefersRetryAfterOverBackoff(t *testing.T) {
	t.Parallel()

	headers := map[string][]string{httpHeaderRetryAfter: {"0"}}

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: HttpResponse{statusCode: http.StatusTooManyRequests, headers: headers}},
		{resp: statusResponse(http.StatusOK)},
	}}

	policy := DefaultMirrorNodeHttpRetryPolicy()
	policy = policy.WithMaxAttempts(2)
	// Backoff that would dominate the test if Retry-After were ignored.
	policy = policy.WithInitialBackoff(5 * time.Second)
	policy = policy.WithMaxBackoff(5 * time.Second)
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy)

	start := time.Now()
	resp, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Less(t, elapsed, time.Second)
}

func TestUnitMirrorHttpClientHonoursRetryAfterAboveMaxBackoff(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: HttpResponse{statusCode: http.StatusServiceUnavailable, headers: map[string][]string{httpHeaderRetryAfter: {"1"}}}},
		{resp: statusResponse(http.StatusOK)},
	}}

	policy := fastRetryPolicy(2)
	policy = policy.WithTotalDeadline(10 * time.Second)
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy)

	start := time.Now()
	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, elapsed, 900*time.Millisecond)
}

func TestUnitMirrorHttpClientFailsFastWhenRetryAfterExceedsDeadline(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: HttpResponse{statusCode: http.StatusTooManyRequests, headers: map[string][]string{httpHeaderRetryAfter: {"3600"}}}},
		{resp: statusResponse(http.StatusOK)},
	}}

	policy := fastRetryPolicy(5)
	policy = policy.WithTotalDeadline(5 * time.Second)
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy)

	start := time.Now()
	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())

	require.ErrorIs(t, err, errMirrorHttpDeadlineExceeded)
	assert.Less(t, time.Since(start), time.Second)
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorHttpClientIgnoresRetryAfterOnRequestTimeoutStatus(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: HttpResponse{statusCode: http.StatusRequestTimeout, headers: map[string][]string{httpHeaderRetryAfter: {"3600"}}}},
		{resp: statusResponse(http.StatusOK)},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(2))

	start := time.Now()
	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())

	require.NoError(t, err)
	assert.Less(t, time.Since(start), time.Second)
}

func TestUnitMirrorHttpClientBoundsTotalWallClock(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusServiceUnavailable), delay: 200 * time.Millisecond},
	}}

	policy := fastRetryPolicy(20)
	policy = policy.WithPerAttemptTimeout(time.Second)
	policy = policy.WithTotalDeadline(300 * time.Millisecond)
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy)

	start := time.Now()
	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	elapsed := time.Since(start)

	require.ErrorIs(t, err, errMirrorHttpDeadlineExceeded)
	require.NotErrorIs(t, err, errMirrorHttpRetriesExhausted)
	assert.Less(t, elapsed, 2*time.Second)
	assert.Less(t, transport.callCount(), 20)
}

func TestUnitMirrorHttpClientTotalDeadlineSpansCall(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusOK), delay: 200 * time.Millisecond},
	}}

	policy := fastRetryPolicy(1)
	policy = policy.WithTotalDeadline(300 * time.Millisecond)
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy)

	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.NoError(t, err)

	_, err = client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.Error(t, err)
}

func TestUnitMirrorHttpClientAppliesPerAttemptTimeout(t *testing.T) {
	t.Parallel()

	const maxAttempts = 2
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusOK), delay: 500 * time.Millisecond},
	}}

	policy := fastRetryPolicy(maxAttempts)
	policy = policy.WithPerAttemptTimeout(50 * time.Millisecond)
	policy = policy.WithTotalDeadline(5 * time.Second)
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy)

	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.Error(t, err)
	assert.Equal(t, maxAttempts, transport.callCount())
}

func TestUnitMirrorHttpClientHonoursCallerCancellation(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusOK), delay: 5 * time.Second},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(3))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationFromContext(ctx))

	require.ErrorIs(t, err, context.Canceled)
	requireHttpErrorID(t, err, HttpTransportCancelledError)
	assert.Less(t, time.Since(start), 2*time.Second)
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorHttpClientCancellationInterruptsBackoff(t *testing.T) {
	t.Parallel()

	// Retry-After makes the wait exactly 5s; a jittered backoff could draw one shorter than the
	// delay before cancelling.
	throttled := HttpResponse{statusCode: http.StatusServiceUnavailable, headers: map[string][]string{httpHeaderRetryAfter: {"5"}}}
	transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: throttled}}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(3))

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(50*time.Millisecond, cancel)

	start := time.Now()
	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationFromContext(ctx))

	requireHttpErrorID(t, err, HttpTransportCancelledError)
	assert.Less(t, time.Since(start), time.Second)
}

func TestUnitMirrorHttpClientRetriesPostLikeGet(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{
		{resp: statusResponse(http.StatusServiceUnavailable)},
		{resp: statusResponse(http.StatusOK)},
	}}
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, fastRetryPolicy(3))

	resp, err := client.post(testMirrorPath(t, "/contracts/call"), "application/json", []byte(`{}`), CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, 2, transport.callCount())
	assert.Equal(t, HttpMethodPost, transport.lastRequest().method)
	assert.Equal(t, []byte(`{}`), transport.lastRequest().body)
}

func TestUnitMirrorHttpClientResolvesPathAgainstBaseURL(t *testing.T) {
	t.Parallel()

	for _, base := range []string{"https://mirror.example.com/api/v1", "https://mirror.example.com/api/v1/"} {
		transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
		client := newMirrorNodeHttpClient(base, transport, fastRetryPolicy(1))

		_, err := client.get(testMirrorPath(t, "/accounts/0.0.1?limit=25"), CancellationNone())
		require.NoError(t, err)
		assert.Equal(t,
			"https://mirror.example.com/api/v1/accounts/0.0.1?limit=25",
			transport.lastRequest().url,
			"both spellings of the base URL must resolve identically")
	}
}

func TestUnitMirrorHttpShouldRetryStatus(t *testing.T) {
	t.Parallel()

	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", &fakeHttpTransport{}, DefaultMirrorNodeHttpRetryPolicy())

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

	policy := DefaultMirrorNodeHttpRetryPolicy()
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", &fakeHttpTransport{}, policy)

	distinct := make(map[time.Duration]struct{})
	for attempt := range 6 {
		for range 40 {
			delay := client.backoffDelay(attempt)
			require.GreaterOrEqual(t, delay, time.Duration(0))
			require.LessOrEqual(t, delay, policy.GetMaxBackoff())
			distinct[delay] = struct{}{}
		}
	}

	assert.Greater(t, len(distinct), 1)
	assert.LessOrEqual(t, client.backoffDelay(0), policy.GetInitialBackoff())
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

			headers := map[string][]string{}
			if tt.value != "" {
				headers[httpHeaderRetryAfter] = []string{tt.value}
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

		headers := map[string][]string{
			httpHeaderRetryAfter: {time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)},
		}

		got, ok := retryAfterDelay(headers)
		require.True(t, ok)
		assert.Greater(t, got, 20*time.Second)
		assert.LessOrEqual(t, got, 30*time.Second)
	})
}

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

	client := newMirrorNodeHttpClient(server.URL+mirrorHttpAPIVersionPrefix, newTestHttpTransport(t, 0), DefaultMirrorNodeHttpRetryPolicy())

	resp, err := client.get(testMirrorPath(t, "/balances?account.id=0.0.1"), CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, `{"balances":[]}`, string(resp.body))
	assert.Equal(t, int32(2), atomic.LoadInt32(&attempts))
	assert.Equal(t, "/api/v1/balances?account.id=0.0.1", gotPath)
}

func TestUnitMirrorHttpClientPassesCancellationThroughUnwrapped(t *testing.T) {
	t.Parallel()

	transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	policy := fastRetryPolicy(1).WithTotalDeadline(time.Minute)
	client := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := client.get(testMirrorPath(t, testMirrorAccountPath), CancellationFromContext(ctx))
	require.NoError(t, err)

	transport.mu.Lock()
	defer transport.mu.Unlock()
	assert.Equal(t, ctx, transport.cancellations[0].GetContext())
}

func TestUnitMirrorHttpClientDerivesEachAttemptDeadline(t *testing.T) {
	t.Parallel()

	path := testMirrorPath(t, testMirrorAccountPath)
	deadlineFor := func(policy MirrorNodeHttpRetryPolicy) time.Duration {
		transport := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
		_, err := newMirrorNodeHttpClient("https://mirror.example.com/api/v1", transport, policy).get(path, CancellationNone())
		require.NoError(t, err)
		return transport.lastRequest().deadline
	}

	assert.Equal(t, time.Second, deadlineFor(fastRetryPolicy(1).WithPerAttemptTimeout(time.Second).WithTotalDeadline(time.Minute)),
		"a per-attempt cap tighter than the remaining total is used as is")

	fromTotal := deadlineFor(fastRetryPolicy(1).WithPerAttemptTimeout(30 * time.Second).WithTotalDeadline(2 * time.Second))
	assert.LessOrEqual(t, fromTotal, 2*time.Second)
	assert.Greater(t, fromTotal, time.Second)

	assert.Equal(t, 2*time.Second, deadlineFor(fastRetryPolicy(1).WithPerAttemptTimeout(0).WithTotalDeadline(2*time.Second)).Round(time.Second),
		"0 is no per-attempt cap, so the remaining total is the whole bound")
	assert.Zero(t, deadlineFor(fastRetryPolicy(1).WithPerAttemptTimeout(0)))
}
