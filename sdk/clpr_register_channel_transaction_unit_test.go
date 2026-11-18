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

// clprAssertSignedRoundTrip freezes tx, signs it with both vector keys through the generic helpers and
// checks that the decoded bytes keep every node signature unchanged.
func clprAssertSignedRoundTrip(t *testing.T, tx TransactionInterface) TransactionInterface {
	t.Helper()
	_, err := TransactionFreezeWith(tx, nil)
	require.NoError(t, err)
	for _, vector := range clprOwnershipVectors {
		_, err = TransactionSign(tx, vector.key(t))
		require.NoError(t, err)
	}
	txBytes, err := TransactionToBytes(tx)
	require.NoError(t, err)

	decoded, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	// Count sig pairs on the wire rather than through GetSignatures.
	signedTransactions := decoded.getBaseTransaction().signedTransactions
	require.Equal(t, 2, signedTransactions._Length())
	for i := 0; i < signedTransactions._Length(); i++ {
		assert.Len(t, signedTransactions._Get(i).(*services.SignedTransaction).GetSigMap().GetSigPair(), 2)
	}
	for _, vector := range clprOwnershipVectors {
		assert.True(t, vector.key(t).PublicKey().VerifyTransaction(decoded))
	}

	roundTrip, err := TransactionToBytes(decoded)
	require.NoError(t, err)
	assert.Equal(t, txBytes, roundTrip)
	return decoded
}

// clprAssertFrozenRejected decodes body as frozen bytes and checks execution fails locally on field.
func clprAssertFrozenRejected(t *testing.T, body *services.TransactionBody, field string) {
	t.Helper()
	frozen, err := TransactionFromBytes(clprFrozenBytes(t, body))
	require.NoError(t, err)

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	_, err = TransactionExecute(frozen, client)
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "'"+field+"'")
}

func TestUnitClprRegisterChannelTransactionSignedRoundTrip(t *testing.T) {
	t.Parallel()

	commitment := clprHex(t, clprOwnershipVectors[0].channelCommitment)
	tx := NewClprRegisterChannelTransaction().SetOwnershipCommitment(commitment)
	clprOffline(t, tx.Transaction)

	decoded := clprAssertSignedRoundTrip(t, tx).(ClprRegisterChannelTransaction)
	assert.True(t, decoded.IsFrozen())
	assert.Equal(t, commitment, decoded.GetOwnershipCommitment())
	assert.Equal(t, []AccountID{{Account: 3}, {Account: 4}}, decoded.GetNodeAccountIDs())
	assert.Equal(t, tx.GetTransactionID().String(), decoded.GetTransactionID().String())
}

func TestUnitClprRegisterChannelTransactionUnsignedRoundTrip(t *testing.T) {
	t.Parallel()

	for _, commitment := range [][]byte{nil, clprBytes(32, 1)} {
		decoded, err := TransactionFromBytes(clprUnsignedBytes(t, &services.TransactionBody{
			Data: &services.TransactionBody_ClprRegisterChannel{
				ClprRegisterChannel: &services.ClprRegisterChannelTransactionBody{OwnershipCommitment: commitment},
			},
		}))
		require.NoError(t, err)
		result := decoded.(ClprRegisterChannelTransaction)
		assert.False(t, result.IsFrozen())
		assert.Equal(t, commitment, result.GetOwnershipCommitment())
	}
}

func TestUnitClprRegisterChannelTransactionRequiredFields(t *testing.T) {
	t.Parallel()

	tx := NewClprRegisterChannelTransaction()
	clprOffline(t, tx.Transaction)
	assert.Nil(t, tx.GetOwnershipCommitment())

	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.EqualError(t, err, "ClprRegisterChannelTransaction: 'ownershipCommitment' must be explicitly set before freezing or executing")
	assert.False(t, tx.IsFrozen())

	for _, length := range []int{0, 31, 33} {
		tx.SetOwnershipCommitment(clprBytes(length, 1))
		_, err = tx.Freeze()
		assert.ErrorContains(t, err, "ownershipCommitment must be exactly 32 bytes")
		assert.False(t, tx.IsFrozen())
	}

	tx.SetOwnershipCommitment(clprBytes(32, 1))
	_, err = tx.Freeze()
	assert.NoError(t, err)

	clprAssertFrozenRejected(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprRegisterChannel{ClprRegisterChannel: &services.ClprRegisterChannelTransactionBody{}},
	}, "ownershipCommitment")
}

func TestUnitClprRegisterChannelTransactionCopiesBytes(t *testing.T) {
	t.Parallel()

	commitment := clprBytes(32, 1)
	tx := NewClprRegisterChannelTransaction().SetOwnershipCommitment(commitment)
	commitment[0] = 0
	tx.GetOwnershipCommitment()[1] = 0

	assert.Equal(t, clprBytes(32, 1), tx.GetOwnershipCommitment())
}

func TestUnitClprRegisterChannelTransactionFrozenSetterKeepsState(t *testing.T) {
	t.Parallel()

	tx := NewClprRegisterChannelTransaction().SetOwnershipCommitment(clprBytes(32, 1))
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(clprOwnershipVectors[1].key(t))
	before, err := tx.ToBytes()
	require.NoError(t, err)

	tx.SetOwnershipCommitment(clprBytes(32, 2))

	assert.Equal(t, clprBytes(32, 1), tx.GetOwnershipCommitment())
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestUnitClprRegisterChannelTransactionCommitmentMatchesVectors(t *testing.T) {
	t.Parallel()

	channelID := clprHex(t, clprVectorChannelID)
	for _, vector := range clprOwnershipVectors {
		commitment, err := ComputeClprChannelOwnershipCommitment(channelID, vector.key(t).PublicKey())
		require.NoError(t, err)

		tx := NewClprRegisterChannelTransaction().SetOwnershipCommitment(commitment)
		assert.Equal(t, vector.channelCommitment, hex.EncodeToString(tx.GetOwnershipCommitment()))
	}
}

func TestUnitClprRegisterChannelTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	_, err := NewClprRegisterChannelTransaction().Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprRegisterChannelTransaction`")
}
