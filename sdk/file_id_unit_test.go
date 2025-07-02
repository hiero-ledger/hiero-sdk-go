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

	fileID = GetAddressBookFileIDFor(3, 5)
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

	fileID = GetFeeScheduleFileIDFor(3, 5)
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

	fileID = GetExchangeRatesFileIDFor(3, 5)
	assert.Equal(t, uint64(3), fileID.Shard)
	assert.Equal(t, uint64(5), fileID.Realm)
	assert.Equal(t, uint64(112), fileID.File)
}

func TestUnitFileIDFromEvmAddressIncorrectAddress(t *testing.T) {
	t.Parallel()

	// Test with an EVM address that's too short
	_, err := FileIDFromEvmAddress(0, 0, "abc123")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")

	// Test with an EVM address that's too long
	_, err = FileIDFromEvmAddress(0, 0, "0123456789abcdef0123456789abcdef0123456789abcdef")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")

	// Test with a 0x prefix that gets removed but then is too short
	_, err = FileIDFromEvmAddress(0, 0, "0xabc123")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")

	// Test with non-long-zero address
	_, err = FileIDFromEvmAddress(0, 0, "742d35Cc6634C0532925a3b844Bc454e4438f44e")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")
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
