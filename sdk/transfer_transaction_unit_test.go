//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitTransferTransactionSetTokenTransferWithDecimals(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	senderAccountID := AccountID{Account: 2}
	amount := int64(10)
	decimals := uint32(5)

	transaction := NewTransferTransaction().
		AddTokenTransferWithDecimals(tokenID, senderAccountID, amount, decimals)

	require.Equal(t, transaction.GetTokenIDDecimals()[tokenID], decimals)
}

func TestUnitTransferTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	transfer := NewTransferTransaction().
		AddHbarTransfer(accountID, HbarFromTinybar(1))

	err = transfer.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	transfer := NewTransferTransaction().
		AddHbarTransfer(accountID, HbarFromTinybar(1))

	err = transfer.validateNetworkOnIDs(client)
	require.Error(t, err)
}

func TestUnitTransferTransactionAddHbarTransferWithHook(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 123}
	amount := NewHbar(5)
	evmCall := *NewEvmHookCall().SetData([]byte{0x01, 0x02}).SetGasLimit(25000)
	hookCall := NewFungibleHookCallWithHookId(1, evmCall, PRE_HOOK)

	transaction := NewTransferTransaction().
		AddHbarTransferWithHook(accountID, amount, *hookCall)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.hbarTransfers))
	require.Equal(t, accountID, *transaction.hbarTransfers[0].accountID)
	require.Equal(t, amount, transaction.hbarTransfers[0].amount)
	require.NotNil(t, transaction.hbarTransfers[0].hookCall)
	require.Equal(t, int64(1), transaction.hbarTransfers[0].hookCall.GetHookId())
}

func TestUnitTransferTransactionAddTokenTransferWithHook(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 100}
	accountID := AccountID{Account: 200}
	amount := int64(1000)
	evmCall := *NewEvmHookCall().SetData([]byte{0x03, 0x04}).SetGasLimit(30000)
	hookCall := NewFungibleHookCallWithHookId(2, evmCall, PRE_POST_HOOK)

	transaction := NewTransferTransaction().
		AddTokenTransferWithHook(tokenID, accountID, amount, *hookCall)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.tokenTransfers))
	require.NotNil(t, transaction.tokenTransfers[tokenID])
	require.Equal(t, 1, len(transaction.tokenTransfers[tokenID].Transfers))
	require.NotNil(t, transaction.tokenTransfers[tokenID].Transfers[0].hookCall)
	require.Equal(t, int64(2), transaction.tokenTransfers[tokenID].Transfers[0].hookCall.GetHookId())
}

func TestUnitTransferTransactionAddNftTransferWithHook(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 300}
	nftID := tokenID.Nft(1)
	sender := AccountID{Account: 400}
	receiver := AccountID{Account: 500}
	evmCall := *NewEvmHookCall().SetData([]byte{0x05, 0x06}).SetGasLimit(35000)
	senderHookCall := NewNftHookCallWithHookId(3, evmCall, PRE_HOOK_SENDER)
	receiverHookCall := NewNftHookCallWithHookId(4, evmCall, PRE_HOOK_RECEIVER)

	transaction := NewTransferTransaction().
		AddNftTransferWitHook(nftID, sender, receiver, senderHookCall, receiverHookCall)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.nftTransfers))
	require.NotNil(t, transaction.nftTransfers[nftID.TokenID])
	require.Equal(t, 1, len(transaction.nftTransfers[nftID.TokenID]))
	require.NotNil(t, transaction.nftTransfers[nftID.TokenID][0].SenderHookCall)
	require.NotNil(t, transaction.nftTransfers[nftID.TokenID][0].ReceiverHookCall)
	require.Equal(t, int64(3), transaction.nftTransfers[nftID.TokenID][0].SenderHookCall.GetHookId())
	require.Equal(t, int64(4), transaction.nftTransfers[nftID.TokenID][0].ReceiverHookCall.GetHookId())
}

func TestUnitTransferTransactionValidateHbarTransferWithHook(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	contractID, err := ContractIDFromString("0.0.456-fuxra")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookIdFull := NewHookId(*entityId, 789)
	evmCall := *NewEvmHookCall().SetData([]byte{0x01}).SetGasLimit(20000)
	hookCall := NewFungibleHookCallWithHookIdFull(*hookIdFull, evmCall, PRE_HOOK)

	transaction := NewTransferTransaction().
		AddHbarTransferWithHook(accountID, NewHbar(1), *hookCall)

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateHbarTransferWithHookWrongChecksum(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	// Wrong checksum for contract ID
	contractID, err := ContractIDFromString("0.0.456-rmkykd")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookIdFull := NewHookId(*entityId, 789)
	evmCall := *NewEvmHookCall().SetData([]byte{0x01}).SetGasLimit(20000)
	hookCall := NewFungibleHookCallWithHookIdFull(*hookIdFull, evmCall, PRE_HOOK)

	transaction := NewTransferTransaction().
		AddHbarTransferWithHook(accountID, NewHbar(1), *hookCall)

	err = transaction.validateNetworkOnIDs(client)
	require.Error(t, err)
}

func TestUnitTransferTransactionValidateTokenTransferWithHook(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	tokenID, err := TokenIDFromString("0.0.100-quros")
	require.NoError(t, err)

	accountID, err := AccountIDFromString("0.0.200-tyrmb")
	require.NoError(t, err)

	contractID, err := ContractIDFromString("0.0.300-xcrjk")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookIdFull := NewHookId(*entityId, 123)
	evmCall := *NewEvmHookCall().SetData([]byte{0x02}).SetGasLimit(25000)
	hookCall := NewFungibleHookCallWithHookIdFull(*hookIdFull, evmCall, PRE_POST_HOOK)

	transaction := NewTransferTransaction().
		AddTokenTransferWithHook(tokenID, accountID, 1000, *hookCall)

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateTokenTransferWithHookWrongTokenChecksum(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	// Wrong checksum for token ID
	tokenID, err := TokenIDFromString("0.0.100-rmkykd")
	require.NoError(t, err)

	accountID, err := AccountIDFromString("0.0.200-tyrmb")
	require.NoError(t, err)

	evmCall := *NewEvmHookCall().SetData([]byte{0x02}).SetGasLimit(25000)
	hookCall := NewFungibleHookCallWithHookId(1, evmCall, PRE_HOOK)

	transaction := NewTransferTransaction().
		AddTokenTransferWithHook(tokenID, accountID, 1000, *hookCall)

	err = transaction.validateNetworkOnIDs(client)
	require.Error(t, err)
}

func TestUnitTransferTransactionValidateNftTransferWithHooks(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	tokenID, err := TokenIDFromString("0.0.100-quros")
	require.NoError(t, err)

	sender, err := AccountIDFromString("0.0.200-tyrmb")
	require.NoError(t, err)

	receiver, err := AccountIDFromString("0.0.300-xcrjk")
	require.NoError(t, err)

	contractID1, err := ContractIDFromString("0.0.400-agrgt")
	require.NoError(t, err)

	contractID2, err := ContractIDFromString("0.0.500-dkrec")
	require.NoError(t, err)

	entityId1 := NewHookEntityIdWithContractId(contractID1)
	hookIdFull1 := NewHookId(*entityId1, 111)
	evmCall1 := *NewEvmHookCall().SetData([]byte{0x03}).SetGasLimit(30000)
	senderHookCall := NewNftHookCallWithHookIdFull(*hookIdFull1, evmCall1, PRE_HOOK_SENDER)

	entityId2 := NewHookEntityIdWithContractId(contractID2)
	hookIdFull2 := NewHookId(*entityId2, 222)
	evmCall2 := *NewEvmHookCall().SetData([]byte{0x04}).SetGasLimit(35000)
	receiverHookCall := NewNftHookCallWithHookIdFull(*hookIdFull2, evmCall2, PRE_HOOK_RECEIVER)

	nftID := tokenID.Nft(1)

	transaction := NewTransferTransaction().
		AddNftTransferWitHook(nftID, sender, receiver, senderHookCall, receiverHookCall)

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateNftTransferWithHooksWrongSenderChecksum(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	tokenID, err := TokenIDFromString("0.0.100-quros")
	require.NoError(t, err)

	// Wrong checksum for sender
	sender, err := AccountIDFromString("0.0.200-rmkykd")
	require.NoError(t, err)

	receiver, err := AccountIDFromString("0.0.300-xcrjk")
	require.NoError(t, err)

	evmCall := *NewEvmHookCall().SetData([]byte{0x03}).SetGasLimit(30000)
	senderHookCall := NewNftHookCallWithHookId(1, evmCall, PRE_HOOK_SENDER)

	nftID := tokenID.Nft(1)

	transaction := NewTransferTransaction().
		AddNftTransferWitHook(nftID, sender, receiver, senderHookCall, nil)

	err = transaction.validateNetworkOnIDs(client)
	require.Error(t, err)
}

func TestUnitTransferTransactionValidateNftTransferWithHooksWrongHookChecksum(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	tokenID, err := TokenIDFromString("0.0.100-quros")
	require.NoError(t, err)

	sender, err := AccountIDFromString("0.0.200-tyrmb")
	require.NoError(t, err)

	receiver, err := AccountIDFromString("0.0.300-xcrjk")
	require.NoError(t, err)

	// Wrong checksum for contract ID in hook
	contractID, err := ContractIDFromString("0.0.400-rmkykd")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookIdFull := NewHookId(*entityId, 111)
	evmCall := *NewEvmHookCall().SetData([]byte{0x03}).SetGasLimit(30000)
	senderHookCall := NewNftHookCallWithHookIdFull(*hookIdFull, evmCall, PRE_HOOK_SENDER)

	nftID := tokenID.Nft(1)

	transaction := NewTransferTransaction().
		AddNftTransferWitHook(nftID, sender, receiver, senderHookCall, nil)

	err = transaction.validateNetworkOnIDs(client)
	require.Error(t, err)
}

func TestUnitTransferTransactionValidateMixedTransfersWithHooks(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	// Setup accounts
	account1, err := AccountIDFromString("0.0.100-quros")
	require.NoError(t, err)

	account2, err := AccountIDFromString("0.0.200-tyrmb")
	require.NoError(t, err)

	// Setup token
	tokenID, err := TokenIDFromString("0.0.300-xcrjk")
	require.NoError(t, err)

	// Setup contract for hook
	contractID, err := ContractIDFromString("0.0.400-agrgt")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookIdFull := NewHookId(*entityId, 123)
	evmCall := *NewEvmHookCall().SetData([]byte{0x05}).SetGasLimit(40000)
	hookCall := NewFungibleHookCallWithHookIdFull(*hookIdFull, evmCall, PRE_HOOK)

	// Create transaction with mixed transfers (HBAR and Token with hooks)
	transaction := NewTransferTransaction().
		AddHbarTransferWithHook(account1, NewHbar(-2), *hookCall).
		AddHbarTransferWithHook(account2, NewHbar(2), *hookCall).
		AddTokenTransferWithHook(tokenID, account1, -1000, *hookCall).
		AddTokenTransferWithHook(tokenID, account2, 1000, *hookCall)

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateWithNoClient(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 123}
	evmCall := *NewEvmHookCall().SetData([]byte{0x01}).SetGasLimit(20000)
	hookCall := NewFungibleHookCallWithHookId(1, evmCall, PRE_HOOK)

	transaction := NewTransferTransaction().
		AddHbarTransferWithHook(accountID, NewHbar(1), *hookCall)

	err := transaction.validateNetworkOnIDs(nil)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateWithChecksumDisabled(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(false) // Disable checksum validation

	// Use wrong checksum but it should pass because validation is disabled
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	evmCall := *NewEvmHookCall().SetData([]byte{0x01}).SetGasLimit(20000)
	hookCall := NewFungibleHookCallWithHookId(1, evmCall, PRE_HOOK)

	transaction := NewTransferTransaction().
		AddHbarTransferWithHook(accountID, NewHbar(1), *hookCall)

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionGetHbarTransfers(t *testing.T) {
	t.Parallel()

	account1 := AccountID{Account: 100}
	account2 := AccountID{Account: 200}
	amount1 := NewHbar(5)
	amount2 := NewHbar(-5)

	transaction := NewTransferTransaction().
		AddHbarTransfer(account1, amount1).
		AddHbarTransfer(account2, amount2)

	transfers := transaction.GetHbarTransfers()
	require.Equal(t, 2, len(transfers))
	require.Equal(t, amount1, transfers[account1])
	require.Equal(t, amount2, transfers[account2])
}

func TestUnitTransferTransactionGetTokenTransfers(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 100}
	account1 := AccountID{Account: 200}
	account2 := AccountID{Account: 300}
	amount1 := int64(1000)
	amount2 := int64(-1000)

	transaction := NewTransferTransaction().
		AddTokenTransfer(tokenID, account1, amount1).
		AddTokenTransfer(tokenID, account2, amount2)

	transfers := transaction.GetTokenTransfers()
	require.Equal(t, 1, len(transfers))
	require.Equal(t, 2, len(transfers[tokenID]))
}

func TestUnitTransferTransactionGetNftTransfers(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 100}
	nftID := tokenID.Nft(1)
	sender := AccountID{Account: 200}
	receiver := AccountID{Account: 300}

	transaction := NewTransferTransaction().
		AddNftTransfer(nftID, sender, receiver)

	transfers := transaction.GetNftTransfers()
	require.Equal(t, 1, len(transfers))
	require.Equal(t, 1, len(transfers[nftID.TokenID]))
	require.Equal(t, sender, transfers[nftID.TokenID][0].SenderAccountID)
	require.Equal(t, receiver, transfers[nftID.TokenID][0].ReceiverAccountID)
}

func TestUnitTransferTransactionSetHbarTransferApproval(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 100}
	amount := NewHbar(5)

	transaction := NewTransferTransaction().
		AddHbarTransfer(accountID, amount).
		SetHbarTransferApproval(accountID, true)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.hbarTransfers))
	require.True(t, transaction.hbarTransfers[0].isApproved)
}

func TestUnitTransferTransactionSetTokenTransferApproval(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 100}
	accountID := AccountID{Account: 200}
	amount := int64(1000)

	transaction := NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID, amount).
		SetTokenTransferApproval(tokenID, accountID, true)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.tokenTransfers))
	require.True(t, transaction.tokenTransfers[tokenID].Transfers[0].isApproved)
}

func TestUnitTransferTransactionSetNftTransferApproval(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 100}
	nftID := tokenID.Nft(1)
	sender := AccountID{Account: 200}
	receiver := AccountID{Account: 300}

	transaction := NewTransferTransaction().
		AddNftTransfer(nftID, sender, receiver).
		SetNftTransferApproval(nftID, true)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.nftTransfers))
	require.True(t, transaction.nftTransfers[nftID.TokenID][0].IsApproved)
}

func TestUnitTransferTransactionAddApprovedHbarTransfer(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 100}
	amount := NewHbar(5)

	transaction := NewTransferTransaction().
		AddApprovedHbarTransfer(accountID, amount, true)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.hbarTransfers))
	require.True(t, transaction.hbarTransfers[0].isApproved)
}

func TestUnitTransferTransactionAddApprovedTokenTransfer(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 100}
	accountID := AccountID{Account: 200}
	amount := int64(1000)

	transaction := NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, accountID, amount, true)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.tokenTransfers))
	require.True(t, transaction.tokenTransfers[tokenID].Transfers[0].isApproved)
}

func TestUnitTransferTransactionAddApprovedNftTransfer(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 100}
	nftID := tokenID.Nft(1)
	sender := AccountID{Account: 200}
	receiver := AccountID{Account: 300}

	transaction := NewTransferTransaction().
		AddApprovedNftTransfer(nftID, sender, receiver, true)

	require.NotNil(t, transaction)
	require.Equal(t, 1, len(transaction.nftTransfers))
	require.True(t, transaction.nftTransfers[nftID.TokenID][0].IsApproved)
}

func TestUnitTransferTransactionGetName(t *testing.T) {
	t.Parallel()

	transaction := NewTransferTransaction()
	require.Equal(t, "TransferTransaction", transaction.getName())
}

func TestUnitTransferTransactionAddHbarTransferProtobuf(t *testing.T) {
	t.Parallel()

	account1 := AccountID{Account: 100}
	account2 := AccountID{Account: 200}
	amount1 := NewHbar(5)
	amount2 := NewHbar(-5)

	transaction := NewTransferTransaction().
		AddHbarTransfer(account1, amount1).
		AddHbarTransfer(account2, amount2)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.NotNil(t, body.Transfers)
	require.Equal(t, 2, len(body.Transfers.AccountAmounts))

	// Verify first transfer
	require.Equal(t, int64(100), body.Transfers.AccountAmounts[0].AccountID.GetAccountNum())
	require.Equal(t, amount1.AsTinybar(), body.Transfers.AccountAmounts[0].Amount)
	require.False(t, body.Transfers.AccountAmounts[0].IsApproval)
	require.Nil(t, body.Transfers.AccountAmounts[0].HookCall)

	// Verify second transfer
	require.Equal(t, int64(200), body.Transfers.AccountAmounts[1].AccountID.GetAccountNum())
	require.Equal(t, amount2.AsTinybar(), body.Transfers.AccountAmounts[1].Amount)
	require.False(t, body.Transfers.AccountAmounts[1].IsApproval)
	require.Nil(t, body.Transfers.AccountAmounts[1].HookCall)
}

func TestUnitTransferTransactionAddHbarTransferWithHookProtobuf(t *testing.T) {
	t.Parallel()

	account := AccountID{Account: 123}
	amount := NewHbar(10)
	evmCall := *NewEvmHookCall().SetData([]byte{0x01, 0x02}).SetGasLimit(25000)
	hookCall := NewFungibleHookCallWithHookId(5, evmCall, PRE_HOOK)

	transaction := NewTransferTransaction().
		AddHbarTransferWithHook(account, amount, *hookCall)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.NotNil(t, body.Transfers)
	require.Equal(t, 1, len(body.Transfers.AccountAmounts))

	// Verify transfer with hook
	require.Equal(t, int64(123), body.Transfers.AccountAmounts[0].AccountID.GetAccountNum())
	require.Equal(t, amount.AsTinybar(), body.Transfers.AccountAmounts[0].Amount)
	require.False(t, body.Transfers.AccountAmounts[0].IsApproval)
	require.NotNil(t, body.Transfers.AccountAmounts[0].HookCall)
	protoHookCall := body.Transfers.AccountAmounts[0].GetPreTxAllowanceHook()
	require.NotNil(t, protoHookCall)
	require.Equal(t, int64(5), protoHookCall.GetHookId())
}

func TestUnitTransferTransactionAddApprovedHbarTransferProtobuf(t *testing.T) {
	t.Parallel()

	account := AccountID{Account: 150}
	amount := NewHbar(7)

	transaction := NewTransferTransaction().
		AddApprovedHbarTransfer(account, amount, true)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.NotNil(t, body.Transfers)
	require.Equal(t, 1, len(body.Transfers.AccountAmounts))

	// Verify approved transfer (no hook)
	require.Equal(t, int64(150), body.Transfers.AccountAmounts[0].AccountID.GetAccountNum())
	require.Equal(t, amount.AsTinybar(), body.Transfers.AccountAmounts[0].Amount)
	require.True(t, body.Transfers.AccountAmounts[0].IsApproval)
	require.Nil(t, body.Transfers.AccountAmounts[0].HookCall)
}

func TestUnitTransferTransactionAddTokenTransferProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 300}
	account1 := AccountID{Account: 400}
	account2 := AccountID{Account: 500}
	amount1 := int64(1000)
	amount2 := int64(-1000)

	transaction := NewTransferTransaction().
		AddTokenTransfer(tokenID, account1, amount1).
		AddTokenTransfer(tokenID, account2, amount2)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify token transfer list
	require.Equal(t, int64(300), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 2, len(body.TokenTransfers[0].Transfers))
	require.Nil(t, body.TokenTransfers[0].ExpectedDecimals)

	// Verify first token transfer
	require.Equal(t, int64(400), body.TokenTransfers[0].Transfers[0].AccountID.GetAccountNum())
	require.Equal(t, amount1, body.TokenTransfers[0].Transfers[0].Amount)
	require.False(t, body.TokenTransfers[0].Transfers[0].IsApproval)
	require.Nil(t, body.TokenTransfers[0].Transfers[0].HookCall)

	// Verify second token transfer
	require.Equal(t, int64(500), body.TokenTransfers[0].Transfers[1].AccountID.GetAccountNum())
	require.Equal(t, amount2, body.TokenTransfers[0].Transfers[1].Amount)
	require.False(t, body.TokenTransfers[0].Transfers[1].IsApproval)
	require.Nil(t, body.TokenTransfers[0].Transfers[1].HookCall)
}

func TestUnitTransferTransactionAddTokenTransferWithDecimalsProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 310}
	account := AccountID{Account: 410}
	amount := int64(5000)
	decimals := uint32(6)

	transaction := NewTransferTransaction().
		AddTokenTransferWithDecimals(tokenID, account, amount, decimals)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify token transfer with decimals
	require.Equal(t, int64(310), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 1, len(body.TokenTransfers[0].Transfers))
	require.NotNil(t, body.TokenTransfers[0].ExpectedDecimals)
	require.Equal(t, decimals, body.TokenTransfers[0].ExpectedDecimals.Value)

	// Verify transfer details
	require.Equal(t, int64(410), body.TokenTransfers[0].Transfers[0].AccountID.GetAccountNum())
	require.Equal(t, amount, body.TokenTransfers[0].Transfers[0].Amount)
	require.False(t, body.TokenTransfers[0].Transfers[0].IsApproval)
}

func TestUnitTransferTransactionAddTokenTransferWithHookProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 320}
	account := AccountID{Account: 420}
	amount := int64(2000)
	evmCall := *NewEvmHookCall().SetData([]byte{0x03, 0x04}).SetGasLimit(30000)
	hookCall := NewFungibleHookCallWithHookId(7, evmCall, PRE_POST_HOOK)

	transaction := NewTransferTransaction().
		AddTokenTransferWithHook(tokenID, account, amount, *hookCall)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify token transfer with hook
	require.Equal(t, int64(320), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 1, len(body.TokenTransfers[0].Transfers))

	// Verify transfer with hook
	require.Equal(t, int64(420), body.TokenTransfers[0].Transfers[0].AccountID.GetAccountNum())
	require.Equal(t, amount, body.TokenTransfers[0].Transfers[0].Amount)
	require.False(t, body.TokenTransfers[0].Transfers[0].IsApproval)
	require.NotNil(t, body.TokenTransfers[0].Transfers[0].HookCall)
	protoHookCall := body.TokenTransfers[0].Transfers[0].GetPrePostTxAllowanceHook()
	require.NotNil(t, protoHookCall)
	require.Equal(t, int64(7), protoHookCall.GetHookId())
}

func TestUnitTransferTransactionAddApprovedTokenTransferProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 330}
	account := AccountID{Account: 430}
	amount := int64(3000)

	transaction := NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, account, amount, true)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify approved token transfer (no hook)
	require.Equal(t, int64(330), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 1, len(body.TokenTransfers[0].Transfers))

	// Verify transfer is approved
	require.Equal(t, int64(430), body.TokenTransfers[0].Transfers[0].AccountID.GetAccountNum())
	require.Equal(t, amount, body.TokenTransfers[0].Transfers[0].Amount)
	require.True(t, body.TokenTransfers[0].Transfers[0].IsApproval)
	require.Nil(t, body.TokenTransfers[0].Transfers[0].HookCall)
}

func TestUnitTransferTransactionAddApprovedTokenTransferWithDecimalsProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 340}
	account := AccountID{Account: 440}
	amount := int64(4000)
	decimals := uint32(8)

	transaction := NewTransferTransaction().
		AddApprovedTokenTransferWithDecimals(tokenID, account, amount, decimals, true)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify approved token transfer with decimals (no hook)
	require.Equal(t, int64(340), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 1, len(body.TokenTransfers[0].Transfers))
	require.NotNil(t, body.TokenTransfers[0].ExpectedDecimals)
	require.Equal(t, decimals, body.TokenTransfers[0].ExpectedDecimals.Value)

	// Verify transfer is approved
	require.Equal(t, int64(440), body.TokenTransfers[0].Transfers[0].AccountID.GetAccountNum())
	require.Equal(t, amount, body.TokenTransfers[0].Transfers[0].Amount)
	require.True(t, body.TokenTransfers[0].Transfers[0].IsApproval)
	require.Nil(t, body.TokenTransfers[0].Transfers[0].HookCall)
}

func TestUnitTransferTransactionAddNftTransferProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 500}
	nftID := tokenID.Nft(1)
	sender := AccountID{Account: 600}
	receiver := AccountID{Account: 700}

	transaction := NewTransferTransaction().
		AddNftTransfer(nftID, sender, receiver)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify NFT transfer
	require.Equal(t, int64(500), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 1, len(body.TokenTransfers[0].NftTransfers))

	// Verify NFT transfer details (no hooks)
	require.Equal(t, int64(600), body.TokenTransfers[0].NftTransfers[0].SenderAccountID.GetAccountNum())
	require.Equal(t, int64(700), body.TokenTransfers[0].NftTransfers[0].ReceiverAccountID.GetAccountNum())
	require.Equal(t, int64(1), body.TokenTransfers[0].NftTransfers[0].SerialNumber)
	require.False(t, body.TokenTransfers[0].NftTransfers[0].IsApproval)
	require.Nil(t, body.TokenTransfers[0].NftTransfers[0].SenderAllowanceHookCall)
	require.Nil(t, body.TokenTransfers[0].NftTransfers[0].ReceiverAllowanceHookCall)
}

func TestUnitTransferTransactionAddNftTransferWithHookProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 510}
	nftID := tokenID.Nft(2)
	sender := AccountID{Account: 610}
	receiver := AccountID{Account: 710}
	evmCall := *NewEvmHookCall().SetData([]byte{0x05, 0x06}).SetGasLimit(35000)
	senderHookCall := NewNftHookCallWithHookId(8, evmCall, PRE_HOOK_SENDER)
	receiverHookCall := NewNftHookCallWithHookId(9, evmCall, PRE_HOOK_RECEIVER)

	transaction := NewTransferTransaction().
		AddNftTransferWitHook(nftID, sender, receiver, senderHookCall, receiverHookCall)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify NFT transfer with hooks
	require.Equal(t, int64(510), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 1, len(body.TokenTransfers[0].NftTransfers))

	// Verify NFT transfer details with hooks
	nftTransfer := body.TokenTransfers[0].NftTransfers[0]
	require.Equal(t, int64(610), nftTransfer.SenderAccountID.GetAccountNum())
	require.Equal(t, int64(710), nftTransfer.ReceiverAccountID.GetAccountNum())
	require.Equal(t, int64(2), nftTransfer.SerialNumber)
	require.False(t, nftTransfer.IsApproval)

	// Verify sender hook
	require.NotNil(t, nftTransfer.SenderAllowanceHookCall)
	senderHook := nftTransfer.GetPreTxSenderAllowanceHook()
	require.NotNil(t, senderHook)
	require.Equal(t, int64(8), senderHook.GetHookId())

	// Verify receiver hook
	require.NotNil(t, nftTransfer.ReceiverAllowanceHookCall)
	receiverHook := nftTransfer.GetPreTxReceiverAllowanceHook()
	require.NotNil(t, receiverHook)
	require.Equal(t, int64(9), receiverHook.GetHookId())
}

func TestUnitTransferTransactionAddApprovedNftTransferProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 520}
	nftID := tokenID.Nft(3)
	sender := AccountID{Account: 620}
	receiver := AccountID{Account: 720}

	transaction := NewTransferTransaction().
		AddApprovedNftTransfer(nftID, sender, receiver, true)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))

	// Verify approved NFT transfer (no hooks)
	require.Equal(t, int64(520), body.TokenTransfers[0].Token.TokenNum)
	require.Equal(t, 1, len(body.TokenTransfers[0].NftTransfers))

	// Verify NFT transfer is approved
	nftTransfer := body.TokenTransfers[0].NftTransfers[0]
	require.Equal(t, int64(620), nftTransfer.SenderAccountID.GetAccountNum())
	require.Equal(t, int64(720), nftTransfer.ReceiverAccountID.GetAccountNum())
	require.Equal(t, int64(3), nftTransfer.SerialNumber)
	require.True(t, nftTransfer.IsApproval)
	require.Nil(t, nftTransfer.SenderAllowanceHookCall)
	require.Nil(t, nftTransfer.ReceiverAllowanceHookCall)
}

// Test multiple transfer types in single transaction
func TestUnitTransferTransactionMixedTransfersProtobuf(t *testing.T) {
	t.Parallel()

	// Setup IDs
	account1 := AccountID{Account: 100}
	account2 := AccountID{Account: 200}
	tokenID := TokenID{Token: 300}
	nftID := tokenID.Nft(1)

	// Create mixed transaction
	transaction := NewTransferTransaction().
		AddHbarTransfer(account1, NewHbar(-10)).
		AddHbarTransfer(account2, NewHbar(10)).
		AddTokenTransfer(tokenID, account1, -1000).
		AddTokenTransfer(tokenID, account2, 1000).
		AddNftTransfer(nftID, account1, account2)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)

	// Verify HBAR transfers
	require.Equal(t, 2, len(body.Transfers.AccountAmounts))

	// Verify token transfers (fungible and NFT transfers create separate TokenTransferList entries)
	require.Equal(t, 2, len(body.TokenTransfers))
	// One list has fungible transfers
	require.True(t, len(body.TokenTransfers[0].Transfers) == 2 || len(body.TokenTransfers[1].Transfers) == 2)
	// One list has NFT transfers
	require.True(t, len(body.TokenTransfers[0].NftTransfers) == 1 || len(body.TokenTransfers[1].NftTransfers) == 1)
}

// Test that adding same account multiple times accumulates amounts
func TestUnitTransferTransactionAccumulateHbarTransfers(t *testing.T) {
	t.Parallel()

	account := AccountID{Account: 123}

	transaction := NewTransferTransaction().
		AddHbarTransfer(account, NewHbar(5)).
		AddHbarTransfer(account, NewHbar(3))

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.Transfers.AccountAmounts))
	require.Equal(t, NewHbar(8).AsTinybar(), body.Transfers.AccountAmounts[0].Amount)
}

// Test that adding same token account multiple times accumulates amounts
func TestUnitTransferTransactionAccumulateTokenTransfers(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 300}
	account := AccountID{Account: 400}

	transaction := NewTransferTransaction().
		AddTokenTransfer(tokenID, account, 1000).
		AddTokenTransfer(tokenID, account, 500)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))
	require.Equal(t, 1, len(body.TokenTransfers[0].Transfers))
	require.Equal(t, int64(1500), body.TokenTransfers[0].Transfers[0].Amount)
}

// Test multiple NFT transfers for same token
func TestUnitTransferTransactionMultipleNftTransfersProtobuf(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 500}
	nftID1 := tokenID.Nft(1)
	nftID2 := tokenID.Nft(2)
	sender := AccountID{Account: 600}
	receiver := AccountID{Account: 700}

	transaction := NewTransferTransaction().
		AddNftTransfer(nftID1, sender, receiver).
		AddNftTransfer(nftID2, sender, receiver)

	body := transaction.buildProtoBody()

	require.NotNil(t, body)
	require.Equal(t, 1, len(body.TokenTransfers))
	require.Equal(t, 2, len(body.TokenTransfers[0].NftTransfers))
	require.Equal(t, int64(1), body.TokenTransfers[0].NftTransfers[0].SerialNumber)
	require.Equal(t, int64(2), body.TokenTransfers[0].NftTransfers[1].SerialNumber)
}

// Test that hook and approval are mutually exclusive for HBAR transfers
func TestUnitTransferTransactionHookAndApprovalMutuallyExclusiveHbar(t *testing.T) {
	t.Parallel()

	account1 := AccountID{Account: 100}
	account2 := AccountID{Account: 200}
	evmCall := *NewEvmHookCall().SetData([]byte{0x01}).SetGasLimit(20000)
	hookCall := NewFungibleHookCallWithHookId(1, evmCall, PRE_HOOK)

	// Add transfer with hook
	txWithHook := NewTransferTransaction().
		AddHbarTransferWithHook(account1, NewHbar(5), *hookCall)

	bodyWithHook := txWithHook.buildProtoBody()
	require.NotNil(t, bodyWithHook.Transfers.AccountAmounts[0].HookCall)
	require.False(t, bodyWithHook.Transfers.AccountAmounts[0].IsApproval)

	// Add approved transfer (no hook)
	txWithApproval := NewTransferTransaction().
		AddApprovedHbarTransfer(account2, NewHbar(5), true)

	bodyWithApproval := txWithApproval.buildProtoBody()
	require.Nil(t, bodyWithApproval.Transfers.AccountAmounts[0].HookCall)
	require.True(t, bodyWithApproval.Transfers.AccountAmounts[0].IsApproval)
}

// Test that hook and approval are mutually exclusive for token transfers
func TestUnitTransferTransactionHookAndApprovalMutuallyExclusiveToken(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 300}
	account1 := AccountID{Account: 400}
	account2 := AccountID{Account: 500}
	evmCall := *NewEvmHookCall().SetData([]byte{0x02}).SetGasLimit(25000)
	hookCall := NewFungibleHookCallWithHookId(2, evmCall, PRE_POST_HOOK)

	// Add transfer with hook
	txWithHook := NewTransferTransaction().
		AddTokenTransferWithHook(tokenID, account1, 1000, *hookCall)

	bodyWithHook := txWithHook.buildProtoBody()
	require.NotNil(t, bodyWithHook.TokenTransfers[0].Transfers[0].HookCall)
	require.False(t, bodyWithHook.TokenTransfers[0].Transfers[0].IsApproval)

	// Add approved transfer (no hook)
	txWithApproval := NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID, account2, 1000, true)

	bodyWithApproval := txWithApproval.buildProtoBody()
	require.Nil(t, bodyWithApproval.TokenTransfers[0].Transfers[0].HookCall)
	require.True(t, bodyWithApproval.TokenTransfers[0].Transfers[0].IsApproval)
}

// Test that hook and approval are mutually exclusive for NFT transfers
func TestUnitTransferTransactionHookAndApprovalMutuallyExclusiveNft(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 500}
	nftID1 := tokenID.Nft(1)
	nftID2 := tokenID.Nft(2)
	sender := AccountID{Account: 600}
	receiver := AccountID{Account: 700}
	evmCall := *NewEvmHookCall().SetData([]byte{0x03}).SetGasLimit(30000)
	senderHookCall := NewNftHookCallWithHookId(3, evmCall, PRE_HOOK_SENDER)

	// Add NFT transfer with hook
	txWithHook := NewTransferTransaction().
		AddNftTransferWitHook(nftID1, sender, receiver, senderHookCall, nil)

	bodyWithHook := txWithHook.buildProtoBody()
	require.NotNil(t, bodyWithHook.TokenTransfers[0].NftTransfers[0].SenderAllowanceHookCall)
	require.False(t, bodyWithHook.TokenTransfers[0].NftTransfers[0].IsApproval)

	// Add approved NFT transfer (no hook)
	txWithApproval := NewTransferTransaction().
		AddApprovedNftTransfer(nftID2, sender, receiver, true)

	bodyWithApproval := txWithApproval.buildProtoBody()
	require.Nil(t, bodyWithApproval.TokenTransfers[0].NftTransfers[0].SenderAllowanceHookCall)
	require.Nil(t, bodyWithApproval.TokenTransfers[0].NftTransfers[0].ReceiverAllowanceHookCall)
	require.True(t, bodyWithApproval.TokenTransfers[0].NftTransfers[0].IsApproval)
}
