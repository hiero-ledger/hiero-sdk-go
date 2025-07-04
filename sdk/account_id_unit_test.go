//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"encoding/hex"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	id.ToStringWithChecksum(client)
	id.GetChecksum()
	sol := id.ToSolidityAddress()
	AccountIDFromSolidityAddress(sol)
	err = id.Validate(client)
	require.Error(t, err)
	evmID, err := AccountIDFromEvmAddress(0, 0, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	require.NoError(t, err)
	pb := evmID._ToProtobuf()
	_AccountIDFromProtobuf(pb)

	idByte := id.ToBytes()
	AccountIDFromBytes(idByte)

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	alias := key.ToAccountID(0, 0)
	pb = alias._ToProtobuf()
	_AccountIDFromProtobuf(pb)

	require.NoError(t, err)
	assert.Equal(t, id.Account, uint64(123))
}

func TestUnitAccountIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitAccountIDFromStringAlias(t *testing.T) {
	t.Parallel()

	key, err := GeneratePrivateKey()
	require.NoError(t, err)
	id, err := AccountIDFromString("0.0." + key.PublicKey().String())
	require.NoError(t, err)
	id2 := key.ToAccountID(0, 0)

	assert.Equal(t, id.String(), id2.String())
}

func TestUnitChecksum(t *testing.T) {
	t.Parallel()

	id, err := LedgerIDFromString("01")
	require.NoError(t, err)
	ad1, err := _ChecksumParseAddress(id, "0.0.3")
	require.NoError(t, err)
	id, err = LedgerIDFromString("10")
	require.NoError(t, err)
	ad2, err := _ChecksumParseAddress(id, "0.0.3")
	require.NoError(t, err)

	require.NotEqual(t, ad1.correctChecksum, ad2.correctChecksum)
}

func TestUnitAccountIDEvm(t *testing.T) {
	t.Parallel()

	id, err := AccountIDFromString("0.0.0011223344556677889900112233445566778899")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.0011223344556677889900112233445566778899")
}

func TestUnitAccountIDPopulateFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateAccount(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateAccount(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailWithNoMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDChecksumError(t *testing.T) {
	t.Parallel()

	id, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	client, err := _NewMockClient()
	require.NoError(t, err)
	_, err = id.ToStringWithChecksum(client)
	require.Error(t, err)
}

func TestUnitAccountIDFromEvmAddressIncorrectSize(t *testing.T) {
	t.Parallel()

	// Test with an EVM address that's too short
	_, err := AccountIDFromEvmAddress(0, 0, "abc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Test with an EVM address that's too long
	_, err = AccountIDFromEvmAddress(0, 0, "0123456789abcdef0123456789abcdef0123456789abcdef")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Test with a 0x prefix that gets removed but then is too short
	_, err = AccountIDFromEvmAddress(0, 0, "0xabc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Verify a correct length works
	correctAddress := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	id, err := AccountIDFromEvmAddress(0, 0, correctAddress)
	require.NoError(t, err)
	require.NotNil(t, id.AliasEvmAddress)

	require.Equal(t, strings.ToLower(evmAddress), hex.EncodeToString(*id.AliasEvmAddress))
}

func TestUnitAccountIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal EVM address
	evmAddress := evmAddress
	bytes, err := hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err := AccountIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)

	// Test with a different shard and realm
	id, err = AccountIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)

	// Test with a long zero address
	evmAddress = longZeroAddress
	bytes, err = hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err = AccountIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)

	// Test with a different shard and realm
	id, err = AccountIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)
}

func TestUnitAccountIDToEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal account ID
	id := AccountID{Shard: 0, Realm: 0, Account: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a different shard and realm
	id = AccountID{Shard: 1, Realm: 1, Account: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a long zero address
	bytes, err := hex.DecodeString(longZeroAddress)
	id = AccountID{Shard: 1, Realm: 1, AliasEvmAddress: &bytes}
	require.NoError(t, err)
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a normal EVM address
	emvAddress := evmAddress
	bytes, err = hex.DecodeString(emvAddress)
	id = AccountID{Shard: 0, Realm: 0, AliasEvmAddress: &bytes}
	expected := strings.ToLower(evmAddress)
	require.NoError(t, err)
	require.Equal(t, expected, id.ToEvmAddress())

	// Test with different shard and realm
	id = AccountID{Shard: 1, Realm: 1, AliasEvmAddress: &bytes}
	require.Equal(t, expected, id.ToEvmAddress())
}
