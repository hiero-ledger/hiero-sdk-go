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

	id, err := TopicIDFromSolidityAddress("000000000000000000000000000000000000007b")
	require.NoError(t, err)
	require.Equal(t, TopicID{Shard: 0, Realm: 0, Topic: 123}, id)
}

func TestUnitTopicIDToSolidityAddress(t *testing.T) {
	t.Parallel()

	id := TopicID{Shard: 0, Realm: 0, Topic: 123}
	address := id.ToSolidityAddress()
	require.Equal(t, "000000000000000000000000000000000000007b", address)
}
