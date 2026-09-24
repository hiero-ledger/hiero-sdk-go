//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tokenBalanceClient returns a mock client whose mirror node requests go to a scripted fake transport.
func tokenBalanceClient(t *testing.T, turns ...fakeHttpTurn) (*Client, *fakeHttpTransport) {
	t.Helper()

	client, err := _NewMockClient()
	require.NoError(t, err)
	transport := &fakeHttpTransport{turns: turns}
	client.SetMirrorNodeHttpConfig(DefaultMirrorNodeHttpConfig().WithTransport(transport).WithRetryPolicy(fastRetryPolicy(3)))

	return client, transport
}

func okBody(body string) fakeHttpTurn {
	return fakeHttpTurn{resp: HttpResponse{statusCode: http.StatusOK, body: []byte(body)}}
}

const twoTokenRelationships = `{"tokens":[
	{"token_id":"0.0.1001","balance":250,"decimals":2,"automatic_association":true,"created_timestamp":"1700000000.000000001","freeze_status":"UNFROZEN","kyc_status":"GRANTED"},
	{"token_id":"0.0.1002","balance":3,"decimals":0,"automatic_association":false,"created_timestamp":"1700000001.000000000","freeze_status":"NOT_APPLICABLE","kyc_status":"NOT_APPLICABLE"}
],"links":{"next":null}}`

func TestUnitMirrorNodeTokenBalanceQueryGetSet(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Shard: 1, Realm: 2, Account: 3}
	tokenID := TokenID{Shard: 1, Realm: 2, Token: 4}
	query := NewMirrorNodeTokenBalanceQuery().
		SetAccountID(accountID).
		SetTokenID(tokenID)

	assert.Equal(t, accountID, query.GetAccountID())
	assert.Equal(t, tokenID, query.GetTokenID())

	next := "/api/v1/accounts/0.0.5/tokens?limit=100&order=asc&token.id=gt:0.0.1002"
	assert.Equal(t, next, NewMirrorNodeTokenBalanceQuery().SetNextPage(next).GetNextPage())
}

func TestUnitMirrorNodeTokenBalanceQueryDefaults(t *testing.T) {
	t.Parallel()

	query := NewMirrorNodeTokenBalanceQuery()

	assert.Equal(t, AccountID{}, query.GetAccountID())
	assert.Equal(t, TokenID{}, query.GetTokenID())
	assert.Empty(t, query.GetNextPage())
}

func TestUnitMirrorNodeTokenBalanceQueryMapsEveryField(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, okBody(twoTokenRelationships))

	page, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
	require.NoError(t, err)

	assert.Contains(t, transport.lastRequest().url, "/api/v1/accounts/0.0.5/tokens?limit=100&order=asc")
	assert.Empty(t, page.Next)
	require.Len(t, page.Tokens, 2)

	automatic := true
	assert.Equal(t, MirrorNodeTokenBalance{
		TokenID:              TokenID{Token: 1001},
		Balance:              250,
		Decimals:             2,
		AutomaticAssociation: &automatic,
		CreatedTimestamp:     "1700000000.000000001",
		FreezeStatus:         "UNFROZEN",
		KycStatus:            "GRANTED",
	}, page.Tokens[0])
	assert.Equal(t, uint64(3), page.Tokens[1].Balance)
	assert.Zero(t, page.Tokens[1].Decimals)
}

func TestUnitMirrorNodeTokenBalanceQueryFollowsNextPage(t *testing.T) {
	t.Parallel()

	const next = "/api/v1/accounts/0.0.5/tokens?limit=100&order=asc&token.id=gt:0.0.1002"
	client, transport := tokenBalanceClient(t,
		okBody(`{"tokens":[{"token_id":"0.0.1002","balance":1,"decimals":0,"kyc_status":"NOT_APPLICABLE"}],"links":{"next":"`+next+`"}}`),
		okBody(`{"tokens":[],"links":{"next":null}}`),
	)

	first, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
	require.NoError(t, err)
	assert.Equal(t, next, first.Next)

	last, err := NewMirrorNodeTokenBalanceQuery().SetNextPage(first.Next).Execute(client)
	require.NoError(t, err)
	assert.Empty(t, last.Tokens)
	assert.Empty(t, last.Next)
	assert.Contains(t, transport.lastRequest().url, "/api/v1/accounts/0.0.5/tokens?limit=100&order=asc&token.id=gt:0.0.1002")
}

func TestUnitMirrorNodeTokenBalanceQuerySingleTokenHasNoNextPage(t *testing.T) {
	t.Parallel()

	// The mirror node echoes the request as next when an exact token.id page is full.
	client, transport := tokenBalanceClient(t, okBody(`{"tokens":[{"token_id":"0.0.1001","balance":9,"decimals":1,"kyc_status":"GRANTED"}],
		"links":{"next":"/api/v1/accounts/0.0.5/tokens?limit=100&token.id=0.0.1001"}}`))

	page, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).SetTokenID(TokenID{Token: 1001}).Execute(client)
	require.NoError(t, err)

	assert.Contains(t, transport.lastRequest().url, "/api/v1/accounts/0.0.5/tokens?limit=100&token.id=0.0.1001")
	assert.Empty(t, page.Next)
	require.Len(t, page.Tokens, 1)
	assert.Equal(t, uint64(9), page.Tokens[0].Balance)
}

func TestUnitMirrorNodeTokenBalanceQueryEmptyPagesAreNotErrors(t *testing.T) {
	t.Parallel()

	client, _ := tokenBalanceClient(t, okBody(`{"tokens":[],"links":{"next":null}}`))

	page, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).SetTokenID(TokenID{Token: 7}).Execute(client)
	require.NoError(t, err)
	assert.Empty(t, page.Tokens)
}

func TestUnitMirrorNodeTokenBalanceQueryNonExistentAccountErrors(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, fakeHttpTurn{resp: HttpResponse{
		statusCode: http.StatusNotFound,
		body:       []byte(`{"_status":{"messages":[{"message":"Not found"}]}}`),
	}})

	_, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Shard: 1, Realm: 2, Account: 3}).Execute(client)
	require.ErrorIs(t, err, ErrHederaPreCheckStatus{Status: StatusInvalidAccountID})
	assert.Equal(t, 1, transport.callCount())
}

func TestUnitMirrorNodeTokenBalanceQueryResolvesEvmAddress(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, okBody(`{"tokens":[],"links":{"next":null}}`))
	accountID, err := AccountIDFromEvmAddress(0, 0, "7F2c5e0f4a3b9E8d1C6A5b4E3f2D1c0B9a8E7F6d")
	require.NoError(t, err)

	_, err = NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).Execute(client)
	require.NoError(t, err)
	assert.Contains(t, transport.lastRequest().url, "/accounts/7f2c5e0f4a3b9e8d1c6a5b4e3f2d1c0b9a8e7f6d/tokens")
}

func TestUnitMirrorNodeTokenBalanceQueryRejectsInvalidCombinations(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, okBody(`{"tokens":[],"links":{"next":null}}`))
	const next = "/api/v1/accounts/0.0.5/tokens?limit=100&order=asc&token.id=gt:0.0.1"

	_, err := NewMirrorNodeTokenBalanceQuery().Execute(client)
	require.ErrorIs(t, err, errMirrorNodeTokenBalanceQueryNoAccountID)

	_, err = NewMirrorNodeTokenBalanceQuery().SetTokenID(TokenID{Token: 1}).Execute(client)
	require.ErrorIs(t, err, errMirrorNodeTokenBalanceQueryNoAccountID)

	_, err = NewMirrorNodeTokenBalanceQuery().SetNextPage(next).SetTokenID(TokenID{Token: 1}).Execute(client)
	require.ErrorIs(t, err, errMirrorNodeTokenBalanceQueryNextPageConflict)

	_, err = NewMirrorNodeTokenBalanceQuery().SetNextPage(next).SetAccountID(AccountID{Account: 5}).Execute(client)
	require.ErrorIs(t, err, errMirrorNodeTokenBalanceQueryNextPageConflict)

	assert.Zero(t, transport.callCount())
}

func TestUnitMirrorNodeTokenBalanceQueryRejectsForeignCursors(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, okBody(`{"tokens":[],"links":{"next":null}}`))

	for _, next := range []string{
		"https://evil.example/api/v1/accounts/0.0.5/tokens",
		"//evil.example/api/v1/accounts/0.0.5/tokens",
		"/api/v1/accounts/0.0.5/../../network/nodes",
		"/api/v1/balances?account.id=0.0.5",
		"/api/v1/accounts/0.0.5",
	} {
		_, err := NewMirrorNodeTokenBalanceQuery().SetNextPage(next).Execute(client)
		assert.Error(t, err, next)
	}
	assert.Zero(t, transport.callCount())
}

func TestUnitMirrorNodeTokenBalanceQueryRetriesWithClientPolicy(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, fakeHttpTurn{resp: statusResponse(http.StatusServiceUnavailable)}, okBody(twoTokenRelationships))
	_, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
	require.NoError(t, err)
	assert.Equal(t, 2, transport.callCount())

	down, _ := tokenBalanceClient(t, fakeHttpTurn{resp: statusResponse(http.StatusServiceUnavailable)})
	_, err = NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(down)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "received non-200 response from mirror node: 503")
}

func TestUnitMirrorNodeTokenBalanceQueryLargeBalanceIsLossless(t *testing.T) {
	t.Parallel()

	client, _ := tokenBalanceClient(t, okBody(`{"tokens":[{"token_id":"0.0.1","balance":9007199254740993,"decimals":0,"kyc_status":"GRANTED"}],"links":{"next":null}}`))

	page, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
	require.NoError(t, err)
	assert.Equal(t, uint64(9007199254740993), page.Tokens[0].Balance)
}

func TestUnitMirrorNodeTokenBalanceQueryValidatesChecksums(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, okBody(`{"tokens":[],"links":{"next":null}}`))
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	badAccount, err := AccountIDFromString("0.0.123-abcde")
	require.NoError(t, err)
	_, err = NewMirrorNodeTokenBalanceQuery().SetAccountID(badAccount).Execute(client)
	require.Error(t, err)

	badToken, err := TokenIDFromString("0.0.123-abcde")
	require.NoError(t, err)
	_, err = NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).SetTokenID(badToken).Execute(client)
	require.Error(t, err)
	assert.Zero(t, transport.callCount())

	alias, err := AccountIDFromEvmAddress(0, 0, "7f2c5e0f4a3b9e8d1c6a5b4e3f2d1c0b9a8e7f6d")
	require.NoError(t, err)
	_, err = NewMirrorNodeTokenBalanceQuery().SetAccountID(alias).Execute(client)
	require.NoError(t, err)
}

func TestUnitMirrorNodeTokenBalanceQueryRejectsMalformedResponses(t *testing.T) {
	t.Parallel()

	for body, why := range map[string]string{
		`<html>502 Bad Gateway</html>`: "a proxy page at 200",
		`{"links":{"next":null}}`:      "no tokens array",
		`{"tokens":null}`:              "a null tokens array",
	} {
		client, transport := tokenBalanceClient(t, okBody(body))
		_, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
		assert.Error(t, err, why)
		assert.Equal(t, 1, transport.callCount(), why)
	}
}

func TestUnitMirrorNodeTokenBalanceQueryRequiresMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, transport := tokenBalanceClient(t, okBody(`{"tokens":[],"links":{"next":null}}`))
	client.SetMirrorNetwork([]string{})

	_, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 5}).Execute(client)
	require.ErrorContains(t, err, "mirror node is not set")
	assert.Zero(t, transport.callCount())

	_, err = NewMirrorNodeTokenBalanceQuery().Execute(nil)
	require.ErrorIs(t, err, errNoClientProvided)
}
