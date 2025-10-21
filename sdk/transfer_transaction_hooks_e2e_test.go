//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTransferHbarWithPreTransactionAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	accountId := *receipt.AccountID

	hookCall := NewFungibleHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK)

	transferResponse, err := NewTransferTransaction().
		AddHbarTransfer(env.OperatorID, NewHbar(1)).
		AddHbarTransferWithHook(accountId, NewHbar(-1), *hookCall).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationMultipleAccountsHooksMustAllApprove(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	// Two different hook ids for two different accounts
	hookDetails1 := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	hookDetails2 := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	// Create two recipient accounts, each with its own hook
	key1, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp1, err := NewAccountCreateTransaction().
		SetKey(key1).
		SetInitialBalance(NewHbar(1)).
		AddHook(*hookDetails1).
		Execute(env.Client)
	require.NoError(t, err)

	receipt1, err := resp1.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	acct1 := *receipt1.AccountID

	key2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp2, err := NewAccountCreateTransaction().
		SetKey(key2).
		SetInitialBalance(NewHbar(1)).
		AddHook(*hookDetails2).
		Execute(env.Client)
	require.NoError(t, err)

	receipt2, err := resp2.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	acct2 := *receipt2.AccountID

	hookCall1 := NewFungibleHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK)
	hookCall2 := NewFungibleHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.OperatorID, NewHbar(2)).
		AddHbarTransferWithHook(acct1, NewHbar(-1), *hookCall1).
		AddHbarTransferWithHook(acct2, NewHbar(-1), *hookCall2).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, receipt.Status)
}

func TestIntegrationTransferFungibleTokenWithPreTransactionAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	tokenId, err := createFungibleToken(&env)
	require.NoError(t, err)

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	accountId := *receipt.AccountID

	associateTx, err := NewTokenAssociateTransaction().
		SetAccountID(accountId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTx.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		AddTokenTransfer(tokenId, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenId, accountId, 1000).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	hookCall := NewFungibleHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK)

	transferResponse, err := NewTransferTransaction().
		AddTokenTransfer(tokenId, env.Client.GetOperatorAccountID(), 100).
		AddTokenTransferWithHook(tokenId, accountId, -100, *hookCall).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationTransferNftWithPreTransactionAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	tokenId, err := createNft(&env)
	require.NoError(t, err)

	mintResp, err := NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadatas([][]byte{{1}, {2}, {3}}).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mintResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(3), mintReceipt.SerialNumbers[2])

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	accountId := *receipt.AccountID

	associateTx, err := NewTokenAssociateTransaction().
		SetAccountID(accountId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTx.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		AddNftTransfer(tokenId.Nft(1), env.Client.GetOperatorAccountID(), accountId).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	hookCall := NewNftHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK_SENDER)

	transferResponse, err := NewTransferTransaction().
		AddNftTransferWitHook(tokenId.Nft(1), accountId, env.Client.GetOperatorAccountID(), hookCall, nil).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationTransferNftWithReceiverAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	tokenId, err := createNft(&env)
	require.NoError(t, err)

	mintResp, err := NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadatas([][]byte{{1}, {2}}).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mintResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(2), mintReceipt.SerialNumbers[1])

	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(receiverKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	receiverId := *receipt.AccountID

	associateTx, err := NewTokenAssociateTransaction().
		SetAccountID(receiverId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	receiverHookCall := NewNftHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK_RECEIVER)

	transferResponse, err := NewTransferTransaction().
		AddNftTransfer(tokenId.Nft(1), env.Client.GetOperatorAccountID(), receiverId).
		AddNftTransferWitHook(tokenId.Nft(2), env.Client.GetOperatorAccountID(), receiverId, nil, receiverHookCall).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationTransferNftWithBothSenderAndReceiverHooksSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	tokenId, err := createNft(&env)
	require.NoError(t, err)

	mintResp, err := NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadatas([][]byte{{1}, {2}}).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mintResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.Equal(t, int64(2), mintReceipt.SerialNumbers[1])

	senderKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	senderResp, err := NewAccountCreateTransaction().
		SetKey(senderKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	senderReceipt, err := senderResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	senderId := *senderReceipt.AccountID

	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	receiverResp, err := NewAccountCreateTransaction().
		SetKey(receiverKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receiverReceipt, err := receiverResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	receiverId := *receiverReceipt.AccountID

	senderAssociateTx, err := NewTokenAssociateTransaction().
		SetAccountID(senderId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := senderAssociateTx.Sign(senderKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	receiverAssociateTx, err := NewTokenAssociateTransaction().
		SetAccountID(receiverId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = receiverAssociateTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Transfer NFT to sender account first
	resp, err = NewTransferTransaction().
		AddNftTransfer(tokenId.Nft(1), env.Client.GetOperatorAccountID(), senderId).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Test NFT transfer with both sender and receiver hooks
	senderHookCall := NewNftHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK_SENDER)
	receiverHookCall := NewNftHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_HOOK_RECEIVER)

	transferResponse, err := NewTransferTransaction().
		AddNftTransferWitHook(tokenId.Nft(1), senderId, receiverId, senderHookCall, receiverHookCall).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationTransferWithInvalidGasHook(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	accountId := *receipt.AccountID

	hookCall := NewFungibleHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(0), PRE_HOOK)

	transferResponse, err := NewTransferTransaction().
		AddHbarTransfer(env.OperatorID, NewHbar(1)).
		AddHbarTransferWithHook(accountId, NewHbar(-1), *hookCall).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INSUFFICIENT_GA")
}

func TestIntegrationTransferHbarWithPrePostTransactionAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	accountId := *receipt.AccountID

	hookCall := NewFungibleHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_POST_HOOK)

	transferResponse, err := NewTransferTransaction().
		AddHbarTransfer(env.OperatorID, NewHbar(1)).
		AddHbarTransferWithHook(accountId, NewHbar(-1), *hookCall).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationTransferFungibleTokenWithPrePostTransactionAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	tokenId, err := createFungibleToken(&env)
	require.NoError(t, err)

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	accountId := *receipt.AccountID

	associateTx, err := NewTokenAssociateTransaction().
		SetAccountID(accountId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTx.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		AddTokenTransfer(tokenId, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenId, accountId, 1000).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	hookCall := NewFungibleHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_POST_HOOK)

	transferResponse, err := NewTransferTransaction().
		AddTokenTransfer(tokenId, env.Client.GetOperatorAccountID(), 100).
		AddTokenTransferWithHook(tokenId, accountId, -100, *hookCall).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationTransferNftWithPrePostSenderAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	tokenId, err := createNft(&env)
	require.NoError(t, err)

	mintResp, err := NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadatas([][]byte{{1}, {2}}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = mintResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(accountKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	accountId := *receipt.AccountID

	associateTx, err := NewTokenAssociateTransaction().
		SetAccountID(accountId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTx.Sign(accountKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Transfer NFT to the account first
	resp, err = NewTransferTransaction().
		AddNftTransfer(tokenId.Nft(1), env.Client.GetOperatorAccountID(), accountId).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	hookCall := NewNftHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_POST_HOOK_SENDER)

	transferResponse, err := NewTransferTransaction().
		AddNftTransferWitHook(tokenId.Nft(1), accountId, env.Client.GetOperatorAccountID(), hookCall, nil).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

func TestIntegrationTransferNftWithPrePostReceiverAllowanceHookSucceeds(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	hookContractId := createHookContractId(t, &env)

	hookDetails := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(2).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

	tokenId, err := createNft(&env)
	require.NoError(t, err)

	mintResp, err := NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadatas([][]byte{{1}, {2}}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = mintResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	receiverKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(receiverKey).
		SetInitialBalance(NewHbar(2)).
		AddHook(*hookDetails).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	receiverId := *receipt.AccountID

	associateTx, err := NewTokenAssociateTransaction().
		SetAccountID(receiverId).
		AddTokenID(tokenId).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = associateTx.Sign(receiverKey).Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	receiverHookCall := NewNftHookCall(2, *NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), PRE_POST_HOOK_RECEIVER)

	transferResponse, err := NewTransferTransaction().
		AddNftTransfer(tokenId.Nft(1), env.Client.GetOperatorAccountID(), receiverId).
		AddNftTransferWitHook(tokenId.Nft(2), env.Client.GetOperatorAccountID(), receiverId, nil, receiverHookCall).
		Execute(env.Client)
	require.NoError(t, err)

	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, transferReceipt.Status)
}

// func TestIntegrationTransferWithFullHookId(t *testing.T) {
// 	t.Parallel()
// 	env := NewIntegrationTestEnv(t)
// 	defer CloseIntegrationTestEnv(env, nil)

// 	hookContractId := createHookContractId(t, &env)

// 	hookDetails := NewHookCreationDetails().
// 		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
// 		SetHookId(2).
// 		SetLambdaEvmHook(*NewLambdaEvmHook().SetContractId(hookContractId))

// 	accountKey, err := PrivateKeyGenerateEd25519()
// 	require.NoError(t, err)

// 	resp, err := NewAccountCreateTransaction().
// 		SetKey(accountKey).
// 		SetInitialBalance(NewHbar(2)).
// 		AddHook(*hookDetails).
// 		Execute(env.Client)
// 	require.NoError(t, err)

// 	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
// 	require.NoError(t, err)
// 	accountId := *receipt.AccountID

// 	hookCall := NewFungibleHookCallFull(*NewHookId(*NewHookEntityIdWithAccountId(accountId), 2), *NewEvmHookCall().SetData([]byte{}).SetGasLimit(20000), PRE_HOOK)

// 	transferResponse, err := NewTransferTransaction().
// 		AddHbarTransfer(env.OperatorID, NewHbar(1)).
// 		AddHbarTransferWithHook(accountId, NewHbar(-1), *hookCall).
// 		Execute(env.Client)
// 	require.NoError(t, err)

// 	transferReceipt, err := transferResponse.SetValidateStatus(true).GetReceipt(env.Client)
// 	require.NoError(t, err)
// 	assert.Equal(t, StatusSuccess, transferReceipt.Status)
// }

// Helper function to create bytecode file
func createBytecodeFile(t *testing.T, env *IntegrationTestEnv) FileID {
	const smartContractBytecode = "6080604052348015600e575f5ffd5b506107d18061001c5f395ff3fe608060405260043610610033575f3560e01c8063124d8b301461003757806394112e2f14610067578063bd0dd0b614610097575b5f5ffd5b610051600480360381019061004c91906106f2565b6100c7565b60405161005e9190610782565b60405180910390f35b610081600480360381019061007c91906106f2565b6100d2565b60405161008e9190610782565b60405180910390f35b6100b160048036038101906100ac91906106f2565b6100dd565b6040516100be9190610782565b60405180910390f35b5f6001905092915050565b5f6001905092915050565b5f6001905092915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f60a08284031215610112576101116100f9565b5b81905092915050565b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6101658261011f565b810181811067ffffffffffffffff821117156101845761018361012f565b5b80604052505050565b5f6101966100e8565b90506101a2828261015c565b919050565b5f5ffd5b5f5ffd5b5f67ffffffffffffffff8211156101c9576101c861012f565b5b602082029050602081019050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610207826101de565b9050919050565b610217816101fd565b8114610221575f5ffd5b50565b5f813590506102328161020e565b92915050565b5f8160070b9050919050565b61024d81610238565b8114610257575f5ffd5b50565b5f8135905061026881610244565b92915050565b5f604082840312156102835761028261011b565b5b61028d604061018d565b90505f61029c84828501610224565b5f8301525060206102af8482850161025a565b60208301525092915050565b5f6102cd6102c8846101af565b61018d565b905080838252602082019050604084028301858111156102f0576102ef6101da565b5b835b818110156103195780610305888261026e565b8452602084019350506040810190506102f2565b5050509392505050565b5f82601f830112610337576103366101ab565b5b81356103478482602086016102bb565b91505092915050565b5f67ffffffffffffffff82111561036a5761036961012f565b5b602082029050602081019050919050565b5f67ffffffffffffffff8211156103955761039461012f565b5b602082029050602081019050919050565b5f606082840312156103bb576103ba61011b565b5b6103c5606061018d565b90505f6103d484828501610224565b5f8301525060206103e784828501610224565b60208301525060406103fb8482850161025a565b60408301525092915050565b5f6104196104148461037b565b61018d565b9050808382526020820190506060840283018581111561043c5761043b6101da565b5b835b81811015610465578061045188826103a6565b84526020840193505060608101905061043e565b5050509392505050565b5f82601f830112610483576104826101ab565b5b8135610493848260208601610407565b91505092915050565b5f606082840312156104b1576104b061011b565b5b6104bb606061018d565b90505f6104ca84828501610224565b5f83015250602082013567ffffffffffffffff8111156104ed576104ec6101a7565b5b6104f984828501610323565b602083015250604082013567ffffffffffffffff81111561051d5761051c6101a7565b5b6105298482850161046f565b60408301525092915050565b5f61054761054284610350565b61018d565b9050808382526020820190506020840283018581111561056a576105696101da565b5b835b818110156105b157803567ffffffffffffffff81111561058f5761058e6101ab565b5b80860161059c898261049c565b8552602085019450505060208101905061056c565b5050509392505050565b5f82601f8301126105cf576105ce6101ab565b5b81356105df848260208601610535565b91505092915050565b5f604082840312156105fd576105fc61011b565b5b610607604061018d565b90505f82013567ffffffffffffffff811115610626576106256101a7565b5b61063284828501610323565b5f83015250602082013567ffffffffffffffff811115610655576106546101a7565b5b610661848285016105bb565b60208301525092915050565b5f604082840312156106825761068161011b565b5b61068c604061018d565b90505f82013567ffffffffffffffff8111156106ab576106aa6101a7565b5b6106b7848285016105e8565b5f83015250602082013567ffffffffffffffff8111156106da576106d96101a7565b5b6106e6848285016105e8565b60208301525092915050565b5f5f60408385031215610708576107076100f1565b5b5f83013567ffffffffffffffff811115610725576107246100f5565b5b610731858286016100fd565b925050602083013567ffffffffffffffff811115610752576107516100f5565b5b61075e8582860161066d565b9150509250929050565b5f8115159050919050565b61077c81610768565b82525050565b5f6020820190506107955f830184610773565b9291505056fea26469706673582212207dfe7723f6d6869419b1cb0619758b439da0cf4ffd9520997c40a3946299d4dc64736f6c634300081e0033"
	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey.PublicKey()).
		SetContents([]byte(smartContractBytecode)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.NotNil(t, receipt.FileID)

	return *receipt.FileID
}

// Helper function to create hook contract ID
func createHookContractId(t *testing.T, env *IntegrationTestEnv) *ContractID {
	fileId := createBytecodeFile(t, env)

	response, err := NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey.PublicKey()).
		SetGas(1000000).
		SetBytecodeFileID(fileId).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.NotNil(t, receipt.ContractID)

	return receipt.ContractID
}
