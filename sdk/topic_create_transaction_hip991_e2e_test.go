//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestIntegrationRevenueGeneratingTopicCanCreateAndUpdate(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	exemptKey1, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	exemptKey2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	customFeeTokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	customFeeTokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	customFee1 := NewCustomFixedFee().
		SetAmount(1).
		SetDenominatingTokenID(customFeeTokenID).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	customFee2 := NewCustomFixedFee().
		SetAmount(2).
		SetDenominatingTokenID(customFeeTokenID2).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	// Create a revenue generating topic
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFeeExemptKeys([]Key{exemptKey1.PublicKey(), exemptKey2.PublicKey()}).
		SetCustomFees([]*CustomFixedFee{customFee1, customFee2}).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.FeeScheduleKey.String())
	assert.Equal(t, exemptKey1.PublicKey().String(), info.FeeExemptKeys[0].String())
	assert.Equal(t, exemptKey2.PublicKey().String(), info.FeeExemptKeys[1].String())
	assert.Equal(t, customFee1.String(), info.CustomFees[0].String())
	assert.Equal(t, customFee2.String(), info.CustomFees[1].String())

	newFeeScheduleKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Update the revenue generating topic
	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetFeeScheduleKey(newFeeScheduleKey).
		SetFeeExemptKeys([]Key{exemptKey2.PublicKey()}).
		SetCustomFees([]*CustomFixedFee{customFee2}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())
	assert.Equal(t, newFeeScheduleKey.PublicKey().String(), info.FeeScheduleKey.String())
	assert.Equal(t, exemptKey2.PublicKey().String(), info.FeeExemptKeys[0].String())
	assert.True(t, len(info.FeeExemptKeys) == 1)
	assert.Equal(t, customFee2.String(), info.CustomFees[0].String())
	assert.True(t, len(info.CustomFees) == 1)
}

func TestIntegrationRevenueGeneratingTopicCannotCreateWithInvalidExemptKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	exemptKey1, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Duplicate exempt key - fails with FEE_EXEMPT_KEY_LIST_CONTAINS_DUPLICATED_KEYS
	_, err = NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeExemptKeys([]Key{exemptKey1, exemptKey1}).
		Execute(env.Client)
	require.ErrorContains(t, err, "exceptional precheck status FEE_EXEMPT_KEY_LIST_CONTAINS_DUPLICATED_KEYS")

	// Create key with invalid length
	invalidKey := &PrivateKey{ed25519PrivateKey: &_Ed25519PrivateKey{keyData: make([]byte, 34)}}

	// Invalid exempt key - fails with INVALID_KEY_IN_FEE_EXEMPT_KEY_LIST
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeExemptKeys([]Key{invalidKey}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_KEY_IN_FEE_EXEMPT_KEY_LIST")

	exemptKey2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey3, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey4, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey5, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey6, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey7, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey8, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey9, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey10, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	exemptKey11, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Set 11 keys - fails with MAX_ENTRIES_FOR_FEE_EXEMPT_KEY_LIST_EXCEEDED
	resp, err = NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeExemptKeys([]Key{exemptKey1, exemptKey2, exemptKey3, exemptKey4, exemptKey5, exemptKey6, exemptKey7, exemptKey8, exemptKey9, exemptKey10, exemptKey11}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: MAX_ENTRIES_FOR_FEE_EXEMPT_KEY_LIST_EXCEEDED")
}

func TestIntegrationRevenueGeneratingTopicCannotUpdateFeeScheduleKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a revenue generating topic without fee schedule key
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	newFeeScheduleKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Update the revenue generating topic with new fee schedule key - fails with FEE_SCHEDULE_KEY_CANNOT_BE_UPDATED
	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetFeeScheduleKey(newFeeScheduleKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: FEE_SCHEDULE_KEY_CANNOT_BE_UPDATED")
}

func TestIntegrationRevenueGeneratingTopicCannotUpdateCustomFees(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a revenue generating topic without fee schedule key
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	customFeeTokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	customFeeTokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	customFee1 := NewCustomFixedFee().
		SetAmount(1).
		SetDenominatingTokenID(customFeeTokenID).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	customFee2 := NewCustomFixedFee().
		SetAmount(2).
		SetDenominatingTokenID(customFeeTokenID2).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	// Update the revenue generating topic with new custom fees - fails with FEE_SCHEDULE_KEY_NOT_SET
	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetCustomFees([]*CustomFixedFee{customFee1, customFee2}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: FEE_SCHEDULE_KEY_NOT_SET")
}

func TestIntegrationRevenueGeneratingTopicCanChargeHbarsWithLimit(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	var hbar int64 = 100_000_000
	customFee := NewCustomFixedFee().
		SetAmount(hbar / 2).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	// Create a revenue generating topic
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		AddCustomFee(customFee).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	// Create payer with 1 Hbar
	payerId, payerPrivateKey, err := createAccount(&env)
	require.NoError(t, err)

	customFeeLimit := NewCustomFeeLimit().
		SetPayerId(payerId).
		AddCustomFee(NewCustomFixedFee().SetAmount(hbar))

	// Submit a message to the revenue generating topic with custom fee limit
	env.Client.SetOperator(payerId, payerPrivateKey)
	frozenTxn, err := NewTopicMessageSubmitTransaction().
		SetMessage("message").
		SetTopicID(topicID).
		AddCustomFeeLimit(customFeeLimit).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenTxn.Sign(payerPrivateKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// Verify the custom fee charged
	accountInfo, err := NewAccountInfoQuery().
		SetAccountID(payerId).
		Execute(env.Client)
	require.NoError(t, err)

	assert.True(t, accountInfo.Balance.AsTinybar() < hbar/2)
}

func TestIntegrationRevenueGeneratingTopicCanChargeHbarsWithoutLimit(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	var hbar int64 = 100_000_000
	customFee := NewCustomFixedFee().
		SetAmount(hbar / 2).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	// Create a revenue generating topic
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		AddCustomFee(customFee).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	// Create payer with 1 Hbar
	payerId, payerPrivateKey, err := createAccount(&env)
	require.NoError(t, err)

	// Submit a message to the revenue generating topic with custom fee limit
	env.Client.SetOperator(payerId, payerPrivateKey)
	resp, err = NewTopicMessageSubmitTransaction().
		SetMessage("message").
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// Verify the custom fee charged
	accountInfo, err := NewAccountInfoQuery().
		SetAccountID(payerId).
		Execute(env.Client)
	require.NoError(t, err)

	assert.True(t, accountInfo.Balance.AsTinybar() < hbar/2)
}

func TestIntegrationRevenueGeneratingTopicCanChargeTokensWithLimit(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenId, err := createFungibleToken(&env)
	require.NoError(t, err)

	customFee := NewCustomFixedFee().
		SetAmount(1).
		SetDenominatingTokenID(tokenId).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	// Create a revenue generating topic
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		AddCustomFee(customFee).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	// Create payer
	payerId, payerPrivateKey, err := createAccount(&env, func(transaction *AccountCreateTransaction) {
		transaction.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// Send tokens to payer
	resp, err = NewTransferTransaction().
		AddTokenTransfer(tokenId, env.Client.GetOperatorAccountID(), -1).
		AddTokenTransfer(tokenId, payerId, 1).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFeeLimit := NewCustomFeeLimit().
		SetPayerId(payerId).
		AddCustomFee(NewCustomFixedFee().SetAmount(2).SetDenominatingTokenID(tokenId))

	// Submit a message to the revenue generating topic with custom fee limit
	env.Client.SetOperator(payerId, payerPrivateKey)
	resp, err = NewTopicMessageSubmitTransaction().
		SetMessage("message").
		SetTopicID(topicID).
		AddCustomFeeLimit(customFeeLimit).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// Verify the custom fee charged
	accountBalance, err := NewAccountBalanceQuery().
		SetAccountID(payerId).
		Execute(env.Client)
	require.NoError(t, err)

	assert.True(t, accountBalance.Tokens.Get(tokenId) == 0)
}

func TestIntegrationRevenueGeneratingTopicCanChargeTokensWithoutLimit(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenId, err := createFungibleToken(&env)
	require.NoError(t, err)

	customFee := NewCustomFixedFee().
		SetAmount(1).
		SetDenominatingTokenID(tokenId).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	// Create a revenue generating topic
	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		AddCustomFee(customFee).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	// Create payer
	payerId, payerPrivateKey, err := createAccount(&env, func(transaction *AccountCreateTransaction) {
		transaction.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// Send tokens to payer
	resp, err = NewTransferTransaction().
		AddTokenTransfer(tokenId, env.Client.GetOperatorAccountID(), -1).
		AddTokenTransfer(tokenId, payerId, 1).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Submit a message to the revenue generating topic with custom fee limit
	env.Client.SetOperator(payerId, payerPrivateKey)
	resp, err = NewTopicMessageSubmitTransaction().
		SetMessage("message").
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// Verify the custom fee charged
	accountBalance, err := NewAccountBalanceQuery().
		SetAccountID(payerId).
		Execute(env.Client)
	require.NoError(t, err)

	assert.True(t, accountBalance.Tokens.Get(tokenId) == 0)
}
