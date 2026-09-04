package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"sync"
	"time"
)

// How the mirror HTTP layer attaches to Client. Three concerns, kept separate on purpose:
//
//  1. the transport — one per Client, built on first use, long-lived, shared by every mirror
//     REST call. Expensive: it owns the connection pool.
//  2. the policy — a plain value, resolved per call as
//     package default → Client → per-query override. Cheap.
//  3. the base URL — read per call, because SetMirrorNetwork can change it at any time.
//
// Go has no inheritance, so "the query inherits the client's settings" is a struct copy
// followed by overwriting the fields the query actually set. Because the policy is a value,
// two queries can disagree about attempts while sharing one connection pool, and nothing on
// the hot path needs a lock.

// Grace period Client.Close gives in-flight mirror reads before aborting them. Matches the
// mirror node's own per-attempt default rather than the gRPC shutdown budget.
const mirrorHttpDefaultCloseGrace = 5 * time.Second

// mirrorHttpState is the per-Client HTTP state, in one struct so Client gains one field
// instead of four.
//
// Client holds it by pointer, following the same reason _ManagedNetwork holds
// healthyNodesMutex by pointer: Client has public API that copies it by value
// (GetGrpcDeadline, the deprecated GetNetworkName), and copying a mutex is a vet error. By
// pointer, a Client copy shares this state rather than forking it.
type mirrorHttpState struct {
	mu        sync.Mutex
	transport httpTransport
	// nil until a caller overrides the package defaults.
	options *mirrorHttpRetryOptions
	// Fixed when the transport is built, so it is client-level rather than per call. 0 keeps
	// net/http's own bounds.
	connectTimeout time.Duration
}

// state returns the Client's HTTP state, allocating it if this Client did not come from a
// constructor. Only the error paths in ClientFromConfig return a zero-value Client, and one is
// never shared before it is handed back, so allocating here needs no lock of its own.
func (client *Client) mirrorHttpOrInit() *mirrorHttpState {
	if client.mirrorHttp == nil {
		client.mirrorHttp = &mirrorHttpState{}
	}

	return client.mirrorHttp
}

// settingsLocked must be called with mu held.
func (s *mirrorHttpState) settingsLocked() mirrorHttpRetryOptions {
	if s.options == nil {
		return defaultMirrorHttpRetryOptions()
	}

	return *s.options
}

// mirrorHttpSettings returns the Client's HTTP policy as a value the caller may modify freely.
func (client *Client) mirrorHttpSettings() mirrorHttpRetryOptions {
	if client.mirrorHttp == nil {
		return defaultMirrorHttpRetryOptions()
	}

	client.mirrorHttp.mu.Lock()
	defer client.mirrorHttp.mu.Unlock()

	return client.mirrorHttp.settingsLocked()
}

// setMirrorHttpSettings replaces the Client-level HTTP policy. Unexported while this is a
// prototype; the shipping version exposes it, since a caller cannot otherwise raise a timeout.
func (client *Client) setMirrorHttpSettings(options mirrorHttpRetryOptions) {
	state := client.mirrorHttpOrInit()

	state.mu.Lock()
	defer state.mu.Unlock()

	state.options = &options
}

// setMirrorHttpTransport injects a transport, replacing any already built. The Client takes
// over closing it. This is the seam an application uses for a proxy, mTLS or tracing, and the
// seam tests use for a fake.
func (client *Client) setMirrorHttpTransport(transport httpTransport) {
	state := client.mirrorHttpOrInit()

	state.mu.Lock()
	defer state.mu.Unlock()

	state.transport = transport
}

// setMirrorHttpConnectTimeout sets the dial/handshake bound for the transport. It must be
// called before the first mirror REST call, because the transport is built once. 0 restores
// net/http's own bounds.
func (client *Client) setMirrorHttpConnectTimeout(connectTimeout time.Duration) {
	state := client.mirrorHttpOrInit()

	state.mu.Lock()
	defer state.mu.Unlock()

	state.connectTimeout = connectTimeout
}

// sharedMirrorHttpTransport returns the Client's transport, building it on first use.
//
// Lazy on purpose: most Clients never make a mirror REST call, so one built in a test or a
// short-lived process should not open a pool it will not use. It also means a transport
// injected any time before the first call still wins.
func (client *Client) sharedMirrorHttpTransport() httpTransport {
	state := client.mirrorHttpOrInit()

	state.mu.Lock()
	defer state.mu.Unlock()

	if state.transport == nil {
		state.transport = newDefaultHttpTransport(state.connectTimeout)
	}

	return state.transport
}

// mirrorRestClient pairs the Client's shared transport with a per-call policy and the base URL
// as it stands now. Callers must not close the returned client: the transport belongs to the
// Client and is released by Client.Close.
func (client *Client) mirrorRestClient(options mirrorHttpRetryOptions) (*mirrorHttpClient, error) {
	baseURL, err := mirrorNodeRestBaseURL(client)
	if err != nil {
		return nil, err
	}

	return newMirrorHttpClient(baseURL, client.sharedMirrorHttpTransport(), options), nil
}

// mirrorRestClientForBaseURL is mirrorRestClient with the base URL supplied rather than
// resolved. It exists for the call sites that still rewrite the base for a local network — Go
// currently rewrites it to three different ports across four sites, which is tracked as the
// ingress-standardization question and not fixed here.
func (client *Client) mirrorRestClientForBaseURL(baseURL string, options mirrorHttpRetryOptions) *mirrorHttpClient {
	return newMirrorHttpClient(baseURL, client.sharedMirrorHttpTransport(), options)
}

// mirrorHttpOptionsForQuery resolves the three levels a query's HTTP policy can come from.
//
// queryMaxAttempts of 0 means "not set by the query", matching the existing convention on
// MirrorNodeAccountBalanceQuery.SetMaxAttempts.
//
// The Client.GetMaxAttempts fallback is a compatibility bridge, not a design choice:
// Client.SetMaxAttempts is the *gRPC* retry knob, and the shipped balance query already reaches
// for it, so dropping it here would silently change behaviour for anyone relying on it. It is
// the same knob-reuse this layer argues against upstream — kept only until mirror HTTP has its
// own client-level attempts setting.
func (client *Client) mirrorHttpOptionsForQuery(queryMaxAttempts uint64) mirrorHttpRetryOptions {
	options := client.mirrorHttpSettings()

	switch {
	case queryMaxAttempts > 0:
		options.maxAttempts = int(queryMaxAttempts)
	case client.GetMaxAttempts() > 0:
		options.maxAttempts = client.GetMaxAttempts()
	}

	return options
}

// closeMirrorHttp releases the transport if one was ever built, and is a no-op otherwise. It
// clears the field so a Close/reuse cycle does not hand back a dead transport.
func (client *Client) closeMirrorHttp(grace time.Duration) error {
	if client.mirrorHttp == nil {
		return nil
	}

	client.mirrorHttp.mu.Lock()
	transport := client.mirrorHttp.transport
	client.mirrorHttp.transport = nil
	client.mirrorHttp.mu.Unlock()

	if transport == nil {
		return nil
	}

	return transport.close(grace)
}
