//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUnitRegisteredNodeUpdateTransactionMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
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
						RegisteredNodeId: 7,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(7).
		SetDescription("updated").
		SetTransactionID(tran).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, uint64(7), receipt.RegisteredNodeId)
}

func TestUnitRegisteredNodeUpdateTransactionDescriptionSet(t *testing.T) {
	t.Parallel()

	nodeId := uint64(1)
	tx := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(nodeId).
		SetDescription("Updated Node")

	body := tx.buildProtoBody()
	assert.NotNil(t, body.Description)
	assert.Equal(t, "Updated Node", body.Description.Value)
}

func TestUnitRegisteredNodeUpdateTransactionDescriptionNotSet(t *testing.T) {
	t.Parallel()

	nodeId := uint64(1)
	tx := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(nodeId)

	body := tx.buildProtoBody()
	assert.Nil(t, body.Description)
}

func TestUnitRegisteredNodeUpdateTransactionRoundTrip(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &RpcRelayServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			domainName: "node.example.com",
			port:       443,
		},
	}

	tx := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(42).
		SetAdminKey(key.PublicKey()).
		SetDescription("Updated").
		AddServiceEndpoint(endpoint)

	body := tx.buildProtoBody()
	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_RegisteredNodeUpdate{
			RegisteredNodeUpdate: body,
		},
	}

	restored := _RegisteredNodeUpdateTransactionFromProtobuf(*tx.Transaction, pbBody)
	assert.Equal(t, uint64(42), restored.GetRegisteredNodeId())
	assert.Equal(t, "Updated", restored.GetDescription())
	assert.Len(t, restored.GetServiceEndpoints(), 1)

	_, ok := restored.GetServiceEndpoints()[0].(*RpcRelayServiceEndpoint)
	assert.True(t, ok, "expected *RpcRelayServiceEndpoint after round-trip")
}

func TestUnitRegisteredNodeUpdateTransactionGetName(t *testing.T) {
	t.Parallel()

	tx := NewRegisteredNodeUpdateTransaction()
	assert.Equal(t, "RegisteredNodeUpdateTransaction", tx.getName())
}

func TestUnitRegisteredNodeUpdateTransactionGet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	key, _ := PrivateKeyGenerateEd25519()
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress: []byte{10, 0, 0, 1},
			port:      8080,
		},
		endpointApi: BlockNodeApiStatus,
	}

	transaction, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(1).
		SetAdminKey(key.PublicKey()).
		SetDescription("test").
		AddServiceEndpoint(endpoint).
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegisteredNodeId()
	transaction.GetAdminKey()
	transaction.GetDescription()
	transaction.GetServiceEndpoints()
}

func TestUnitRegisteredNodeUpdateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(0).
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegisteredNodeId()
	transaction.GetAdminKey()
	transaction.GetDescription()
	transaction.GetServiceEndpoints()
}

func TestUnitRegisteredNodeUpdateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	key, _ := PrivateKeyGenerateEd25519()
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	endpoint := &MirrorNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			domainName: "mirror.example.com",
			port:       443,
		},
	}

	transaction, err := NewRegisteredNodeUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetRegisteredNodeId(7).
		SetAdminKey(key.PublicKey()).
		SetDescription("test").
		AddServiceEndpoint(endpoint).
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetRegisteredNodeUpdate()
	require.Equal(t, uint64(7), proto.RegisteredNodeId)
	require.Equal(t, "test", proto.Description.Value)
	require.Equal(t, key.PublicKey()._ToProtoKey(), proto.AdminKey)
	require.Len(t, proto.ServiceEndpoint, 1)
}

func TestUnitRegisteredNodeUpdateTransactionMinimalRoundTrip(t *testing.T) {
	t.Parallel()

	tx := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(99)

	body := tx.buildProtoBody()
	assert.Equal(t, uint64(99), body.RegisteredNodeId)
	assert.Nil(t, body.Description)
	assert.Nil(t, body.AdminKey)
	assert.Empty(t, body.ServiceEndpoint)

	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_RegisteredNodeUpdate{
			RegisteredNodeUpdate: body,
		},
	}

	restored := _RegisteredNodeUpdateTransactionFromProtobuf(*tx.Transaction, pbBody)
	assert.Equal(t, uint64(99), restored.GetRegisteredNodeId())
	assert.Nil(t, restored.description)
}

func TestUnitRegisteredNodeUpdateTransactionFromProtobufWithEmptyDescription(t *testing.T) {
	t.Parallel()

	body := &services.RegisteredNodeUpdateTransactionBody{
		RegisteredNodeId: 5,
		Description:      wrapperspb.String(""),
	}

	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_RegisteredNodeUpdate{
			RegisteredNodeUpdate: body,
		},
	}

	tx := NewRegisteredNodeUpdateTransaction()
	restored := _RegisteredNodeUpdateTransactionFromProtobuf(*tx.Transaction, pbBody)
	assert.NotNil(t, restored.description)
	assert.Equal(t, "", restored.GetDescription())
}

func TestUnitRegisteredNodeUpdateTransactionCoverage(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	trx, err := NewRegisteredNodeUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetRegisteredNodeId(1).
		SetDescription("test").
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetTransactionMemo("no").
		SetTransactionValidDuration(time.Second * 30).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	trx.validateNetworkOnIDs(client)
	_, err = trx.Schedule()
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
	trx.GetMaxTransactionFee()
	trx.GetTransactionMemo()
	trx.GetRegenerateTransactionID()
	trx.GetRegisteredNodeId()
	trx.GetAdminKey()
	trx.GetDescription()
	trx.GetServiceEndpoints()
	_, err = trx.GetSignatures()
	require.NoError(t, err)
	trx.getName()
	switch b := txFromBytes.(type) {
	case RegisteredNodeUpdateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

func TestUnitRegisteredNodeUpdateTransactionFromToBytes(t *testing.T) {
	tx := NewRegisteredNodeUpdateTransaction().
		SetRegisteredNodeId(99).
		SetDescription("test")

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	assert.Equal(t, tx.buildProtoBody(), txFromBytes.(RegisteredNodeUpdateTransaction).buildProtoBody())
}
