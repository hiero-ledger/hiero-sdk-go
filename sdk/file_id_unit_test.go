//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitFileIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := FileIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	id.ToStringWithChecksum(*client)
	sol := id.ToSolidityAddress()
	FileIDFromSolidityAddress(sol)
	id.Validate(client)

	pb := id._ToProtobuf()
	_FileIDFromProtobuf(pb)

	idByte := id.ToBytes()
	FileIDFromBytes(idByte)

	require.Equal(t, FileID{File: 111}.String(), FileIDForFeeSchedule().String())
	require.Equal(t, FileID{File: 102}.String(), FileIDForAddressBook().String())
	require.Equal(t, FileID{File: 112}.String(), FileIDForExchangeRate().String())

	assert.Equal(t, id.File, uint64(123))
}

func TestUnitFileIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitFileIDChecksumError(t *testing.T) {
	t.Parallel()

	id, err := FileIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()

	_, err = id.ToStringWithChecksum(*client)
	require.Error(t, err)
}

func TestUnitGetAddressBookFileIDFor(t *testing.T) {
	t.Parallel()

	fileID := GetAddressBookFileIDFor(0, 0)
	assert.Equal(t, uint64(0), fileID.Shard)
	assert.Equal(t, uint64(0), fileID.Realm)
	assert.Equal(t, uint64(102), fileID.File)
	assert.Equal(t, FileIDForAddressBook(), fileID)

	fileID = GetAddressBookFileIDFor(5, 3)
	assert.Equal(t, uint64(3), fileID.Shard)
	assert.Equal(t, uint64(5), fileID.Realm)
	assert.Equal(t, uint64(102), fileID.File)
}

func TestUnitGetFeeScheduleFileIDFor(t *testing.T) {
	t.Parallel()

	fileID := GetFeeScheduleFileIDFor(0, 0)
	assert.Equal(t, uint64(0), fileID.Shard)
	assert.Equal(t, uint64(0), fileID.Realm)
	assert.Equal(t, uint64(111), fileID.File)
	assert.Equal(t, FileIDForFeeSchedule(), fileID)

	fileID = GetFeeScheduleFileIDFor(5, 3)
	assert.Equal(t, uint64(3), fileID.Shard)
	assert.Equal(t, uint64(5), fileID.Realm)
	assert.Equal(t, uint64(111), fileID.File)
}

func TestUnitGetExchangeRatesFileIDFor(t *testing.T) {
	t.Parallel()

	fileID := GetExchangeRatesFileIDFor(0, 0)
	assert.Equal(t, uint64(0), fileID.Shard)
	assert.Equal(t, uint64(0), fileID.Realm)
	assert.Equal(t, uint64(112), fileID.File)
	assert.Equal(t, FileIDForExchangeRate(), fileID)

	fileID = GetExchangeRatesFileIDFor(5, 3)
	assert.Equal(t, uint64(3), fileID.Shard)
	assert.Equal(t, uint64(5), fileID.Realm)
	assert.Equal(t, uint64(112), fileID.File)
}

func TestUnitFileIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a long zero address representing file 1234
	evmAddress := "00000000000000000000000000000000000004d2"
	id, err := FileIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(1234), id.File)

	// Test with a different shard and realm
	id, err = FileIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(1234), id.File)
}

func TestUnitFileIDToEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal file ID
	id := FileID{Shard: 0, Realm: 0, File: 123}
	require.Equal(t, "000000000000000000000000000000000000007b", id.ToEvmAddress())

	// Test with a different shard and realm
	id = FileID{Shard: 1, Realm: 1, File: 123}
	require.Equal(t, "000000000000000000000000000000000000007b", id.ToEvmAddress())
}
