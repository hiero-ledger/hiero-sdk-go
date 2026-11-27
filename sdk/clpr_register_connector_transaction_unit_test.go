//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitClprRegisterConnectorTransactionCommitmentMatchesVectors(t *testing.T) {
	t.Parallel()

	channelID := clprHex(t, clprVectorChannelID)
	salt := clprHex(t, clprVectorSalt)

	for _, vector := range clprOwnershipVectors {
		t.Run(vector.name, func(t *testing.T) {
			publicKey := vector.key(t).PublicKey()
			connectorID, err := ComputeClprConnectorID(channelID, publicKey, salt)
			require.NoError(t, err)
			commitment, err := ComputeClprConnectorOwnershipCommitment(connectorID, publicKey)
			require.NoError(t, err)

			tx := NewClprRegisterConnectorTransaction().SetCommitment(commitment).SetChannelID(channelID)
			assert.Equal(t, vector.connectorID, hex.EncodeToString(connectorID))
			assert.Equal(t, vector.connectorCommitment, hex.EncodeToString(tx.GetCommitment()))
			assert.Equal(t, clprHex(t, clprVectorChannelID), channelID)
			assert.Equal(t, clprHex(t, clprVectorSalt), salt)
		})
	}
}

func TestUnitClprRegisterConnectorTransactionSignedRoundTrip(t *testing.T) {
	t.Parallel()

	commitment := clprHex(t, clprOwnershipVectors[0].connectorCommitment)
	channelID := clprHex(t, clprVectorChannelID)
	tx := NewClprRegisterConnectorTransaction().SetCommitment(commitment).SetChannelID(channelID)
	clprOffline(t, tx.Transaction)

	decoded := clprAssertSignedRoundTrip(t, tx).(ClprRegisterConnectorTransaction)
	assert.Equal(t, commitment, decoded.GetCommitment())
	assert.Equal(t, channelID, decoded.GetChannelID())
	body := decoded.build().GetClprRegisterConnector()
	assert.Equal(t, commitment, body.GetCommitment())
	assert.Equal(t, channelID, body.GetChannelId())
}

func TestUnitClprRegisterConnectorTransactionRequiredFields(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name       string
		commitment []byte
		channelID  []byte
		missing    string
	}{
		{"none", nil, nil, "commitment"},
		{"channelID only", nil, clprBytes(32, 1), "commitment"},
		{"commitment only", clprBytes(32, 2), nil, "channelID"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := TransactionFromBytes(clprUnsignedBytes(t, &services.TransactionBody{
				Data: &services.TransactionBody_ClprRegisterConnector{
					ClprRegisterConnector: &services.ClprRegisterConnectorTransactionBody{Commitment: tt.commitment, ChannelId: tt.channelID},
				},
			}))
			require.NoError(t, err)
			tx := decoded.(ClprRegisterConnectorTransaction).childTransaction
			assert.False(t, tx.IsFrozen())
			assert.Equal(t, tt.commitment, tx.GetCommitment())
			assert.Equal(t, tt.channelID, tx.GetChannelID())

			_, err = tx.Freeze()
			assert.ErrorIs(t, err, errClprFieldNotSet)
			assert.EqualError(t, err, "ClprRegisterConnectorTransaction: '"+tt.missing+"' must be explicitly set before freezing or executing")

			tx.SetCommitment(clprBytes(32, 2)).SetChannelID(clprBytes(32, 1))
			_, err = tx.Freeze()
			assert.NoError(t, err)

			clprAssertFrozenRejected(t, &services.TransactionBody{
				Data: &services.TransactionBody_ClprRegisterConnector{
					ClprRegisterConnector: &services.ClprRegisterConnectorTransactionBody{Commitment: tt.commitment, ChannelId: tt.channelID},
				},
			}, tt.missing)
		})
	}
}

func TestUnitClprRegisterConnectorTransactionRejectsBadLengths(t *testing.T) {
	t.Parallel()

	for _, length := range []int{0, 31, 33} {
		tx := NewClprRegisterConnectorTransaction().SetCommitment(clprBytes(length, 2)).SetChannelID(clprBytes(32, 1))
		clprOffline(t, tx.Transaction)
		_, err := tx.Freeze()
		assert.ErrorContains(t, err, "commitment must be exactly 32 bytes")

		tx = NewClprRegisterConnectorTransaction().SetCommitment(clprBytes(32, 2)).SetChannelID(clprBytes(length, 1))
		clprOffline(t, tx.Transaction)
		_, err = tx.Freeze()
		assert.ErrorContains(t, err, "channelID must be exactly 32 bytes")
		assert.False(t, tx.IsFrozen())
	}
}

func TestUnitClprRegisterConnectorTransactionCopiesBytes(t *testing.T) {
	t.Parallel()

	commitment := clprBytes(32, 2)
	channelID := clprBytes(32, 1)
	tx := NewClprRegisterConnectorTransaction().SetCommitment(commitment).SetChannelID(channelID)

	for _, bytes := range [][]byte{commitment, channelID, tx.GetCommitment(), tx.GetChannelID()} {
		bytes[0] = 255
	}

	assert.Equal(t, clprBytes(32, 2), tx.GetCommitment())
	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
}

func TestUnitClprRegisterConnectorTransactionFrozenSettersKeepState(t *testing.T) {
	t.Parallel()

	key := clprOwnershipVectors[0].key(t)
	tx := NewClprRegisterConnectorTransaction().SetCommitment(clprBytes(32, 2)).SetChannelID(clprBytes(32, 1))
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(key)
	before, err := tx.ToBytes()
	require.NoError(t, err)

	tx.SetCommitment(clprBytes(32, 8)).SetChannelID(clprBytes(32, 9))

	assert.Equal(t, clprBytes(32, 2), tx.GetCommitment())
	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
	assert.True(t, key.PublicKey().VerifyTransaction(tx))
}

func TestUnitClprRegisterConnectorTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	_, err := NewClprRegisterConnectorTransaction().Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprRegisterConnectorTransaction`")
}
