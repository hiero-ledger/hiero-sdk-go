//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSuccessResponse(networkMultiplier int, nodeBase, serviceBase uint64) string {
	networkSubtotal := nodeBase * uint64(networkMultiplier)
	total := networkSubtotal + nodeBase + serviceBase
	return fmt.Sprintf(`{
  "network": {"multiplier": %d, "subtotal": %d},
  "node": {"base": %d, "extras": []},
  "service": {"base": %d, "extras": []},
  "notes": [],
  "total": %d
}`, networkMultiplier, networkSubtotal, nodeBase, serviceBase, total)
}

func newMockClientForREST() *Client {
	net := _NewNetwork()
	client := _NewClient(net, []string{"localhost:8084"}, nil, false, 0, 0)

	err := client.SetNetwork(map[string]AccountID{
		"127.0.0.1:50211": {Account: 3},
	})
	if err != nil {
		// Continue even if SetNetwork fails
	}

	dummyKey, _ := PrivateKeyGenerateEd25519()
	client.SetOperator(AccountID{Account: 1}, dummyKey)

	return client
}

func TestUnitFeeEstimateQueryRetriesOnUnavailableErrors(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		assert.Equal(t, "/api/v1/network/fees", r.URL.Path, "request path should be /api/v1/network/fees")

		contentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/protobuf", contentType, "Content-Type should be application/protobuf")

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Greater(t, len(body), 0, "request body should not be empty")

		queryParams := r.URL.Query()
		assert.Contains(t, queryParams, "mode", "request should contain mode query parameter")

		// First request returns 503, second returns success
		if requestCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("transient error"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(newSuccessResponse(2, 6, 8)))
		}
	}))
	defer server.Close()

	cleanup := SetupMockTransportForDomain("localhost:8084", server.URL)
	defer cleanup()

	client := newMockClientForREST()
	client.SetMaxBackoff(500 * time.Millisecond)

	tx := NewTransferTransaction()

	query := NewFeeEstimateQuery().
		SetTransaction(tx).
		SetMaxAttempts(3)

	response, err := query.Execute(client)
	require.NoError(t, err)

	assert.Equal(t, uint64(26), response.Total)
	assert.GreaterOrEqual(t, requestCount, 2, "should have retried at least once")
}

func TestUnitFeeEstimateQueryRetriesOnDeadlineExceededErrors(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		assert.Equal(t, "/api/v1/network/fees", r.URL.Path, "request path should be /api/v1/network/fees")

		contentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/protobuf", contentType, "Content-Type should be application/protobuf")

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Greater(t, len(body), 0, "request body should not be empty")

		queryParams := r.URL.Query()
		assert.Contains(t, queryParams, "mode", "request should contain mode query parameter")

		// First request returns 504, second returns success
		if requestCount == 1 {
			w.WriteHeader(http.StatusGatewayTimeout)
			_, _ = w.Write([]byte("gateway timeout"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(newSuccessResponse(4, 8, 20)))
		}
	}))
	defer server.Close()

	cleanup := SetupMockTransportForDomain("localhost:8084", server.URL)
	defer cleanup()

	client := newMockClientForREST()
	client.SetMaxBackoff(500 * time.Millisecond)

	tx := NewTransferTransaction()

	query := NewFeeEstimateQuery().
		SetTransaction(tx).
		SetMaxAttempts(3)

	response, err := query.Execute(client)
	require.NoError(t, err)

	assert.Equal(t, uint64(60), response.Total)
	assert.GreaterOrEqual(t, requestCount, 2, "should have retried at least once")
}

func TestUnitFeeEstimateQuerySucceedsOnFirstAttempt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/network/fees", r.URL.Path, "request path should be /api/v1/network/fees")

		contentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/protobuf", contentType, "Content-Type should be application/protobuf")

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Greater(t, len(body), 0, "request body should not be empty")

		queryParams := r.URL.Query()
		assert.Contains(t, queryParams, "mode", "request should contain mode query parameter")
		assert.Equal(t, "INTRINSIC", queryParams.Get("mode"), "mode should be INTRINSIC")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(newSuccessResponse(3, 10, 20)))
	}))
	defer server.Close()

	cleanup := SetupMockTransportForDomain("localhost:8084", server.URL)
	defer cleanup()

	client := newMockClientForREST()

	tx := NewTransferTransaction()

	query := NewFeeEstimateQuery().
		SetTransaction(tx).
		SetMode(FeeEstimateModeIntrinsic)

	response, err := query.Execute(client)
	require.NoError(t, err)

	assert.Equal(t, uint64(60), response.Total)
}
