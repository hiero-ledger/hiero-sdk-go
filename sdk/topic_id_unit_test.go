//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitTopicInfoChecksumError(t *testing.T) {
	t.Parallel()

	id, err := TopicIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	client, err := _NewMockClient()
	require.NoError(t, err)
	_, err = id.ToStringWithChecksum(*client)
	require.Error(t, err)
}

func TestUnitTopicIDFromString(t *testing.T) {
	t.Parallel()

	id, err := TopicIDFromString("0.0.123")
	require.NoError(t, err)
	require.Equal(t, TopicID{Shard: 0, Realm: 0, Topic: 123}, id)
}

func TestUnitTopicIDToString(t *testing.T) {
	t.Parallel()

	id := TopicID{Shard: 0, Realm: 0, Topic: 123}
	require.Equal(t, "0.0.123", id.String())
}

func TestUnitTopicIDChecksum(t *testing.T) {
	t.Parallel()

	id, err := TopicIDFromString("0.0.123")
	require.NoError(t, err)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDMainnet())
	require.NoError(t, err)
	checksum, err := id.ToStringWithChecksum(*client)
	require.NoError(t, err)
	require.Equal(t, "0.0.123-vfmkw", checksum)
}

func TestUnitTopicIDFromSolidityAddress(t *testing.T) {
	t.Parallel()

	id, err := TopicIDFromSolidityAddress(longZeroAddress)
	require.NoError(t, err)
	require.Equal(t, TopicID{Shard: 0, Realm: 0, Topic: 123}, id)
}

func TestUnitTopicIDToSolidityAddress(t *testing.T) {
	t.Parallel()

	id := TopicID{Shard: 0, Realm: 0, Topic: 123}
	address := id.ToSolidityAddress()
	require.Equal(t, longZeroAddress, address)
}

func TestUnitTopicIDFromEvmAddressIncorrectAddress(t *testing.T) {
	t.Parallel()

	// Test with an EVM address that's too short
	_, err := TopicIDFromEvmAddress(0, 0, "abc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotALongZeroAddress)

	// Test with an EVM address that's too long
	_, err = TopicIDFromEvmAddress(0, 0, "0123456789abcdef0123456789abcdef0123456789abcdef")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotALongZeroAddress)

	// Test with a 0x prefix that gets removed but then is too short
	_, err = TopicIDFromEvmAddress(0, 0, "0xabc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotALongZeroAddress)

	// Test with non-long-zero address
	_, err = TopicIDFromEvmAddress(0, 0, evmAddress)
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotALongZeroAddress)
}

func TestUnitTopicIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a long zero address representing topic 1234
	id, err := TopicIDFromEvmAddress(0, 0, longZeroAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(123), id.Topic)

	// Test with a different shard and realm
	id, err = TopicIDFromEvmAddress(1, 1, longZeroAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(123), id.Topic)
}

func TestUnitTopicIDToEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal topic ID
	id := TopicID{Shard: 0, Realm: 0, Topic: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a different shard and realm
	id = TopicID{Shard: 1, Realm: 1, Topic: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())
}
