//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestUnitClprUpdateLedgerConfigurationTransactionSignedRoundTrip(t *testing.T) {
	t.Parallel()

	limits := clprThrottleLimits()
	limits.MaxSyncBytes = math.MaxInt64
	limits.MaxLocalEndpoints = 0
	limits.MaxPeerEndpoints = 0

	tx := NewClprUpdateLedgerConfigurationTransaction().SetThrottles(limits)
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)
	for _, vector := range clprOwnershipVectors {
		tx.Sign(vector.key(t))
	}

	body := tx.buildProtoBody()
	assert.True(t, proto.Equal(&services.ClprLedgerConfiguration{Throttles: limits._ToProtobuf()}, body.GetConfiguration()))

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)
	decoded, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	result := decoded.(ClprUpdateLedgerConfigurationTransaction)

	throttles, ok := result.GetThrottles()
	assert.True(t, ok)
	assert.Equal(t, limits, throttles)
	for _, vector := range clprOwnershipVectors {
		assert.True(t, vector.key(t).PublicKey().VerifyTransaction(result.childTransaction), vector.name)
	}
	roundTrip, err := result.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, txBytes, roundTrip)
}

func TestUnitClprUpdateLedgerConfigurationTransactionRequiresValidThrottles(t *testing.T) {
	t.Parallel()

	tx := NewClprUpdateLedgerConfigurationTransaction()
	clprOffline(t, tx.Transaction)
	_, ok := tx.GetThrottles()
	assert.False(t, ok)

	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "ClprUpdateLedgerConfigurationTransaction: 'throttles'")

	invalid := clprThrottleLimits()
	invalid.MaxQueueDepth = 0
	tx.SetThrottles(invalid)
	_, err = tx.Freeze()
	assert.ErrorIs(t, err, errClprThrottleOutOfRange)
	assert.False(t, tx.IsFrozen())

	tx.SetThrottles(clprThrottleLimits())
	_, err = tx.Freeze()
	assert.NoError(t, err)
}

func TestUnitClprUpdateLedgerConfigurationTransactionDecodedConfiguration(t *testing.T) {
	t.Parallel()

	configuration := clprLedgerConfigurationResponse(clprAnswerHeader(), nil).GetClprGetLedgerConfiguration().GetConfiguration()
	body := &services.TransactionBody{
		Data: &services.TransactionBody_ClprUpdateLedgerConfiguration{
			ClprUpdateLedgerConfiguration: &services.ClprUpdateLedgerConfigurationTransactionBody{Configuration: configuration},
		},
	}

	frozenBytes := clprSignedWireBytes(t, body)
	frozen, err := TransactionFromBytes(frozenBytes)
	require.NoError(t, err)
	frozenResult := frozen.(ClprUpdateLedgerConfigurationTransaction)
	roundTrip, err := frozenResult.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, frozenBytes, roundTrip)

	unfrozen, err := TransactionFromBytes(clprUnsignedBytes(t, body))
	require.NoError(t, err)
	tx := unfrozen.(ClprUpdateLedgerConfigurationTransaction).childTransaction
	assert.False(t, tx.IsFrozen())
	throttles, ok := tx.GetThrottles()
	assert.True(t, ok)
	assert.Equal(t, clprThrottleLimits(), throttles)
	assert.True(t, proto.Equal(&services.ClprLedgerConfiguration{Throttles: clprThrottleLimits()._ToProtobuf()}, tx.buildProtoBody().GetConfiguration()))
}

func TestUnitClprUpdateLedgerConfigurationTransactionReadModifySubmit(t *testing.T) {
	t.Parallel()

	response := clprLedgerConfigurationResponse(clprAnswerHeader(), nil)
	response.GetClprGetLedgerConfiguration().GetConfiguration().GetThrottles().MaxSyncBytes = math.MaxInt64
	info := _ClprLedgerConfigurationInfoFromProtobuf(response.GetClprGetLedgerConfiguration())

	throttles := info.Throttles
	throttles.MaxQueueDepth = 17
	tx := NewClprUpdateLedgerConfigurationTransaction().SetThrottles(throttles)
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)

	sent := _ClprThrottlesFromProtobuf(tx.buildProtoBody().GetConfiguration().GetThrottles())
	expected := clprThrottleLimits()
	expected.MaxSyncBytes = math.MaxInt64
	expected.MaxQueueDepth = 17
	assert.Equal(t, expected, sent)
	assert.Equal(t, uint32(16), info.Throttles.MaxQueueDepth)
}

func TestUnitClprUpdateLedgerConfigurationTransactionFrozenSetterKeepsState(t *testing.T) {
	t.Parallel()

	tx := NewClprUpdateLedgerConfigurationTransaction().SetThrottles(clprThrottleLimits())
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)
	before, err := tx.ToBytes()
	require.NoError(t, err)

	replacement := clprThrottleLimits()
	replacement.MaxQueueDepth = 17
	tx.SetThrottles(replacement)

	throttles, _ := tx.GetThrottles()
	assert.Equal(t, clprThrottleLimits(), throttles)
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestUnitClprUpdateLedgerConfigurationTransactionFrozenBytesRejectedBeforeExecute(t *testing.T) {
	t.Parallel()

	frozen, err := TransactionFromBytes(clprFrozenBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprUpdateLedgerConfiguration{
			ClprUpdateLedgerConfiguration: &services.ClprUpdateLedgerConfigurationTransactionBody{},
		},
	}))
	require.NoError(t, err)

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	_, err = TransactionExecute(frozen, client)
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "'throttles'")

	invalid, err := TransactionFromBytes(clprFrozenBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprUpdateLedgerConfiguration{
			ClprUpdateLedgerConfiguration: &services.ClprUpdateLedgerConfigurationTransactionBody{
				Configuration: &services.ClprLedgerConfiguration{Throttles: &services.ClprThrottles{}},
			},
		},
	}))
	require.NoError(t, err)

	_, err = TransactionExecute(invalid, client)
	assert.ErrorIs(t, err, errClprThrottleOutOfRange)
}

func TestUnitClprUpdateLedgerConfigurationTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	_, err := NewClprUpdateLedgerConfigurationTransaction().SetThrottles(clprThrottleLimits()).Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprUpdateLedgerConfigurationTransaction`")
}
