//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitMirrorNodeAccountBalanceQueryGetSet(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 1234}
	contractID := ContractID{Contract: 5678}

	q := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(accountID).
		SetMaxAttempts(5)

	assert.Equal(t, accountID, q.GetAccountID())
	assert.Equal(t, uint64(5), q.GetMaxAttempts())

	// Setting a contract ID clears the account ID (they are mutually exclusive).
	q.SetContractID(contractID)
	assert.Equal(t, contractID, q.GetContractID())
	assert.Equal(t, AccountID{}, q.GetAccountID())
}

func TestUnitMirrorNodeAccountBalanceQueryNilClient(t *testing.T) {
	t.Parallel()

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 1234}).
		Execute(nil)
	require.ErrorIs(t, err, errNoClientProvided)
}

func TestUnitMirrorNodeAccountBalanceQueryNoIDSet(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetMirrorNetwork([]string{"mirror.example.com:443"})

	_, err = NewMirrorNodeAccountBalanceQuery().Execute(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "either an account ID or a contract ID must be set")
}

func TestUnitMirrorNodeAccountBalanceQuerySuccess(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balance.example.com:443"

	var accountPath, tokensPath, tokensLimit string
	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/tokens") {
			tokensPath = r.URL.Path
			tokensLimit = r.URL.Query().Get("limit")
			require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
				"tokens": []map[string]interface{}{
					{"token_id": "0.0.1001", "balance": 42, "decimals": 2},
					{"token_id": "0.0.1002", "balance": 7, "decimals": 0},
				},
				"links": map[string]interface{}{"next": nil},
			}))
			return
		}
		accountPath = r.URL.Path
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
			"account": "0.0.1234",
			"balance": map[string]interface{}{
				"balance": 500000000,
				// The embedded token list is truncated by the mirror node and must be ignored;
				// return a bogus entry to prove it is not used.
				"tokens": []map[string]interface{}{
					{"token_id": "0.0.9999", "balance": 999},
				},
			},
		}))
	})

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 1234}).
		Execute(client)
	require.NoError(t, err)

	assert.True(t, strings.HasSuffix(accountPath, "/accounts/0.0.1234"), "unexpected account path: %s", accountPath)
	assert.True(t, strings.HasSuffix(tokensPath, "/accounts/0.0.1234/tokens"), "unexpected tokens path: %s", tokensPath)
	// The maximum page size must be requested; the default of 25 would cap the walk at 2,500 tokens.
	assert.Equal(t, "100", tokensLimit)

	assert.Equal(t, HbarFromTinybar(500000000).AsTinybar(), balance.Hbars.AsTinybar())

	// Token balances come from the /tokens endpoint, not the truncated embedded list.
	assert.Equal(t, uint64(42), balance.Tokens.Get(TokenID{Token: 1001}))
	assert.Equal(t, uint64(7), balance.Tokens.Get(TokenID{Token: 1002}))
	assert.Equal(t, uint64(42), balance.Token[TokenID{Token: 1001}])
	assert.Equal(t, uint64(0), balance.Tokens.Get(TokenID{Token: 9999}))

	// Decimals are populated from the /tokens endpoint.
	assert.Equal(t, uint64(2), balance.TokenDecimals.Get(TokenID{Token: 1001}))
	assert.Equal(t, uint64(0), balance.TokenDecimals.Get(TokenID{Token: 1002}))
}

func TestUnitMirrorNodeAccountBalanceQueryPagination(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balancepaged.example.com:443"

	tokensCalls := 0
	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !strings.HasSuffix(r.URL.Path, "/tokens") {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
				"balance": map[string]interface{}{"balance": 100},
			}))
			return
		}

		// Second page is requested once the first page hands back a links.next.
		if r.URL.Query().Get("page") == "2" {
			tokensCalls++
			require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
				"tokens": []map[string]interface{}{
					{"token_id": "0.0.2002", "balance": 20, "decimals": 4},
				},
				"links": map[string]interface{}{"next": nil},
			}))
			return
		}

		tokensCalls++
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
			"tokens": []map[string]interface{}{
				{"token_id": "0.0.2001", "balance": 10, "decimals": 1},
			},
			"links": map[string]interface{}{
				"next": "/api/v1/accounts/0.0.1234/tokens?page=2",
			},
		}))
	})

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 1234}).
		Execute(client)
	require.NoError(t, err)

	assert.Equal(t, 2, tokensCalls, "expected both token pages to be fetched")
	assert.Equal(t, uint64(10), balance.Tokens.Get(TokenID{Token: 2001}))
	assert.Equal(t, uint64(20), balance.Tokens.Get(TokenID{Token: 2002}))
	assert.Equal(t, uint64(1), balance.TokenDecimals.Get(TokenID{Token: 2001}))
	assert.Equal(t, uint64(4), balance.TokenDecimals.Get(TokenID{Token: 2002}))
}

// A wrong-network checksum must fail Execute before any network call; no transport is mocked here,
// so reaching the network would fail differently.
func TestUnitMirrorNodeAccountBalanceQueryValidateChecksum(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)
	client.SetMirrorNetwork([]string{"mirror.example.com:443"})

	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	require.NoError(t, NewMirrorNodeAccountBalanceQuery().
		SetAccountID(accountID).
		validateNetworkOnIDs(client))

	wrongAccountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	_, err = NewMirrorNodeAccountBalanceQuery().
		SetAccountID(wrongAccountID).
		Execute(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network mismatch or wrong checksum given")

	wrongContractID, err := ContractIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	_, err = NewMirrorNodeAccountBalanceQuery().
		SetContractID(wrongContractID).
		Execute(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network mismatch or wrong checksum given")
}

// A transient 503 must be retried, with the query succeeding on the next attempt.
func TestUnitMirrorNodeAccountBalanceQueryRetriesTransient503(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balanceretry.example.com:443"

	accountCalls := 0
	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/tokens") {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
				"tokens": []map[string]interface{}{},
				"links":  map[string]interface{}{"next": nil},
			}))
			return
		}
		accountCalls++
		if accountCalls == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
			"balance": map[string]interface{}{"balance": 100},
		}))
	})

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 1234}).
		Execute(client)
	require.NoError(t, err)

	assert.Equal(t, 2, accountCalls, "expected the initial 503 to be retried once")
	assert.Equal(t, int64(100), balance.Hbars.AsTinybar())
}

func TestUnitMirrorNodeAccountBalanceQueryNon200(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balancenotfound.example.com:443"

	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"_status":{"messages":[{"message":"Not found"}]}}`))
	})

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 9999}).
		Execute(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-200")
}

func TestUnitMirrorNodeAccountBalanceQueryTokensNon200(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balancetokenserr.example.com:443"

	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/tokens") {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"_status":{"messages":[{"message":"Not found"}]}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
			"balance": map[string]interface{}{"balance": 100},
		}))
	})

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 1234}).
		Execute(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-200")
}

func TestUnitMirrorNodeAccountBalanceQueryContractUsesAccountsEndpoint(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balancecontract.example.com:443"

	var accountPath, tokensPath string
	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/tokens") {
			tokensPath = r.URL.Path
			require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
				"tokens": []map[string]interface{}{},
				"links":  map[string]interface{}{"next": nil},
			}))
			return
		}
		accountPath = r.URL.Path
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
			"balance": map[string]interface{}{"balance": 123},
		}))
	})

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetContractID(ContractID{Contract: 5678}).
		Execute(client)
	require.NoError(t, err)

	// Contracts are routed through the accounts endpoint, not the contracts endpoint.
	assert.True(t, strings.HasSuffix(accountPath, "/accounts/0.0.5678"), "unexpected account path: %s", accountPath)
	assert.True(t, strings.HasSuffix(tokensPath, "/accounts/0.0.5678/tokens"), "unexpected tokens path: %s", tokensPath)
	assert.Equal(t, HbarFromTinybar(123).AsTinybar(), balance.Hbars.AsTinybar())
}

// A contract built from an EVM address must be queried by bare hex; ContractID.String() would
// render 0.0.<hex>, which the mirror node does not resolve.
func TestUnitMirrorNodeAccountBalanceQueryContractByEvmAddress(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balanceevmcontract.example.com:443"
	const evmAddress = "742d35cc6634c0532925a3b844bc454e4438f44e"

	var accountPath string
	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/tokens") {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
				"tokens": []map[string]interface{}{},
				"links":  map[string]interface{}{"next": nil},
			}))
			return
		}
		accountPath = r.URL.Path
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
			"balance": map[string]interface{}{"balance": 77},
		}))
	})

	contractID, err := ContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetContractID(contractID).
		Execute(client)
	require.NoError(t, err)

	assert.True(t, strings.HasSuffix(accountPath, "/accounts/"+evmAddress), "unexpected account path: %s", accountPath)
	assert.Equal(t, int64(77), balance.Hbars.AsTinybar())
}

// An unparseable token ID must fail the query rather than silently drop the token and report a
// partial balance as complete.
func TestUnitMirrorNodeAccountBalanceQueryInvalidTokenID(t *testing.T) {
	// Not parallel: SetupMockTransportForDomain mutates http.DefaultTransport.
	const domain = "balancebadtoken.example.com:443"

	client := newMockMirrorClient(t, domain, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/tokens") {
			require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
				"tokens": []map[string]interface{}{
					{"token_id": "not-a-token-id", "balance": 5, "decimals": 0},
				},
				"links": map[string]interface{}{"next": nil},
			}))
			return
		}
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{
			"balance": map[string]interface{}{"balance": 100},
		}))
	})

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 1234}).
		Execute(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token ID")
}

// A mirror network that resolves to an unusable URL must fail immediately rather than burn the
// whole retry budget on a request that can never succeed.
func TestUnitMirrorNodeAccountBalanceQueryRejectsUnusableURL(t *testing.T) {
	t.Parallel()

	require.NoError(t, mirrorNodeValidateURL("https://mirror.example.com:443/api/v1/accounts/0.0.1"))
	require.NoError(t, mirrorNodeValidateURL("http://localhost:5551/api/v1"))

	require.ErrorContains(t, mirrorNodeValidateURL("ftp://mirror.example.com/api/v1"), "unsupported mirror node URL scheme")
	// A bare path has no scheme at all, so it trips the scheme check first.
	require.ErrorContains(t, mirrorNodeValidateURL("/api/v1/accounts/0.0.1"), "unsupported mirror node URL scheme")
	require.ErrorContains(t, mirrorNodeValidateURL("http:///api/v1"), "has no host")
	require.ErrorContains(t, mirrorNodeValidateURL("http://[::1/api"), "invalid mirror node URL")
}
