//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitLiveHash_BytesRoundTrip(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	keys := NewKeyList()
	keys.Add(key.PublicKey())

	original := LiveHash{
		AccountID:        AccountID{Shard: 0, Realm: 0, Account: 123},
		Hash:             []byte{1, 2, 3},
		Keys:             *keys,
		LiveHashDuration: 30 * time.Second,
	}

	bytes := original.ToBytes()

	recovered, err := LiveHashFromBytes(bytes)
	require.NoError(t, err)

	assert.Equal(t, original.AccountID, recovered.AccountID)
	assert.Equal(t, original.Hash, recovered.Hash)
	assert.Equal(t, original.Keys, recovered.Keys)
	assert.Equal(t, original.LiveHashDuration, recovered.LiveHashDuration)
}

func TestUnitLiveHash_FromProtobufNil(t *testing.T) {
	t.Parallel()

	result, err := _LiveHashFromProtobuf(nil)

	assert.Equal(t, LiveHash{}, result)
	require.ErrorIs(t, err, errParameterNull)
}

func TestUnitLiveHash_FromBytesNil(t *testing.T) {
	t.Parallel()

	result, err := LiveHashFromBytes(nil)

	assert.Equal(t, LiveHash{}, result)
	require.ErrorIs(t, err, errByteArrayNull)
}

func TestUnitLiveHash_ProtobufRoundTripNilKeys(t *testing.T) {
	t.Parallel()

	pb := &services.LiveHash{
		AccountId: &services.AccountID{
			ShardNum: 0,
			RealmNum: 0,
			Account: &services.AccountID_AccountNum{
				AccountNum: 123,
			},
		},
		Duration: &services.Duration{
			Seconds: 30,
		},
		Keys: nil,
	}

	result, err := _LiveHashFromProtobuf(pb)
	require.NoError(t, err)

	assert.Equal(t, KeyList{}, result.Keys)
	assert.Equal(t, AccountID{Shard: 0, Realm: 0, Account: 123}, result.AccountID)
	assert.Equal(t, 30*time.Second, result.LiveHashDuration)
}

func TestUnitLiveHash_DeprecatedDurationNotSerialised(t *testing.T) {
	t.Parallel()

	lh := LiveHash{
		AccountID:        AccountID{Shard: 0, Realm: 0, Account: 123},
		Hash:             []byte{1, 2, 3},
		Keys:             KeyList{},
		Duration:         time.Unix(1700000000, 0),
		LiveHashDuration: 30 * time.Second,
	}

	bytes := lh.ToBytes()

	result, err := LiveHashFromBytes(bytes)
	require.NoError(t, err)

	assert.True(t, result.Duration.IsZero())
	assert.Equal(t, 30*time.Second, result.LiveHashDuration)
}
