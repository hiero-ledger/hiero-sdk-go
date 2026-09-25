//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requireHttpErrorID asserts that err is an HttpTransportError with id.
func requireHttpErrorID(t *testing.T, err error, id HttpTransportErrorID) {
	t.Helper()

	var httpErr *HttpTransportError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, id, httpErr.ID, "got %s", httpErr.ID)
}

func newTestHttpTransport(t *testing.T, connectTimeout time.Duration) *DefaultHttpTransport {
	t.Helper()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().WithConnectTimeout(connectTimeout))
	t.Cleanup(func() {
		transport.Close(time.Second)
	})

	return transport
}

func TestUnitHttpTransportReturnsNon2xxAsResponse(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusBadRequest, http.StatusNotFound, http.StatusServiceUnavailable} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"_status":{"messages":[{"detail":"nope"}]}}`))
		}))
		defer server.Close()

		transport := newTestHttpTransport(t, time.Second)

		resp, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
		require.NoError(t, err)
		assert.Equal(t, status, resp.statusCode)
		assert.Contains(t, string(resp.body), "nope")
	}
}

func TestUnitHttpTransportLowercasesResponseHeaders(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "5")
		w.Header().Add("X-Mirror-Trace", "a")
		w.Header().Add("X-Mirror-Trace", "b")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	resp, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, []string{"5"}, resp.headers["retry-after"])
	assert.Equal(t, []string{"a", "b"}, resp.headers["x-mirror-trace"])
	assert.NotContains(t, resp.headers, "Retry-After")
}

func TestUnitHttpTransportSendsXUserAgentAndContentType(t *testing.T) {
	t.Parallel()

	var gotAgent, gotUserAgent, gotContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAgent = r.Header.Get(httpHeaderXUserAgent)
		gotUserAgent = r.Header.Get("User-Agent")
		gotContentType = r.Header.Get(httpHeaderContentType)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.RoundTrip(HttpRequest{
		method:      HttpMethodPost,
		url:         server.URL,
		body:        []byte(`{}`),
		contentType: "application/json",
	}, CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, getUserAgent(), gotAgent)
	assert.NotContains(t, gotUserAgent, "hiero-sdk-go")
	assert.Equal(t, "application/json", gotContentType)
}

func TestUnitHttpTransportXUserAgentCannotBeReplaced(t *testing.T) {
	t.Parallel()

	var gotAgent []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAgent = r.Header.Values(httpHeaderXUserAgent)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.RoundTrip(HttpRequest{
		method:  HttpMethodGet,
		url:     server.URL,
		headers: map[string]string{"X-User-Agent": "my-app/1.0"},
	}, CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, []string{getUserAgent()}, gotAgent)
}

func TestUnitHttpTransportClassifiesUnresolvableHostAsTerminal(t *testing.T) {
	t.Parallel()

	base := &http.Transport{DialContext: func(context.Context, string, string) (net.Conn, error) {
		return nil, &net.DNSError{Err: "no such host", Name: "mirror.example.com", IsNotFound: true}
	}}
	transport := newHttpTransportOver(base, DefaultHttpTransportConfiguration())
	defer transport.Close(0)

	_, err := transport.RoundTrip(HttpRequest{
		method: HttpMethodGet,
		url:    "https://mirror.example.com/api/v1/accounts/0.0.1",
	}, CancellationNone())
	require.Error(t, err)

	requireHttpErrorID(t, err, HttpTransportUnknownHostError)
	kind, _ := httpErrorKindOf(err)
	assert.Equal(t, httpTerminal, kind)
}

func TestUnitHttpTransportClassifiesUntrustedCertificateAsTerminal(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
	require.Error(t, err)

	requireHttpErrorID(t, err, HttpTransportTLSError)
	kind, _ := httpErrorKindOf(err)
	assert.Equal(t, httpTerminal, kind)
}

func TestUnitHttpTransportClassifiesRefusedConnectionAsTransient(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	url := server.URL
	server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: url}, CancellationNone())
	require.Error(t, err)

	requireHttpErrorID(t, err, HttpTransportConnectionError)
	kind, _ := httpErrorKindOf(err)
	assert.Equal(t, httpTransient, kind)
}

func TestUnitHttpTransportHonoursConnectTimeout(t *testing.T) {
	t.Parallel()

	// Accept the TCP connection but never complete the TLS handshake, so the handshake
	// deadline — which is set from connectTimeout — is what ends the attempt.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			// Hold it open and say nothing.
			t.Cleanup(func() { _ = conn.Close() })
		}
	}()

	transport := newTestHttpTransport(t, 150*time.Millisecond)

	start := time.Now()
	_, err = transport.RoundTrip(HttpRequest{
		method: HttpMethodGet,
		url:    fmt.Sprintf("https://%s/api/v1/accounts", listener.Addr().String()),
	}, CancellationNone())
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Less(t, elapsed, 3*time.Second)
}

func TestUnitHttpTransportBoundsRedirects(t *testing.T) {
	t.Parallel()

	var hops int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hops, 1)
		http.Redirect(w, r, "/again", http.StatusFound)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "stopped after 5 redirects")
	assert.LessOrEqual(t, atomic.LoadInt32(&hops), int32(httpDefaultMaxRedirects+1))
}

func TestUnitHttpTransportDropsCallerHeadersCrossOrigin(t *testing.T) {
	t.Parallel()

	var secondHopAuth, secondHopKey, secondHopAgent string
	var sawSecondHop atomic.Bool
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secondHopAuth = r.Header.Get("Authorization")
		secondHopKey = r.Header.Get("X-Api-Key")
		secondHopAgent = r.Header.Get(httpHeaderXUserAgent)
		sawSecondHop.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer destination.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, destination.URL, http.StatusFound)
	}))
	defer origin.Close()

	transport := newTestHttpTransport(t, time.Second)

	resp, err := transport.RoundTrip(HttpRequest{
		method:  HttpMethodGet,
		url:     origin.URL,
		headers: map[string]string{"Authorization": "Bearer super-secret", "x-api-key": "also-secret"},
	}, CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	require.True(t, sawSecondHop.Load())
	assert.Empty(t, secondHopAuth)
	assert.Empty(t, secondHopKey)
	assert.Equal(t, getUserAgent(), secondHopAgent)
}

func TestUnitHttpTransportAppliesConfiguredBounds(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		size, _ := strconv.Atoi(r.URL.Query().Get("size"))
		_, _ = w.Write(make([]byte, size))
	}))
	defer server.Close()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().WithMaxResponseBytes(1024).WithMaxRedirects(0))
	t.Cleanup(func() { transport.Close(time.Second) })

	resp, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL + "/?size=1024"}, CancellationNone())
	require.NoError(t, err)
	assert.Len(t, resp.body, 1024)

	_, err = transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL + "/?size=1025"}, CancellationNone())
	requireHttpErrorID(t, err, HttpTransportResponseTooLargeError)

	_, err = transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL + "/redirect"}, CancellationNone())
	require.Error(t, err)
}

func TestUnitHttpTransportHeaderPrecedence(t *testing.T) {
	t.Parallel()

	var got http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().
		WithDefaultHeader("X-Trace", "default").
		WithDefaultHeader("authorization", "default").
		WithDefaultHeader("content-type", "text/plain").
		WithDefaultHeader(httpHeaderXUserAgent, "default"))
	t.Cleanup(func() { transport.Close(time.Second) })

	_, err := transport.RoundTrip(HttpRequest{
		method:      HttpMethodPost,
		url:         server.URL,
		body:        []byte(`{}`),
		contentType: "application/json",
		headers:     map[string]string{"authorization": "caller", "content-type": "text/html"},
	}, CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, "default", got.Get("X-Trace"))
	assert.Equal(t, "caller", got.Get("Authorization"))
	assert.Equal(t, "application/json", got.Get(httpHeaderContentType))
	assert.Equal(t, getUserAgent(), got.Get(httpHeaderXUserAgent))
}

func TestUnitHttpTransportRejectsWorkAfterClose(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().WithConnectTimeout(time.Second))
	transport.Close(time.Second)

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
	require.Error(t, err)
	require.ErrorIs(t, err, errHttpTransportClosed)

	requireHttpErrorID(t, err, HttpTransportClientClosedError)
	kind, _ := httpErrorKindOf(err)
	assert.Equal(t, httpClosed, kind)
}

func TestUnitHttpTransportCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().WithConnectTimeout(time.Second))

	transport.Close(time.Second)
	transport.Close(time.Second)
}

func TestUnitHttpTransportCloseWaitsForInFlightRequest(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().WithConnectTimeout(time.Second))

	done := make(chan error, 1)
	go func() {
		_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
		done <- err
	}()

	// Give the request time to register as in-flight before closing.
	time.Sleep(100 * time.Millisecond)

	closed := make(chan struct{})
	go func() {
		transport.Close(5 * time.Second)
		close(closed)
	}()

	select {
	case <-closed:
		t.Fatal("close returned before the in-flight request finished")
	case <-time.After(150 * time.Millisecond):
	}

	close(release)

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("in-flight request never completed")
	}

	select {
	case <-closed:
	case <-time.After(3 * time.Second):
		t.Fatal("close never returned after the request drained")
	}
}

func TestUnitHttpTransportCloseAbortsAfterGracePeriod(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	// Unblock the handler before httptest.Server.Close waits on it: defers run LIFO.
	defer close(release)

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().WithConnectTimeout(time.Second))

	done := make(chan error, 1)
	go func() {
		_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
		done <- err
	}()

	time.Sleep(100 * time.Millisecond)

	start := time.Now()
	transport.Close(200 * time.Millisecond)
	assert.Less(t, time.Since(start), 2*time.Second)

	select {
	case err := <-done:
		require.Error(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("abort did not reach the in-flight request")
	}
}

func TestUnitHttpTransportPropagatesCallerCancellation(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	// Unblock the handler before httptest.Server.Close waits on it: defers run LIFO.
	defer close(release)

	transport := newTestHttpTransport(t, time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationFromContext(ctx))
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)

	requireHttpErrorID(t, err, HttpTransportCancelledError)
	kind, _ := httpErrorKindOf(err)
	assert.Equal(t, httpTerminal, kind)
}

func TestUnitHttpTransportRejectsMalformedURLAsTerminal(t *testing.T) {
	t.Parallel()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: "://nope"}, CancellationNone())
	require.Error(t, err)

	kind, ok := httpErrorKindOf(err)
	assert.False(t, ok)
	assert.Equal(t, httpTerminal, kind)
}

func TestUnitHttpErrorKindOf(t *testing.T) {
	t.Parallel()

	kind, ok := httpErrorKindOf(&HttpTransportError{ID: HttpTransportTLSError, Err: errors.New("boom")})
	assert.True(t, ok)
	assert.Equal(t, httpTerminal, kind)

	// Wrapping must not lose the classification.
	wrapped := errors.Join(errors.New("context"), &HttpTransportError{ID: HttpTransportClientClosedError, Err: errHttpTransportClosed})
	kind, ok = httpErrorKindOf(wrapped)
	assert.True(t, ok)
	assert.Equal(t, httpClosed, kind)

	// The case an injected transport produces: a bare error nothing identified.
	kind, ok = httpErrorKindOf(errors.New("plain"))
	assert.False(t, ok)
	assert.Equal(t, httpTerminal, kind)
}

func TestUnitHttpTransportUsesEnvironmentProxy(t *testing.T) {
	t.Parallel()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration())
	t.Cleanup(func() { transport.Close(time.Second) })

	inner, ok := transport.httpClient.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, inner.Proxy)
}

func TestUnitHttpTransportRejectsOversizedBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		chunk := make([]byte, 1<<20)
		for range (httpDefaultMaxResponseBytes >> 20) + 1 {
			if _, err := w.Write(chunk); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, 0)

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
	require.ErrorIs(t, err, errHttpResponseTooLarge)

	requireHttpErrorID(t, err, HttpTransportResponseTooLargeError)
	kind, _ := httpErrorKindOf(err)
	assert.Equal(t, httpTerminal, kind)
}

func TestUnitHttpTransportMaxResponseBytesAtMaxInt64(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("body"))
	}))
	defer server.Close()

	transport := NewDefaultHttpTransport(DefaultHttpTransportConfiguration().WithMaxResponseBytes(math.MaxInt64))
	defer transport.Close(0)

	resp, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL}, CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, []byte("body"), resp.GetBody())
}

func TestUnitHttpTransportServesConcurrentRequests(t *testing.T) {
	t.Parallel()

	const callers = 24

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		caller, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/caller-"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = fmt.Fprintf(w, "caller-%d", caller)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, 0)

	var wg sync.WaitGroup
	bodies := make([]string, callers)
	errs := make([]error, callers)

	for i := range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("%s/caller-%d", server.URL, i)
			resp, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: url}, CancellationNone())
			bodies[i], errs[i] = string(resp.body), err
		}()
	}
	wg.Wait()

	for i := range callers {
		require.NoError(t, errs[i])
		assert.Equal(t, fmt.Sprintf("caller-%d", i), bodies[i])
	}
}

func TestUnitHttpErrorIDRetryVerdict(t *testing.T) {
	t.Parallel()

	cases := []struct {
		id   HttpTransportErrorID
		name string
		kind httpErrorKind
	}{
		{HttpTransportConnectionError, "connection-error", httpTransient},
		{HttpTransportTimeoutError, "timeout-error", httpTransient},
		{HttpTransportUnknownHostError, "unknown-host-error", httpTerminal},
		{HttpTransportTLSError, "tls-error", httpTerminal},
		{HttpTransportClientClosedError, "client-closed-error", httpClosed},
		{HttpTransportCancelledError, "cancelled-error", httpTerminal},
		{HttpTransportResponseTooLargeError, "response-too-large-error", httpTerminal},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.name, tc.id.String())
		assert.Equal(t, tc.kind, tc.id.kind(), tc.name)
	}
}

func TestUnitClassifyTransportError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		id   HttpTransportErrorID
	}{
		{"caller cancelled", fmt.Errorf("get: %w", context.Canceled), HttpTransportCancelledError},
		{"deadline", fmt.Errorf("get: %w", context.DeadlineExceeded), HttpTransportTimeoutError},
		{"temporary DNS failure", &net.DNSError{Err: "server misbehaving", IsTemporary: true}, HttpTransportUnknownHostError},
		{"refused", &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}, HttpTransportConnectionError},
		{"reset", fmt.Errorf("read: %w", syscall.ECONNRESET), HttpTransportConnectionError},
		{"server hung up", fmt.Errorf("get: %w", io.EOF), HttpTransportConnectionError},
		{"truncated body", io.ErrUnexpectedEOF, HttpTransportConnectionError},
		{"tls alert", tls.AlertError(40), HttpTransportTLSError},
	}

	for _, tc := range cases {
		id, ok := classifyTransportError(tc.err)
		assert.True(t, ok, tc.name)
		assert.Equal(t, tc.id, id, "%s: got %s", tc.name, id)
	}

	_, ok := classifyTransportError(errors.New("proxy returned something odd"))
	assert.False(t, ok)
}

func TestUnitHttpTransportDeadlineBoundsWholeBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		flusher, _ := w.(http.Flusher)
		for range 50 {
			select {
			case <-r.Context().Done():
				return
			case <-time.After(100 * time.Millisecond):
			}
			_, _ = w.Write([]byte("x"))
			flusher.Flush()
		}
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	start := time.Now()
	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL, deadline: 300 * time.Millisecond}, CancellationNone())

	requireHttpErrorID(t, err, HttpTransportTimeoutError)
	assert.Less(t, time.Since(start), 2*time.Second)
}

func TestUnitHttpTransportCancellationIsNotTimeout(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	defer close(release)

	transport := newTestHttpTransport(t, time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(50*time.Millisecond, cancel)

	_, err := transport.RoundTrip(HttpRequest{method: HttpMethodGet, url: server.URL, deadline: 5 * time.Second}, CancellationFromContext(ctx))
	requireHttpErrorID(t, err, HttpTransportCancelledError)
}

func TestUnitHttpTransportReleasesCancellationRegistrations(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)
	source := newCountingContext()

	for i := range 20 {
		req := HttpRequest{method: HttpMethodGet, url: server.URL}
		if i%2 == 0 {
			req.deadline = time.Second
		}
		_, err := transport.RoundTrip(req, CancellationFromContext(source))
		require.NoError(t, err)
	}

	assert.Zero(t, source.liveRegistrations())
}

func TestUnitNewHttpRequestValidatesURL(t *testing.T) {
	t.Parallel()

	req, err := NewHttpRequest(HttpMethodPost, "https://mirror.example.com/api/v1/contracts/call")
	require.NoError(t, err)
	derived := req.WithBody("application/json", []byte(`{}`)).WithHeader("X-Trace", "1").WithDeadline(time.Second)

	assert.Empty(t, req.GetBody())
	assert.Equal(t, HttpMethodPost, derived.GetMethod())
	assert.Equal(t, "application/json", derived.GetContentType())
	assert.Equal(t, map[string]string{"x-trace": "1"}, derived.GetHeaders())
	assert.Equal(t, time.Second, derived.GetDeadline())

	for _, bad := range []string{"", "/api/v1/accounts", "ftp://mirror.example.com", "https://mirror.example.com/a b"} {
		_, err := NewHttpRequest(HttpMethodGet, bad)
		assert.Error(t, err, bad)
	}
	assert.Panics(t, func() { req.WithDeadline(-time.Second) })
}

func TestUnitNewHttpResponseLowercasesHeaderNames(t *testing.T) {
	t.Parallel()

	resp := NewHttpResponse(http.StatusTooManyRequests, []byte("slow down"), map[string][]string{"Retry-After": {"5"}})

	assert.Equal(t, uint16(http.StatusTooManyRequests), resp.GetStatusCode())
	assert.Equal(t, []string{"5"}, resp.GetHeaders()["retry-after"])

	read := resp.GetHeaders()
	read["retry-after"][0] = "0"
	assert.Equal(t, []string{"5"}, resp.GetHeaders()["retry-after"])
}

func TestUnitHttpTransportErrorVerdicts(t *testing.T) {
	t.Parallel()

	assert.False(t, (&HttpTransportError{Err: errors.New("boom")}).IsRetryable())
	assert.True(t, (&HttpTransportError{ID: HttpTransportConnectionError}).IsRetryable())
	assert.False(t, (&HttpTransportError{ID: HttpTransportTLSError}).IsRetryable())

	kind, ok := httpErrorKindOf(fmt.Errorf("proxy hiccup: %w", ErrHttpTransportRetryable))
	assert.True(t, ok)
	assert.Equal(t, httpTransient, kind)
}
