//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// waitForTokenRelationships polls until the mirror node shows at least want relationships for accountID.
func waitForTokenRelationships(t *testing.T, client *Client, accountID AccountID, want int) MirrorNodeTokenBalancePage {
	t.Helper()

	var page MirrorNodeTokenBalancePage
	var err error
	for range 20 {
		page, err = NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).Execute(client)
		if err == nil && len(page.Tokens) >= want {
			return page
		}
		time.Sleep(1500 * time.Millisecond)
	}
	require.NoError(t, err)
	require.FailNow(t, "mirror node did not ingest the token relationships in time", "have %d, want %d", len(page.Tokens), want)

	return page
}

func autoAssociating(tx *AccountCreateTransaction) {
	tx.SetMaxAutomaticTokenAssociations(-1)
}

func TestIntegrationMirrorNodeTokenBalanceQueryReturnsFungibleAndNftBalances(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	accountID, _, err := createAccount(&env, autoAssociating)
	require.NoError(t, err)
	fungible, err := createFungibleToken(&env, func(tx *TokenCreateTransaction) { tx.SetDecimals(3) })
	require.NoError(t, err)
	nft, err := createNft(&env)
	require.NoError(t, err)
	unrelated, err := createFungibleToken(&env)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(nft).SetMetadatas([][]byte{{1}, {2}}).Execute(env.Client)
	require.NoError(t, err)
	minted, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	transfer, err := NewTransferTransaction().
		AddTokenTransfer(fungible, env.OperatorID, -1500).
		AddTokenTransfer(fungible, accountID, 1500).
		AddNftTransfer(nft.Nft(minted.SerialNumbers[0]), env.OperatorID, accountID).
		AddNftTransfer(nft.Nft(minted.SerialNumbers[1]), env.OperatorID, accountID).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	page := waitForTokenRelationships(t, env.Client, accountID, 2)
	require.Len(t, page.Tokens, 2)
	assert.Empty(t, page.Next)
	assert.Less(t, page.Tokens[0].TokenID.Token, page.Tokens[1].TokenID.Token)

	byToken := map[TokenID]MirrorNodeTokenBalance{}
	for _, balance := range page.Tokens {
		byToken[balance.TokenID] = balance
		assert.NotEmpty(t, balance.KycStatus)
		assert.NotEmpty(t, balance.CreatedTimestamp)
		require.NotNil(t, balance.AutomaticAssociation)
		assert.True(t, *balance.AutomaticAssociation)
	}
	assert.Equal(t, uint64(1500), byToken[fungible].Balance)
	assert.Equal(t, uint64(3), byToken[fungible].Decimals)
	assert.Equal(t, uint64(2), byToken[nft].Balance)
	assert.Zero(t, byToken[nft].Decimals)

	single, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).SetTokenID(fungible).Execute(env.Client)
	require.NoError(t, err)
	require.Len(t, single.Tokens, 1)
	assert.Equal(t, uint64(1500), single.Tokens[0].Balance)
	assert.Empty(t, single.Next)

	none, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).SetTokenID(unrelated).Execute(env.Client)
	require.NoError(t, err)
	assert.Empty(t, none.Tokens)
}

func TestIntegrationMirrorNodeTokenBalanceQueryNonExistentAccountErrors(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(AccountID{Account: 9_999_999}).Execute(env.Client)
	require.ErrorIs(t, err, ErrHederaPreCheckStatus{Status: StatusInvalidAccountID})
}

func TestIntegrationMirrorNodeTokenBalanceQueryKeepsAssociationsToDeletedTokens(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	accountID, _, err := createAccount(&env, autoAssociating)
	require.NoError(t, err)
	token, err := createFungibleToken(&env)
	require.NoError(t, err)

	transfer, err := NewTransferTransaction().
		AddTokenTransfer(token, env.OperatorID, -100).
		AddTokenTransfer(token, accountID, 100).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = transfer.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	waitForTokenRelationships(t, env.Client, accountID, 1)

	deletion, err := NewTokenDeleteTransaction().SetTokenID(token).Execute(env.Client)
	require.NoError(t, err)
	_, err = deletion.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	waitForMirrorTokenDeleted(t, env.Client, token)

	page, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).Execute(env.Client)
	require.NoError(t, err)
	single, err := NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).SetTokenID(token).Execute(env.Client)
	require.NoError(t, err)

	for name, got := range map[string]MirrorNodeTokenBalancePage{"all tokens": page, "single token": single} {
		require.Len(t, got.Tokens, 1, name)
		assert.Equal(t, token, got.Tokens[0].TokenID, name)
		assert.Equal(t, uint64(100), got.Tokens[0].Balance, name)
	}
}

// waitForMirrorTokenDeleted polls until the mirror node reports the token deleted.
func waitForMirrorTokenDeleted(t *testing.T, client *Client, token TokenID) {
	t.Helper()

	path, err := newMirrorNodeRestPath("/tokens/" + token.String())
	require.NoError(t, err)
	for range 20 {
		result, err := mirrorNodeGetJSONObject(client, path)
		if err == nil && result["deleted"] == true {
			return
		}
		time.Sleep(1500 * time.Millisecond)
	}
	require.FailNow(t, "mirror node never reported the token deleted", "token %s", token)
}
