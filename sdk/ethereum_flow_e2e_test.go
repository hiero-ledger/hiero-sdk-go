//go:build all || e2e
// +build all e2e

package hiero

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

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
