//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitClientMirrorHttpTransportIsLazyAndShared(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	assert.Nil(t, client.mirrorHttp.transport, "a Client that never made a mirror REST call owns no pool")

	first := client.sharedMirrorHttpTransport()
	require.NotNil(t, first)
	assert.Same(t, first, client.sharedMirrorHttpTransport(), "every mirror call shares one transport")

	require.NoError(t, client.closeMirrorHttp(time.Second))
}

func TestUnitClientMirrorHttpInjectedTransportWins(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	client.setMirrorHttpTransport(injected)

	assert.Same(t, injected, client.sharedMirrorHttpTransport(), "injection is the seam for proxy, mTLS, tracing and tests")
}

// The "inheritance" chain, which is a struct copy plus an overwrite rather than embedding.
func TestUnitClientMirrorHttpOptionsForQueryLayering(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	require.Equal(t, -1, client.GetMaxAttempts(), "an unset client reports -1, not a usable count")

	assert.Equal(t, mirrorNodeDefaultMaxAttempts, client.mirrorHttpOptionsForQuery(0).maxAttempts,
		"nothing set anywhere falls back to the package default")

	client.SetMaxAttempts(4)
	assert.Equal(t, 4, client.mirrorHttpOptionsForQuery(0).maxAttempts, "the client fills in for a query that said nothing")
	assert.Equal(t, 2, client.mirrorHttpOptionsForQuery(2).maxAttempts, "the query wins over the client")

	settings := client.mirrorHttpSettings()
	settings.perAttemptTimeout = 90 * time.Second
	client.setMirrorHttpSettings(settings)
	assert.Equal(t, 90*time.Second, client.mirrorHttpOptionsForQuery(2).perAttemptTimeout,
		"client-level policy carries into a query that only overrode attempts")
}

func TestUnitClientMirrorHttpSettingsAreCopies(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	mutated := client.mirrorHttpSettings()
	mutated.maxAttempts = 99
	mutated.retryableStatusCodes = append(mutated.retryableStatusCodes, http.StatusTeapot)

	assert.NotEqual(t, 99, client.mirrorHttpSettings().maxAttempts,
		"a caller holding the value must not be able to reach into the Client")
}

func TestUnitClientMirrorRestClientResolvesBaseURL(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{"mirror.example.com:443"})

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	client.setMirrorHttpTransport(injected)

	restClient, err := client.mirrorRestClient(client.mirrorHttpOptionsForQuery(0))
	require.NoError(t, err)
	assert.Equal(t, "https://mirror.example.com:443/api/v1", restClient.baseURL)

	_, err = restClient.get(context.Background(), testMirrorPath(t, testMirrorAccountPath))
	require.NoError(t, err)
	assert.Equal(t, "https://mirror.example.com:443/api/v1/accounts/0.0.1", injected.lastRequest().url)
}

func TestUnitClientMirrorRestClientRequiresMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{})

	_, err = client.mirrorRestClient(defaultMirrorHttpRetryOptions())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node is not set")
}

// Two queries can disagree about policy while sharing one connection pool — the reason the
// transport is per-Client and the options are per-call.
func TestUnitClientMirrorRestClientSharesTransportAcrossPolicies(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	client.setMirrorHttpTransport(injected)

	patient, err := client.mirrorRestClient(client.mirrorHttpOptionsForQuery(10))
	require.NoError(t, err)
	impatient, err := client.mirrorRestClient(client.mirrorHttpOptionsForQuery(1))
	require.NoError(t, err)

	assert.Equal(t, 10, patient.options.maxAttempts)
	assert.Equal(t, 1, impatient.options.maxAttempts)
	assert.Same(t, patient.transport, impatient.transport)
}

func TestUnitClientCloseReleasesMirrorHttpTransport(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	injected := &fakeHttpTransport{turns: []fakeHttpTurn{{resp: statusResponse(http.StatusOK)}}}
	client.setMirrorHttpTransport(injected)

	require.NoError(t, client.Close())
	assert.Equal(t, 1, injected.closeCalls, "Client.Close owns the transport's lifetime")
	assert.Nil(t, client.mirrorHttp.transport, "a closed Client must not hand back a dead transport")
}

func TestUnitClientCloseWithoutMirrorHttpIsNoOp(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	require.NoError(t, client.Close(), "closing a Client that never made a mirror REST call must not fail")
}

// The constructors' error paths return &Client{}, so the accessors must not panic on one.
func TestUnitClientMirrorHttpZeroValueClientIsSafe(t *testing.T) {
	t.Parallel()

	client := &Client{}

	assert.Equal(t, defaultMirrorHttpRetryOptions().maxAttempts, client.mirrorHttpSettings().maxAttempts)
	require.NoError(t, client.closeMirrorHttp(time.Second), "nothing was built, so there is nothing to release")

	_, err := client.mirrorRestClient(defaultMirrorHttpRetryOptions())
	require.Error(t, err, "a zero-value Client has no mirror network to resolve")
}

// The connect timeout is a new capability, not a new default: unset, the transport keeps
// net/http's own dial and handshake bounds, so attaching this layer to an existing call site
// cannot tighten a limit the caller was already living with.
func TestUnitClientMirrorHttpConnectTimeoutIsOptIn(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	assert.Zero(t, client.mirrorHttpOrInit().connectTimeout, "unset means net/http's own bounds")

	client.setMirrorHttpConnectTimeout(2 * time.Second)
	assert.Equal(t, 2*time.Second, client.mirrorHttpOrInit().connectTimeout)

	require.NoError(t, client.closeMirrorHttp(time.Second))
}

// Nothing bounded total wall clock before this layer, so the bound ships off by default.
func TestUnitMirrorHttpTotalDeadlineIsOptIn(t *testing.T) {
	t.Parallel()

	assert.Zero(t, defaultMirrorHttpRetryOptions().totalDeadline,
		"turning this on is a behaviour change that belongs to the cross-SDK policy decision")
	assert.Equal(t, mirrorNodeDefaultMaxAttempts, defaultMirrorHttpRetryOptions().maxAttempts,
		"attempts match the shipped mirror node default")
	assert.Equal(t, mirrorNodeDefaultTimeout, defaultMirrorHttpRetryOptions().perAttemptTimeout,
		"the per-attempt bound matches the shipped mirror node default")
}
