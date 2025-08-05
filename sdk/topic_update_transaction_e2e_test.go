//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const oldTopicMemo = "go-sdk::TestConsensusTopicUpdateTransaction_Execute::initial"

func TestIntegrationTopicUpdateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(oldTopicMemo).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	require.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, oldTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	newTopicMemo := "go-sdk::TestConsensusTopicUpdateTransaction_Execute::updated"

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTopicMemo(newTopicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, newTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTopicUpdateTransactionNoMemo(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(oldTopicMemo).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	require.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, oldTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTopicUpdateTransactionNoTopicID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(oldTopicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	require.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, oldTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	_, err = NewTopicUpdateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.Error(t, err)
	if err != nil {
		assert.ErrorContains(t, err, "exceptional precheck status INVALID_TOPIC_ID")
	}

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTopicUpdateTransactionClearFeeExemptKeys(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	feeExemptKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		AddFeeExemptKey(feeExemptKey).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	require.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, feeExemptKey.PublicKey().String(), info.FeeExemptKeys[0].String())

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		ClearFeeExemptKeys().
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Nil(t, info.FeeExemptKeys)
}

func TestIntegrationTopicUpdateTransactionClearCustomFees(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	customFee := NewCustomFixedFee().
		SetAmount(1).
		SetFeeCollectorAccountID(env.Client.GetOperatorAccountID())

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]*CustomFixedFee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	require.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, customFee, info.CustomFees[0])

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		ClearCustomFees().
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err = NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Nil(t, info.CustomFees)

}

// HIP-1139: Immutable Topics Tests
func TestIntegrationTopicUpdateTransactionPreventMessageSubmissionWhenSubmitKeyZeroKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a private topic with both Admin and Submit Keys
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	submitKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		SetSubmitKey(submitKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Verify initial message submission works
	frozenSubmit, err := NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Test message before zero key").
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenSubmit.Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update Submit Key to zero key using valid Admin Key signature
	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	frozenUpdate, err := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetSubmitKey(zeroKey).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify that no further messages can be submitted
	resp, err = NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Test message after zero key").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTopicUpdateTransactionAllowMessageSubmissionButPreventAdminUpdatesWhenAdminKeyEmptyList(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a private topic with both Admin and Submit Keys
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	submitKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		SetSubmitKey(submitKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Update Admin Key to empty key list using valid Admin Key signature
	zeroAccount, err := AccountIDFromString("0.0.0")
	require.NoError(t, err)

	frozenUpdate, err := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAdminKey(NewKeyList()).
		SetAutoRenewAccountID(zeroAccount).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify messages can still be submitted with the submit key
	frozenSubmit, err := NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Message after admin key dead").
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenSubmit.Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify that no further administrative updates can be made
	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetTopicMemo("Cannot update memo").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "UNAUTHORIZED")
}

func TestIntegrationTopicUpdateTransactionMakeTopicFullyImmutableWhenBothKeysZeroKeys(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a private topic with both Admin and Submit Keys
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	submitKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		SetSubmitKey(submitKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Update both Submit Key and Admin Key to zero keys with valid Admin Key signature
	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	zeroAccount, err := AccountIDFromString("0.0.0")
	require.NoError(t, err)

	frozenUpdate, err := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetSubmitKey(zeroKey).
		SetAdminKey(NewKeyList()).
		SetAutoRenewAccountID(zeroAccount).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify that message submission fails
	resp, err = NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Message should fail").
		Execute(env.Client)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")

	// Verify that administrative updates fail
	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetTopicMemo("Should fail").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "UNAUTHORIZED")
}

func TestIntegrationTopicUpdateTransactionSuccessfullyUpdateSubmitKeyToZeroKeyWhenTopicHasOnlySubmitKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a private topic with only Submit Key (no Admin Key)
	submitKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetSubmitKey(submitKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Verify initial message submission works
	frozenSubmit, err := NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Test message before zero key").
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenSubmit.Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update Submit Key to zero key with valid Submit Key signature
	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	frozenUpdate, err := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetSubmitKey(zeroKey).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify that no more messages can be submitted
	resp, err = NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Message should fail").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTopicUpdateTransactionMakePublicTopicAdministrativelyImmutableWhenAdminKeyEmptyList(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a public topic with Admin Key but no Submit Key
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Verify initial message submission works (no submit key required)
	resp, err = NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Public message before dead admin key").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Update Admin Key to empty list with valid Admin Key signature
	zeroAccount, err := AccountIDFromString("0.0.0")
	require.NoError(t, err)

	frozenUpdate, err := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAdminKey(NewKeyList()).
		SetAutoRenewAccountID(zeroAccount).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify message submission still works (topic remains public)
	resp, err = NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Public message after dead admin key").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify that administrative updates fail
	frozenUpdate, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetTopicMemo("Should fail").
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "UNAUTHORIZED")
}

func TestIntegrationTopicUpdateTransactionFailMessageSubmissionWhenSubmitKeyZero(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a topic with zero Submit Key
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		SetSubmitKey(zeroKey).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	resp, err = NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Should fail").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTopicUpdateTransactionFailToUpdateSubmitKeyToZeroKeyWithoutValidSubmitKeySignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a topic with Submit Key
	submitKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetSubmitKey(submitKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Attempt to update Submit Key without proper signature
	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetSubmitKey(zeroKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTopicUpdateTransactionFailToUpdateAdminKeyToZeroKeyWithoutValidAdminKeySignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a topic with Admin Key
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	zeroAccount, err := AccountIDFromString("0.0.0")
	require.NoError(t, err)

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAdminKey(NewKeyList()).
		SetAutoRenewAccountID(zeroAccount).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTopicUpdateTransactionSuccessfullyUpdateSubmitKeyToZeroKeyWithValidAdminKeySignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a topic with both Admin and Submit Keys
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	submitKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		SetSubmitKey(submitKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Sign(submitKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Update Submit Key to zero key using Admin Key signature (should succeed)
	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	frozenUpdate, err := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetSubmitKey(zeroKey).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the update was successful by checking topic info
	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, zeroKey.String(), info.SubmitKey.String())
}

func TestIntegrationTopicUpdateTransactionSuccessfullyUpdateSubmitKeyFromZeroKeyToValidKeyWithAdminKeySignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Create a topic with Admin Key and zero Submit Key
	adminKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	frozen, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		SetSubmitKey(zeroKey).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozen.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID

	// Update Submit Key from zero key to valid key using Admin Key signature
	newSubmitKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	frozenUpdate, err := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetSubmitKey(newSubmitKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenUpdate.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify the update was successful by submitting a message
	frozenSubmit, err := NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage("Message with restored submit key").
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = frozenSubmit.Sign(newSubmitKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify topic info shows the new key
	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, newSubmitKey.PublicKey().String(), info.SubmitKey.String())
}
