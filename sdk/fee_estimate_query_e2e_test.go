//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mirrorSyncDelayMillis = 2 * time.Second

func waitForMirrorNodeSync() {
	time.Sleep(mirrorSyncDelayMillis)
}

func subtotal(estimate FeeEstimate) uint64 {
	total := estimate.Base
	for _, extra := range estimate.Extras {
		total += extra.Subtotal
	}
	return total
}

func assertFeeComponentsPresent(t *testing.T, response FeeEstimateResponse) {
	require.NotNil(t, response)

	require.NotNil(t, response.NetworkFee)
	assert.Greater(t, response.NetworkFee.Multiplier, uint32(0))
	assert.GreaterOrEqual(t, response.NetworkFee.Subtotal, uint64(0))

	require.NotNil(t, response.NodeFee)
	assert.GreaterOrEqual(t, response.NodeFee.Base, uint64(0))
	require.NotNil(t, response.NodeFee.Extras)

	require.NotNil(t, response.ServiceFee)
	assert.GreaterOrEqual(t, response.ServiceFee.Base, uint64(0))
	require.NotNil(t, response.ServiceFee.Extras)

	require.NotNil(t, response.Notes)
	assert.Greater(t, response.Total, uint64(0))
}

func assertComponentTotalsConsistent(t *testing.T, response FeeEstimateResponse) {
	network := response.NetworkFee
	node := response.NodeFee
	service := response.ServiceFee

	nodeSubtotal := subtotal(node)
	serviceSubtotal := subtotal(service)

	assert.Equal(t, network.Subtotal, nodeSubtotal*uint64(network.Multiplier))
	assert.Equal(t, response.Total, network.Subtotal+nodeSubtotal+serviceSubtotal)
}

func TestIntegrationFeeEstimateQueryTokenCreateTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTokenCreateTransaction().
		SetTokenName("Test Token").
		SetTokenSymbol("TEST").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeState).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryTransferTransactionStateMode(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeState).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryTransferTransactionIntrinsicMode(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeIntrinsic).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryTransferTransactionDefaultModeIsState(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryTokenMintTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTokenMintTransaction().
		SetTokenID(TokenID{Token: 1234}).
		SetAmount(10).
		FreezeWith(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeIntrinsic).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	require.NotNil(t, response.NodeFee.Extras)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryTopicCreateTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTopicCreateTransaction().
		SetTopicMemo("integration test topic").
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeState).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryContractCreateTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewContractCreateTransaction().
		SetBytecode([]byte{1, 2, 3}).
		SetGas(1000).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeState).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryFileCreateTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetContents([]byte("integration test file")).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = transaction.SignWithOperator(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeState).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryFileAppendTransactionAggregatesChunks(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewFileAppendTransaction().
		SetFileID(FileID{File: 1234}).
		SetContents(make([]byte, 5000)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeIntrinsic).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryTopicMessageSubmitSingleChunk(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTopicMessageSubmitTransaction().
		SetTopicID(TopicID{Topic: 1234}).
		SetMessage(make([]byte, 128)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeIntrinsic).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryTopicMessageSubmitMultipleChunk(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	transaction, err := NewTopicMessageSubmitTransaction().
		SetTopicID(TopicID{Topic: 1234}).
		SetMessage(make([]byte, 5000)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	waitForMirrorNodeSync()

	response, err := NewFeeEstimateQuery().
		SetTransaction(transaction).
		SetMode(FeeEstimateModeIntrinsic).
		Execute(env.Client)
	require.NoError(t, err)

	assertFeeComponentsPresent(t, response)
	assertComponentTotalsConsistent(t, response)
}

func TestIntegrationFeeEstimateQueryWithoutTransactionThrowsError(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewFeeEstimateQuery().
		SetMode(FeeEstimateModeState).
		Execute(env.Client)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction is required")
}
