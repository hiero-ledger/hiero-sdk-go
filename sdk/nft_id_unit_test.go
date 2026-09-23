//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitNftIDFromStringWithSlash(t *testing.T) {
	t.Parallel()

	id, err := NftIDFromString("0.0.5005/10")
	require.NoError(t, err)
	require.Equal(t, TokenID{Shard: 0, Realm: 0, Token: 5005}, id.TokenID)
	require.Equal(t, int64(10), id.SerialNumber)
	require.Equal(t, "0.0.5005/10", id.String())
}

func TestUnitNftIDFromStringWithLegacyAtSeparator(t *testing.T) {
	t.Parallel()

	legacy, err := NftIDFromString("0.0.5005@10")
	require.NoError(t, err)

	canonical, err := NftIDFromString("0.0.5005/10")
	require.NoError(t, err)

	require.Equal(t, canonical, legacy)
	require.Equal(t, "0.0.5005/10", legacy.String())
}

func TestUnitNftIDFromStringPreservesChecksum(t *testing.T) {
	t.Parallel()

	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	expected := tokenID.Nft(2)

	for _, s := range []string{"0.0.123-esxsf/2", "0.0.123-esxsf@2"} {
		id, err := NftIDFromString(s)
		require.NoError(t, err)
		require.Equal(t, expected, id)
	}
}

func TestUnitNftIDFromStringInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"no separator", "0.0.5005"},
		{"missing token with slash", "/10"},
		{"missing token with at", "@10"},
		{"missing serial with slash", "0.0.5005/"},
		{"missing serial with at", "0.0.5005@"},
		{"too many slashes", "0.0.5005/10/2"},
		{"too many ats", "0.0.5005@10@2"},
		{"mixed separators", "0.0.5005@10/2"},
		{"serial not a number", "0.0.5005/notanumber"},
		{"serial out of range", "0.0.5005/9223372036854775808"},
		{"token not an id", "abc/10"},
		{"token missing parts", "0.0/10"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var err error
			require.NotPanics(t, func() {
				_, err = NftIDFromString(tt.input)
			})
			require.Error(t, err)
		})
	}
}

func TestUnitNftIDFromStringMaxSerialNumber(t *testing.T) {
	t.Parallel()

	id, err := NftIDFromString("0.0.5005/9223372036854775807")
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt64), id.SerialNumber)
	require.Equal(t, "0.0.5005/9223372036854775807", id.String())
}

func TestUnitNftIDToStringWithChecksum(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	id, err := NftIDFromString("0.0.123/2")
	require.NoError(t, err)

	s, err := id.ToStringWithChecksum(*client)
	require.NoError(t, err)
	require.Equal(t, "0.0.123-esxsf/2", s)
}

func TestUnitNftIDProtobuf(t *testing.T) {
	t.Parallel()

	id := NftID{TokenID: TokenID{Shard: 1, Realm: 2, Token: 5005}, SerialNumber: 10}
	require.Equal(t, id, _NftIDFromProtobuf(id._ToProtobuf()))
	require.Equal(t, NftID{}, _NftIDFromProtobuf(nil))
}

func TestUnitNftIDBytes(t *testing.T) {
	t.Parallel()

	id := NftID{TokenID: TokenID{Shard: 1, Realm: 2, Token: 5005}, SerialNumber: 10}
	decoded, err := NftIDFromBytes(id.ToBytes())
	require.NoError(t, err)
	require.Equal(t, id, decoded)
}

func TestUnitNftIDIsZero(t *testing.T) {
	t.Parallel()

	assert.True(t, NftID{}._IsZero())
	assert.False(t, NftID{SerialNumber: 1}._IsZero())
	assert.False(t, NftID{TokenID: TokenID{Token: 5005}}._IsZero())
}
