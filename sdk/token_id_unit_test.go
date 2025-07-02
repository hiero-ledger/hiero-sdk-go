//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenIDFromString(t *testing.T) {
	t.Parallel()

	tokID := TokenID{
		Shard: 1,
		Realm: 2,
		Token: 3,
	}

	gotTokID, err := TokenIDFromString(tokID.String())
	require.NoError(t, err)
	assert.Equal(t, tokID.Token, gotTokID.Token)
}

func TestUnitTokenIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := TokenIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	id.ToStringWithChecksum(*client)
	sol := id.ToSolidityAddress()
	TokenIDFromSolidityAddress(sol)
	id.Validate(client)
	pb := id._ToProtobuf()
	_TokenIDFromProtobuf(pb)

	idByte := id.ToBytes()
	TokenIDFromBytes(idByte)

	id.Compare(TokenID{Token: 32})

	assert.Equal(t, id.Token, uint64(123))
}

func TestUnitTokenIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitTokenIDFromStringEVM(t *testing.T) {
	t.Parallel()

	id, err := TokenIDFromString("0.0.434")
	require.NoError(t, err)

	require.Equal(t, "0.0.434", id.String())
}

func TestUnitTokenIDProtobuf(t *testing.T) {
	t.Parallel()

	id, err := TokenIDFromString("0.0.434")
	require.NoError(t, err)

	pb := id._ToProtobuf()

	require.Equal(t, pb, &services.TokenID{
		ShardNum: 0,
		RealmNum: 0,
		TokenNum: 434,
	})

	pbFrom := _TokenIDFromProtobuf(pb)

	require.Equal(t, id, *pbFrom)
}

func TestUnitTokenIDChecksumError(t *testing.T) {
	t.Parallel()

	id, err := TokenIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()

	_, err = id.ToStringWithChecksum(*client)
	require.Error(t, err)
}

func TestUnitTokenIDFromEvmAddressIncorrectAddress(t *testing.T) {
	t.Parallel()

	// Test with an EVM address that's too short
	_, err := TokenIDFromEvmAddress(0, 0, "abc123")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")

	// Test with an EVM address that's too long
	_, err = TokenIDFromEvmAddress(0, 0, "0123456789abcdef0123456789abcdef0123456789abcdef")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")

	// Test with a 0x prefix that gets removed but then is too short
	_, err = TokenIDFromEvmAddress(0, 0, "0xabc123")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")

	// Test with non-long-zero address
	_, err = TokenIDFromEvmAddress(0, 0, "742d35Cc6634C0532925a3b844Bc454e4438f44e")
	require.Error(t, err)
	require.Contains(t, err.Error(), "EVM address is not a correct long zero addres")
}

func TestUnitTokenIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a long zero address representing token 1234
	evmAddress := "00000000000000000000000000000000000004d2"
	id, err := TokenIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(1234), id.Token)

	// Test with a different shard and realm
	id, err = TokenIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(1234), id.Token)
}

func TestUnitTokenIDToEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal token ID
	id := TokenID{Shard: 0, Realm: 0, Token: 123}
	require.Equal(t, "000000000000000000000000000000000000007b", id.ToEvmAddress())

	// Test with a different shard and realm
	id = TokenID{Shard: 1, Realm: 1, Token: 123}
	require.Equal(t, "000000000000000000000000000000000000007b", id.ToEvmAddress())
}
