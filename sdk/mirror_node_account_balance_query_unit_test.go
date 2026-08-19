//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// balancesHandler serves a /balances response, recording the account.id it was asked for.
func balancesHandler(t *testing.T, gotAccountID *string, body string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/balances", r.URL.Path)
		*gotAccountID = r.URL.Query().Get("account.id")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}
}

func TestUnitMirrorNodeAccountBalanceQueryGetSet(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Shard: 1, Realm: 2, Account: 3}
	query := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(accountID).
		SetMaxAttempts(7)

	assert.Equal(t, accountID, query.GetAccountID())
	assert.Equal(t, uint64(7), query.GetMaxAttempts())
}

func TestUnitMirrorNodeAccountBalanceQueryDefaults(t *testing.T) {
	t.Parallel()

	query := NewMirrorNodeAccountBalanceQuery()

	assert.Equal(t, AccountID{}, query.GetAccountID())
	// Zero means unset; resolveAttempts supplies the effective value.
	assert.Zero(t, query.GetMaxAttempts())
}

func TestUnitMirrorNodeAccountBalanceQueryReturnsHbarBalance(t *testing.T) {
	var gotAccountID string
	client := newMockMirrorClient(t, "balance.example.com:443", balancesHandler(t, &gotAccountID,
		`{"timestamp":"1234567890.000000000","balances":[{"account":"0.0.12345","balance":123456789}],"links":{"next":null}}`))

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 12345}).
		Execute(client)

	require.NoError(t, err)
	assert.Equal(t, HbarFromTinybar(123456789), balance.Hbars)
	assert.Equal(t, "0.0.12345", gotAccountID)
}

// Tinybars must survive values a float64 would round.
func TestUnitMirrorNodeAccountBalanceQueryLargeBalanceIsLossless(t *testing.T) {
	var gotAccountID string
	// 2^53 + 1, the first integer float64 cannot represent exactly.
	const tinybars int64 = 9007199254740993
	client := newMockMirrorClient(t, "large.example.com:443", balancesHandler(t, &gotAccountID,
		fmt.Sprintf(`{"balances":[{"account":"0.0.7","balance":%d}]}`, tinybars)))

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 7}).
		Execute(client)

	require.NoError(t, err)
	assert.Equal(t, tinybars, balance.Hbars.AsTinybar())
}

func TestUnitMirrorNodeAccountBalanceQueryResolvesEvmAddress(t *testing.T) {
	var gotAccountID string
	client := newMockMirrorClient(t, "evm.example.com:443", balancesHandler(t, &gotAccountID,
		`{"balances":[{"account":"0.0.1001","balance":500}]}`))

	accountID, err := AccountIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(accountID).
		Execute(client)

	require.NoError(t, err)
	assert.Equal(t, HbarFromTinybar(500), balance.Hbars)
	assert.Equal(t, "742d35cc6634c0532925a3b844bc454e4438f44e", gotAccountID)
}

func TestUnitMirrorNodeAccountBalanceQueryResolvesPublicKeyAlias(t *testing.T) {
	var gotAccountID string
	client := newMockMirrorClient(t, "alias.example.com:443", balancesHandler(t, &gotAccountID,
		`{"balances":[{"account":"0.0.2002","balance":42}]}`))

	key, err := PrivateKeyFromStringEd25519(mockPrivateKey)
	require.NoError(t, err)
	publicKey := key.PublicKey()

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{AliasKey: &publicKey}).
		Execute(client)

	require.NoError(t, err)
	assert.Equal(t, HbarFromTinybar(42), balance.Hbars)
	assert.NotEmpty(t, gotAccountID)
	// Base32, never the DER hex AccountID.String() emits.
	assert.NotContains(t, gotAccountID, "302a300506032b6570")
	assert.Equal(t, gotAccountID, url.QueryEscape(gotAccountID), "alias must need no escaping")
}

func TestUnitMirrorNodeAccountBalanceQueryNonExistentAccountErrors(t *testing.T) {
	var gotAccountID string
	client := newMockMirrorClient(t, "missing.example.com:443", balancesHandler(t, &gotAccountID,
		`{"timestamp":null,"balances":[],"links":{"next":null}}`))

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 999999999}).
		Execute(client)

	var status ErrHederaPreCheckStatus
	require.ErrorAs(t, err, &status)
	assert.Equal(t, StatusInvalidAccountID, status.Status)
	assert.Contains(t, err.Error(), "INVALID_ACCOUNT_ID")
}

func TestUnitMirrorNodeAccountBalanceQueryZeroBalanceIsNotMissing(t *testing.T) {
	var gotAccountID string
	client := newMockMirrorClient(t, "zerobalance.example.com:443", balancesHandler(t, &gotAccountID,
		`{"balances":[{"account":"0.0.42","balance":0}]}`))

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 42}).
		Execute(client)

	require.NoError(t, err)
	assert.Zero(t, balance.Hbars.AsTinybar())
}

func TestUnitMirrorNodeAccountBalanceQueryMissingBalancesArrayIsMalformed(t *testing.T) {
	var gotAccountID string
	client := newMockMirrorClient(t, "nobalances.example.com:443", balancesHandler(t, &gotAccountID,
		`{"timestamp":null,"links":{"next":null}}`))

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 7}).
		Execute(client)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no balances array")
	var status ErrHederaPreCheckStatus
	assert.False(t, errors.As(err, &status), "a malformed payload must not read as a missing account")
}

func TestUnitMirrorNodeAccountBalanceQueryNoAccountIDErrorsBeforeRequest(t *testing.T) {
	var called atomic.Bool
	client := newMockMirrorClient(t, "unused.example.com:443", func(w http.ResponseWriter, r *http.Request) {
		called.Store(true)
		w.WriteHeader(http.StatusOK)
	})

	_, err := NewMirrorNodeAccountBalanceQuery().Execute(client)

	require.ErrorIs(t, err, errMirrorNodeAccountBalanceQueryNoAccountID)
	assert.False(t, called.Load(), "no request may be sent when the account ID is unset")
}

func TestUnitMirrorNodeAccountBalanceQueryNilClientErrors(t *testing.T) {
	t.Parallel()

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 3}).
		Execute(nil)

	require.ErrorIs(t, err, errNoClientProvided)
}

func TestUnitMirrorNodeAccountBalanceQueryRetriesTransientError(t *testing.T) {
	var attempts atomic.Int32
	client := newMockMirrorClient(t, "retry.example.com:443", func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"balances":[{"account":"0.0.5","balance":900}]}`))
	})

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 5}).
		Execute(client)

	require.NoError(t, err)
	assert.Equal(t, HbarFromTinybar(900), balance.Hbars)
	assert.Equal(t, int32(2), attempts.Load(), "the 503 should be retried exactly once")
}

// A 4xx is the intended answer and must not burn the retry budget.
func TestUnitMirrorNodeAccountBalanceQueryDoesNotRetryClientError(t *testing.T) {
	var attempts atomic.Int32
	client := newMockMirrorClient(t, "badrequest.example.com:443", func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	})

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 5}).
		Execute(client)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Equal(t, int32(1), attempts.Load(), "a 400 must not be retried")
}

func TestUnitMirrorNodeAccountBalanceQueryMalformedJSONErrors(t *testing.T) {
	client := newMockMirrorClient(t, "garbage.example.com:443", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"balances":`))
	})

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 5}).
		Execute(client)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
}

func TestUnitMirrorNodeAccountBalanceQueryNoMirrorNetworkErrors(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{})

	_, err = NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 3}).
		Execute(client)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node is not set")
}

// Checksum validation matches the consensus-node queries' validateNetworkOnIDs shape.
func TestUnitMirrorNodeAccountBalanceQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	query := NewMirrorNodeAccountBalanceQuery().SetAccountID(accountID)

	require.NoError(t, query.validateNetworkOnIDs(client))
}

func TestUnitMirrorNodeAccountBalanceQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	query := NewMirrorNodeAccountBalanceQuery().SetAccountID(accountID)

	err = query.validateNetworkOnIDs(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network mismatch or wrong checksum given")
}

// A bad checksum must stop Execute before any request reaches the mirror node.
func TestUnitMirrorNodeAccountBalanceQueryBadChecksumErrorsBeforeRequest(t *testing.T) {
	var called atomic.Bool
	client := newMockMirrorClient(t, "checksum.example.com:443", func(w http.ResponseWriter, r *http.Request) {
		called.Store(true)
		w.WriteHeader(http.StatusOK)
	})
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	_, err = NewMirrorNodeAccountBalanceQuery().SetAccountID(accountID).Execute(client)

	require.Error(t, err)
	assert.False(t, called.Load(), "a checksum mismatch must not reach the network")
}

// The fallback must not be a single attempt; 5xx responses still have to be retried.
func TestUnitMirrorNodeAccountBalanceQueryResolveAttempts(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	require.Equal(t, -1, client.GetMaxAttempts(), "an unset client reports -1, not a usable count")

	query := NewMirrorNodeAccountBalanceQuery()
	assert.Equal(t, uint64(maxAttempts), query.resolveAttempts(client), "falls back to the SDK default")

	client.SetMaxAttempts(4)
	assert.Equal(t, uint64(4), query.resolveAttempts(client), "client setting is used when the query has none")

	query.SetMaxAttempts(2)
	assert.Equal(t, uint64(2), query.resolveAttempts(client), "the query setting wins")
}

// A client-level budget must reach the retry loop, not just resolveAttempts.
func TestUnitMirrorNodeAccountBalanceQueryHonoursClientMaxAttempts(t *testing.T) {
	var attempts atomic.Int32
	client := newMockMirrorClient(t, "clientattempts.example.com:443", func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	client.SetMaxAttempts(2)

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 5}).
		Execute(client)

	require.Error(t, err)
	assert.Equal(t, int32(2), attempts.Load(), "the client's budget of 2 bounds the retries")
}

func TestUnitMirrorNodeAccountBalanceQueryBuildURL(t *testing.T) {
	t.Parallel()

	query := NewMirrorNodeAccountBalanceQuery().SetAccountID(AccountID{Shard: 1, Realm: 2, Account: 3})

	assert.Equal(t, "https://mirror.example.com/api/v1/balances?account.id=1.2.3",
		query.buildURL("https://mirror.example.com/api/v1"))
}

func TestUnitMirrorNodeAccountBalanceQueryWrapsTransportFailure(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork([]string{"127.0.0.1:1"})

	_, err = NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 3}).
		SetMaxAttempts(1).
		Execute(client)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send request")
}
