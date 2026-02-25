//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTransactionResponseRecordQueryPinnedToSubmittingNode(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)

	// Record query should be pinned to submitting node only
	recordQuery := resp.GetRecordQuery(env.Client)
	nodeAccountIDs := recordQuery.GetNodeAccountIDs()

	require.Len(t, nodeAccountIDs, 1)
	assert.Equal(t, resp.NodeID, nodeAccountIDs[0])

	// Verify record can still be obtained
	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, record.Receipt.AccountID)
}

func TestIntegrationTransactionResponseRecordWithFailoverEnabled(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Enable receipt node failover
	env.Client.SetAllowReceiptNodeFailover(true)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	recordQuery := resp.GetRecordQuery(env.Client)
	nodeAccountIDs := recordQuery.GetNodeAccountIDs()

	require.GreaterOrEqual(t, len(nodeAccountIDs), 1)
	assert.Equal(t, resp.NodeID, nodeAccountIDs[0])

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, record.Receipt.AccountID)
}
