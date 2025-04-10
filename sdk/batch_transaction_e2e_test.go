//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationBatchTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	innerTransaction, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetInitialBalance(newBalance).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxAutomaticTokenAssociations(100).
		Batchify(env.OperatorKey, env.Client)
	require.NoError(t, err)

	batchTransaction := NewBatchTransaction().
		AddInnerTransaction(innerTransaction)

	innerTransactionIDs := batchTransaction.GetInnerTransactionIDs()

	resp, err := batchTransaction.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	receipt, err := NewTransactionReceiptQuery().
		SetTransactionID(innerTransactionIDs[0]).Execute(env.Client)

	require.NoError(t, err)
	fmt.Println("accountId", receipt.AccountID.String())
}
