//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitClientMirrorHttpTransportIsLazyAndShared(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	assert.Nil(t, client.mirrorHttp.built)

	first := client.sharedMirrorHttpTransport()
	require.NotNil(t, first)
	assert.Same(t, first, client.sharedMirrorHttpTransport())

	client.closeMirrorHttp(time.Second)
}

func TestUnitClientMirrorHttpInjectedTransportWins(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	setMirrorHttpTransport(client, injected)

	assert.Same(t, injected, client.sharedMirrorHttpTransport())
}

func TestUnitClientMirrorHttpPolicyForQueryLayering(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	assert.Equal(t, uint16(mirrorNodeHttpDefaultMaxAttempts), client.mirrorHttpPolicyForQuery(0).GetMaxAttempts(),
		"nothing set anywhere falls back to the package default")

	client.SetMaxAttempts(4)
	assert.Equal(t, uint16(mirrorNodeHttpDefaultMaxAttempts), client.mirrorHttpPolicyForQuery(0).GetMaxAttempts(),
		"Client.SetMaxAttempts is the gRPC knob and no longer reaches mirror REST")

	client.SetMirrorNodeHttpConfig(client.GetMirrorNodeHttpConfig().
		WithRetryPolicy(client.mirrorHttpPolicy().WithMaxAttempts(6).WithPerAttemptTimeout(90 * time.Second)))
	assert.Equal(t, uint16(6), client.mirrorHttpPolicyForQuery(0).GetMaxAttempts())

	policy := client.mirrorHttpPolicyForQuery(2)
	assert.Equal(t, uint16(2), policy.GetMaxAttempts())
	assert.Equal(t, 90*time.Second, policy.GetPerAttemptTimeout())
}

func TestUnitClientMirrorHttpSettingsAreCopies(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	_ = client.mirrorHttpPolicy().WithMaxAttempts(99)
	codes := client.mirrorHttpPolicy().GetRetryableStatusCodes()
	codes[0] = http.StatusTeapot

	assert.Equal(t, uint16(mirrorNodeHttpDefaultMaxAttempts), client.mirrorHttpPolicy().GetMaxAttempts(),
		"a caller holding the value must not be able to reach into the Client")
	assert.NotContains(t, client.mirrorHttpPolicy().GetRetryableStatusCodes(), uint16(http.StatusTeapot))
}

func TestUnitClientMirrorRestClientResolvesBaseURL(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{"mirror.example.com:443"})

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	setMirrorHttpTransport(client, injected)

	restClient, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), client.mirrorHttpPolicyForQuery(0))
	require.NoError(t, err)
	assert.Equal(t, "https://mirror.example.com:443/api/v1", restClient.baseURL)

	_, err = restClient.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.NoError(t, err)
	assert.Equal(t, "https://mirror.example.com:443/api/v1/accounts/0.0.1", injected.lastRequest().url)
}

func TestUnitClientMirrorRestClientRequiresMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{})

	_, err = client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), DefaultMirrorNodeHttpRetryPolicy())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node is not set")
}

func TestUnitClientMirrorRestClientSharesTransportAcrossPolicies(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	setMirrorHttpTransport(client, injected)

	patient, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), client.mirrorHttpPolicyForQuery(10))
	require.NoError(t, err)
	impatient, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), client.mirrorHttpPolicyForQuery(1))
	require.NoError(t, err)

	assert.Equal(t, uint16(10), patient.policy.GetMaxAttempts())
	assert.Equal(t, uint16(1), impatient.policy.GetMaxAttempts())
	assert.Same(t, patient.transport, impatient.transport)
}

func TestUnitClientCloseReleasesBuiltMirrorHttpTransport(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	built, ok := client.sharedMirrorHttpTransport().(*DefaultHttpTransport)
	require.True(t, ok)

	require.NoError(t, client.Close())
	assert.True(t, built.isClosed())
}

func TestUnitClientCloseLeavesInjectedMirrorHttpTransportOpen(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	setMirrorHttpTransport(client, injected)

	require.NoError(t, client.Close())
	assert.Zero(t, injected.closeCalls)
}

func TestUnitClientCloseInterruptsMirrorHttpBackoff(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	// Retry-After gives a fixed wait; a jittered backoff could end before Close.
	throttled := HttpResponse{statusCode: http.StatusServiceUnavailable, headers: map[string][]string{httpHeaderRetryAfter: {"5"}}}
	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: throttled}}}
	setMirrorHttpTransport(client, injected)

	policy := fastRetryPolicy(5)
	restClient, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), policy)
	require.NoError(t, err)

	time.AfterFunc(100*time.Millisecond, func() { client.closeMirrorHttp(time.Second) })

	start := time.Now()
	_, err = restClient.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())

	requireHttpErrorID(t, err, HttpTransportClientClosedError)
	assert.Less(t, time.Since(start), time.Second)
	assert.Equal(t, 1, injected.callCount())
	assert.Zero(t, injected.closeCalls)
}

func TestUnitClientCloseRefusesNewMirrorHttpCalls(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	client.closeMirrorHttp(time.Second)
	client.closeMirrorHttp(time.Second)

	restClient, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), DefaultMirrorNodeHttpRetryPolicy())
	require.NoError(t, err)

	_, err = restClient.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	requireHttpErrorID(t, err, HttpTransportClientClosedError)
	assert.Nil(t, client.mirrorHttp.built)
}

func TestUnitClientMirrorNodeHttpConfigGetReturnsSetConfig(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	require.NotNil(t, client.sharedMirrorHttpTransport())
	assert.Nil(t, client.GetMirrorNodeHttpConfig().GetTransport())

	client.SetMirrorNodeHttpConfig(client.GetMirrorNodeHttpConfig())
	assert.Equal(t, DefaultMirrorNodeHttpConfig(), client.GetMirrorNodeHttpConfig())

	client.closeMirrorHttp(time.Second)
}

func TestUnitClientMirrorNodeHttpConfigSetReplaces(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithRequestHeader("authorization", "Bearer a"))
	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithRetryPolicy(DefaultMirrorNodeHttpRetryPolicy().WithMaxAttempts(2)))

	assert.Empty(t, client.GetMirrorNodeHttpConfig().GetRequestHeaders())
	assert.Equal(t, uint16(2), client.GetMirrorNodeHttpConfig().GetRetryPolicy().GetMaxAttempts())
}

func TestUnitClientMirrorNodeHttpConfigSendsRequestHeaders(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().
		WithTransport(injected).
		WithRequestHeader("Authorization", "Bearer test"))

	restClient, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), client.mirrorHttpPolicy())
	require.NoError(t, err)
	_, err = restClient.get(testMirrorPath(t, testMirrorAccountPath), CancellationNone())
	require.NoError(t, err)

	assert.Equal(t, "Bearer test", injected.lastRequest().headers["authorization"])
}

func TestUnitClientMirrorHttpInjectingAfterBuildDoesNotLeak(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	built, ok := client.sharedMirrorHttpTransport().(*DefaultHttpTransport)
	require.True(t, ok)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	setMirrorHttpTransport(client, injected)
	assert.Same(t, injected, client.sharedMirrorHttpTransport())

	client.closeMirrorHttp(time.Second)
	assert.True(t, built.isClosed())
	assert.Zero(t, injected.closeCalls)
}

func TestUnitClientCloseWithoutMirrorHttpIsNoOp(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	require.NoError(t, client.Close())
}

func TestUnitClientMirrorHttpNilClientReturnsError(t *testing.T) {
	t.Parallel()

	var client *Client
	assert.Equal(t, DefaultMirrorNodeHttpConfig(), client.GetMirrorNodeHttpConfig())

	evmAddress := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	calls := map[string]func() error{
		"MirrorNodeContractCallQuery": func() error {
			_, err := NewMirrorNodeContractCallQuery().SetContractID(ContractID{Contract: 5}).Execute(client)
			return err
		},
		"MirrorNodeContractEstimateGasQuery": func() error {
			_, err := NewMirrorNodeContractEstimateGasQuery().SetContractID(ContractID{Contract: 5}).Execute(client)
			return err
		},
		"MirrorNodeAccountBalanceQuery": func() error {
			_, err := NewMirrorNodeAccountBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
			return err
		},
		"RegisteredNodeAddressBookQuery": func() error {
			_, err := NewRegisteredNodeAddressBookQuery().Execute(client)
			return err
		},
		"FeeEstimateQuery": func() error {
			_, err := NewFeeEstimateQuery().SetTransaction(NewTransferTransaction()).Execute(client)
			return err
		},
		"AccountID.PopulateAccount": func() error {
			return (&AccountID{AliasEvmAddress: &evmAddress}).PopulateAccount(client)
		},
		"AccountID.PopulateEvmAddress": func() error {
			return (&AccountID{Account: 5}).PopulateEvmAddress(client)
		},
		"ContractID.PopulateContract": func() error {
			return (&ContractID{EvmAddress: evmAddress}).PopulateContract(client)
		},
	}

	for name, call := range calls {
		require.ErrorIs(t, call(), errNoClientProvided, name)
	}
}

// The constructors' error paths return &Client{}, so the accessors must not panic on one.
func TestUnitClientMirrorHttpZeroValueClientIsSafe(t *testing.T) {
	t.Parallel()

	client := &Client{}

	assert.Equal(t, DefaultMirrorNodeHttpRetryPolicy(), client.mirrorHttpPolicy())
	client.closeMirrorHttp(time.Second)

	_, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), DefaultMirrorNodeHttpRetryPolicy())
	require.Error(t, err)
}

func TestUnitClientMirrorHttpConnectTimeoutIsOptIn(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	built, ok := client.sharedMirrorHttpTransport().(*DefaultHttpTransport)
	require.True(t, ok)
	inner, ok := built.httpClient.Transport.(*http.Transport)
	require.True(t, ok)
	assert.Equal(t, httpDefaultHandshakeTimeout, inner.TLSHandshakeTimeout)

	client.closeMirrorHttp(time.Second)
}

func TestUnitClientMirrorHttpBuildsTransportFromConfiguration(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithTransportConfiguration(DefaultHttpTransportConfiguration().
		WithConnectTimeout(2 * time.Second).
		WithMaxRedirects(1).
		WithMaxResponseBytes(1024)))

	built, ok := client.sharedMirrorHttpTransport().(*DefaultHttpTransport)
	require.True(t, ok)
	inner, ok := built.httpClient.Transport.(*http.Transport)
	require.True(t, ok)
	assert.Equal(t, 2*time.Second, inner.TLSHandshakeTimeout)
	assert.Equal(t, 1, built.maxRedirects)
	assert.Equal(t, int64(1024), built.maxResponseBytes)

	client.closeMirrorHttp(time.Second)
}

func TestUnitClientMirrorHttpTotalDeadlineInheritsRequestTimeout(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	restClient, err := client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), DefaultMirrorNodeHttpRetryPolicy())
	require.NoError(t, err)
	assert.Equal(t, client.GetRequestTimeout(), restClient.policy.GetTotalDeadline())

	client.SetRequestTimeout(30 * time.Second)
	restClient, err = client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), DefaultMirrorNodeHttpRetryPolicy())
	require.NoError(t, err)
	assert.Equal(t, 30*time.Second, restClient.policy.GetTotalDeadline())

	explicit := DefaultMirrorNodeHttpRetryPolicy()
	explicit = explicit.WithTotalDeadline(5 * time.Second)
	restClient, err = client.mirrorRestClient(testMirrorPath(t, testMirrorAccountPath), explicit)
	require.NoError(t, err)
	assert.Equal(t, 5*time.Second, restClient.policy.GetTotalDeadline())
}

func TestUnitClientMirrorRestClientResolvesBaseURLByEndpointFamily(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mirror string
		path   string
		want   string
	}{
		{name: "local REST stays on the REST port", mirror: "localhost:5600", path: "/balances", want: "http://localhost:38081/api/v1"},
		{name: "local registered nodes go to 8084", mirror: "localhost:5600", path: "/network/registered-nodes?limit=1", want: "http://localhost:8084/api/v1"},
		{name: "local loopback IP goes to 8084", mirror: "127.0.0.1:5600", path: "/network/registered-nodes", want: "http://localhost:8084/api/v1"},
		{name: "local fees go to 8084", mirror: "localhost:5600", path: "/network/fees?mode=STATE", want: "http://localhost:8084/api/v1"},
		{name: "local contract calls go to 8545", mirror: "127.0.0.1:5600", path: "/contracts/call", want: "http://localhost:8545/api/v1"},
		{name: "remote is left alone", mirror: "mirror.example.com:443", path: "/contracts/call", want: "https://mirror.example.com:443/api/v1"},
		{name: "a host merely containing localhost is remote", mirror: "mylocalhost.example.com:443", path: "/network/registered-nodes", want: "https://mylocalhost.example.com:443/api/v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := _NewMockClient()
			require.NoError(t, err)
			client.SetMirrorNetwork([]string{tt.mirror})

			restClient, err := client.mirrorRestClient(testMirrorPath(t, tt.path), DefaultMirrorNodeHttpRetryPolicy())
			require.NoError(t, err)
			assert.Equal(t, tt.want, restClient.baseURL)
		})
	}
}

func TestUnitIsLoopbackMirrorURL(t *testing.T) {
	t.Parallel()

	for _, local := range []string{"http://localhost:38081/api/v1", "http://LOCALHOST:1", "http://127.0.0.1:1", "http://[::1]:1", "http://0.0.0.0:1"} {
		assert.True(t, isLoopbackMirrorURL(local), local)
	}
	for _, remote := range []string{"https://mylocalhost.example.com/api/v1", "https://localhost.example.com", "http://127.0.0.2:1", "://nope"} {
		assert.False(t, isLoopbackMirrorURL(remote), remote)
	}
}

func TestUnitClientMirrorHttpLoopbackDefaults(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{"localhost:5600"})
	path := testMirrorPath(t, testMirrorAccountPath)

	restClient, err := client.mirrorRestClient(path, client.mirrorHttpPolicyForQuery(0))
	require.NoError(t, err)
	assert.Equal(t, uint16(15), restClient.policy.GetMaxAttempts())
	assert.Equal(t, 90*time.Second, restClient.policy.GetTotalDeadline())
	assert.Equal(t, 30*time.Second, restClient.policy.GetPerAttemptTimeout())

	restClient, err = client.mirrorRestClient(path, client.mirrorHttpPolicyForQuery(2))
	require.NoError(t, err)
	assert.Equal(t, uint16(2), restClient.policy.GetMaxAttempts())
	assert.Equal(t, 90*time.Second, restClient.policy.GetTotalDeadline())

	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithRetryPolicy(DefaultMirrorNodeHttpRetryPolicy().WithTotalDeadline(0)))
	restClient, err = client.mirrorRestClient(path, client.mirrorHttpPolicyForQuery(0))
	require.NoError(t, err)
	assert.Equal(t, client.GetRequestTimeout(), restClient.policy.GetTotalDeadline())

	remote, err := _NewMockClient()
	require.NoError(t, err)
	restClient, err = remote.mirrorRestClient(path, remote.mirrorHttpPolicyForQuery(0))
	require.NoError(t, err)
	assert.Equal(t, uint16(5), restClient.policy.GetMaxAttempts())
}

// applicationTransport is written against the exported API alone, the way an application wraps its
// own HTTP stack.
type applicationTransport struct {
	mu     sync.Mutex
	urls   []string
	closed bool
}

func (a *applicationTransport) RoundTrip(req HttpRequest, _ Cancellation) (HttpResponse, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.urls = append(a.urls, req.GetURL())

	return NewHttpResponse(http.StatusOK, []byte(`{"balances":[{"account":"0.0.5","balance":7}],"links":{"next":null}}`), nil), nil
}

func (a *applicationTransport) Close(time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.closed = true
}

func TestUnitClientRunsQueriesOverInjectedTransport(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	app := &applicationTransport{}
	client.SetMirrorNodeHttpConfig(client.GetMirrorNodeHttpConfig().WithTransport(app))

	balance, err := NewMirrorNodeAccountBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
	require.NoError(t, err)
	assert.Equal(t, HbarFromTinybar(7), balance.Hbars)

	require.NoError(t, client.Close())
	app.mu.Lock()
	defer app.mu.Unlock()
	require.Len(t, app.urls, 1)
	assert.Contains(t, app.urls[0], "/api/v1/balances?account.id=0.0.5")
	assert.False(t, app.closed)
}
