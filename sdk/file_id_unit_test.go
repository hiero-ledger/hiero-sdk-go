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

	fileID, err := GetAddressBookFileIDFor(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), fileID.Shard)
	assert.Equal(t, uint64(0), fileID.Realm)
	assert.Equal(t, uint64(102), fileID.File)
	assert.Equal(t, FileIDForAddressBook(), fileID)

	fileID, err = GetAddressBookFileIDFor(5, 3)
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), fileID.Shard)
	assert.Equal(t, uint64(5), fileID.Realm)
	assert.Equal(t, uint64(102), fileID.File)

	_, err = GetAddressBookFileIDFor(-1, 0)
	assert.Error(t, err)
	_, err = GetAddressBookFileIDFor(0, -1)
	assert.Error(t, err)
	_, err = GetAddressBookFileIDFor(-1, -1)
	assert.Error(t, err)
}

func TestUnitGetFeeScheduleFileIDFor(t *testing.T) {
	t.Parallel()

	fileID, err := GetFeeScheduleFileIDFor(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), fileID.Shard)
	assert.Equal(t, uint64(0), fileID.Realm)
	assert.Equal(t, uint64(111), fileID.File)
	assert.Equal(t, FileIDForFeeSchedule(), fileID)

	fileID, err = GetFeeScheduleFileIDFor(5, 3)
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), fileID.Shard)
	assert.Equal(t, uint64(5), fileID.Realm)
	assert.Equal(t, uint64(111), fileID.File)

	_, err = GetFeeScheduleFileIDFor(-1, 0)
	assert.Error(t, err)
	_, err = GetFeeScheduleFileIDFor(0, -1)
	assert.Error(t, err)
	_, err = GetFeeScheduleFileIDFor(-1, -1)
	assert.Error(t, err)
}

func TestUnitGetExchangeRatesFileIDFor(t *testing.T) {
	t.Parallel()

	fileID, err := GetExchangeRatesFileIDFor(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), fileID.Shard)
	assert.Equal(t, uint64(0), fileID.Realm)
	assert.Equal(t, uint64(112), fileID.File)
	assert.Equal(t, FileIDForExchangeRate(), fileID)

	fileID, err = GetExchangeRatesFileIDFor(5, 3)
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), fileID.Shard)
	assert.Equal(t, uint64(5), fileID.Realm)
	assert.Equal(t, uint64(112), fileID.File)

	_, err = GetAddressBookFileIDFor(-1, 0)
	assert.Error(t, err)
	_, err = GetFeeScheduleFileIDFor(0, -1)
	assert.Error(t, err)
	_, err = GetExchangeRatesFileIDFor(-1, -1)
	assert.Error(t, err)
}
