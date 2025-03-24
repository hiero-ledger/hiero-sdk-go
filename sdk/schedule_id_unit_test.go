//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitScheduleIDFromString(t *testing.T) {
	t.Parallel()

	id, err := ScheduleIDFromString("0.0.123")
	require.NoError(t, err)
	require.Equal(t, ScheduleID{Shard: 0, Realm: 0, Schedule: 123}, id)
}

func TestUnitScheduleIDToString(t *testing.T) {
	t.Parallel()

	id := ScheduleID{Shard: 0, Realm: 0, Schedule: 123}
	require.Equal(t, "0.0.123", id.String())
}

func TestUnitScheduleIDChecksum(t *testing.T) {
	t.Parallel()

	id, err := ScheduleIDFromString("0.0.123")
	require.NoError(t, err)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	checksum, err := id.ToStringWithChecksum(*client)
	require.NoError(t, err)
	require.Equal(t, "0.0.123-esxsf", checksum)
}

func TestUnitScheduleIDChecksumError(t *testing.T) {
	t.Parallel()

	id, err := ScheduleIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()

	_, err = id.ToStringWithChecksum(*client)
	require.Error(t, err)
}
