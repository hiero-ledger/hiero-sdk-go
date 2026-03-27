//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registeredNodeDeleteMockResponses() [][]interface{} {
	return [][]interface{}{{
		grpcStatus.New(codes.Unavailable, "node is UNAVAILABLE").Err(),
		grpcStatus.New(codes.Internal, "Received RST_STREAM with code 0").Err(),
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
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
						Status: services.ResponseCodeEnum_RECEIPT_NOT_FOUND,
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
						RegisteredNodeId: 42,
					},
				},
			},
		},
	}}
}

func TestUnitRegisteredNodeDeleteTransactionMock(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer(registeredNodeDeleteMockResponses())
	defer server.Close()

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(42).
		SetTransactionID(tran).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, uint64(42), receipt.RegisteredNodeId)
}

func TestUnitRegisteredNodeDeleteTransactionBuild(t *testing.T) {
	t.Parallel()

	tx := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(5)

	body := tx.buildProtoBody()
	assert.Equal(t, uint64(5), body.RegisteredNodeId)
}

func TestUnitRegisteredNodeDeleteTransactionRoundTrip(t *testing.T) {
	t.Parallel()

	tx := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(42)

	body := tx.buildProtoBody()
	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_RegisteredNodeDelete{
			RegisteredNodeDelete: body,
		},
	}

	restored := _RegisteredNodeDeleteTransactionFromProtobuf(*tx.Transaction, pbBody)
	assert.Equal(t, uint64(42), restored.GetRegisteredNodeId())
}

func TestUnitRegisteredNodeDeleteTransactionGetName(t *testing.T) {
	t.Parallel()

	tx := NewRegisteredNodeDeleteTransaction()
	assert.Equal(t, "RegisteredNodeDeleteTransaction", tx.getName())
}

func TestUnitRegisteredNodeDeleteTransactionGet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(1).
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
}

func TestUnitRegisteredNodeDeleteTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewRegisteredNodeDeleteTransaction().
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
}

func TestUnitRegisteredNodeDeleteTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewRegisteredNodeDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetRegisteredNodeId(7).
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetRegisteredNodeDelete()
	require.Equal(t, proto.RegisteredNodeId, uint64(7))
}

func buildFrozenRegisteredNodeDeleteTransaction(t *testing.T) (*RegisteredNodeDeleteTransaction, *Client, PrivateKey) {
	t.Helper()
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	trx, err := NewRegisteredNodeDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetRegisteredNodeId(1).
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

	return trx, client, key
}

func TestUnitRegisteredNodeDeleteTransactionCoverage(t *testing.T) {
	t.Parallel()

	trx, client, key := buildFrozenRegisteredNodeDeleteTransaction(t)

	trx.validateNetworkOnIDs(client)
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
	trx.GetMaxTransactionFee()
	trx.GetTransactionMemo()
	trx.GetRegenerateTransactionID()
	trx.GetRegisteredNodeId()
	_, err = trx.GetSignatures()
	require.NoError(t, err)
	trx.getName()
	switch b := txFromBytes.(type) {
	case RegisteredNodeDeleteTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

func TestUnitRegisteredNodeDeleteTransactionFromToBytes(t *testing.T) {
	tx := NewRegisteredNodeDeleteTransaction().
		SetRegisteredNodeId(99)

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	assert.Equal(t, tx.buildProtoBody(), txFromBytes.(RegisteredNodeDeleteTransaction).buildProtoBody())
}

