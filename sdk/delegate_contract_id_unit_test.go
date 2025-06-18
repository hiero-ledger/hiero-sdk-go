//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitDelegatableContractIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	err = id.ValidateChecksum(client)
	require.Error(t, err)
	require.Equal(t, id.Contract, uint64(123))
	strChecksum, err := id.ToStringWithChecksum(*client)
	require.NoError(t, err)
	// different checksum because of different network
	require.Equal(t, strChecksum, "0.0.123-esxsf")
}

func TestUnitDelegatableContractIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := DelegatableContractID{
		Shard:    50,
		Realm:    150,
		Contract: 520,
	}
	require.Equal(t, "50.150.520", id.String())
}

func TestUnitDelegatableContractIDFromStringEVM(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, "0.0.0011223344556677889900112233445577889900", id.String())
}

func TestUnitDelegatableContractIDProtobuf(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	pb := id._ToProtobuf()

	decoded, err := hex.DecodeString("0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: decoded},
	})

	pbFrom := _DelegatableContractIDFromProtobuf(pb)

	require.Equal(t, id, *pbFrom)
}

func TestUnitDelegatableContractIDEvm(t *testing.T) {
	t.Parallel()

	hexString, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	id, err := DelegatableContractIDFromString(fmt.Sprintf("0.0.%s", hexString.PublicKey().String()))
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(id.EvmAddress), hexString.PublicKey().String())

	pb := id._ToProtobuf()
	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: id.EvmAddress},
	})

	id, err = DelegatableContractIDFromString("0.0.123")
	require.NoError(t, err)
	require.Equal(t, id.Contract, uint64(123))
	require.Nil(t, id.EvmAddress)

	pb = id._ToProtobuf()
	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_ContractNum{ContractNum: 123},
	})
}

func TestUnitDelegatableContractIDToFromBytes(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123")
	require.NoError(t, err)
	require.Equal(t, id.Contract, uint64(123))
	require.Nil(t, id.EvmAddress)

	idBytes := id.ToBytes()
	idFromBytes, err := DelegatableContractIDFromBytes(idBytes)
	require.NoError(t, err)
	require.Equal(t, id, idFromBytes)
}

func TestUnitDelegatableContractIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal EVM address
	evmAddress := "742d35Cc6634C0532925a3b844Bc454e4438f44e"
	bytes, err := hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err := DelegatableContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)

	// Test with a different shard and realm
	id, err = DelegatableContractIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)

	// Test with a long zero address
	evmAddress = "00000000000000000000000000000000000004d2"
	bytes, err = hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err = DelegatableContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)

	// Test with a different shard and realm
	id, err = DelegatableContractIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)
}

func TestUnitDelegatableContractIDFromSolidityAddress(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	sol := id.ToSolidityAddress()
	idFromSolidity, err := DelegatableContractIDFromSolidityAddress(sol)
	require.NoError(t, err)
	require.Equal(t, idFromSolidity.Contract, uint64(123))
}

func TestUnitDelegatableContractIDToProtoKey(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	pb := id._ToProtoKey()
	require.Equal(t, pb.GetContractID().GetContractNum(), int64(123))
}

func TestUnitDelegatableContractIDChecksumError(t *testing.T) {
	t.Parallel()

	id, err := DelegatableContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()

	_, err = id.ToStringWithChecksum(*client)
	require.Error(t, err)
}

func TestUnitDelegatableContractIDToEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal contract ID
	id := DelegatableContractID{Shard: 0, Realm: 0, Contract: 123}
	require.Equal(t, "000000000000000000000000000000000000007b", id.ToEvmAddress())

	// Test with a different shard and realm
	id = DelegatableContractID{Shard: 1, Realm: 1, Contract: 123}
	require.Equal(t, "000000000000000000000000000000000000007b", id.ToEvmAddress())

	// Test with a long zero address
	longZeroAddress := "00000000000000000000000000000000000004d2"
	bytes, err := hex.DecodeString(longZeroAddress)
	id = DelegatableContractID{Shard: 1, Realm: 1, EvmAddress: bytes}
	require.NoError(t, err)
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a normal EVM address
	evmAddress := "742d35Cc6634C0532925a3b844Bc454e4438f44e"
	bytes, err = hex.DecodeString(evmAddress)
	id = DelegatableContractID{Shard: 0, Realm: 0, EvmAddress: bytes}
	expected := strings.ToLower("742d35Cc6634C0532925a3b844Bc454e4438f44e")
	require.NoError(t, err)
	require.Equal(t, expected, id.ToEvmAddress())

	// Test with different shard and realm
	id = DelegatableContractID{Shard: 1, Realm: 1, EvmAddress: bytes}
	require.Equal(t, expected, id.ToEvmAddress())
}
