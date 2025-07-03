//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationEthereumFlowCanCreateLargeContract(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	ecdsaPrivateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	aliasAccountId := ecdsaPrivateKey.ToAccountID(0, 0)

	// Create a shallow account for the ECDSA key
	resp, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(*aliasAccountId, NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	chainId, err := hex.DecodeString("012a")
	maxPriorityGas, err := hex.DecodeString("00")
	nonce, err := hex.DecodeString("00")
	maxGas, err := hex.DecodeString("B71B00")        // 12mil
	gasLimitBytes, err := hex.DecodeString("B71B00") // 12mil
	contractBytes, err := hex.DecodeString("00")
	value, err := hex.DecodeString("00")
	callDataBytes, err := hex.DecodeString(LARGE_SMART_CONTRACT_BYTECODE)
	require.NoError(t, err)

	messageBytes, err := getCallData(chainId, nonce, maxPriorityGas, maxGas, gasLimitBytes, contractBytes, value, callDataBytes, ecdsaPrivateKey)
	require.NoError(t, err)

	response, err := NewEthereumFlow().
		SetEthereumDataBytes(messageBytes).
		SetMaxGasAllowance(HbarFrom(10, HbarUnits.Hbar)).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := response.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, record.CallResult.SignerNonce, int64(1))
}

func TestIntegrationEthereumFlowGetGetNodeIdFromResponseOnError(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	ecdsaPrivateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	aliasAccountId := ecdsaPrivateKey.ToAccountID(0, 0)

	// Create a shallow account for the ECDSA key
	resp, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(*aliasAccountId, NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	chainId, err := hex.DecodeString("012a")
	maxPriorityGas, err := hex.DecodeString("00")
	nonce, err := hex.DecodeString("00")
	maxGas, err := hex.DecodeString("00")
	gasLimitBytes, err := hex.DecodeString("00")
	contractBytes, err := hex.DecodeString("00")
	value, err := hex.DecodeString("00")
	callDataBytes, err := hex.DecodeString(LARGE_SMART_CONTRACT_BYTECODE)
	require.NoError(t, err)

	messageBytes, err := getCallData(chainId, nonce, maxPriorityGas, maxGas, gasLimitBytes, contractBytes, value, callDataBytes, ecdsaPrivateKey)
	require.NoError(t, err)

	tx := NewEthereumFlow().
		SetEthereumDataBytes(messageBytes).
		SetMaxGasAllowance(HbarFrom(10, HbarUnits.Hbar))
	require.NotNil(t, tx.GetEthereumData())
	require.NotNil(t, tx.GetEthereumData().GetData())

	response, err := tx.Execute(env.Client)
	require.Error(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.NodeID)
}

func TestIntegrationEthereumFlowJumboTransactionBelowTheLimit(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	ecdsaPrivateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	aliasAccountId := ecdsaPrivateKey.ToAccountID(0, 0)

	// Create a shallow account for the ECDSA key
	resp, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-100)).
		AddHbarTransfer(*aliasAccountId, NewHbar(100)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Create file with the contract bytecode
	resp, err = NewFileCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKeys(env.OperatorKey.PublicKey()).
		SetContents([]byte(SMART_CONTRACT_BYTECODE_JUMBO)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	// Create contract to be called by EthereumTransaction
	resp, err = NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetGas(300_000).
		SetBytecodeFileID(fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	var largeCalldata = make([]byte, 1024*50)
	callData := NewContractFunctionParameters().AddBytes(largeCalldata)
	function := "consumeLargeCalldata"

	chainId, err := hex.DecodeString("012a")
	maxPriorityGas, err := hex.DecodeString("00")
	nonce, err := hex.DecodeString("00")
	maxGas, err := hex.DecodeString("d1385c7bf0")
	gasLimitBytes, err := hex.DecodeString("35E59C")
	contractBytes, err := hex.DecodeString(contractID.ToSolidityAddress())
	value, err := hex.DecodeString("00")
	callDataBytes := callData._Build(&function)
	require.NoError(t, err)

	messageBytes, err := getCallData(chainId, nonce, maxPriorityGas, maxGas, gasLimitBytes, contractBytes, value, callDataBytes, ecdsaPrivateKey)
	require.NoError(t, err)

	response, err := NewEthereumFlow().
		SetEthereumDataBytes(messageBytes).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := response.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, record.CallResult.SignerNonce, int64(1))
}

func TestIntegrationEthereumFlowJumboTransactionAboveTheLimitFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	ecdsaPrivateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	aliasAccountId := ecdsaPrivateKey.ToAccountID(0, 0)

	// Create a shallow account for the ECDSA key
	resp, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-100)).
		AddHbarTransfer(*aliasAccountId, NewHbar(100)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Create file with the contract bytecode
	resp, err = NewFileCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKeys(env.OperatorKey.PublicKey()).
		SetContents([]byte(SMART_CONTRACT_BYTECODE_JUMBO)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	// Create contract to be called by EthereumTransaction
	resp, err = NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetGas(300_000).
		SetBytecodeFileID(fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	var largeCalldata = make([]byte, 1024*140)
	callData := NewContractFunctionParameters().AddBytes(largeCalldata)
	function := "consumeLargeCalldata"

	chainId, err := hex.DecodeString("012a")
	maxPriorityGas, err := hex.DecodeString("00")
	nonce, err := hex.DecodeString("00")
	maxGas, err := hex.DecodeString("d1385c7bf0")
	gasLimitBytes, err := hex.DecodeString("3567E0")
	contractBytes, err := hex.DecodeString(contractID.ToSolidityAddress())
	value, err := hex.DecodeString("00")
	callDataBytes := callData._Build(&function)
	require.NoError(t, err)

	messageBytes, err := getCallData(chainId, nonce, maxPriorityGas, maxGas, gasLimitBytes, contractBytes, value, callDataBytes, ecdsaPrivateKey)
	require.NoError(t, err)

	_, err = NewEthereumFlow().
		SetEthereumDataBytes(messageBytes).
		Execute(env.Client)
	require.ErrorContains(t, err, "TRANSACTION_OVERSIZE")
}
