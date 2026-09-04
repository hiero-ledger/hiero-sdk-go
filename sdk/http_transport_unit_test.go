//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHttpTransport(t *testing.T, connectTimeout time.Duration) *defaultHttpTransport {
	t.Helper()

	transport := newDefaultHttpTransport(connectTimeout)
	t.Cleanup(func() {
		require.NoError(t, transport.close(time.Second))
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

		resp, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: server.URL})
		require.NoError(t, err, "a completed exchange is not an error, whatever its status")
		assert.Equal(t, status, resp.statusCode)
		assert.Contains(t, string(resp.body), "nope", "the body is available for the caller to build its own message")
	}
}

func TestUnitHttpTransportSendsUserAgentAndContentType(t *testing.T) {
	t.Parallel()

	var gotAgent, gotContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAgent = r.Header.Get(httpHeaderUserAgent)
		gotContentType = r.Header.Get(httpHeaderContentType)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.roundTrip(context.Background(), httpRequest{
		method:      httpMethodPost,
		url:         server.URL,
		body:        []byte(`{}`),
		contentType: "application/json",
	})
	require.NoError(t, err)
	assert.Contains(t, gotAgent, "hiero-sdk-go")
	assert.Equal(t, "application/json", gotContentType)
}

func TestUnitHttpTransportAppendsCallerUserAgent(t *testing.T) {
	t.Parallel()

	var gotAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAgent = r.Header.Get(httpHeaderUserAgent)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.roundTrip(context.Background(), httpRequest{
		method:  httpMethodGet,
		url:     server.URL,
		headers: map[string]string{httpHeaderUserAgent: "my-app/1.0"},
	})
	require.NoError(t, err)
	assert.Contains(t, gotAgent, "hiero-sdk-go", "the SDK token is always present")
	assert.Contains(t, gotAgent, "my-app/1.0", "the caller's agent is appended, not replaced")
}

func TestUnitHttpTransportClassifiesUnresolvableHostAsTerminal(t *testing.T) {
	t.Parallel()

	transport := newTestHttpTransport(t, 2*time.Second)

	// .invalid is reserved and never resolves, so this is NXDOMAIN rather than a timeout.
	_, err := transport.roundTrip(context.Background(), httpRequest{
		method: httpMethodGet,
		url:    "https://mirror-node.invalid/api/v1/accounts/0.0.1",
	})
	require.Error(t, err)

	kind, ok := httpErrorKindOf(err)
	require.True(t, ok, "the transport must classify what it returns")
	assert.Equal(t, httpTerminal, kind, "retrying cannot make a name resolve")
}

func TestUnitHttpTransportClassifiesUntrustedCertificateAsTerminal(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: server.URL})
	require.Error(t, err)

	kind, ok := httpErrorKindOf(err)
	require.True(t, ok)
	assert.Equal(t, httpTerminal, kind, "re-validating a bad certificate is not a retry strategy")
}

func TestUnitHttpTransportClassifiesRefusedConnectionAsTransient(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	url := server.URL
	server.Close()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: url})
	require.Error(t, err)

	kind, ok := httpErrorKindOf(err)
	require.True(t, ok)
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
	_, err = transport.roundTrip(context.Background(), httpRequest{
		method: httpMethodGet,
		url:    fmt.Sprintf("https://%s/api/v1/accounts", listener.Addr().String()),
	})
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Less(t, elapsed, 3*time.Second, "connect must not wait for the per-attempt deadline")
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

	_, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: server.URL})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "stopped after 5 redirects")
	assert.LessOrEqual(t, atomic.LoadInt32(&hops), int32(httpMaxRedirects+1))
}

func TestUnitHttpTransportDropsAuthorizationCrossOrigin(t *testing.T) {
	t.Parallel()

	var secondHopAuth string
	var sawSecondHop atomic.Bool
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secondHopAuth = r.Header.Get(httpHeaderAuthorization)
		sawSecondHop.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer destination.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, destination.URL, http.StatusFound)
	}))
	defer origin.Close()

	transport := newTestHttpTransport(t, time.Second)

	resp, err := transport.roundTrip(context.Background(), httpRequest{
		method:  httpMethodGet,
		url:     origin.URL,
		headers: map[string]string{httpHeaderAuthorization: "Bearer super-secret"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.statusCode)
	require.True(t, sawSecondHop.Load(), "the redirect should have been followed")
	assert.Empty(t, secondHopAuth, "credentials must not cross an origin boundary")
}

func TestUnitHttpTransportRejectsWorkAfterClose(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newDefaultHttpTransport(time.Second)
	require.NoError(t, transport.close(time.Second))

	_, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: server.URL})
	require.Error(t, err)
	require.ErrorIs(t, err, errHttpTransportClosed)

	kind, ok := httpErrorKindOf(err)
	require.True(t, ok)
	assert.Equal(t, httpClosed, kind, "a closed transport is permanent, never retried")
}

func TestUnitHttpTransportCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	transport := newDefaultHttpTransport(time.Second)

	require.NoError(t, transport.close(time.Second))
	require.NoError(t, transport.close(time.Second), "closing twice is a no-op, not a failure")
}

func TestUnitHttpTransportCloseWaitsForInFlightRequest(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := newDefaultHttpTransport(time.Second)

	done := make(chan error, 1)
	go func() {
		_, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: server.URL})
		done <- err
	}()

	// Give the request time to register as in-flight before closing.
	time.Sleep(100 * time.Millisecond)

	closed := make(chan struct{})
	go func() {
		require.NoError(t, transport.close(5*time.Second))
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
		require.NoError(t, err, "an exchange already in flight is allowed to finish")
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

	transport := newDefaultHttpTransport(time.Second)

	done := make(chan error, 1)
	go func() {
		_, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: server.URL})
		done <- err
	}()

	time.Sleep(100 * time.Millisecond)

	start := time.Now()
	require.NoError(t, transport.close(200*time.Millisecond))
	assert.Less(t, time.Since(start), 2*time.Second, "close is bounded by its grace period")

	select {
	case err := <-done:
		require.Error(t, err, "an exchange past the grace period is aborted, not left hanging")
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

	_, err := transport.roundTrip(ctx, httpRequest{method: httpMethodGet, url: server.URL})
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled, "a caller can cancel a single mirror read")

	kind, ok := httpErrorKindOf(err)
	require.True(t, ok)
	assert.Equal(t, httpTerminal, kind, "the caller asked to stop; do not retry against their wishes")
}

func TestUnitHttpTransportRejectsMalformedURLAsTerminal(t *testing.T) {
	t.Parallel()

	transport := newTestHttpTransport(t, time.Second)

	_, err := transport.roundTrip(context.Background(), httpRequest{method: httpMethodGet, url: "://nope"})
	require.Error(t, err)

	kind, ok := httpErrorKindOf(err)
	require.True(t, ok)
	assert.Equal(t, httpTerminal, kind)
}

func TestUnitHttpErrorKindOf(t *testing.T) {
	t.Parallel()

	kind, ok := httpErrorKindOf(&httpError{kind: httpTerminal, err: errors.New("boom")})
	assert.True(t, ok)
	assert.Equal(t, httpTerminal, kind)

	// Wrapping must not lose the classification.
	wrapped := errors.Join(errors.New("context"), &httpError{kind: httpClosed, err: errHttpTransportClosed})
	kind, ok = httpErrorKindOf(wrapped)
	assert.True(t, ok)
	assert.Equal(t, httpClosed, kind)

	kind, ok = httpErrorKindOf(errors.New("plain"))
	assert.False(t, ok)
	assert.Equal(t, httpTransient, kind, "an unclassified error defaults to retryable")
}
