//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ADDRESS = "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"
)

func TestMirrorNodeContractQueryCanSimulateTransaction(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SIMPLE_SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)

	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)

	result, err := NewContractCallQuery().
		SetContractID(*contractID).
		SetGas(gas).
		SetFunction("getOwner", nil).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	simulationResult, err := NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%x", result.GetAddress(0)), simulationResult[26:])

	param, err := NewContractFunctionParameters().AddAddress(ADDRESS)
	require.NoError(t, err)

	gas, err = NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetFunction("addOwner", param).
		Execute(env.Client)
	require.NoError(t, err)

	response, err = NewContractExecuteTransaction().
		SetContractID(*contractID).
		SetGas(gas).
		SetFunction("addOwner", param).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("addOwner", param).
		Execute(env.Client)
	require.NoError(t, err)
}

func TestMirrorNodeContractQueryReturnsDefaultGasWhenContractIsNotDeployed(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	contractID := ContractID{Contract: 1231456}
	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, uint64(22892), gas)
}

func TestMirrorNodeContractQueryFailWhenGasLimitIsLow(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SIMPLE_SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)

	_, err = NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetGasLimit(100).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetGasLimit(100).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")
}

func TestMirrorNodeContractQueryFailWhenSenderIsNotSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SIMPLE_SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)
	param, err := NewContractFunctionParameters().AddAddress(ADDRESS)

	_, err = NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetFunction("addOwnerAndTransfer", param).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("addOwnerAndTransfer", param).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")
}
func TestMirrorNodeContractQueryCanSimulateWithSenderSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SIMPLE_SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	receiverId, _, err := createAccount(&env, func(transaction *AccountCreateTransaction) {
		pk, _ := PrivateKeyGenerateEd25519()
		transaction.SetKeyWithoutAlias(pk)
	})
	require.NoError(t, err)
	receiverEvmAddress := receiverId.ToSolidityAddress()

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)
	param, err := NewContractFunctionParameters().AddAddress(receiverEvmAddress)

	owner, err := NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)

	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetSenderEvmAddress(owner[26:]).
		SetFunction("addOwnerAndTransfer", param).
		SetValue(123).
		SetGasLimit(1_000_000).
		Execute(env.Client)
	require.NoError(t, err)

	resp, err := NewContractExecuteTransaction().
		SetContractID(*contractID).
		SetGas(gas).
		SetPayableAmount(NewHbar(1)).
		SetFunction("addOwnerAndTransfer", param).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetSenderEvmAddress(owner[26:]).
		SetFunction("addOwnerAndTransfer", param).
		SetValue(123).
		SetGasLimit(1_000_000).
		Execute(env.Client)
	require.NoError(t, err)
}
