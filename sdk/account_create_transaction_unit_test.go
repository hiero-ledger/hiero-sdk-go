//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountCreateTransactionMock(t *testing.T) {
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

	newBalance := NewHbar(2)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetKeyWithoutAlias(newKey).
		SetTransactionID(tran).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, receipt.AccountID, &AccountID{Account: 234})
}

func TestUnitAccountCreateTransactionGet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	key, err := PrivateKeyGenerateEd25519()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKeyWithoutAlias(key).
		SetAccountMemo("").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		AddHook(*NewHookCreationDetails()).
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
	transaction.GetAccountMemo()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetInitialBalance()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
	transaction.GetHooks()
}

func TestUnitAccountCreateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountCreateTransaction().
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
	transaction.GetAccountMemo()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetInitialBalance()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
}

func TestUnitAccountCreateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	stackedAccountID := AccountID{Account: 5}

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	alias := "5c562e90feaf0eebd33ea75d21024f249d451417"

	transaction, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKeyWithoutAlias(key).
		SetInitialBalance(NewHbar(3)).
		SetAccountMemo("ty").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetStakedAccountID(stackedAccountID).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(60 * time.Second).
		AddHook(*NewHookCreationDetails().SetHookId(1).SetAdminKey(key)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetAlias(alias).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetCryptoCreateAccount()
	require.Equal(t, proto.Key.String(), key._ToProtoKey().String())
	require.Equal(t, proto.InitialBalance, uint64(NewHbar(3).AsTinybar()))
	require.Equal(t, proto.Memo, "ty")
	require.Equal(t, proto.ReceiverSigRequired, true)
	require.Equal(t, proto.MaxAutomaticTokenAssociations, int32(2))
	require.Equal(t, proto.StakedId.(*services.CryptoCreateTransactionBody_StakedAccountId).StakedAccountId.String(),
		stackedAccountID._ToProtobuf().String())
	require.Equal(t, proto.DeclineReward, true)
	require.Equal(t, proto.AutoRenewPeriod.String(), _DurationToProtobuf(60*time.Second).String())
	require.Equal(t, hex.EncodeToString(proto.Alias), alias)
	require.Len(t, proto.HookCreationDetails, 1)
	require.Equal(t, proto.HookCreationDetails[0].HookId, int64(1))
	require.Equal(t, proto.HookCreationDetails[0].AdminKey, key._ToProtoKey())
}

func TestUnitAccountCreateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	alias := "5c562e90feaf0eebd33ea75d21024f249d451417"

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	trx, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKeyWithoutAlias(key).
		SetInitialBalance(NewHbar(3)).
		SetAccountMemo("ty").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetStakedAccountID(account).
		SetStakedNodeID(4).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetTransactionMemo("no").
		SetTransactionValidDuration(time.Second * 30).
		SetRegenerateTransactionID(false).
		SetAlias(alias).
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
	trx.GetStakedAccountID()
	trx.GetStakedNodeID()
	trx.GetDeclineStakingReward()
	trx.GetAlias()
	_, err = trx.GetSignatures()
	require.NoError(t, err)
	trx.getName()
	switch b := txFromBytes.(type) {
	case AccountCreateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

func TestUnitAccountCreateSetStakedNodeID(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	tx := NewAccountCreateTransaction().
		SetStakedAccountID(account).
		SetStakedNodeID(4)

	require.Equal(t, AccountID{}, tx.GetStakedAccountID())
	require.Equal(t, int64(4), tx.GetStakedNodeID())
}
func TestUnitAccountCreateSetStakedAccountID(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	tx := NewAccountCreateTransaction().
		SetStakedNodeID(4).
		SetStakedAccountID(account)

	require.Equal(t, int64(0), tx.GetStakedNodeID())
	require.Equal(t, account, tx.GetStakedAccountID())
}

func TestUnitAccountCreateTransactionFromToBytes(t *testing.T) {
	tx := NewAccountCreateTransaction()

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	assert.Equal(t, tx.buildProtoBody(), txFromBytes.(AccountCreateTransaction).buildProtoBody())
}

type invalidKey struct{}

func (k invalidKey) _ToProtoKey() *services.Key {
	return nil
}

func (k invalidKey) String() string {
	return "invalidKey"
}

func TestUnitAccountCreateSetECDSAKeyWithAliasInvalidKey(t *testing.T) {
	t.Parallel()

	// Test with invalid key
	ecdsaPrivateKey := invalidKey{}
	tx := NewAccountCreateTransaction()
	tx.SetECDSAKeyWithAlias(ecdsaPrivateKey)
	require.Error(t, tx.freezeError)
	tx.SetKeyWithAlias(ecdsaPrivateKey, ecdsaPrivateKey)
	require.Error(t, tx.freezeError)
}

func TestUnitAccountCreateSetECDSAKeyWithAlias(t *testing.T) {
	t.Parallel()

	// Test with ECDSA private key
	ecdsaPrivateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	expectedEvmAddress := ecdsaPrivateKey.PublicKey().ToEvmAddress()

	tx := NewAccountCreateTransaction()
	tx.SetECDSAKeyWithAlias(ecdsaPrivateKey)
	require.NoError(t, tx.freezeError)
	key, err := tx.GetKey()
	require.NoError(t, err)
	require.Equal(t, ecdsaPrivateKey, key)

	// Verify the alias is set correctly
	evmAddressBytes, err := hex.DecodeString(strings.TrimPrefix(expectedEvmAddress, "0x"))
	require.NoError(t, err)
	require.Equal(t, evmAddressBytes, tx.GetAlias())

	// Test with ECDSA public key
	ecdsaPublicKey := ecdsaPrivateKey.PublicKey()
	tx = NewAccountCreateTransaction()
	tx.SetECDSAKeyWithAlias(ecdsaPublicKey)
	require.NoError(t, tx.freezeError)
	key, err = tx.GetKey()
	require.NoError(t, err)
	require.Equal(t, ecdsaPublicKey, key)
	require.Equal(t, evmAddressBytes, tx.GetAlias())

	// Test with non-ECDSA private key
	edKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx = NewAccountCreateTransaction()
	tx.SetECDSAKeyWithAlias(edKey)
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Private key is not ECDSA")

	// Test with non-ECDSA public key
	edPublicKey := edKey.PublicKey()
	tx = NewAccountCreateTransaction()
	tx.SetECDSAKeyWithAlias(edPublicKey)
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Public key is not ECDSA")
}

func TestUnitAccountCreateSetKeyWithAlias(t *testing.T) {
	t.Parallel()

	// Generate test keys
	edKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ecdsaKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	expectedEvmAddress := ecdsaKey.PublicKey().ToEvmAddress()

	// Test with private keys
	tx := NewAccountCreateTransaction()
	tx.SetKeyWithAlias(edKey, ecdsaKey)
	require.NoError(t, tx.freezeError)
	key, err := tx.GetKey()
	require.NoError(t, err)
	require.Equal(t, edKey, key)

	// Verify the alias is set correctly
	evmAddressBytes, err := hex.DecodeString(strings.TrimPrefix(expectedEvmAddress, "0x"))
	require.NoError(t, err)
	require.Equal(t, evmAddressBytes, tx.GetAlias())

	// Test with public key for alias
	tx = NewAccountCreateTransaction()
	tx.SetKeyWithAlias(edKey, ecdsaKey.PublicKey())
	require.NoError(t, tx.freezeError)
	key, err = tx.GetKey()
	require.NoError(t, err)
	require.Equal(t, edKey, key)

	// Test with non-ECDSA private key for alias
	tx = NewAccountCreateTransaction()
	tx.SetKeyWithAlias(edKey, edKey)
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Private key is not ECDSA")

	// Test with non-ECDSA public key for alias
	tx = NewAccountCreateTransaction()
	tx.SetKeyWithAlias(edKey, edKey.PublicKey())
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Public key is not ECDSA")
}

func TestUnitAccountCreateSetHooks(t *testing.T) {
	tx := NewAccountCreateTransaction()

	hook1 := NewHookCreationDetails()
	hook2 := NewHookCreationDetails()

	tx.AddHook(*hook1)
	tx.AddHook(*hook2)

	require.Equal(t, 2, len(tx.GetHooks()))
	require.Equal(t, *hook1, tx.GetHooks()[0])
	require.Equal(t, *hook2, tx.GetHooks()[1])

	tx.SetHooks([]HookCreationDetails{*hook1, *hook2})
	require.Equal(t, 2, len(tx.GetHooks()))
	require.Equal(t, *hook1, tx.GetHooks()[0])
	require.Equal(t, *hook2, tx.GetHooks()[1])
}

func TestUnitAccountCreateToProtoHooks(t *testing.T) {
	tx := NewAccountCreateTransaction()
	proto := tx.buildProtoBody()
	require.Equal(t, 0, len(proto.HookCreationDetails))

	hook := NewHookCreationDetails()
	tx.AddHook(*hook)
	proto = tx.buildProtoBody()
	require.Equal(t, 1, len(proto.HookCreationDetails))
	require.Equal(t, hook.toProtobuf(), proto.HookCreationDetails[0])
}

func TestUnitAccountCreateBytesHooks(t *testing.T) {
	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)
	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()

	hook := NewHookCreationDetails().SetHookId(1).
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetEvmHook(*NewEvmHook().SetContractId(&contractID)).
		SetAdminKey(ed25519PublicKey)
	tx := NewAccountCreateTransaction().SetHooks([]HookCreationDetails{*hook})
	byt, err := tx.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	accountCreateTx := txFromBytes.(AccountCreateTransaction)
	require.NoError(t, err)
	require.Equal(t, *hook, accountCreateTx.GetHooks()[0])
}

func TestUnitAccountCreateSetDelegationAddress(t *testing.T) {
	t.Parallel()

	// Test with hex string with 0x prefix
	delegationAddr1 := "0x1111111111111111111111111111111111111111"
	delegationAddrBytes1, err := hex.DecodeString("1111111111111111111111111111111111111111")
	require.NoError(t, err)

	tx := NewAccountCreateTransaction()
	tx.SetDelegationAddress(delegationAddr1)
	require.NoError(t, tx.freezeError)

	retrievedAddr := tx.GetDelegationAddress()
	require.NotNil(t, retrievedAddr)
	require.Equal(t, delegationAddrBytes1, retrievedAddr)

	// Test with hex string without 0x prefix
	delegationAddr2 := "2222222222222222222222222222222222222222"
	delegationAddrBytes2, err := hex.DecodeString(delegationAddr2)
	require.NoError(t, err)

	tx = NewAccountCreateTransaction()
	tx.SetDelegationAddress(delegationAddr2)
	require.NoError(t, tx.freezeError)

	retrievedAddr = tx.GetDelegationAddress()
	require.NotNil(t, retrievedAddr)
	require.Equal(t, delegationAddrBytes2, retrievedAddr)

	// Test with bytes
	delegationAddrBytes3 := []byte{0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
		0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33}

	tx = NewAccountCreateTransaction()
	tx.SetDelegationAddress(delegationAddrBytes3)
	require.NoError(t, tx.freezeError)

	retrievedAddr = tx.GetDelegationAddress()
	require.NotNil(t, retrievedAddr)
	require.Equal(t, delegationAddrBytes3, retrievedAddr)

	// Test without delegation address (should return nil)
	tx = NewAccountCreateTransaction()
	retrievedAddr = tx.GetDelegationAddress()
	require.Nil(t, retrievedAddr)
}

func TestUnitAccountCreateSetDelegationAddressInvalid(t *testing.T) {
	t.Parallel()

	// Test with invalid hex string (wrong length)
	invalidAddr := "0x12345"
	tx := NewAccountCreateTransaction()
	tx.SetDelegationAddress(invalidAddr)
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Invalid delegation address format")

	// Test with invalid bytes (wrong size)
	invalidBytes := []byte{0x01, 0x02, 0x03} // Only 3 bytes, should be 20
	tx = NewAccountCreateTransaction()
	tx.SetDelegationAddress(invalidBytes)
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Delegation address must be exactly 20 bytes")

	// Test with invalid type
	tx = NewAccountCreateTransaction()
	tx.SetDelegationAddress(12345) // Invalid type
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Delegation address must be a string or []byte")
}

func TestUnitAccountCreateDelegationAddressProto(t *testing.T) {
	t.Parallel()

	delegationAddr := "0x4444444444444444444444444444444444444444"
	delegationAddrBytes, err := hex.DecodeString("4444444444444444444444444444444444444444")
	require.NoError(t, err)

	tx := NewAccountCreateTransaction()
	tx.SetDelegationAddress(delegationAddr)

	proto := tx.buildProtoBody()
	require.Equal(t, delegationAddrBytes, proto.DelegationAddress)

	// Test without delegation address
	tx2 := NewAccountCreateTransaction()
	proto2 := tx2.buildProtoBody()
	require.Nil(t, proto2.DelegationAddress)
}

func TestUnitAccountCreateDelegationAddressBytes(t *testing.T) {
	t.Parallel()

	delegationAddr := "0x5555555555555555555555555555555555555555"
	delegationAddrBytes, err := hex.DecodeString("5555555555555555555555555555555555555555")
	require.NoError(t, err)

	tx := NewAccountCreateTransaction().SetDelegationAddress(delegationAddr)
	byt, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)

	accountCreateTx := txFromBytes.(AccountCreateTransaction)
	require.Equal(t, delegationAddrBytes, accountCreateTx.GetDelegationAddress())
}
