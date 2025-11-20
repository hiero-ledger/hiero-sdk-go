//go:build all || e2e || testnets

package hiero

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestIntegrationTopicMessageSubmitTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	resp, err = NewTopicMessageSubmitTransaction().
		SetMessage(strings.Repeat("a", 100)).
		SetChunkSize(50).
		SetMaxChunks(2).
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTopicMessageSubmitTransactionInvalidChunkSize(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	message := strings.Repeat("a", 100)

	resp, err = NewTopicMessageSubmitTransaction().
		SetMessage(message).
		SetChunkSize(10).
		SetMaxChunks(2).
		SetTopicID(topicID).
		Execute(env.Client)
	require.Error(t, err)

	resp, err = NewTopicMessageSubmitTransaction().
		SetMessage(message).
		SetChunkSize(0).
		SetTopicID(topicID).
		Execute(env.Client)
	require.Error(t, err)

	resp, err = NewTopicMessageSubmitTransaction().
		SetMessage(message).
		SetMaxChunks(0).
		SetTopicID(topicID).
		Execute(env.Client)
	require.Error(t, err)
}

func TestIntegrationTopicMessageSubmitTransactinoWrongMessageType(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	resp, err = NewTopicMessageSubmitTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMessage(1234). // wrong message type
		SetTopicID(topicID).
		Execute(env.Client)
	require.ErrorContains(t, err, "no transactions to execute")
}
