//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitRegisteredNodeCreateTransactionBuild(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress:   []byte{10, 0, 0, 1},
			port:        8080,
			requiresTls: true,
		},
		endpointApis: []BlockNodeApi{BlockNodeApiStatus},
	}

	tx := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("Test Node").
		AddServiceEndpoint(endpoint)

	body := tx.buildProtoBody()
	assert.NotNil(t, body.AdminKey)
	assert.Equal(t, "Test Node", body.Description)
	assert.Len(t, body.ServiceEndpoint, 1)
}

func TestUnitRegisteredNodeCreateTransactionRoundTrip(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &MirrorNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress:   []byte{192, 168, 1, 1},
			port:        443,
			requiresTls: true,
		},
	}

	tx := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("My Node").
		AddServiceEndpoint(endpoint)

	body := tx.buildProtoBody()

	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_RegisteredNodeCreate{
			RegisteredNodeCreate: body,
		},
	}

	restored := _RegisteredNodeCreateTransactionFromProtobuf(*tx.Transaction, pbBody)
	assert.Equal(t, "My Node", restored.GetDescription())
	assert.Len(t, restored.GetServiceEndpoints(), 1)

	_, ok := restored.GetServiceEndpoints()[0].(*MirrorNodeServiceEndpoint)
	assert.True(t, ok, "expected *MirrorNodeServiceEndpoint after round-trip")
}

func TestUnitRegisteredNodeCreateTransactionGetMethod(t *testing.T) {
	t.Parallel()

	tx := NewRegisteredNodeCreateTransaction()
	assert.Equal(t, "RegisteredNodeCreateTransaction", tx.getName())
}

func TestUnitRegisteredNodeCreateTransactionGettersSetters(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &RpcRelayServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			domainName: "node.example.com",
			port:       443,
		},
	}

	tx := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("Test").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint})

	assert.Equal(t, key.PublicKey(), tx.GetAdminKey())
	assert.Equal(t, "Test", tx.GetDescription())
	assert.Len(t, tx.GetServiceEndpoints(), 1)
}

func registeredNodeCreateMockResponses() [][]interface{} {
	return [][]interface{}{{
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_COST_ANSWER,
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status:           services.ResponseCodeEnum_SUCCESS,
						RegisteredNodeId: 11,
					},
				},
			},
		},
	}}
}

func TestUnitRegisteredNodeCreateTransactionMock(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer(registeredNodeCreateMockResponses())
	defer server.Close()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress: []byte{10, 0, 0, 1},
			port:      8080,
		},
		endpointApis: []BlockNodeApi{BlockNodeApiPublish},
	}

	tran := TransactionIDGenerate(AccountID{Account: 3})
	resp, err := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("mock node").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		SetTransactionID(tran).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, uint64(11), receipt.RegisteredNodeId)
}

func TestUnitRegisteredNodeCreateTransactionScheduledBuild(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &MirrorNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress: []byte{192, 168, 1, 1},
			port:      443,
		},
	}

	tx := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("scheduled").
		AddServiceEndpoint(endpoint)

	scheduled, err := tx.buildScheduled()
	require.NoError(t, err)
	require.NotNil(t, scheduled.GetRegisteredNodeCreate())

	again, err := tx.constructScheduleProtobuf()
	require.NoError(t, err)
	require.NotNil(t, again.GetRegisteredNodeCreate())
}

func buildFrozenRegisteredNodeCreateTransaction(t *testing.T) (*RegisteredNodeCreateTransaction, *Client, PrivateKey) {
	t.Helper()
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress: []byte{10, 0, 0, 1},
			port:      8080,
		},
		endpointApis: []BlockNodeApi{BlockNodeApiStatus},
	}

	trx, err := NewRegisteredNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(key.PublicKey()).
		SetDescription("coverage").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint}).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	return trx, client, key
}

func TestUnitRegisteredNodeCreateTransactionCoverage(t *testing.T) {
	t.Parallel()

	trx, client, key := buildFrozenRegisteredNodeCreateTransaction(t)

	require.NoError(t, trx.validateNetworkOnIDs(client))
	_, err := trx.Schedule()
	require.NoError(t, err)
	trx.GetTransactionID()
	trx.GetNodeAccountIDs()
	trx.GetMaxRetry()
	trx.GetMaxTransactionFee()
	trx.GetMaxBackoff()
	trx.GetMinBackoff()
	trx.GetRegenerateTransactionID()

	byt, err := trx.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)

	sig, err := key.SignTransaction(trx)
	require.NoError(t, err)

	_, err = trx.GetTransactionHash()
	require.NoError(t, err)
	_, err = trx.GetSignatures()
	require.NoError(t, err)

	switch b := txFromBytes.(type) {
	case RegisteredNodeCreateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}
