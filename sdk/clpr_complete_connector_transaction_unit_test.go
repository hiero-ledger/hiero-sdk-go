//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"math"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// clprPayerSignedBytes serializes a body signed by the Ed25519 vector key; ToBytes needs a sigMap.
func clprPayerSignedBytes(t *testing.T, body *services.TransactionBody) []byte {
	t.Helper()
	key := clprOwnershipVectors[1].key(t)
	bodyBytes, err := proto.Marshal(clprWireBody(body))
	require.NoError(t, err)
	signed, err := proto.Marshal(&services.SignedTransaction{
		BodyBytes: bodyBytes,
		SigMap: &services.SignatureMap{SigPair: []*services.SignaturePair{{
			PubKeyPrefix: key.PublicKey().BytesRaw(),
			Signature:    &services.SignaturePair_Ed25519{Ed25519: key.Sign(bodyBytes)},
		}}},
	})
	require.NoError(t, err)
	list, err := proto.Marshal(&sdk.TransactionList{TransactionList: []*services.Transaction{{SignedTransactionBytes: signed}}})
	require.NoError(t, err)
	return list
}

func clprConnectorAdminKey(t *testing.T) *KeyList {
	t.Helper()
	return KeyListWithThreshold(1).
		Add(clprOwnershipVectors[0].key(t).PublicKey()).
		Add(clprOwnershipVectors[1].key(t).PublicKey())
}

// clprBoundCompleteConnector sets the fields the ownership proof derives from.
func clprBoundCompleteConnector(t *testing.T, vector clprOwnershipVector) *ClprCompleteConnectorTransaction {
	t.Helper()
	tx := NewClprCompleteConnectorTransaction().
		SetChannelID(clprHex(t, clprVectorChannelID)).
		SetConnectorID(clprHex(t, vector.connectorID)).
		SetSalt(clprHex(t, clprVectorSalt))
	clprOffline(t, tx.Transaction)
	return tx
}

func clprReadyCompleteConnector(t *testing.T, vector clprOwnershipVector) *ClprCompleteConnectorTransaction {
	t.Helper()
	tx := clprBoundCompleteConnector(t, vector).
		SetConnectorContract(ContractID{Shard: 1, Realm: 2, Contract: 345}).
		SetAdminKey(clprConnectorAdminKey(t))
	_, err := tx.SetOwnershipProof(vector.key(t))
	require.NoError(t, err)
	return tx
}

func clprCompleteConnectorWire(t *testing.T, vector clprOwnershipVector, lockedStake uint64) *services.TransactionBody {
	t.Helper()
	return &services.TransactionBody{
		Data: &services.TransactionBody_ClprCompleteConnector{
			ClprCompleteConnector: &services.ClprCompleteConnectorTransactionBody{
				ChannelId:         clprHex(t, clprVectorChannelID),
				ConnectorId:       clprHex(t, vector.connectorID),
				Salt:              clprHex(t, clprVectorSalt),
				PublicKey:         clprHex(t, vector.publicKey),
				Signature:         clprHex(t, vector.connectorSignature),
				SignatureScheme:   services.ClprSignatureScheme(vector.scheme),
				ConnectorContract: ContractID{Shard: 1, Realm: 2, Contract: 345}._ToProtobuf(),
				AdminKey:          clprConnectorAdminKey(t)._ToProtoKey(),
				LockedStake:       lockedStake,
			},
		},
	}
}

func TestUnitClprCompleteConnectorTransactionOwnershipProof(t *testing.T) {
	t.Parallel()

	for _, vector := range clprOwnershipVectors {
		t.Run(vector.name, func(t *testing.T) {
			tx := clprReadyCompleteConnector(t, vector)

			assert.Equal(t, vector.connectorDigest, hex.EncodeToString(_ClprKeccak(_ClprConnectorOwnershipPreimage(clprHex(t, vector.connectorID)))))
			assert.Equal(t, vector.publicKey, hex.EncodeToString(tx.GetOwnershipPublicKey()))
			assert.Equal(t, vector.connectorSignature, hex.EncodeToString(tx.GetOwnershipSignature()))
			scheme, ok := tx.GetOwnershipSignatureScheme()
			assert.True(t, ok)
			assert.Equal(t, vector.scheme, scheme)
			assertClprSignatureVerifies(t, vector, clprHex(t, vector.connectorDigest), tx.GetOwnershipSignature())
		})
	}
}

func TestUnitClprCompleteConnectorTransactionRejectsMismatchedConnectorID(t *testing.T) {
	t.Parallel()

	ecdsa, ed25519 := clprOwnershipVectors[0], clprOwnershipVectors[1]

	wrongSalt := clprBoundCompleteConnector(t, ed25519).SetSalt(clprBytes(32, 0))
	_, err := wrongSalt.SetOwnershipProof(ed25519.key(t))
	assert.ErrorIs(t, err, errClprConnectorIDMismatch)
	assert.Nil(t, wrongSalt.GetOwnershipSignature())

	otherKey := clprBoundCompleteConnector(t, ecdsa)
	_, err = otherKey.SetOwnershipProof(ed25519.key(t))
	assert.ErrorIs(t, err, errClprConnectorIDMismatch)
	assert.Nil(t, otherKey.GetOwnershipPublicKey())
	_, ok := otherKey.GetOwnershipSignatureScheme()
	assert.False(t, ok)

	_, err = otherKey.SetOwnershipProof(PrivateKey{})
	assert.ErrorIs(t, err, errClprPrivateKeyRequired)

	_, err = otherKey.SetOwnershipProof(ecdsa.key(t))
	assert.NoError(t, err)
}

func TestUnitClprCompleteConnectorTransactionBoundFields(t *testing.T) {
	t.Parallel()

	vector := clprOwnershipVectors[1]
	fields := []struct {
		name  string
		value []byte
		set   func(*ClprCompleteConnectorTransaction, []byte) *ClprCompleteConnectorTransaction
	}{
		{"channelID", clprHex(t, clprVectorChannelID), (*ClprCompleteConnectorTransaction).SetChannelID},
		{"connectorID", clprHex(t, vector.connectorID), (*ClprCompleteConnectorTransaction).SetConnectorID},
		{"salt", clprHex(t, clprVectorSalt), (*ClprCompleteConnectorTransaction).SetSalt},
	}

	for i, field := range fields {
		t.Run(field.name, func(t *testing.T) {
			missing := NewClprCompleteConnectorTransaction()
			for j, other := range fields {
				if j != i {
					other.set(missing, other.value)
				}
			}
			_, err := missing.SetOwnershipProof(vector.key(t))
			assert.ErrorIs(t, err, errClprFieldNotSet)
			assert.ErrorContains(t, err, "'"+field.name+"'")
			assert.Nil(t, missing.GetOwnershipSignature())

			tx := clprReadyCompleteConnector(t, vector)
			signature := tx.GetOwnershipSignature()

			field.set(tx, field.value)
			assert.Equal(t, signature, tx.GetOwnershipSignature())

			field.set(tx, clprBytes(32, 8))
			assert.Nil(t, tx.GetOwnershipPublicKey())
			assert.Nil(t, tx.GetOwnershipSignature())
			_, ok := tx.GetOwnershipSignatureScheme()
			assert.False(t, ok)

			field.set(tx, field.value)
			_, err = tx.Freeze()
			assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)
			assert.False(t, tx.IsFrozen())

			_, err = tx.SetOwnershipProof(vector.key(t))
			require.NoError(t, err)
			_, err = tx.Freeze()
			assert.NoError(t, err)
		})
	}
}

func TestUnitClprCompleteConnectorTransactionSignedRoundTrip(t *testing.T) {
	t.Parallel()

	for _, vector := range clprOwnershipVectors {
		t.Run(vector.name, func(t *testing.T) {
			tx := clprReadyCompleteConnector(t, vector).
				SetConnectorContract(ContractID{Shard: 1, Realm: 2, Contract: 345}).
				SetLockedStake(HbarFromTinybar(math.MaxInt64))
			_, err := tx.Freeze()
			require.NoError(t, err)
			tx.Sign(clprOwnershipVectors[1].key(t))

			txBytes, err := tx.ToBytes()
			require.NoError(t, err)
			decoded, err := TransactionFromBytes(txBytes)
			require.NoError(t, err)
			result := decoded.(ClprCompleteConnectorTransaction)

			assert.True(t, result.IsFrozen())
			assert.Equal(t, clprHex(t, clprVectorChannelID), result.GetChannelID())
			assert.Equal(t, clprHex(t, vector.connectorID), result.GetConnectorID())
			assert.Equal(t, clprHex(t, clprVectorSalt), result.GetSalt())
			assert.Equal(t, "1.2.345", result.GetConnectorContract().String())
			assert.Equal(t, uint64(math.MaxInt64), result.GetLockedStakeTinybars())
			assert.Equal(t, vector.connectorSignature, hex.EncodeToString(result.GetOwnershipSignature()))

			adminKey, ok := result.GetAdminKey().(*KeyList)
			require.True(t, ok)
			assert.Equal(t, 1, adminKey.threshold)
			assert.Equal(t, clprConnectorAdminKey(t).String(), adminKey.String())

			roundTrip, err := result.ToBytes()
			require.NoError(t, err)
			assert.Equal(t, txBytes, roundTrip)

			signatures, err := result.GetSignatures()
			require.NoError(t, err)
			assert.Len(t, signatures, 2)
			assert.True(t, clprOwnershipVectors[1].key(t).PublicKey().VerifyTransaction(result))
		})
	}
}

func TestUnitClprCompleteConnectorTransactionDecodedStakeKeepsUint64Range(t *testing.T) {
	t.Parallel()

	for _, stake := range []uint64{1 << 63, math.MaxUint64} {
		txBytes := clprPayerSignedBytes(t, clprCompleteConnectorWire(t, clprOwnershipVectors[1], stake))
		decoded, err := TransactionFromBytes(txBytes)
		require.NoError(t, err)
		result := decoded.(ClprCompleteConnectorTransaction)

		assert.Equal(t, stake, result.GetLockedStakeTinybars())
		assert.Equal(t, HbarFromTinybar(math.MaxInt64), result.GetLockedStake())

		roundTrip, err := result.ToBytes()
		require.NoError(t, err)
		assert.Equal(t, txBytes, roundTrip)
		assert.True(t, clprOwnershipVectors[1].key(t).PublicKey().VerifyTransaction(result))
	}
}

func TestUnitClprCompleteConnectorTransactionLockedStake(t *testing.T) {
	t.Parallel()

	vector := clprOwnershipVectors[1]
	tx := clprReadyCompleteConnector(t, vector)
	assert.Equal(t, uint64(0), tx.GetLockedStakeTinybars())
	assert.Equal(t, uint64(0), tx.buildProtoBody().GetLockedStake())

	tx.SetLockedStake(HbarFromTinybar(7)).SetLockedStake(HbarFromTinybar(-1))
	assert.Equal(t, HbarFromTinybar(7), tx.GetLockedStake())
	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprLockedStakeNegative)
	assert.ErrorContains(t, err, "ClprCompleteConnectorTransaction")

	tx.SetLockedStake(HbarFromTinybar(math.MaxInt64))
	_, err = tx.Freeze()
	require.NoError(t, err)
	assert.Equal(t, uint64(math.MaxInt64), tx.buildProtoBody().GetLockedStake())
}

func TestUnitClprCompleteConnectorTransactionAdminKeyKinds(t *testing.T) {
	t.Parallel()

	ecdsaKey := clprOwnershipVectors[0].key(t).PublicKey()
	ed25519Key := clprOwnershipVectors[1].key(t).PublicKey()
	evmDelegate, err := DelegatableContractIDFromEvmAddress(1, 2, "0011223344556677889900112233445577889900")
	require.NoError(t, err)

	for _, adminKey := range []Key{
		ed25519Key,
		ecdsaKey,
		clprConnectorAdminKey(t),
		NewKeyList().Add(ed25519Key).Add(ecdsaKey),
		ContractID{Shard: 1, Realm: 2, Contract: 345},
		DelegatableContractID{Shard: 1, Realm: 2, Contract: 346},
		evmDelegate,
		KeyListWithThreshold(1).Add(ed25519Key).Add(evmDelegate),
	} {
		tx := NewClprCompleteConnectorTransaction().SetAdminKey(adminKey)
		decoded, err := TransactionFromBytes(clprUnsignedBytes(t, &services.TransactionBody{
			Data: &services.TransactionBody_ClprCompleteConnector{ClprCompleteConnector: tx.buildProtoBody()},
		}))
		require.NoError(t, err)
		result := decoded.(ClprCompleteConnectorTransaction)

		assert.Equal(t, adminKey._ToProtoKey().String(), result.GetAdminKey()._ToProtoKey().String())
		assert.Equal(t, adminKey._ToProtoKey().String(), tx.buildProtoBody().GetAdminKey().String())
	}
}

func TestUnitClprCompleteConnectorTransactionCopiesInputs(t *testing.T) {
	t.Parallel()

	adminKey := clprConnectorAdminKey(t)
	contract := ContractID{Shard: 1, Realm: 2, EvmAddress: clprHex(t, "0011223344556677889900112233445577889900")}
	tx := clprReadyCompleteConnector(t, clprOwnershipVectors[1]).
		SetAdminKey(adminKey).
		SetConnectorContract(contract)

	adminKey.Add(clprOwnershipVectors[0].key(t).PublicKey())
	returned := tx.GetAdminKey().(*KeyList)
	returned.Add(clprOwnershipVectors[0].key(t).PublicKey())
	for _, bytes := range [][]byte{contract.EvmAddress, tx.GetConnectorContract().EvmAddress,
		tx.GetChannelID(), tx.GetConnectorID(), tx.GetSalt(), tx.GetOwnershipPublicKey(), tx.GetOwnershipSignature()} {
		for i := range bytes {
			bytes[i] = 0
		}
	}

	assert.Equal(t, clprConnectorAdminKey(t).String(), tx.GetAdminKey().String())
	assert.Equal(t, "0011223344556677889900112233445577889900", hex.EncodeToString(tx.GetConnectorContract().EvmAddress))
	assert.Equal(t, "0011223344556677889900112233445577889900", hex.EncodeToString(tx.buildProtoBody().GetConnectorContract().GetEvmAddress()))
	assert.Equal(t, clprHex(t, clprVectorChannelID), tx.GetChannelID())
	assert.Equal(t, clprHex(t, clprOwnershipVectors[1].connectorID), tx.GetConnectorID())
	assert.Equal(t, clprHex(t, clprVectorSalt), tx.GetSalt())
	assert.Equal(t, clprOwnershipVectors[1].connectorSignature, hex.EncodeToString(tx.GetOwnershipSignature()))
}

func TestUnitClprCompleteConnectorTransactionRequiredFields(t *testing.T) {
	t.Parallel()

	vector := clprOwnershipVectors[1]
	tx := NewClprCompleteConnectorTransaction()
	clprOffline(t, tx.Transaction)

	steps := []struct {
		field string
		set   func()
	}{
		{"'channelID'", func() { tx.SetChannelID(clprHex(t, clprVectorChannelID)) }},
		{"'connectorID'", func() { tx.SetConnectorID(clprHex(t, vector.connectorID)) }},
		{"'salt'", func() { tx.SetSalt(clprHex(t, clprVectorSalt)) }},
		{"'connectorContract'", func() { tx.SetConnectorContract(ContractID{Contract: 345}) }},
		{"'adminKey'", func() { tx.SetAdminKey(clprConnectorAdminKey(t)) }},
	}
	for _, step := range steps {
		_, err := tx.Freeze()
		assert.ErrorIs(t, err, errClprFieldNotSet)
		assert.ErrorContains(t, err, "ClprCompleteConnectorTransaction: "+step.field)
		step.set()
	}

	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)

	_, err = tx.SetOwnershipProof(vector.key(t))
	require.NoError(t, err)
	_, err = tx.Freeze()
	assert.NoError(t, err)
}

func TestUnitClprCompleteConnectorTransactionRejectsBadIDLength(t *testing.T) {
	t.Parallel()

	for _, length := range []int{0, 31, 33} {
		tx := clprReadyCompleteConnector(t, clprOwnershipVectors[1])
		salt := tx.GetSalt()
		signature := tx.GetOwnershipSignature()

		tx.SetSalt(clprBytes(length, 1))
		assert.Equal(t, salt, tx.GetSalt())
		assert.Equal(t, signature, tx.GetOwnershipSignature())
		_, err := tx.Freeze()
		assert.ErrorContains(t, err, "salt must be exactly 32 bytes")

		tx.SetSalt(salt)
		_, err = tx.Freeze()
		assert.NoError(t, err)
	}
}

func TestUnitClprCompleteConnectorTransactionIncompleteDecodedProof(t *testing.T) {
	t.Parallel()

	vector := clprOwnershipVectors[1]
	for _, clear := range []func(*services.ClprCompleteConnectorTransactionBody){
		func(body *services.ClprCompleteConnectorTransactionBody) { body.PublicKey = nil },
		func(body *services.ClprCompleteConnectorTransactionBody) { body.Signature = nil },
		func(body *services.ClprCompleteConnectorTransactionBody) { body.PublicKey, body.Signature = nil, nil },
	} {
		wire := clprCompleteConnectorWire(t, vector, 0)
		clear(wire.GetClprCompleteConnector())
		decoded, err := TransactionFromBytes(clprUnsignedBytes(t, wire))
		require.NoError(t, err)
		tx := decoded.(ClprCompleteConnectorTransaction).childTransaction
		clprOffline(t, tx.Transaction)

		_, err = tx.Freeze()
		assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)

		_, err = tx.SetOwnershipProof(vector.key(t))
		require.NoError(t, err)
		_, err = tx.Freeze()
		assert.NoError(t, err)
	}
}

func TestUnitClprCompleteConnectorTransactionFrozenSettersKeepState(t *testing.T) {
	t.Parallel()

	vector := clprOwnershipVectors[1]
	tx := clprReadyCompleteConnector(t, vector).SetLockedStake(HbarFromTinybar(7))
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(vector.key(t))
	before, err := tx.ToBytes()
	require.NoError(t, err)

	tx.SetChannelID(clprBytes(32, 8)).
		SetConnectorID(clprBytes(32, 8)).
		SetSalt(clprBytes(32, 8)).
		SetAdminKey(clprOwnershipVectors[0].key(t).PublicKey()).
		SetConnectorContract(ContractID{Contract: 1}).
		SetLockedStake(HbarFromTinybar(-1))
	_, err = tx.SetOwnershipProof(vector.key(t))
	assert.ErrorIs(t, err, errTransactionIsFrozen)

	assert.Equal(t, clprHex(t, clprVectorChannelID), tx.GetChannelID())
	assert.Equal(t, clprHex(t, vector.connectorID), tx.GetConnectorID())
	assert.Equal(t, clprHex(t, clprVectorSalt), tx.GetSalt())
	assert.Equal(t, clprConnectorAdminKey(t).String(), tx.GetAdminKey().String())
	assert.Equal(t, ContractID{Shard: 1, Realm: 2, Contract: 345}, tx.GetConnectorContract())
	assert.Equal(t, HbarFromTinybar(7), tx.GetLockedStake())
	assert.Equal(t, vector.connectorSignature, hex.EncodeToString(tx.GetOwnershipSignature()))
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestUnitClprCompleteConnectorTransactionFrozenBytesRejectedBeforeExecute(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	for _, clear := range []func(*services.ClprCompleteConnectorTransactionBody){
		func(body *services.ClprCompleteConnectorTransactionBody) { body.PublicKey, body.Signature = nil, nil },
		func(body *services.ClprCompleteConnectorTransactionBody) { body.PublicKey = nil },
		func(body *services.ClprCompleteConnectorTransactionBody) { body.Signature = nil },
	} {
		wire := clprCompleteConnectorWire(t, clprOwnershipVectors[1], 0)
		clear(wire.GetClprCompleteConnector())
		frozen, err := TransactionFromBytes(clprFrozenBytes(t, wire))
		require.NoError(t, err)
		assert.True(t, frozen.(ClprCompleteConnectorTransaction).IsFrozen())

		_, err = TransactionExecute(frozen, client)
		assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)
	}
}

func TestUnitClprCompleteConnectorTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	tx := clprReadyCompleteConnector(t, clprOwnershipVectors[1])
	_, err := tx.Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprCompleteConnectorTransaction`")
}

func TestUnitClprCompleteConnectorTransactionRejectsEmptyAdminKey(t *testing.T) {
	t.Parallel()

	tx := clprReadyCompleteConnector(t, clprOwnershipVectors[1])
	adminKey := tx.GetAdminKey()

	for _, empty := range []Key{nil, PublicKey{}} {
		tx.SetAdminKey(empty)
		assert.Equal(t, adminKey.String(), tx.GetAdminKey().String())
		_, err := tx.Freeze()
		assert.ErrorIs(t, err, errClprAdminKeyInvalid)
	}

	tx.SetAdminKey(adminKey)
	_, err := tx.Freeze()
	assert.NoError(t, err)
}

func TestUnitClprCompleteConnectorTransactionDecodedInvalidAdminKey(t *testing.T) {
	t.Parallel()

	body := clprCompleteConnectorWire(t, clprOwnershipVectors[1], 7)
	body.GetClprCompleteConnector().AdminKey = &services.Key{Key: &services.Key_Ed25519{Ed25519: []byte{1, 2, 3}}}

	frozen, err := TransactionFromBytes(clprFrozenBytes(t, body))
	require.NoError(t, err)

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	_, err = TransactionExecute(frozen, client)
	assert.ErrorIs(t, err, errClprAdminKeyInvalid)
}

func TestUnitClprCompleteConnectorTransactionAdminKeyKeepsType(t *testing.T) {
	t.Parallel()

	contract := ContractID{Contract: 5, EvmAddress: clprBytes(20, 1)}
	tx := NewClprCompleteConnectorTransaction().SetAdminKey(contract)
	contract.EvmAddress[0] = 9

	got, ok := tx.GetAdminKey().(ContractID)
	require.True(t, ok)
	assert.Equal(t, clprBytes(20, 1), got.EvmAddress)
	got.EvmAddress[0] = 9

	again, ok := tx.GetAdminKey().(ContractID)
	require.True(t, ok)
	assert.Equal(t, clprBytes(20, 1), again.EvmAddress)

	_, ok = NewClprCompleteConnectorTransaction().SetAdminKey(clprConnectorAdminKey(t)).GetAdminKey().(*KeyList)
	assert.True(t, ok)
}
