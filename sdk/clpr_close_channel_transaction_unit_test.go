//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitClprCloseChannelTransactionSignedRoundTrip(t *testing.T) {
	t.Parallel()

	channelID := clprHex(t, clprVectorChannelID)
	commitment := clprHex(t, clprOwnershipVectors[0].channelCommitment)

	for _, withCommitment := range []bool{false, true} {
		tx := NewClprCloseChannelTransaction().SetChannelID(channelID)
		if withCommitment {
			tx.SetOwnershipCommitment(commitment)
		}
		clprOffline(t, tx.Transaction)

		decoded := clprAssertSignedRoundTrip(t, tx).(ClprCloseChannelTransaction)
		assert.Equal(t, channelID, decoded.GetChannelID())
		if withCommitment {
			assert.Equal(t, commitment, decoded.GetOwnershipCommitment())
		} else {
			assert.Nil(t, decoded.GetOwnershipCommitment())
		}
	}
}

func TestUnitClprCloseChannelTransactionOmitsUnsetCommitment(t *testing.T) {
	t.Parallel()

	tx := NewClprCloseChannelTransaction().SetChannelID(clprBytes(32, 1))
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)

	body := tx.build().GetClprCloseChannel()
	assert.Nil(t, body.GetOwnershipCommitment())
}

func TestUnitClprCloseChannelTransactionUnsignedRoundTrip(t *testing.T) {
	t.Parallel()

	decoded, err := TransactionFromBytes(clprUnsignedBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprCloseChannel{
			ClprCloseChannel: &services.ClprCloseChannelTransactionBody{OwnershipCommitment: clprBytes(32, 2)},
		},
	}))
	require.NoError(t, err)
	result := decoded.(ClprCloseChannelTransaction)
	assert.False(t, result.IsFrozen())
	assert.Nil(t, result.GetChannelID())
	assert.Equal(t, clprBytes(32, 2), result.GetOwnershipCommitment())
}

func TestUnitClprCloseChannelTransactionRequiredFields(t *testing.T) {
	t.Parallel()

	tx := NewClprCloseChannelTransaction().SetOwnershipCommitment(clprBytes(32, 2))
	clprOffline(t, tx.Transaction)

	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.EqualError(t, err, "ClprCloseChannelTransaction: 'channelID' must be explicitly set before freezing or executing")
	assert.False(t, tx.IsFrozen())

	tx.SetChannelID(clprBytes(32, 1))
	_, err = tx.Freeze()
	assert.NoError(t, err)

	clprAssertFrozenRejected(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprCloseChannel{ClprCloseChannel: &services.ClprCloseChannelTransactionBody{}},
	}, "channelID")
}

func TestUnitClprCloseChannelTransactionRejectsBadLengths(t *testing.T) {
	t.Parallel()

	for _, length := range []int{0, 31, 33} {
		tx := NewClprCloseChannelTransaction().SetChannelID(clprBytes(length, 1))
		clprOffline(t, tx.Transaction)
		_, err := tx.Freeze()
		assert.ErrorContains(t, err, "channelID must be exactly 32 bytes")

		tx = NewClprCloseChannelTransaction().SetChannelID(clprBytes(32, 1)).SetOwnershipCommitment(clprBytes(length, 1))
		clprOffline(t, tx.Transaction)
		_, err = tx.Freeze()
		assert.ErrorContains(t, err, "ownershipCommitment must be exactly 32 bytes")
		assert.False(t, tx.IsFrozen())
	}
}

func TestUnitClprCloseChannelTransactionCopiesBytes(t *testing.T) {
	t.Parallel()

	channelID := clprBytes(32, 1)
	commitment := clprBytes(32, 2)
	tx := NewClprCloseChannelTransaction().SetChannelID(channelID).SetOwnershipCommitment(commitment)

	for _, bytes := range [][]byte{channelID, commitment, tx.GetChannelID(), tx.GetOwnershipCommitment()} {
		bytes[0] = 255
	}

	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
	assert.Equal(t, clprBytes(32, 2), tx.GetOwnershipCommitment())
}

func TestUnitClprCloseChannelTransactionFrozenSettersKeepState(t *testing.T) {
	t.Parallel()

	key := clprOwnershipVectors[1].key(t)
	tx := NewClprCloseChannelTransaction().SetChannelID(clprBytes(32, 1)).SetOwnershipCommitment(clprBytes(32, 2))
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(key)
	before, err := tx.ToBytes()
	require.NoError(t, err)

	tx.SetChannelID(clprBytes(32, 8)).SetOwnershipCommitment(clprBytes(32, 9))

	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
	assert.Equal(t, clprBytes(32, 2), tx.GetOwnershipCommitment())
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
	assert.True(t, key.PublicKey().VerifyTransaction(tx))
}

func TestUnitClprCloseChannelTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	_, err := NewClprCloseChannelTransaction().SetChannelID(clprBytes(32, 1)).Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprCloseChannelTransaction`")
}
