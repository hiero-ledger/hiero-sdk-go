//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitMirrorNodePostWithRetryRejectsZeroMaxAttempts(t *testing.T) {
	// Note: Not running in parallel since sibling tests modify global http.DefaultTransport.
	client, err := _NewMockClient()
	require.NoError(t, err)

	resp, err := mirrorNodePostWithRetry(client, "https://example.com", "application/json", []byte("{}"), 0, time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maxAttempts must be at least 1")
	assert.Nil(t, resp)
}

func TestUnitMirrorNodePostWithRetryReturnsLastResponseAfterExhaustion(t *testing.T) {
	// Note: Not running in parallel since sibling tests modify global http.DefaultTransport.
	const maxAttempts = uint64(2)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		// Always transient: the helper retries until attempts are exhausted, then
		// returns the final (still non-200) response with a nil error so the caller
		// can format its own error.
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client, err := _NewMockClient()
	require.NoError(t, err)

	resp, err := mirrorNodePostWithRetry(client, server.URL, "application/json", []byte("{}"), maxAttempts, time.Second)
	require.NoError(t, err, "an exhausted retryable response is returned without a transport error")
	require.NotNil(t, resp)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	assert.Equal(t, int32(maxAttempts), atomic.LoadInt32(&attempts), "every attempt should be used before giving up")
}

func TestUnitMirrorNodePostWithRetryRetriesOnTooManyRequests(t *testing.T) {
	// Note: Not running in parallel since sibling tests modify global http.DefaultTransport.
	const maxAttempts = uint64(2)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		// 429 is transient: the helper must retry it (unlike a genuine 4xx) and, once
		// attempts are exhausted, return the final 429 response with a nil error.
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client, err := _NewMockClient()
	require.NoError(t, err)

	resp, err := mirrorNodePostWithRetry(client, server.URL, "application/json", []byte("{}"), maxAttempts, time.Second)
	require.NoError(t, err, "an exhausted retryable response is returned without a transport error")
	require.NotNil(t, resp)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, int32(maxAttempts), atomic.LoadInt32(&attempts), "429 should be retried up to maxAttempts")
}

func TestUnitMirrorNodePostWithRetryReturnsTransportErrorAfterExhaustion(t *testing.T) {
	// Note: Not running in parallel since sibling tests modify global http.DefaultTransport.
	const maxAttempts = uint64(2)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		// Drop the connection every time, surfacing as a transport error (EOF) on
		// every attempt so the helper exhausts its retries and returns the error.
		hj, ok := w.(http.Hijacker)
		require.True(t, ok, "test server must support connection hijacking")
		conn, _, hijackErr := hj.Hijack()
		require.NoError(t, hijackErr)
		_ = conn.Close()
	}))
	defer server.Close()

	client, err := _NewMockClient()
	require.NoError(t, err)

	resp, err := mirrorNodePostWithRetry(client, server.URL, "application/json", []byte("{}"), maxAttempts, time.Second)
	require.Error(t, err, "a transport failure on every attempt is surfaced as an error")
	assert.Nil(t, resp)
	assert.Equal(t, int32(maxAttempts), atomic.LoadInt32(&attempts), "every attempt should be used before giving up")
}

func TestUnitMirrorNodePostWithRetryStopsOnNetworkUpdateCancellation(t *testing.T) {
	// Note: Not running in parallel since sibling tests modify global http.DefaultTransport.
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client, err := _NewMockClient()
	require.NoError(t, err)

	// Cancelling the network-update context aborts the backoff wait between attempts.
	client.cancelNetworkUpdate()

	resp, err := mirrorNodePostWithRetry(client, server.URL, "application/json", []byte("{}"), 3, time.Second)
	require.Error(t, err)
	assert.ErrorIs(t, err, client.networkUpdateContext.Err())
	assert.Nil(t, resp)
	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts), "backoff cancellation should stop before a second attempt")
}
