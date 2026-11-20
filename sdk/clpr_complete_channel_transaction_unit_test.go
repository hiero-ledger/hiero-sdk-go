//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func clprReadyCompleteChannel(t *testing.T, key PrivateKey) *ClprCompleteChannelTransaction {
	t.Helper()
	tx := NewClprCompleteChannelTransaction().
		SetChannelID(clprHex(t, clprVectorChannelID)).
		SetVerifierContract(ContractID{Shard: 1, Realm: 2, Contract: 345}).
		SetConfigProofBytes([]byte{255, 0, 7}).
		SetEndpointManifestProofBytes([]byte{128, 0, 9, 255})
	_, err := tx.SetOwnershipProof(key)
	require.NoError(t, err)
	clprOffline(t, tx.Transaction)
	return tx
}

// clprWireBody fills the transaction header of a hand-built CLPR body.
func clprWireBody(body *services.TransactionBody) *services.TransactionBody {
	body.TransactionID = TransactionIDGenerate(AccountID{Account: 5006})._ToProtobuf()
	body.NodeAccountID = AccountID{Account: 3}._ToProtobuf()
	body.TransactionFee = uint64(NewHbar(1).AsTinybar())
	body.TransactionValidDuration = &services.Duration{Seconds: 120}
	return body
}

// clprUnsignedBytes serializes a body in the legacy body-bytes form, which decodes unfrozen.
func clprUnsignedBytes(t *testing.T, body *services.TransactionBody) []byte {
	t.Helper()
	bodyBytes, err := proto.Marshal(clprWireBody(body))
	require.NoError(t, err)
	list, err := proto.Marshal(&sdk.TransactionList{TransactionList: []*services.Transaction{{BodyBytes: bodyBytes}}}) // nolint
	require.NoError(t, err)
	return list
}

// clprFrozenBytes serializes a body as signed-transaction bytes, which decode frozen.
func clprFrozenBytes(t *testing.T, body *services.TransactionBody) []byte {
	t.Helper()
	bodyBytes, err := proto.Marshal(clprWireBody(body))
	require.NoError(t, err)
	signed, err := proto.Marshal(&services.SignedTransaction{BodyBytes: bodyBytes, SigMap: &services.SignatureMap{}})
	require.NoError(t, err)
	list, err := proto.Marshal(&sdk.TransactionList{TransactionList: []*services.Transaction{{SignedTransactionBytes: signed}}})
	require.NoError(t, err)
	return list
}

func TestUnitClprCompleteChannelTransactionOwnershipProof(t *testing.T) {
	t.Parallel()

	for _, vector := range clprOwnershipVectors {
		t.Run(vector.name, func(t *testing.T) {
			tx := clprReadyCompleteChannel(t, vector.key(t))

			assert.Equal(t, vector.publicKey, hex.EncodeToString(tx.GetOwnershipPublicKey()))
			assert.Equal(t, vector.channelSignature, hex.EncodeToString(tx.GetOwnershipSignature()))
			scheme, ok := tx.GetOwnershipSignatureScheme()
			assert.True(t, ok)
			assert.Equal(t, vector.scheme, scheme)
			assertClprSignatureVerifies(t, vector, clprHex(t, vector.channelDigest), tx.GetOwnershipSignature())
		})
	}
}

func TestUnitClprCompleteChannelTransactionSignedRoundTrip(t *testing.T) {
	t.Parallel()

	for _, vector := range clprOwnershipVectors {
		t.Run(vector.name, func(t *testing.T) {
			payer := clprOwnershipVectors[1].key(t)
			tx := clprReadyCompleteChannel(t, vector.key(t))
			_, err := tx.Freeze()
			require.NoError(t, err)
			tx.Sign(payer)

			txBytes, err := tx.ToBytes()
			require.NoError(t, err)

			decoded, err := TransactionFromBytes(txBytes)
			require.NoError(t, err)
			result := decoded.(ClprCompleteChannelTransaction)

			assert.True(t, result.IsFrozen())
			assert.Equal(t, clprHex(t, clprVectorChannelID), result.GetChannelID())
			assert.Equal(t, "1.2.345", result.GetVerifierContract().String())
			assert.Equal(t, []byte{255, 0, 7}, result.GetConfigProofBytes())
			assert.Equal(t, []byte{128, 0, 9, 255}, result.GetEndpointManifestProofBytes())
			assert.Equal(t, vector.channelSignature, hex.EncodeToString(result.GetOwnershipSignature()))
			scheme, ok := result.GetOwnershipSignatureScheme()
			assert.True(t, ok)
			assert.Equal(t, vector.scheme, scheme)

			roundTrip, err := result.ToBytes()
			require.NoError(t, err)
			assert.Equal(t, txBytes, roundTrip)

			signatures, err := result.GetSignatures()
			require.NoError(t, err)
			assert.Len(t, signatures, 2)
		})
	}
}

func TestUnitClprCompleteChannelTransactionAbsentProof(t *testing.T) {
	t.Parallel()

	tx := NewClprCompleteChannelTransaction().SetChannelID(clprBytes(32, 1))
	assert.Nil(t, tx.GetOwnershipPublicKey())
	assert.Nil(t, tx.GetOwnershipSignature())
	_, ok := tx.GetOwnershipSignatureScheme()
	assert.False(t, ok)

	decoded, err := TransactionFromBytes(clprUnsignedBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprCompleteChannel{
			ClprCompleteChannel: &services.ClprCompleteChannelTransactionBody{ChannelId: clprBytes(32, 1)},
		},
	}))
	require.NoError(t, err)
	result := decoded.(ClprCompleteChannelTransaction)
	assert.Nil(t, result.GetOwnershipPublicKey())
	_, ok = result.GetOwnershipSignatureScheme()
	assert.False(t, ok)
}

func TestUnitClprCompleteChannelTransactionUnknownScheme(t *testing.T) {
	t.Parallel()

	decoded, err := TransactionFromBytes(clprUnsignedBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprCompleteChannel{
			ClprCompleteChannel: &services.ClprCompleteChannelTransactionBody{
				ChannelId:       clprBytes(32, 1),
				PublicKey:       []byte{1},
				Signature:       []byte{2},
				SignatureScheme: 7,
			},
		},
	}))
	require.NoError(t, err)

	result := decoded.(ClprCompleteChannelTransaction)
	scheme, ok := result.GetOwnershipSignatureScheme()
	assert.True(t, ok)
	assert.Equal(t, "UNKNOWN (7)", scheme.String())
}

func TestUnitClprCompleteChannelTransactionRequiredFields(t *testing.T) {
	t.Parallel()

	key := clprOwnershipVectors[1].key(t)

	tx := NewClprCompleteChannelTransaction()
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "ClprCompleteChannelTransaction: 'channelID'")

	_, err = tx.SetOwnershipProof(key)
	assert.ErrorContains(t, err, "'channelID'")
	assert.Nil(t, tx.GetOwnershipPublicKey())

	tx.SetChannelID(clprBytes(32, 1))
	_, err = tx.Freeze()
	assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)

	_, err = tx.SetOwnershipProof(key)
	require.NoError(t, err)
	_, err = tx.Freeze()
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "ClprCompleteChannelTransaction: 'verifierContract'")

	tx.SetVerifierContract(ContractID{Contract: 345})
	_, err = tx.Freeze()
	assert.ErrorIs(t, err, errClprConfigProofEmpty)

	tx.SetConfigProofBytes([]byte{})
	_, err = tx.Freeze()
	assert.ErrorIs(t, err, errClprConfigProofEmpty)

	tx.SetConfigProofBytes([]byte{1})
	_, err = tx.Freeze()
	require.NoError(t, err)
	assert.Empty(t, tx.GetEndpointManifestProofBytes())
}

func TestUnitClprCompleteChannelTransactionRejectsBadChannelIDLength(t *testing.T) {
	t.Parallel()

	for _, length := range []int{0, 31, 33} {
		tx := NewClprCompleteChannelTransaction().SetChannelID(clprBytes(length, 1))
		clprOffline(t, tx.Transaction)
		assert.Nil(t, tx.GetChannelID())
		_, err := tx.Freeze()
		assert.ErrorContains(t, err, "channelID must be exactly 32 bytes")
		assert.False(t, tx.IsFrozen())

		_, err = tx.SetOwnershipProof(clprOwnershipVectors[1].key(t))
		assert.ErrorIs(t, err, errClprFieldNotSet)

		ready := clprReadyCompleteChannel(t, clprOwnershipVectors[1].key(t))
		signature := ready.GetOwnershipSignature()
		ready.SetChannelID(clprBytes(length, 1))
		assert.Equal(t, clprHex(t, clprVectorChannelID), ready.GetChannelID())
		assert.Equal(t, signature, ready.GetOwnershipSignature())
	}
}

func TestUnitClprCompleteChannelTransactionIncompleteDecodedProof(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct{ publicKey, signature []byte }{{nil, []byte{2}}, {[]byte{1}, nil}} {
		decoded, err := TransactionFromBytes(clprUnsignedBytes(t, &services.TransactionBody{
			Data: &services.TransactionBody_ClprCompleteChannel{
				ClprCompleteChannel: &services.ClprCompleteChannelTransactionBody{
					ChannelId:        clprBytes(32, 1),
					VerifierContract: ContractID{Contract: 345}._ToProtobuf(),
					ConfigProofBytes: []byte{1},
					PublicKey:        tt.publicKey,
					Signature:        tt.signature,
					SignatureScheme:  services.ClprSignatureScheme_ED25519,
				},
			},
		}))
		require.NoError(t, err)
		// The decoded value is a copy; its base transaction drives the live child.
		tx := decoded.(ClprCompleteChannelTransaction).childTransaction
		clprOffline(t, tx.Transaction)

		_, err = tx.Freeze()
		assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)

		_, err = tx.SetOwnershipProof(clprOwnershipVectors[1].key(t))
		require.NoError(t, err)
		_, err = tx.Freeze()
		assert.NoError(t, err)
	}
}

func TestUnitClprCompleteChannelTransactionChannelIDChangeClearsProof(t *testing.T) {
	t.Parallel()

	key := clprOwnershipVectors[1].key(t)
	tx := clprReadyCompleteChannel(t, key)
	signature := tx.GetOwnershipSignature()

	tx.SetChannelID(clprHex(t, clprVectorChannelID))
	assert.Equal(t, signature, tx.GetOwnershipSignature())

	tx.SetChannelID(clprBytes(32, 8))
	assert.Nil(t, tx.GetOwnershipPublicKey())
	assert.Nil(t, tx.GetOwnershipSignature())
	_, ok := tx.GetOwnershipSignatureScheme()
	assert.False(t, ok)
	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)

	tx.SetChannelID(clprHex(t, clprVectorChannelID))
	_, err = tx.Freeze()
	assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)

	_, err = tx.SetOwnershipProof(key)
	require.NoError(t, err)
	_, err = tx.Freeze()
	assert.NoError(t, err)
}

func TestUnitClprCompleteChannelTransactionCopiesBytes(t *testing.T) {
	t.Parallel()

	channelID := clprBytes(32, 1)
	configProof := []byte{1, 2}
	manifestProof := []byte{3, 4}
	verifier := ContractID{EvmAddress: clprBytes(20, 5)}
	tx := NewClprCompleteChannelTransaction().
		SetChannelID(channelID).
		SetConfigProofBytes(configProof).
		SetEndpointManifestProofBytes(manifestProof).
		SetVerifierContract(verifier)
	_, err := tx.SetOwnershipProof(clprOwnershipVectors[1].key(t))
	require.NoError(t, err)

	for _, bytes := range [][]byte{channelID, configProof, manifestProof, verifier.EvmAddress,
		tx.GetChannelID(), tx.GetConfigProofBytes(), tx.GetEndpointManifestProofBytes(),
		tx.GetVerifierContract().EvmAddress, tx.GetOwnershipPublicKey(), tx.GetOwnershipSignature()} {
		for i := range bytes {
			bytes[i] = 0
		}
	}

	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
	assert.Equal(t, []byte{1, 2}, tx.GetConfigProofBytes())
	assert.Equal(t, []byte{3, 4}, tx.GetEndpointManifestProofBytes())
	assert.Equal(t, clprBytes(20, 5), tx.GetVerifierContract().EvmAddress)
	assert.Equal(t, clprOwnershipVectors[1].publicKey, hex.EncodeToString(tx.GetOwnershipPublicKey()))
}

func TestUnitClprCompleteChannelTransactionFrozenSettersKeepState(t *testing.T) {
	t.Parallel()

	key := clprOwnershipVectors[1].key(t)
	tx := clprReadyCompleteChannel(t, key)
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(key)
	before, err := tx.ToBytes()
	require.NoError(t, err)

	tx.SetChannelID(clprBytes(32, 8)).
		SetVerifierContract(ContractID{Contract: 1}).
		SetConfigProofBytes([]byte{9}).
		SetEndpointManifestProofBytes([]byte{9})
	_, err = tx.SetOwnershipProof(clprOwnershipVectors[0].key(t))
	assert.ErrorIs(t, err, errTransactionIsFrozen)

	assert.Equal(t, clprHex(t, clprVectorChannelID), tx.GetChannelID())
	assert.Equal(t, "1.2.345", tx.GetVerifierContract().String())
	assert.Equal(t, []byte{255, 0, 7}, tx.GetConfigProofBytes())
	assert.Equal(t, clprOwnershipVectors[1].publicKey, hex.EncodeToString(tx.GetOwnershipPublicKey()))
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestUnitClprCompleteChannelTransactionFrozenBytesRejectedBeforeExecute(t *testing.T) {
	t.Parallel()

	frozen, err := TransactionFromBytes(clprFrozenBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprCompleteChannel{
			ClprCompleteChannel: &services.ClprCompleteChannelTransactionBody{ChannelId: clprBytes(32, 1), ConfigProofBytes: []byte{1}},
		},
	}))
	require.NoError(t, err)
	assert.True(t, frozen.(ClprCompleteChannelTransaction).IsFrozen())

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	_, err = TransactionExecute(frozen, client)
	assert.ErrorIs(t, err, errClprOwnershipProofIncomplete)
}

func TestUnitClprCompleteChannelTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	tx := clprReadyCompleteChannel(t, clprOwnershipVectors[1].key(t))
	_, err := tx.Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprCompleteChannelTransaction`")
}
