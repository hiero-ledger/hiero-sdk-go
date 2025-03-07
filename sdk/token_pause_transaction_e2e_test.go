//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenPauseTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createFungibleToken(&env)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)

	require.NotNil(t, info.PauseStatus)
	require.False(t, *info.PauseStatus)

	resp, err := NewTokenPauseTransaction().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	require.NotNil(t, info.PauseStatus)
	require.True(t, *info.PauseStatus)

	//Unpause token to avoid error in CloseIntegrationTestEnv
	resp, err = NewTokenUnpauseTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(env.Client)

	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenPauseTransactionFromToBytes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createFungibleToken(&env)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(env.Client)

	require.NotNil(t, info.PauseStatus)
	require.False(t, *info.PauseStatus)

	tokenPause := NewTokenPauseTransaction().
		SetTokenID(tokenID)
	require.NoError(t, err)

	tokenPauseBytes, err := tokenPause.ToBytes()
	require.NoError(t, err)

	tokenPauseFromBytes, err := TransactionFromBytes(tokenPauseBytes)
	require.NoError(t, err)

	tokenPauseTransaction := tokenPauseFromBytes.(TokenPauseTransaction)
	require.Equal(t, tokenPauseTransaction.GetTokenID(), tokenID)
	resp, err := tokenPauseTransaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	require.NotNil(t, info.PauseStatus)
	require.True(t, *info.PauseStatus)

	tokenUnpause := NewTokenUnpauseTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID)

	tokenUnpauseBytes, err := tokenUnpause.ToBytes()
	require.NoError(t, err)

	tokenUnpauseFromBytes, err := TransactionFromBytes(tokenUnpauseBytes)
	require.NoError(t, err)

	tokenUnpauseTransaction := tokenUnpauseFromBytes.(TokenUnpauseTransaction)
	require.Equal(t, tokenUnpauseTransaction.GetTokenID(), tokenID)
	require.Equal(t, len(tokenUnpauseTransaction.GetNodeAccountIDs()), 1)
	resp, err = tokenUnpauseTransaction.Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}
