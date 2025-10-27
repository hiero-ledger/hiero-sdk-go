//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitNodeCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tx := NewNodeCreateTransaction().
		SetAccountID(accountID)

	err = tx.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitNodeCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tx := NewNodeCreateTransaction().
		SetAccountID(accountID)

	err = tx.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func endpoints(offsets ...uint) []Endpoint {
	endpoints := make([]Endpoint, 0)

	for _, offset := range offsets {
		endpoints = append(endpoints, Endpoint{
			address: []byte{byte(offset), byte(offset), byte(offset), byte(offset)},
		})
	}

	return endpoints
}

func TestUnitNodeCreateTransactionMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		status.New(codes.Unavailable, "node is UNAVAILABLE").Err(),
		status.New(codes.Internal, "Received RST_STREAM with code 0").Err(),
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
						Status: services.ResponseCodeEnum_SUCCESS,
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
							AccountNum: 234,
						}},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewNodeCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAdminKey(newKey).
		SetDescription("test").
		SetGossipEndpoints(endpoints(0, 1, 2)).
		SetServiceEndpoints(endpoints(3, 4, 5)).
		SetGossipCaCertificate([]byte{111}).
		SetGrpcCertificateHash([]byte{222}).
		SetTransactionID(tran).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, receipt.AccountID, &AccountID{Account: 234})
}

func TestUnitNodeCreateTransactionGet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	key, err := PrivateKeyGenerateEd25519()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(key).
		SetTransactionMemo("").
		SetDescription("test").
		SetGossipEndpoints(endpoints(0, 1, 2)).
		SetServiceEndpoints(endpoints(3, 4, 5)).
		SetGossipCaCertificate([]byte{111}).
		SetGrpcCertificateHash([]byte{222}).
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
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetAccountID()
	transaction.GetDescription()
	transaction.GetGossipEndpoints()
	transaction.GetServiceEndpoints()
	transaction.GetGossipCaCertificate()
	transaction.GetGrpcCertificateHash()
	transaction.GetAdminKey()
}

func TestUnitNodeCreateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeCreateTransaction().
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
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetAccountID()
	transaction.GetDescription()
	transaction.GetGossipEndpoints()
	transaction.GetServiceEndpoints()
	transaction.GetGossipCaCertificate()
	transaction.GetGrpcCertificateHash()
	transaction.GetAdminKey()
}

func TestUnitNodeCreateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	stackedAccountID := AccountID{Account: 5}

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	gossipEndpoints := endpoints(1, 2, 3)
	serviceEndpoints := endpoints(3, 4, 5)
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAccountID(stackedAccountID).
		SetAdminKey(key).
		SetTransactionMemo("").
		SetDescription("test").
		SetGossipEndpoints(gossipEndpoints).
		SetServiceEndpoints(serviceEndpoints).
		SetGossipCaCertificate([]byte{111}).
		SetGrpcCertificateHash([]byte{222}).
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetNodeCreate()
	require.Equal(t, proto.AccountId.String(), stackedAccountID._ToProtobuf().String())
	require.Equal(t, proto.Description, "test")
	require.Equal(t, proto.GossipEndpoint[0], gossipEndpoints[0]._ToProtobuf())
	require.Equal(t, proto.ServiceEndpoint[0], serviceEndpoints[0]._ToProtobuf())
	require.Equal(t, proto.GossipCaCertificate, []byte{111})
	require.Equal(t, proto.GrpcCertificateHash, []byte{222})
	require.Equal(t, proto.AdminKey, key._ToProtoKey())
}

func TestUnitNodeCreateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	trx, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(key).
		SetAccountID(account).
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
	trx.GetAccountID()
	trx.GetDescription()
	trx.GetGossipEndpoints()
	trx.GetServiceEndpoints()
	trx.GetGossipCaCertificate()
	trx.GetGrpcCertificateHash()
	trx.GetAdminKey()
	_, err = trx.GetSignatures()
	require.NoError(t, err)
	trx.getName()
	switch b := txFromBytes.(type) {
	case NodeCreateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

func TestUnitNodeCreateTransactionFromToBytes(t *testing.T) {
	tx := NewNodeCreateTransaction()

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	assert.Equal(t, tx.buildProtoBody(), txFromBytes.(NodeCreateTransaction).buildProtoBody())
}

func TestUnitNodeCreateTransactionDeclineReward(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetDeclineReward(true).
		Freeze()
	require.NoError(t, err)

	assert.True(t, transaction.GetDeclineReward())

	proto := transaction.build().GetNodeCreate()
	assert.True(t, proto.DeclineReward)

	transaction2, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetDeclineReward(false).
		Freeze()
	require.NoError(t, err)

	assert.False(t, transaction2.GetDeclineReward())
	proto2 := transaction2.build().GetNodeCreate()
	assert.False(t, proto2.DeclineReward)
}

func TestUnitNodeCreateTransactionGrpcProxyEndpoint(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	proxyEndpoint := Endpoint{
		address: []byte{1, 2, 3, 4},
	}

	transaction, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetGrpcWebProxyEndpoint(proxyEndpoint).
		Freeze()
	require.NoError(t, err)

	gotEndpoint := transaction.GetGrpcWebProxyEndpoint()
	assert.Equal(t, proxyEndpoint.address, gotEndpoint.address)

	proto := transaction.build().GetNodeCreate()
	assert.NotNil(t, proto.GrpcProxyEndpoint)
	assert.Equal(t, proxyEndpoint._ToProtobuf().String(), proto.GrpcProxyEndpoint.String())

	transaction2, err := NewNodeCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	proto2 := transaction2.build().GetNodeCreate()
	assert.Nil(t, proto2.GrpcProxyEndpoint)
}

func TestUnitNodeCreateTransactionGossipEndpointsValidation(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	t.Run("TooManyGossipEndpoints", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Add 11 gossip endpoints (more than the limit of 10)
		for i := 0; i < 11; i++ {
			endpoint := Endpoint{}
			endpoint.SetAddress([]byte{192, 168, 1, byte(i + 1)})
			endpoint.SetPort(8080)
			tx.AddGossipEndpoint(endpoint)
		}

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errTooManyGossipEndpoints)
	})

	t.Run("GossipEndpointWithNoAddressOrDomainName", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Add an invalid endpoint (no address or domain name)
		invalidEndpoint := Endpoint{}
		tx.AddGossipEndpoint(invalidEndpoint)

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errEndpointMustHaveAddressOrDomainName)
	})

	t.Run("GossipEndpointWithBothAddressOrDomainName", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Add an invalid endpoint (no address or domain name)
		invalidEndpoint := Endpoint{}
		invalidEndpoint.SetAddress([]byte{192, 168, 1, 1})
		invalidEndpoint.SetDomainName("example.com")
		tx.AddGossipEndpoint(invalidEndpoint)

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errEndpointCannotHaveBothAddressAndDomainName)
	})

	t.Run("ValidGossipEndpoints", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Add valid endpoints (within limit and properly configured)
		for i := 0; i < 5; i++ {
			endpoint := Endpoint{}
			endpoint.SetAddress([]byte{192, 168, 1, byte(i + 1)})
			endpoint.SetPort(8080)
			tx.AddGossipEndpoint(endpoint)
		}

		// Should not error with valid endpoints
		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})
}

func TestUnitNodeCreateTransactionServiceEndpointsValidation(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	t.Run("TooManyServiceEndpoints", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		for i := 0; i < 9; i++ {
			endpoint := Endpoint{}
			endpoint.SetAddress([]byte{10, 0, 0, byte(i + 1)})
			endpoint.SetPort(9090)
			tx.AddServiceEndpoint(endpoint)
		}

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errTooManyServiceEndpoints)
	})

	t.Run("ServiceEndpointWithNoAddressOrDomainName", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Add an invalid endpoint (no address or domain name)
		invalidEndpoint := Endpoint{}
		tx.AddServiceEndpoint(invalidEndpoint)

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errEndpointMustHaveAddressOrDomainName)
	})

	t.Run("ServiceEndpointWithBothAddressOrDomainName", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Add an invalid endpoint (both address and domain name set)
		invalidEndpoint := Endpoint{}
		invalidEndpoint.SetAddress([]byte{172, 16, 0, 1})
		invalidEndpoint.SetDomainName("example.com")
		tx.AddServiceEndpoint(invalidEndpoint)

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errEndpointCannotHaveBothAddressAndDomainName)
	})

	t.Run("ValidServiceEndpoints", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Add valid endpoints (within limit and properly configured)
		for i := 0; i < 4; i++ {
			endpoint := Endpoint{}
			endpoint.SetDomainName("service" + string(rune('a'+i)) + ".example.com")
			endpoint.SetPort(9090)
			tx.AddServiceEndpoint(endpoint)
		}

		// Should not error with valid endpoints
		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})
}

func TestUnitNodeCreateTransactionGrpcWebProxyEndpointValidation(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	t.Run("GrpcWebProxyEndpointWithNoAddressOrDomainName", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Set an invalid gRPC web proxy endpoint
		invalidEndpoint := Endpoint{}
		// No address or domain name set, making it invalid
		tx.SetGrpcWebProxyEndpoint(invalidEndpoint)

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errEndpointMustHaveAddressOrDomainName)
	})

	t.Run("GrpcWebProxyEndpointWithBothAddressOrDomainName", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Set an invalid gRPC web proxy endpoint
		invalidEndpoint := Endpoint{}
		invalidEndpoint.SetAddress([]byte{192, 168, 1, 1})
		invalidEndpoint.SetDomainName("example.com")
		tx.SetGrpcWebProxyEndpoint(invalidEndpoint)
		_, err := tx.Freeze()
		require.Error(t, err)
		assert.ErrorIs(t, err, errEndpointCannotHaveBothAddressAndDomainName)
	})

	t.Run("ValidGrpcWebProxyEndpoint", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Set a valid gRPC web proxy endpoint
		validEndpoint := Endpoint{}
		validEndpoint.SetDomainName("proxy.example.com")
		validEndpoint.SetPort(8080)
		tx.SetGrpcWebProxyEndpoint(validEndpoint)

		// Should not error with valid endpoint
		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})

	t.Run("NoGrpcWebProxyEndpoint", func(t *testing.T) {
		nodeAccountID := []AccountID{{Account: 10}}
		transactionID := TransactionIDGenerate(AccountID{Account: 324})
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Don't set gRPC web proxy endpoint (should be fine)
		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})
}

func TestUnitNodeCreateTransactionGossipCaCertificateValidation(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	t.Run("EmptyGossipCaCertificate", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Set an empty gossip CA certificate
		tx.SetGossipCaCertificate([]byte{})

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errGossipCaCertificateEmpty)
	})

	t.Run("ValidGossipCaCertificate", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Set a valid gossip CA certificate
		validCert := []byte{0x30, 0x82, 0x01, 0x22}
		tx.SetGossipCaCertificate(validCert)

		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})

	t.Run("NoGossipCaCertificate", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Don't set gossip CA certificate (should be fine)
		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})
}

func TestUnitNodeCreateTransactionDescriptionValidation(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	t.Run("DescriptionTooLong", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Set a description longer than 100 characters
		longDescription := strings.Repeat("a", 101)
		tx.SetDescription(longDescription)

		_, err := tx.Freeze()

		require.Error(t, err)
		assert.ErrorIs(t, err, errDescriptionTooLong)
	})

	t.Run("DescriptionAtLimit", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		// Set a description exactly 100 characters
		limitDescription := strings.Repeat("a", 100)
		tx.SetDescription(limitDescription)

		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})

	t.Run("ValidDescription", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		tx.SetDescription("Valid node description")

		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})

	t.Run("EmptyDescription", func(t *testing.T) {
		tx := NewNodeCreateTransaction().
			SetTransactionID(transactionID).
			SetNodeAccountIDs(nodeAccountID)

		tx.SetDescription("")

		frozen, err := tx.Freeze()
		require.NoError(t, err)
		assert.NotNil(t, frozen)
	})
}
