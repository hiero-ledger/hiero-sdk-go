//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func clprExecuteCases(t *testing.T) []struct {
	method string
	build  func() TransactionInterface
} {
	key := clprOwnershipVectors[1].key(t)
	channelID := clprHex(t, clprVectorChannelID)
	salt := clprHex(t, clprVectorSalt)
	connectorID := clprHex(t, clprOwnershipVectors[1].connectorID)

	return []struct {
		method string
		build  func() TransactionInterface
	}{
		{"registerChannel", func() TransactionInterface {
			return NewClprRegisterChannelTransaction().SetOwnershipCommitment(clprBytes(32, 1))
		}},
		{"completeChannel", func() TransactionInterface {
			tx, err := NewClprCompleteChannelTransaction().SetChannelID(channelID).SetVerifierContract(ContractID{Contract: 345}).
				SetConfigProofBytes([]byte{1}).SetOwnershipProof(key)
			require.NoError(t, err)
			return tx
		}},
		{"closeChannel", func() TransactionInterface {
			return NewClprCloseChannelTransaction().SetChannelID(channelID)
		}},
		{"registerConnector", func() TransactionInterface {
			return NewClprRegisterConnectorTransaction().SetCommitment(clprBytes(32, 1)).SetChannelID(channelID)
		}},
		{"completeConnector", func() TransactionInterface {
			tx, err := NewClprCompleteConnectorTransaction().SetChannelID(channelID).SetConnectorID(connectorID).
				SetSalt(salt).SetConnectorContract(ContractID{Contract: 345}).SetAdminKey(key.PublicKey()).SetOwnershipProof(key)
			require.NoError(t, err)
			return tx
		}},
		{"deregisterConnector", func() TransactionInterface {
			return NewClprDeregisterConnectorTransaction().SetChannelID(channelID).SetConnectorID(connectorID).
				SetStakeRecipient(AccountID{Account: 1001})
		}},
		{"submitBundle", func() TransactionInterface {
			return NewClprSubmitBundleTransaction().SetChannelID(channelID).SetBundlePayload([]byte{1})
		}},
		{"updateLedgerConfiguration", func() TransactionInterface {
			return NewClprUpdateLedgerConfigurationTransaction().SetThrottles(clprThrottleLimits())
		}},
	}
}

func clprExecuteOnMock(t *testing.T, tx func(*Client) TransactionInterface) string {
	t.Helper()
	var method string
	handler := func(ctx context.Context, _ *services.Transaction) *services.TransactionResponse {
		method, _ = grpc.Method(ctx)
		return &services.TransactionResponse{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK}
	}
	client, server := NewMockClientAndServer([][]interface{}{{MockTransactionHandlerFunc(handler)}})
	defer server.Close()

	_, err := TransactionExecute(tx(client), client)
	require.NoError(t, err)
	return method
}

func TestUnitClprTransactionsExecuteThroughClprService(t *testing.T) {
	t.Parallel()

	for _, tt := range clprExecuteCases(t) {
		t.Run(tt.method, func(t *testing.T) {
			method := clprExecuteOnMock(t, func(*Client) TransactionInterface { return tt.build() })
			assert.Equal(t, "/proto.ClprService/"+tt.method, method)

			decodedMethod := clprExecuteOnMock(t, func(client *Client) TransactionInterface {
				tx := tt.build()
				_, err := TransactionFreezeWith(tx, client)
				require.NoError(t, err)
				txBytes, err := TransactionToBytes(tx)
				require.NoError(t, err)
				decoded, err := TransactionFromBytes(txBytes)
				require.NoError(t, err)
				return decoded
			})
			assert.Equal(t, "/proto.ClprService/"+tt.method, decodedMethod)
		})
	}
}
