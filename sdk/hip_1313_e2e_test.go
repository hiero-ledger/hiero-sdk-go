//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationHIP1313HighVolumeAccountCreate(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		SetHighVolume(true).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NotEqual(t, AccountID{}, accountID)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, record.HighVolumePricingMultiplier, uint64(1000))
}

func TestIntegrationHIP1313HighVolumeWithMaxTransactionFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		SetHighVolume(true).
		SetMaxTransactionFee(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NotEqual(t, AccountID{}, accountID)

	// Verify fee charged does not exceed the max transaction fee
	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)
	assert.True(t, record.TransactionFee.AsTinybar() <= NewHbar(2).AsTinybar())
}

func TestIntegrationHIP1313HighVolumeInsufficientFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	_, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		SetHighVolume(true).
		SetMaxTransactionFee(HbarFromTinybar(1)).
		Execute(env.Client)
	require.Error(t, err)
	assert.Equal(t, "exceptional precheck status INSUFFICIENT_TX_FEE", err.Error())
}
