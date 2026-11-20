//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

var clprBundlePayload = []byte{0, 255, 128, 0, 13, 10, 195, 40, 1}

// clprSignedWireBytes serializes a hand-built body signed by the Ed25519 vector key, as a node would receive it.
func clprSignedWireBytes(t *testing.T, body *services.TransactionBody) []byte {
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

func TestUnitClprSubmitBundleTransactionSignedRoundTrip(t *testing.T) {
	t.Parallel()

	for _, payload := range [][]byte{{}, clprBundlePayload} {
		tx := NewClprSubmitBundleTransaction().
			SetChannelID(clprHex(t, clprVectorChannelID)).
			SetBundlePayload(payload)
		clprOffline(t, tx.Transaction)
		_, err := tx.Freeze()
		require.NoError(t, err)
		for _, vector := range clprOwnershipVectors {
			tx.Sign(vector.key(t))
		}

		txBytes, err := tx.ToBytes()
		require.NoError(t, err)
		decoded, err := TransactionFromBytes(txBytes)
		require.NoError(t, err)
		result := decoded.(ClprSubmitBundleTransaction)

		assert.True(t, result.IsFrozen())
		assert.Equal(t, clprHex(t, clprVectorChannelID), result.GetChannelID())
		assert.Equal(t, len(payload), len(result.GetBundlePayload()))
		assert.Equal(t, string(payload), string(result.GetBundlePayload()))
		for _, vector := range clprOwnershipVectors {
			assert.True(t, vector.key(t).PublicKey().VerifyTransaction(result.childTransaction), vector.name)
		}
		roundTrip, err := result.ToBytes()
		require.NoError(t, err)
		assert.Equal(t, txBytes, roundTrip)
	}
}

func TestUnitClprSubmitBundleTransactionRequiresOnlyChannelID(t *testing.T) {
	t.Parallel()

	tx := NewClprSubmitBundleTransaction().SetBundlePayload(clprBundlePayload)
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "ClprSubmitBundleTransaction: 'channelID'")

	for _, length := range []int{0, 31, 33} {
		tx.SetChannelID(clprBytes(length, 1))
		_, err = tx.Freeze()
		assert.ErrorContains(t, err, "channelID must be exactly 32 bytes")
		assert.False(t, tx.IsFrozen())
	}

	withoutPayload := NewClprSubmitBundleTransaction().SetChannelID(clprBytes(32, 1))
	clprOffline(t, withoutPayload.Transaction)
	_, err = withoutPayload.Freeze()
	assert.NoError(t, err)
	assert.Empty(t, withoutPayload.buildProtoBody().GetBundlePayload())
}

func TestUnitClprSubmitBundleTransactionCopiesBytes(t *testing.T) {
	t.Parallel()

	channelID := clprBytes(32, 1)
	payload := append([]byte{}, clprBundlePayload...)
	tx := NewClprSubmitBundleTransaction().SetChannelID(channelID).SetBundlePayload(payload)

	for _, bytes := range [][]byte{channelID, payload, tx.GetChannelID(), tx.GetBundlePayload()} {
		for i := range bytes {
			bytes[i] = 127
		}
	}

	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
	assert.Equal(t, clprBundlePayload, tx.GetBundlePayload())
}

func TestUnitClprSubmitBundleTransactionFrozenSettersKeepState(t *testing.T) {
	t.Parallel()

	tx := NewClprSubmitBundleTransaction().SetChannelID(clprBytes(32, 1)).SetBundlePayload(clprBundlePayload)
	clprOffline(t, tx.Transaction)
	_, err := tx.Freeze()
	require.NoError(t, err)
	tx.Sign(clprOwnershipVectors[1].key(t))
	before, err := tx.ToBytes()
	require.NoError(t, err)

	tx.SetChannelID(clprBytes(32, 8)).SetBundlePayload([]byte{9})

	assert.Equal(t, clprBytes(32, 1), tx.GetChannelID())
	assert.Equal(t, clprBundlePayload, tx.GetBundlePayload())
	after, err := tx.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestUnitClprSubmitBundleTransactionDeprecatedFields(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"EndpointNodeID", "EndpointSignature"} {
		_, hasGetter := reflect.TypeOf(&ClprSubmitBundleTransaction{}).MethodByName("Get" + name)
		_, hasSetter := reflect.TypeOf(&ClprSubmitBundleTransaction{}).MethodByName("Set" + name)
		assert.False(t, hasGetter || hasSetter, name)
	}

	frozenBytes := clprSignedWireBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprSubmitBundle{
			ClprSubmitBundle: &services.ClprSubmitBundleTransactionBody{
				ChannelId:         clprBytes(32, 1),
				BundlePayload:     clprBundlePayload,
				EndpointNodeId:    9007199254740993,                   // nolint
				EndpointSignature: []byte(strings.Repeat("\x11", 65)), // nolint
			},
		},
	})
	decoded, err := TransactionFromBytes(frozenBytes)
	require.NoError(t, err)
	result := decoded.(ClprSubmitBundleTransaction)

	assert.True(t, result.IsFrozen())
	roundTrip, err := result.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, frozenBytes, roundTrip)
	assert.True(t, clprOwnershipVectors[1].key(t).PublicKey().VerifyTransaction(result.childTransaction))

	rebuilt := result.buildProtoBody()
	assert.Zero(t, rebuilt.GetEndpointNodeId())     // nolint
	assert.Empty(t, rebuilt.GetEndpointSignature()) // nolint
}

func TestUnitClprSubmitBundleTransactionFrozenBytesRejectedBeforeExecute(t *testing.T) {
	t.Parallel()

	frozen, err := TransactionFromBytes(clprFrozenBytes(t, &services.TransactionBody{
		Data: &services.TransactionBody_ClprSubmitBundle{ClprSubmitBundle: &services.ClprSubmitBundleTransactionBody{}},
	}))
	require.NoError(t, err)

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	_, err = TransactionExecute(frozen, client)
	assert.ErrorIs(t, err, errClprFieldNotSet)
	assert.ErrorContains(t, err, "'channelID'")
}

func TestUnitClprSubmitBundleTransactionCannotSchedule(t *testing.T) {
	t.Parallel()

	_, err := NewClprSubmitBundleTransaction().SetChannelID(clprBytes(32, 1)).Schedule()
	assert.EqualError(t, err, "cannot schedule `ClprSubmitBundleTransaction`")
}
