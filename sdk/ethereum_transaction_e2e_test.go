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

func TestIntegrationEthereumTransaction(t *testing.T) {
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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Create file with the contract bytecode
	resp, err = NewFileCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKeys(env.OperatorKey.PublicKey()).
		SetContents([]byte(ETHEREUM_SMART_CONTRACT_BYTECODE)).
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
		SetGas(1000000).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hiero")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hiero-sdk-go::TestContractCreateTransaction_Execute").
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	callData := NewContractFunctionParameters().AddString("new message")
	function := "setMessage"

	chainId, err := hex.DecodeString("012a")
	maxPriorityGas, err := hex.DecodeString("00")
	nonce, err := hex.DecodeString("00")
	maxGas, err := hex.DecodeString("d1385c7bf0")
	gasLimitBytes, err := hex.DecodeString("0249f0") // 150k
	contractBytes, err := hex.DecodeString(contractID.ToEvmAddress())
	value, err := hex.DecodeString("00")
	callDataBytes := callData._Build(&function)
	require.NoError(t, err)

	messageBytes, err := getCallData(chainId, nonce, maxPriorityGas, maxGas, gasLimitBytes, contractBytes, value, callDataBytes, ecdsaPrivateKey)
	require.NoError(t, err)

	resp, err = NewEthereumTransaction().SetEthereumData(messageBytes).Execute(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, int64(1), record.CallResult.SignerNonce)

	resp, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationEthereumTransactionJumboTransaction(t *testing.T) {
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

	var largeCalldata = make([]byte, 1024*120)
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

	response, err := NewEthereumTransaction().
		SetEthereumData(messageBytes).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := response.SetValidateStatus(true).GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, record.CallResult.SignerNonce, int64(1))
}

func getCallData(chainId, nonce, maxPriorityGas, maxGas, gasLimitBytes, contractBytes, value, callDataBytes []byte, ecdsaPrivateKey PrivateKey) ([]byte, error) {
	objectsList := &RLPItem{}
	objectsList.AssignList()
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(chainId))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(nonce))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(maxPriorityGas))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(maxGas))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(gasLimitBytes))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(contractBytes))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(value))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(callDataBytes))
	objectsList.PushBack(NewRLPItem(LIST_TYPE))

	messageBytes, err := objectsList.Write()
	if err != nil {
		return nil, err
	}
	messageBytes = append([]byte{0x02}, messageBytes...)

	sig := ecdsaPrivateKey.Sign(messageBytes)

	r := sig[0:32]
	s := sig[32:64]
	v := ecdsaPrivateKey.GetRecoveryId(r, s, messageBytes)
	recIdBytes := []byte{byte(v)}

	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(recIdBytes))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(r))
	objectsList.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(s))

	messageBytes, err = objectsList.Write()
	if err != nil {
		return nil, err
	}
	messageBytes = append([]byte{0x02}, messageBytes...)

	return messageBytes, nil
}
