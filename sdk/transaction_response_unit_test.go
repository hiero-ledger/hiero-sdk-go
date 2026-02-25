//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitTransactionResponseReceiptQueryPinnedToSubmittingNode(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{}, {}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	nodeID := AccountID{Account: 3}
	response := TransactionResponse{
		TransactionID: testTransactionID,
		NodeID:        nodeID,
	}

	query := response.GetReceiptQuery(client)
	nodeAccountIDs := query.GetNodeAccountIDs()

	require.Len(t, nodeAccountIDs, 1)
	assert.Equal(t, nodeID, nodeAccountIDs[0])
}

func TestUnitTransactionResponseRecordQueryPinnedToSubmittingNode(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{}, {}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	nodeID := AccountID{Account: 3}
	response := TransactionResponse{
		TransactionID: testTransactionID,
		NodeID:        nodeID,
	}

	query := response.GetRecordQuery(client)
	nodeAccountIDs := query.GetNodeAccountIDs()

	require.Len(t, nodeAccountIDs, 1)
	assert.Equal(t, nodeID, nodeAccountIDs[0])
}

func TestUnitTransactionResponseFailoverEnabledSubmittingNodeFirst(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{}, {}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	client.SetAllowReceiptNodeFailover(true)

	nodeID := AccountID{Account: 3}
	response := TransactionResponse{
		TransactionID: testTransactionID,
		NodeID:        nodeID,
	}

	receiptQuery := response.GetReceiptQuery(client).GetNodeAccountIDs()
	require.Greater(t, len(receiptQuery), 1)
	assert.Equal(t, nodeID, receiptQuery[0])
}

func TestUnitTransactionResponseFailoverEnabledNoDuplicateNodes(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{}, {}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	client.SetAllowReceiptNodeFailover(true)

	nodeID := AccountID{Account: 3}
	response := TransactionResponse{
		TransactionID: testTransactionID,
		NodeID:        nodeID,
	}

	query := response.GetReceiptQuery(client)
	nodeAccountIDs := query.GetNodeAccountIDs()

	seen := make(map[AccountID]bool)
	for _, id := range nodeAccountIDs {
		assert.False(t, seen[id], "duplicate node account ID: %v", id)
		seen[id] = true
	}
}

func TestUnitTransactionResponseGetNodeAccountIDsNilClient(t *testing.T) {
	t.Parallel()

	nodeID := AccountID{Account: 3}
	response := TransactionResponse{
		TransactionID: testTransactionID,
		NodeID:        nodeID,
	}

	receiptQuery := response.GetReceiptQuery(nil).GetNodeAccountIDs()
	require.Len(t, receiptQuery, 1)
	assert.Equal(t, nodeID, receiptQuery[0])
}

func TestUnitTransactionResponseFailoverContainsAllHealthyNodes(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{}, {}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	client.SetAllowReceiptNodeFailover(true)

	// Node 0.0.3 is the submitting node, 0.0.4 is the other node in the mock network
	nodeID := AccountID{Account: 3}
	otherNodeID := AccountID{Account: 4}
	response := TransactionResponse{
		TransactionID: testTransactionID,
		NodeID:        nodeID,
	}

	query := response.GetReceiptQuery(client)
	nodeAccountIDs := query.GetNodeAccountIDs()

	require.Len(t, nodeAccountIDs, 2)
	assert.Equal(t, nodeID, nodeAccountIDs[0])
	assert.Equal(t, otherNodeID, nodeAccountIDs[1])
}
