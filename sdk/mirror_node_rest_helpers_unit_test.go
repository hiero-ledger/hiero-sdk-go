//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestUnitMirrorNodeRestBaseURLRejectsNilClient(t *testing.T) {
	t.Parallel()

	_, err := mirrorNodeRestBaseURL(nil)
	require.ErrorIs(t, err, errNoClientProvided)
}

func TestUnitMirrorNodeRestBaseURLRejectsEmptyMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{})

	_, err = mirrorNodeRestBaseURL(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node is not set")
}

func TestUnitMirrorNodeValidateURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{name: "https", url: "https://mirror.example.com/api/v1"},
		{name: "http", url: "http://localhost:5551/api/v1"},
		{name: "unparseable", url: "http://[::1", wantErr: "invalid mirror node URL"},
		{name: "unsupported scheme", url: "ftp://mirror.example.com", wantErr: "unsupported mirror node URL scheme"},
		{name: "no scheme", url: "mirror.example.com/api/v1", wantErr: "unsupported mirror node URL scheme"},
		{name: "no host", url: "https:///api/v1", wantErr: "has no host"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := mirrorNodeValidateURL(tt.url)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// A URL no request could satisfy must fail before the retry loop spends its budget on it.
func TestUnitMirrorNodeGetAndPostRejectInvalidURLWithoutRequesting(t *testing.T) {
	client, err := _NewMockClient()
	require.NoError(t, err)

	_, err = mirrorNodeGetWithRetry(client, "ftp://mirror.example.com", 3, time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported mirror node URL scheme")

	_, err = mirrorNodePostWithRetry(client, "ftp://mirror.example.com", "application/json", []byte("{}"), 3, time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported mirror node URL scheme")
}

func TestUnitMirrorNodeReadBodyRejectsNilResponse(t *testing.T) {
	t.Parallel()

	_, err := mirrorNodeReadBody(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "received nil response from mirror node")
}

func TestUnitMirrorNodeReadBodyReportsNon200WithDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"_status":"invalid parameter"}`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL) // #nosec
	require.NoError(t, err)

	_, err = mirrorNodeReadBody(resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "invalid parameter")
}

func TestUnitMirrorNodeReadBodyReturnsBodyOn200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL) // #nosec
	require.NoError(t, err)

	body, err := mirrorNodeReadBody(resp)
	require.NoError(t, err)
	assert.JSONEq(t, `{"ok":true}`, string(body))
}

// The happy paths live in registered_node_address_book_query_unit_test.go; this is the error branch.
func TestUnitResolveNextURLRejectsUnparseableLink(t *testing.T) {
	t.Parallel()

	base, err := url.Parse("https://mirror.example.com/api/v1/things")
	require.NoError(t, err)

	_, err = resolveNextURL(base, "://nope")
	require.Error(t, err)
}

func TestUnitMirrorNodeWalkPagesStopsWhenNextIsAbsent(t *testing.T) {
	t.Parallel()

	var fetched []string
	err := mirrorNodeWalkPages("https://mirror.example.com/api/v1/things", 10,
		func(pageURL string) ([]byte, error) {
			fetched = append(fetched, pageURL)
			return []byte(`{}`), nil
		},
		func(body []byte) (*string, error) { return nil, nil },
	)

	require.NoError(t, err)
	assert.Equal(t, []string{"https://mirror.example.com/api/v1/things"}, fetched)
}

func TestUnitMirrorNodeWalkPagesFollowsNextAcrossPages(t *testing.T) {
	t.Parallel()

	next := []string{"/api/v1/things?page=2", "/api/v1/things?page=3", ""}
	var fetched []string
	page := 0

	err := mirrorNodeWalkPages("https://mirror.example.com/api/v1/things", 10,
		func(pageURL string) ([]byte, error) {
			fetched = append(fetched, pageURL)
			return []byte(`{}`), nil
		},
		func(body []byte) (*string, error) {
			n := next[page]
			page++
			return &n, nil
		},
	)

	require.NoError(t, err)
	assert.Equal(t, []string{
		"https://mirror.example.com/api/v1/things",
		"https://mirror.example.com/api/v1/things?page=2",
		"https://mirror.example.com/api/v1/things?page=3",
	}, fetched, "each links.next is resolved against the base and fetched in order")
}

// The cap is what stops a mirror node that always returns a next link from looping forever.
func TestUnitMirrorNodeWalkPagesStopsAtPageCap(t *testing.T) {
	t.Parallel()

	const cap = 3
	pages := 0
	forever := "/api/v1/things?page=next"

	err := mirrorNodeWalkPages("https://mirror.example.com/api/v1/things", cap,
		func(pageURL string) ([]byte, error) { pages++; return []byte(`{}`), nil },
		func(body []byte) (*string, error) { return &forever, nil },
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeded pagination cap of 3 pages")
	assert.Equal(t, cap, pages, "no more than the cap may be fetched")
}

func TestUnitMirrorNodeWalkPagesPropagatesErrors(t *testing.T) {
	t.Parallel()

	t.Run("invalid start URL", func(t *testing.T) {
		t.Parallel()
		err := mirrorNodeWalkPages("://nope", 10,
			func(string) ([]byte, error) { return nil, nil },
			func([]byte) (*string, error) { return nil, nil },
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid mirror node page URL")
	})

	t.Run("fetch failure", func(t *testing.T) {
		t.Parallel()
		err := mirrorNodeWalkPages("https://mirror.example.com/api/v1/things", 10,
			func(string) ([]byte, error) { return nil, errors.New("boom") },
			func([]byte) (*string, error) { return nil, nil },
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "boom")
	})

	t.Run("handler failure", func(t *testing.T) {
		t.Parallel()
		err := mirrorNodeWalkPages("https://mirror.example.com/api/v1/things", 10,
			func(string) ([]byte, error) { return []byte(`{}`), nil },
			func([]byte) (*string, error) { return nil, errors.New("bad page") },
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bad page")
	})

	t.Run("invalid next link", func(t *testing.T) {
		t.Parallel()
		bad := "://nope"
		err := mirrorNodeWalkPages("https://mirror.example.com/api/v1/things", 10,
			func(string) ([]byte, error) { return []byte(`{}`), nil },
			func([]byte) (*string, error) { return &bad, nil },
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid pagination next link")
	})
}
