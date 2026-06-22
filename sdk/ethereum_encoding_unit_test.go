//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitUint64EthBytesRoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    uint64
		want []byte
	}{
		{"zero is empty", 0, []byte{}},
		{"one byte", 1, []byte{0x01}},
		{"0xFF", 0xFF, []byte{0xFF}},
		{"two bytes", 0x0100, []byte{0x01, 0x00}},
		{"eight bytes", 0xFFFFFFFFFFFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := _uint64ToEthBytes(tc.v)
			assert.Equal(t, tc.want, got)
			back := _ethBytesToUint64(got)
			assert.Equal(t, tc.v, back)
		})
	}
}

func TestUnitBigIntEthBytesRoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    *big.Int
	}{
		{"nil", nil},
		{"zero", big.NewInt(0)},
		{"small", big.NewInt(0x42)},
		{"large", new(big.Int).Lsh(big.NewInt(1), 200)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded := _bigIntToEthBytes(tc.v)
			decoded := _ethBytesToBigInt(encoded)
			if tc.v == nil || tc.v.Sign() == 0 {
				assert.Equal(t, 0, decoded.Sign())
				assert.Equal(t, 0, len(encoded))
			} else {
				assert.Equal(t, 0, tc.v.Cmp(decoded))
			}
		})
	}
}

func TestUnitEthBytesEncodingEdgeCases(t *testing.T) {
	t.Parallel()

	// Decoding truncates inputs longer than 8 bytes to the low 8.
	assert.Equal(t, uint64(0x0102030405060708),
		_ethBytesToUint64([]byte{0xFF, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}))

	// A negative big.Int encodes its absolute value.
	assert.Equal(t, []byte{0x05}, _bigIntToEthBytes(big.NewInt(-5)))
}

func TestUnitNumberAccessorsAgreeWithRawBytes(t *testing.T) {
	t.Parallel()

	tx := &EthereumEIP1559Transaction{}
	tx.SetChainId(0xCAFE)
	assert.Equal(t, []byte{0xCA, 0xFE}, tx.GetChainIdBytes())
	assert.Equal(t, uint64(0xCAFE), tx.GetChainId())

	// Writing raw bytes is reflected in the numeric getter.
	tx.ChainId = []byte{0x00, 0x01} // non-canonical leading zero
	assert.Equal(t, uint64(1), tx.GetChainId())
}
