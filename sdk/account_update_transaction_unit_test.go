//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountUpdateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	accountUpdate := NewAccountUpdateTransaction().
		SetProxyAccountID(accountID)

	err = accountUpdate.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountUpdateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	accountUpdate := NewAccountUpdateTransaction().
		SetProxyAccountID(accountID)

	err = accountUpdate.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountUpdateTransactionMock(t *testing.T) {
	t.Parallel()

	call := func(request *services.Transaction) *services.TransactionResponse {
		require.NotEmpty(t, request.SignedTransactionBytes)
		signedTransaction := services.SignedTransaction{}
		_ = protobuf.Unmarshal(request.SignedTransactionBytes, &signedTransaction)

		require.NotEmpty(t, signedTransaction.BodyBytes)
		transactionBody := services.TransactionBody{}
		_ = protobuf.Unmarshal(signedTransaction.BodyBytes, &transactionBody)

		require.NotNil(t, transactionBody.TransactionID)
		transactionId := transactionBody.TransactionID.String()
		require.NotEqual(t, "", transactionId)

		sigMap := signedTransaction.GetSigMap()
		require.NotNil(t, sigMap)

		for _, sigPair := range sigMap.SigPair {
			verified := false

			switch k := sigPair.Signature.(type) {
			case *services.SignaturePair_Ed25519:
				pbTemp, _ := PublicKeyFromBytesEd25519(sigPair.PubKeyPrefix)
				verified = pbTemp.VerifySignedMessage(signedTransaction.BodyBytes, k.Ed25519)
			case *services.SignaturePair_ECDSASecp256K1:
				pbTemp, _ := PublicKeyFromBytesECDSA(sigPair.PubKeyPrefix)
				verified = pbTemp.VerifySignedMessage(signedTransaction.BodyBytes, k.ECDSASecp256K1)
			}
			require.True(t, verified)
		}

		if bod, ok := transactionBody.Data.(*services.TransactionBody_CryptoUpdateAccount); ok {
			require.Equal(t, bod.CryptoUpdateAccount.Memo.Value, "no")
			require.Equal(t, bod.CryptoUpdateAccount.AccountIDToUpdate.GetAccountNum(), int64(123))
			//alias := services.Key{}
			//_ = protobuf.Unmarshal(bod.CryptoUpdateAccount.Alias, &alias)
			//require.Equal(t, hex.EncodeToString(alias.GetEd25519()), "1480272863d39c42f902bc11601a968eaf30ad662694e3044c86d5df46fabfd2")
		}

		return &services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		}
	}
	responses := [][]interface{}{{
		call,
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()
	//302a300506032b65700321001480272863d39c42f902bc11601a968eaf30ad662694e3044c86d5df46fabfd2
	newKey, err := PrivateKeyFromStringEd25519("302e020100300506032b657004220420278184257eb568d0e5fcfc1df99828b039b4776da05855dc5af105996e6200d1")
	require.NoError(t, err)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	_, err = NewAccountUpdateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(tran).
		SetAccountMemo("no").
		SetAccountID(AccountID{Account: 123}).
		SetAliasKey(newKey.PublicKey()).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitAccountUpdateTransactionGet(t *testing.T) {
	t.Parallel()

	spenderAccountID1 := AccountID{Account: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	key, err := PrivateKeyGenerateEd25519()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAccountID(spenderAccountID1).
		SetKey(key).
		SetProxyAccountID(spenderAccountID1).
		SetAccountMemo("").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAccountID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
}

func TestUnitAccountUpdateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAccountID()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
}

func TestUnitAccountUpdateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	stackedAccountID := AccountID{Account: 5}
	accountID := AccountID{Account: 6}

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKey(key).
		SetAccountID(accountID).
		SetAccountMemo("ty").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetStakedAccountID(stackedAccountID).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(60 * time.Second).
		SetExpirationTime(time.Unix(34, 56)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetCryptoUpdateAccount()
	require.Equal(t, proto.AccountIDToUpdate.String(), accountID._ToProtobuf().String())
	require.Equal(t, proto.Key.String(), key._ToProtoKey().String())
	require.Equal(t, proto.Memo.Value, "ty")
	require.Equal(t, proto.ReceiverSigRequiredField.(*services.CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper).ReceiverSigRequiredWrapper.Value, true)
	require.Equal(t, proto.MaxAutomaticTokenAssociations.GetValue(), int32(2))
	require.Equal(t, proto.StakedId.(*services.CryptoUpdateTransactionBody_StakedAccountId).StakedAccountId.String(),
		stackedAccountID._ToProtobuf().String())
	require.Equal(t, proto.DeclineReward.Value, true)
	require.Equal(t, proto.AutoRenewPeriod.String(), _DurationToProtobuf(60*time.Second).String())
	require.Equal(t, proto.ExpirationTime.String(), _TimeToProtobuf(time.Unix(34, 56)).String())
}

func TestUnitAccountUpdateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	grpc := time.Second * 30
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	transaction, err := NewAccountUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKey(key).
		SetAccountID(account).
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
		SetGrpcDeadline(&grpc).
		Freeze()
	require.NoError(t, err)

	transaction.validateNetworkOnIDs(client)

	_, err = transaction.Schedule()
	require.NoError(t, err)
	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()
	transaction.GetMaxRetry()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxBackoff()
	transaction.GetMinBackoff()
	transaction.GetRegenerateTransactionID()
	byt, err := transaction.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)
	sig, err := key.SignTransaction(transaction)
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	transaction.GetStakedAccountID()
	transaction.GetStakedNodeID()
	transaction.ClearStakedAccountID()
	transaction.ClearStakedNodeID()
	transaction.GetDeclineStakingReward()
	transaction.GetExpirationTime()
	transaction.GetAccountMemo()

	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.getName()
	switch b := txFromBytes.(type) {
	case *AccountUpdateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

func TestUnitAccountUpdateTransactionFromToBytes(t *testing.T) {
	tx := NewAccountUpdateTransaction()

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	assert.Equal(t, tx.buildProtoBody(), txFromBytes.(AccountUpdateTransaction).buildProtoBody())
}

func TestUnitAccountUpdateSetHooks(t *testing.T) {
	tx := NewAccountUpdateTransaction()

	hook1 := NewHookCreationDetails()
	hook2 := NewHookCreationDetails()

	tx.AddHookToCreate(*hook1).AddHookToCreate(*hook2)
	require.Equal(t, 2, len(tx.GetHooksToCreate()))
	require.Equal(t, *hook1, tx.GetHooksToCreate()[0])
	require.Equal(t, *hook2, tx.GetHooksToCreate()[1])

	tx.SetHooksToCreate([]HookCreationDetails{*hook1, *hook2})
	require.Equal(t, 2, len(tx.GetHooksToCreate()))
	require.Equal(t, *hook1, tx.GetHooksToCreate()[0])
	require.Equal(t, *hook2, tx.GetHooksToCreate()[1])

	tx.AddHookToDelete(1)
	require.Equal(t, []int64{1}, tx.GetHooksToDelete())

	tx.SetHooksToDelete([]int64{1, 2})
	require.Equal(t, []int64{1, 2}, tx.GetHooksToDelete())
}

func TestUnitAccountUpdateTransactionToProtoHooks(t *testing.T) {
	tx := NewAccountUpdateTransaction()
	proto := tx.buildProtoBody()
	require.Equal(t, 0, len(proto.HookCreationDetails))
	require.Equal(t, 0, len(proto.HookIdsToDelete))

	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()

	hook := NewHookCreationDetails().
		SetHookId(1).
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetEvmHook(*NewEvmHook().SetContractId(&ContractID{Contract: 1})).
		SetAdminKey(ed25519PublicKey)

	tx.AddHookToCreate(*hook).AddHookToDelete(1)
	proto = tx.buildProtoBody()
	require.Equal(t, hook.toProtobuf(), proto.HookCreationDetails[0])
	require.Equal(t, []int64{1}, proto.HookIdsToDelete)
}

func TestUnitAccountUpdateTransactionBytesHooks(t *testing.T) {
	tx := NewAccountUpdateTransaction()
	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)
	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()
	hook := NewHookCreationDetails().
		SetHookId(1).SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetEvmHook(*NewEvmHook().SetContractId(&contractID)).
		SetAdminKey(ed25519PublicKey)
	tx.AddHookToCreate(*hook)
	tx.AddHookToDelete(1)
	txBytes, err := tx.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	accountUpdateTx := txFromBytes.(AccountUpdateTransaction)
	require.Equal(t, *hook, accountUpdateTx.GetHooksToCreate()[0])
	require.Equal(t, []int64{1}, accountUpdateTx.GetHooksToDelete())
}

func TestUnitAccountUpdateSetDelegationAddress(t *testing.T) {
	t.Parallel()

	// Test with hex string with 0x prefix
	delegationAddr1 := "0x1111111111111111111111111111111111111111"
	delegationAddrBytes1, err := hex.DecodeString("1111111111111111111111111111111111111111")
	require.NoError(t, err)

	tx := NewAccountUpdateTransaction()
	tx.SetDelegationAddress(delegationAddr1)
	require.NoError(t, tx.freezeError)

	retrievedAddr := tx.GetDelegationAddress()
	require.NotNil(t, retrievedAddr)
	require.Equal(t, delegationAddrBytes1, retrievedAddr)

	// Test with hex string without 0x prefix
	delegationAddr2 := "2222222222222222222222222222222222222222"
	delegationAddrBytes2, err := hex.DecodeString(delegationAddr2)
	require.NoError(t, err)

	tx = NewAccountUpdateTransaction()
	tx.SetDelegationAddress(delegationAddr2)
	require.NoError(t, tx.freezeError)

	retrievedAddr = tx.GetDelegationAddress()
	require.NotNil(t, retrievedAddr)
	require.Equal(t, delegationAddrBytes2, retrievedAddr)

	// Test with bytes
	delegationAddrBytes3 := []byte{0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
		0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33}

	tx = NewAccountUpdateTransaction()
	tx.SetDelegationAddress(delegationAddrBytes3)
	require.NoError(t, tx.freezeError)

	retrievedAddr = tx.GetDelegationAddress()
	require.NotNil(t, retrievedAddr)
	require.Equal(t, delegationAddrBytes3, retrievedAddr)

	// Test with nil (clears delegation)
	tx = NewAccountUpdateTransaction()
	tx.SetDelegationAddress(nil)
	require.NoError(t, tx.freezeError)

	retrievedAddr = tx.GetDelegationAddress()
	require.Nil(t, retrievedAddr)

	// Test with zero address (clears delegation)
	zeroAddress := "0x0000000000000000000000000000000000000000"
	tx = NewAccountUpdateTransaction()
	tx.SetDelegationAddress(zeroAddress)
	require.NoError(t, tx.freezeError)

	retrievedAddr = tx.GetDelegationAddress()
	require.Nil(t, retrievedAddr)

	// Test without delegation address (should return nil)
	tx = NewAccountUpdateTransaction()
	retrievedAddr = tx.GetDelegationAddress()
	require.Nil(t, retrievedAddr)
}

func TestUnitAccountUpdateSetDelegationAddressInvalid(t *testing.T) {
	t.Parallel()

	// Test with invalid hex string (wrong length)
	invalidAddr := "0x12345"
	tx := NewAccountUpdateTransaction()
	tx.SetDelegationAddress(invalidAddr)
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Invalid delegation address format")

	// Test with invalid bytes (wrong size)
	invalidBytes := []byte{0x01, 0x02, 0x03} // Only 3 bytes, should be 20
	tx = NewAccountUpdateTransaction()
	tx.SetDelegationAddress(invalidBytes)
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Delegation address must be exactly 20 bytes")

	// Test with invalid type
	tx = NewAccountUpdateTransaction()
	tx.SetDelegationAddress(12345) // Invalid type
	require.Error(t, tx.freezeError)
	require.Contains(t, tx.freezeError.Error(), "Delegation address must be a string, []byte, or nil")
}

func TestUnitAccountUpdateDelegationAddressProto(t *testing.T) {
	t.Parallel()

	delegationAddr := "0x4444444444444444444444444444444444444444"
	delegationAddrBytes, err := hex.DecodeString("4444444444444444444444444444444444444444")
	require.NoError(t, err)

	tx := NewAccountUpdateTransaction()
	tx.SetDelegationAddress(delegationAddr)

	proto := tx.buildProtoBody()
	require.Equal(t, delegationAddrBytes, proto.DelegationAddress)

	// Test without delegation address (should set empty slice to clear)
	tx2 := NewAccountUpdateTransaction()
	proto2 := tx2.buildProtoBody()
	require.Equal(t, []byte{}, proto2.DelegationAddress)

	// Test with nil (should set empty slice to clear)
	tx3 := NewAccountUpdateTransaction()
	tx3.SetDelegationAddress(nil)
	proto3 := tx3.buildProtoBody()
	require.Equal(t, []byte{}, proto3.DelegationAddress)
}

func TestUnitAccountUpdateDelegationAddressBytes(t *testing.T) {
	t.Parallel()

	delegationAddr := "0x5555555555555555555555555555555555555555"
	delegationAddrBytes, err := hex.DecodeString("5555555555555555555555555555555555555555")
	require.NoError(t, err)

	tx := NewAccountUpdateTransaction().SetDelegationAddress(delegationAddr)
	byt, err := tx.ToBytes()
	require.NoError(t, err)

	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)

	accountUpdateTx := txFromBytes.(AccountUpdateTransaction)
	require.Equal(t, delegationAddrBytes, accountUpdateTx.GetDelegationAddress())
}
