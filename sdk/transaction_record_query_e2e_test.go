//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTransactionRecordQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	require.NoError(t, err)

	resp, err := tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)

	accountID := *record.Receipt.AccountID
	assert.NotNil(t, accountID)

	recordBytes := record.ToBytes()

	_, err = TransactionRecordFromBytes(recordBytes)
	require.NoError(t, err)

	transaction, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTransactionRecordQueryReceiptPaymentZero(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	require.NoError(t, err)

	resp, err := tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)

	accountID := *record.Receipt.AccountID
	assert.NotNil(t, accountID)

	transaction, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTransactionRecordQueryInsufficientFee(t *testing.T) {

	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	require.NoError(t, err)

	resp, err := tx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetIncludeChildren(true).GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = NewTransactionRecordQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(HbarFromTinybar(99999)).
		SetQueryPayment(HbarFromTinybar(1)).
		Execute(env.Client)
	assert.Error(t, err)

	accountID := receipt.AccountID
	assert.NotNil(t, accountID)

	transaction, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func DisabledTestIntegrationTokenTransferRecordsQuery(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetCustomFees([]Fee{CustomFractionalFee{
			CustomFee: CustomFee{
				FeeCollectorAccountID: &env.OperatorID,
			},
			Numerator:     1,
			Denominator:   20,
			MinimumAmount: 1,
			MaximumAmount: 10,
		}})
	})
	require.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddTokenTransfer(tokenID, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, len(record.TokenTransfers), 1)
	assert.Equal(t, len(record.AssessedCustomFees), 0)

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenNftTransferRecordQuery(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokenID, err := createNft(&env)
	require.NoError(t, err)

	metaData := [][]byte{{50}, {50}}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadatas(metaData).
		Execute(env.Client)
	require.NoError(t, err)

	mintReceipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, len(record.NftTransfers), 1)

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetSerialNumbers(mintReceipt.SerialNumbers).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// deployReverterContract deploys a contract with two always-reverting functions:
//
//	error InsufficientBalance(uint256 available, uint256 required);
//	function revertWithReason() external pure { revert("This should fail"); }
//	function revertWithCustomError() external pure { revert InsufficientBalance(5, 10); }
//
// Compiled with solc 0.8.30 --optimize.
func deployReverterContract(t *testing.T, env IntegrationTestEnv) (ContractID, AccountID) {
	reverterBytecode := []byte(`6080604052348015600e575f5ffd5b5060da80601a5f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806346fc4bb11460345780635b2dd10014603c575b5f5ffd5b603a6042565b005b603a606a565b60405163cf47918160e01b815260056004820152600a60248201526044015b60405180910390fd5b60405162461bcd60e51b815260206004820152601060248201526f151a1a5cc81cda1bdd5b190819985a5b60821b6044820152606401606156fea26469706673582212204ff1c43253f237119979541db455673b8e183cb981eee17e42a787a4f1b35fcb64736f6c634300081e0033`)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(reverterBytecode).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetGas(contractDeployGas).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetBytecodeFileID(*receipt.FileID).
		SetContractMemo("hiero-sdk-go::TestTransactionRecordQueryContractRevert").
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	return *receipt.ContractID, resp.NodeID
}

func TestIntegrationTransactionRecordQueryContractRevertReturnsRecord(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	contractID, nodeID := deployReverterContract(t, env)

	resp, err := NewContractExecuteTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetGas(contractDeployGas).
		SetFunction("revertWithReason", nil).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	var receiptErr ErrHederaReceiptStatus
	require.ErrorAs(t, err, &receiptErr)
	assert.Equal(t, StatusContractRevertExecuted, receiptErr.Status)

	require.Equal(t, StatusContractRevertExecuted, record.Receipt.Status)
	result, err := record.GetContractExecuteResult()
	require.NoError(t, err)

	// ErrorMessage is the hex-encoded Error(string) revert data (selector 0x08c379a0).
	require.True(t, strings.HasPrefix(result.ErrorMessage, "0x08c379a0"))
	payload, err := hex.DecodeString(strings.TrimPrefix(result.ErrorMessage, "0x"))
	require.NoError(t, err)
	assert.Contains(t, string(payload), "This should fail")

	// A direct TransactionRecordQuery behaves the same.
	record, err = NewTransactionRecordQuery().
		SetTransactionID(resp.TransactionID).
		SetNodeAccountIDs([]AccountID{nodeID}).
		Execute(env.Client)
	require.Error(t, err)
	require.NotNil(t, record.CallResult)
	assert.Equal(t, result.ErrorMessage, record.CallResult.ErrorMessage)
}

func TestIntegrationTransactionRecordQueryContractRevertValidateStatusFalse(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	contractID, nodeID := deployReverterContract(t, env)

	resp, err := NewContractExecuteTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetGas(contractDeployGas).
		SetFunction("revertWithCustomError", nil).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := resp.SetValidateStatus(false).GetRecord(env.Client)
	require.NoError(t, err)

	require.Equal(t, StatusContractRevertExecuted, record.Receipt.Status)
	result, err := record.GetContractExecuteResult()
	require.NoError(t, err)

	// Custom error InsufficientBalance(uint256,uint256): selector 0xcf479181 + args (5, 10).
	require.True(t, strings.HasPrefix(result.ErrorMessage, "0xcf479181"))
	payload, err := hex.DecodeString(strings.TrimPrefix(result.ErrorMessage, "0x"))
	require.NoError(t, err)
	require.Equal(t, 4+32+32, len(payload))
	assert.Equal(t, byte(5), payload[4+31])
	assert.Equal(t, byte(10), payload[4+63])
}

func TestIntegrationTransactionRecordQueryGetScheduleRef(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	key, err := GeneratePrivateKey()
	require.NoError(t, err)

	createResponse, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(key.PublicKey()).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)

	transactionReceipt, err := createResponse.GetReceipt(env.Client)
	require.NoError(t, err)

	accountId := *transactionReceipt.AccountID

	// Create the transaction
	transaction := NewTransferTransaction().
		AddHbarTransfer(accountId, NewHbar(1).Negated()).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(1)).
		SetMaxTransactionFee(NewHbar(10))

	scheduled, err := transaction.Schedule()
	require.NoError(t, err)

	// Schedule the transaction
	resp, err := scheduled.
		Execute(env.Client)
	require.NoError(t, err)

	record, err := resp.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	// Get the schedule reference
	scheduleRef := record.ScheduleRef
	require.NotNil(t, scheduleRef)
}
