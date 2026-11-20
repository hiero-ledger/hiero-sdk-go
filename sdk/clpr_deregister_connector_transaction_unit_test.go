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

func clprReadyDeregisterConnector(t *testing.T) *ClprDeregisterConnectorTransaction {
	t.Helper()
	tx := NewClprDeregisterConnectorTransaction().
		SetChannelID(clprHex(t, clprVectorChannelID)).
		SetConnectorID(clprHex(t, clprOwnershipVectors[0].connectorID)).
		SetStakeRecipient(AccountID{Shard: 2, Realm: 3, Account: 9007199254740993})
	clprOffline(t, tx.Transaction)
	return tx
}

func TestUnitClprDeregisterConnectorTransactionRoundTrip(t *testing.T) {
	t.Parallel()

	tx := clprReadyDeregisterConnector(t)
	body := tx.buildProtoBody()
	assert.Equal(t, int64(9007199254740993), body.GetStakeRecipient().GetAccountNum())

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)
	decoded, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	result := decoded.(ClprDeregisterConnectorTransaction)

	assert.False(t, result.IsFrozen())
	assert.Equal(t, clprHex(t, clprVectorChannelID), result.GetChannelID())
	assert.Equal(t, clprHex(t, clprOwnershipVectors[0].connectorID), result.GetConnectorID())
	assert.Equal(t, "2.3.9007199254740993", result.GetStakeRecipient().String())

	roundTrip, err := result.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, txBytes, roundTrip)
}

func TestUnitClprDeregisterConnectorTransactionAdminAndRecipientSignatures(t *testing.T) {
	t.Parallel()

	admin, recipient := clprOwnershipVectors[0].key(t), clprOwnershipVectors[1].key(t)
	tx := clprReadyDeregisterConnector(t)
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(admin).Sign(recipient)

	txBytes, err := tx.ToBytes()
	require.NoError(t, err)
	decoded, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)
	result := decoded.(ClprDeregisterConnectorTransaction)

	assert.True(t, result.IsFrozen())
	assert.True(t, admin.PublicKey().VerifyTransaction(result))
	assert.True(t, recipient.PublicKey().VerifyTransaction(result))
	roundTrip, err := result.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, txBytes, roundTrip)
}

func TestUnitClprDeregisterConnectorTransactionEvmAliasRecipient(t *testing.T) {
	t.Parallel()

	const evmAddress = "0011223344556677889900112233445577889900"
	recipient, err := AccountIDFromEvmAddress(2, 3, evmAddress)
	require.NoError(t, err)
	tx := clprReadyDeregisterConnector(t).SetStakeRecipient(recipient)

	for i := range *recipient.AliasEvmAddress {
		(*recipient.AliasEvmAddress)[i] = 0xff
	}
	returned := tx.GetStakeRecipient()
	for i := range *returned.AliasEvmAddress {
		(*returned.AliasEvmAddress)[i] = 0xff
	}

	assert.Equal(t, evmAddress, hex.EncodeToString(*tx.GetStakeRecipient().AliasEvmAddress))
	assert.Equal(t, evmAddress, hex.EncodeToString(tx.buildProtoBody().GetStakeRecipient().GetAlias()))
}

func TestUnitClprDeregisterConnectorTransactionCopiesInputs(t *testing.T) {
	t.Parallel()

	channelID := clprBytes(32, 1)
	connectorID := clprBytes(32, 2)
	tx := NewClprDeregisterConnectorTransaction().SetChannelID(channelID).SetConnectorID(connectorID)

	for _, bytes := range [][]byte{channelID, connectorID, tx.GetChannelID(), tx.GetConnectorID()} {
		for i := range bytes {
			bytes[i] = 0
		}
	}

	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
	assert.Equal(t, clprBytes(32, 2), tx.GetConnectorID())
}

func TestUnitClprDeregisterConnectorTransactionRequiredFields(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	tx := NewClprDeregisterConnectorTransaction()
	clprOffline(t, tx.Transaction)

	steps := []struct {
		field string
		set   func()
	}{
		{"'channelID'", func() { tx.SetChannelID(clprBytes(32, 1)) }},
		{"'connectorID'", func() { tx.SetConnectorID(clprBytes(32, 2)) }},
		{"'stakeRecipient'", func() { tx.SetStakeRecipient(AccountID{Account: 9}) }},
	}
	for _, step := range steps {
		_, err := tx.FreezeWith(client)
		assert.ErrorIs(t, err, errClprFieldNotSet)
		assert.ErrorContains(t, err, "ClprDeregisterConnectorTransaction: "+step.field)
		assert.Equal(t, AccountID{}, tx.GetStakeRecipient())
		step.set()
	}

	_, err = tx.FreezeWith(client)
	require.NoError(t, err)
	assert.Equal(t, "0.0.9", tx.GetStakeRecipient().String())
}

func TestUnitClprDeregisterConnectorTransactionRejectsBadIDLength(t *testing.T) {
	t.Parallel()

	for _, length := range []int{0, 31, 33} {
		tx := clprReadyDeregisterConnector(t).SetConnectorID(clprBytes(length, 1))
		_, err := tx.Freeze()
		assert.ErrorContains(t, err, "connectorID must be exactly 32 bytes")
		assert.False(t, tx.IsFrozen())
	}
}

func TestUnitClprDeregisterConnectorTransactionFrozenSettersKeepState(t *testing.T) {
	t.Parallel()

	tx := clprReadyDeregisterConnector(t)
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(clprOwnershipVectors[1].key(t))
	before, err := tx.ToBytes()
	require.NoError(t, err)

	tx.SetChannelID(clprBytes(32, 8)).
		SetConnectorID(clprBytes(32, 8)).
		SetStakeRecipient(AccountID{Account: 1})

	assert.Equal(t, clprHex(t, clprVectorChannelID), tx.GetChannelID())
	assert.Equal(t, clprHex(t, clprOwnershipVectors[0].connectorID), tx.GetConnectorID())
	assert.Equal(t, "2.3.9007199254740993", tx.GetStakeRecipient().String())
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestUnitClprDeregisterConnectorTransactionFrozenBytesRejectedBeforeExecute(t *testing.T) {
	t.Parallel()

	frozen, err := TransactionFromBytes(clprFrozenBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprDeregisterConnector{
			ClprDeregisterConnector: &services.ClprDeregisterConnectorTransactionBody{
				ChannelId:   clprBytes(32, 0),
				ConnectorId: clprBytes(32, 0),
			},
		},
	}))
	require.NoError(t, err)
	assert.True(t, frozen.(ClprDeregisterConnectorTransaction).IsFrozen())

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	_, err = TransactionExecute(frozen, client)
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "'stakeRecipient'")
}

func TestUnitClprDeregisterConnectorTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	_, err := clprReadyDeregisterConnector(t).Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprDeregisterConnectorTransaction`")
}
