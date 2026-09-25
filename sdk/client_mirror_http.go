package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"net/url"
	"strings"
	"sync"
	"time"
)

// mirrorHttpDefaultCloseGrace is how long Client.Close waits for in-flight mirror node requests.
const mirrorHttpDefaultCloseGrace = 5 * time.Second

// mirrorHttpState is the Client's mirror HTTP state. It is held by pointer so copies of a Client share it.
type mirrorHttpState struct {
	mu     sync.Mutex
	config MirrorNodeHttpConfig
	// built is the transport the Client created, and the only one it closes.
	built HttpTransport
	// closed is closed by Client.Close to stop retries in progress.
	closed chan struct{}
}

func newMirrorHttpState() *mirrorHttpState {
	return &mirrorHttpState{closed: make(chan struct{})}
}

// isClosedLocked must be called with mu held.
func (s *mirrorHttpState) isClosedLocked() bool {
	select {
	case <-s.closed:
		return true
	default:
		return false
	}
}

// mirrorHttpOrInit returns the mirror HTTP state, allocating it for a Client not built by a constructor.
func (client *Client) mirrorHttpOrInit() *mirrorHttpState {
	if client.mirrorHttp == nil {
		client.mirrorHttp = newMirrorHttpState()
	}

	return client.mirrorHttp
}

// SetMirrorNodeHttpConfig replaces the configuration for mirror node REST requests.
// A transport configuration only takes effect if set before the first mirror node request.
func (client *Client) SetMirrorNodeHttpConfig(config MirrorNodeHttpConfig) *Client {
	state := client.mirrorHttpOrInit()

	state.mu.Lock()
	defer state.mu.Unlock()

	state.config = config

	return client
}

// GetMirrorNodeHttpConfig returns the configuration set with SetMirrorNodeHttpConfig.
func (client *Client) GetMirrorNodeHttpConfig() MirrorNodeHttpConfig {
	if client == nil || client.mirrorHttp == nil {
		return DefaultMirrorNodeHttpConfig()
	}

	client.mirrorHttp.mu.Lock()
	defer client.mirrorHttp.mu.Unlock()

	return client.mirrorHttp.config
}

func (client *Client) mirrorHttpPolicy() MirrorNodeHttpRetryPolicy {
	return client.GetMirrorNodeHttpConfig().GetRetryPolicy()
}

// sharedMirrorHttpTransport returns the injected transport, or the Client's own, building it on first use.
func (client *Client) sharedMirrorHttpTransport() HttpTransport {
	state := client.mirrorHttpOrInit()

	state.mu.Lock()
	defer state.mu.Unlock()

	if state.config.transport != nil {
		return state.config.transport
	}
	if state.built == nil && !state.isClosedLocked() {
		state.built = NewDefaultHttpTransport(state.config.transportConfiguration)
	}

	return state.built
}

// mirrorRestClient returns a client for one call. It picks one mirror node, so every request of
// the call, pages included, goes to that node.
func (client *Client) mirrorRestClient(path mirrorNodeRestPath, policy MirrorNodeHttpRetryPolicy) (*mirrorNodeHttpClient, error) {
	baseURL, err := mirrorNodeRestBaseURL(client)
	if err != nil {
		return nil, err
	}
	if isLoopbackMirrorURL(baseURL) {
		policy = policy.withLoopbackDefaults()
	}

	restClient := newMirrorNodeHttpClient(mirrorNodeBaseURLForPath(baseURL, path), client.sharedMirrorHttpTransport(), client.inheritTotalDeadline(policy))
	restClient.headers = client.GetMirrorNodeHttpConfig().requestHeaders
	restClient.clientClosed = client.mirrorHttpOrInit().closed

	return restClient, nil
}

// Local mirror nodes serve these endpoints on their own ports.
var mirrorNodeLocalModules = []struct {
	pathPrefix string
	baseURL    string
}{
	{pathPrefix: "/network/registered-nodes", baseURL: "http://localhost:8084" + mirrorHttpAPIVersionPrefix},
	{pathPrefix: "/network/fees", baseURL: "http://localhost:8084" + mirrorHttpAPIVersionPrefix},
	{pathPrefix: "/contracts/call", baseURL: "http://localhost:8545" + mirrorHttpAPIVersionPrefix},
}

// mirrorNodeBaseURLForPath returns the local base URL that serves path, or baseURL if it is not loopback.
func mirrorNodeBaseURLForPath(baseURL string, path mirrorNodeRestPath) string {
	if !isLoopbackMirrorURL(baseURL) {
		return baseURL
	}

	for _, module := range mirrorNodeLocalModules {
		if strings.HasPrefix(path.String(), module.pathPrefix) {
			return module.baseURL
		}
	}

	return baseURL
}

// isLoopbackMirrorURL reports whether baseURL's host is a loopback address.
func isLoopbackMirrorURL(baseURL string) bool {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	switch strings.ToLower(parsed.Hostname()) {
	case localhostName, "127.0.0.1", "::1", "0.0.0.0":
		return true
	default:
		return false
	}
}

// inheritTotalDeadline sets an unset total deadline to the Client's request timeout.
func (client *Client) inheritTotalDeadline(policy MirrorNodeHttpRetryPolicy) MirrorNodeHttpRetryPolicy {
	if policy.GetTotalDeadline() == 0 && client.requestTimeout > 0 {
		return policy.WithTotalDeadline(client.requestTimeout)
	}

	return policy
}

// mirrorHttpPolicyForQuery returns the Client's policy with the query's max attempts applied, if set.
func (client *Client) mirrorHttpPolicyForQuery(queryMaxAttempts uint64) MirrorNodeHttpRetryPolicy {
	policy := client.mirrorHttpPolicy()
	if queryMaxAttempts > 0 {
		policy = policy.WithMaxAttempts(uint16(min(queryMaxAttempts, math.MaxUint16)))
	}

	return policy
}

// closeMirrorHttp stops further mirror node requests and closes the transport the Client built.
func (client *Client) closeMirrorHttp(grace time.Duration) {
	if client.mirrorHttp == nil {
		return
	}

	state := client.mirrorHttp
	state.mu.Lock()
	if state.isClosedLocked() {
		state.mu.Unlock()
		return
	}
	close(state.closed)
	built := state.built
	state.mu.Unlock()

	if built != nil {
		built.Close(grace)
	}
}
