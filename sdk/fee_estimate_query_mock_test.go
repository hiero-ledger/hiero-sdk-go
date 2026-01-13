//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubResponse struct {
	status int
	body   string
}

type stubMirrorRestServer struct {
	responses        []stubResponse
	observedRequests int
	server           *http.Server
	listener         net.Listener
	baseURL          string
}

// newStubMirrorRestServer creates a new mock HTTP server.
// Note: These tests cannot run in parallel because the mirror node code hardcodes port 5551 for localhost.
func newStubMirrorRestServer(t *testing.T) *stubMirrorRestServer {
	stub := &stubMirrorRestServer{
		responses: make([]stubResponse, 0),
	}

	listener, err := net.Listen("tcp", "127.0.0.1:5551")
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stub.observedRequests++

		assert.Equal(t, "/api/v1/network/fees", r.URL.Path, "request path should be /api/v1/network/fees")

		contentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/json", contentType, "Content-Type should be application/json")

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Greater(t, len(body), 0, "request body should not be empty")

		var payload map[string]interface{}
		err = json.Unmarshal(body, &payload)
		require.NoError(t, err)
		assert.Contains(t, payload, "transaction", "request should contain transaction field")
		assert.Contains(t, payload, "mode", "request should contain mode field")

		if len(stub.responses) == 0 {
			http.Error(w, "no response queued", http.StatusInternalServerError)
			return
		}

		response := stub.responses[0]
		stub.responses = stub.responses[1:]

		w.WriteHeader(response.status)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response.body))
	})

	stub.server = &http.Server{Handler: handler}
	stub.listener = listener
	stub.baseURL = fmt.Sprintf("http://%s", listener.Addr().String())

	go func() {
		_ = stub.server.Serve(listener)
	}()

	return stub
}

func (s *stubMirrorRestServer) enqueue(response stubResponse) {
	s.responses = append(s.responses, response)
}

func (s *stubMirrorRestServer) stop() {
	if s.server != nil {
		_ = s.server.Close()
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

func (s *stubMirrorRestServer) requestCount() int {
	return s.observedRequests
}

func (s *stubMirrorRestServer) getURL() string {
	return s.baseURL
}

func (s *stubMirrorRestServer) getMirrorNetworkAddress() string {
	parsedURL, err := url.Parse(s.baseURL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test server URL: %v", err))
	}
	return parsedURL.Host
}

func (s *stubMirrorRestServer) verify(t *testing.T) {
	assert.Empty(t, s.responses, "all queued responses should have been served")
}

func newMockClientForREST(mirrorNetworkAddress string) *Client {
	net := _NewNetwork()
	client := _NewClient(net, []string{mirrorNetworkAddress}, nil, false, 0, 0)

	err := client.SetNetwork(map[string]AccountID{
		"127.0.0.1:50211": {Account: 3},
	})
	if err != nil {
		// Continue even if SetNetwork fails
	}

	// Set a dummy operator so transactions can be frozen (required by FreezeWith)
	dummyKey, _ := PrivateKeyGenerateEd25519()
	client.SetOperator(AccountID{Account: 1}, dummyKey)

	return client
}

func newSuccessResponse(mode FeeEstimateMode, networkMultiplier int, nodeBase, serviceBase uint64) string {
	networkSubtotal := nodeBase * uint64(networkMultiplier)
	total := networkSubtotal + nodeBase + serviceBase
	return fmt.Sprintf(`{
  "mode": "%s",
  "network": {"multiplier": %d, "subtotal": %d},
  "node": {"base": %d, "extras": []},
  "service": {"base": %d, "extras": []},
  "notes": [],
  "total": %d
}`, mode.String(), networkMultiplier, networkSubtotal, nodeBase, serviceBase, total)
}
