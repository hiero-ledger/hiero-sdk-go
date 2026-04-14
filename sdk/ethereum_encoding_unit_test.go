//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestUnitEIP1559NumberAccessors(t *testing.T) {
	t.Parallel()

	gasPrice := new(big.Int).SetUint64(900_000_000_000) // 0xd1385c7bf0
	tx := (&EthereumEIP1559Transaction{}).
		SetChainId(298).
		SetNonce(42).
		SetGasLimit(150_000).
		SetMaxGas(gasPrice).
		SetMaxPriorityGas(big.NewInt(0)).
		SetValue(big.NewInt(1_000_000)).
		SetRecoveryId(1)

	assert.Equal(t, uint64(298), tx.GetChainId())
	assert.Equal(t, uint64(42), tx.GetNonce())
	assert.Equal(t, uint64(150_000), tx.GetGasLimit())
	assert.Equal(t, 0, gasPrice.Cmp(tx.GetMaxGas()))
	assert.Equal(t, 0, tx.GetMaxPriorityGas().Sign())
	assert.Equal(t, int64(1_000_000), tx.GetValue().Int64())
	assert.Equal(t, 1, tx.GetRecoveryId())

	// Underlying bytes match canonical encoding.
	require.Equal(t, []byte{0x01, 0x2a}, tx.ChainId)
	require.Equal(t, []byte{0x2a}, tx.Nonce)
	require.Equal(t, []byte{}, tx.MaxPriorityGas)
}

func TestUnitEIP2930NumberAccessors(t *testing.T) {
	t.Parallel()

	tx := (&EthereumEIP2930Transaction{}).
		SetChainId(1).
		SetNonce(7).
		SetGasPrice(big.NewInt(20_000_000_000)).
		SetGasLimit(21_000).
		SetValue(big.NewInt(0)).
		SetRecoveryId(0)

	assert.Equal(t, uint64(1), tx.GetChainId())
	assert.Equal(t, uint64(7), tx.GetNonce())
	assert.Equal(t, int64(20_000_000_000), tx.GetGasPrice().Int64())
	assert.Equal(t, uint64(21_000), tx.GetGasLimit())
	assert.Equal(t, 0, tx.GetValue().Sign())
	assert.Equal(t, 0, tx.GetRecoveryId())
	assert.Equal(t, []byte{}, tx.Value)
}

func TestUnitEIP7702NumberAccessors(t *testing.T) {
	t.Parallel()

	tx := (&EthereumEIP7702Transaction{}).
		SetChainId(298).
		SetNonce(0).
		SetMaxGas(big.NewInt(0x05)).
		SetMaxPriorityGas(big.NewInt(0x02)).
		SetGasLimit(150_000).
		SetValue(big.NewInt(0))

	assert.Equal(t, uint64(298), tx.GetChainId())
	assert.Equal(t, uint64(0), tx.GetNonce())
	assert.Equal(t, int64(5), tx.GetMaxGas().Int64())
	assert.Equal(t, int64(2), tx.GetMaxPriorityGas().Int64())
	assert.Equal(t, uint64(150_000), tx.GetGasLimit())
}

func TestUnitLegacyNumberAccessors(t *testing.T) {
	t.Parallel()

	tx := (&EthereumLegacyTransaction{}).
		SetNonce(3).
		SetGasPrice(big.NewInt(20_000_000_000)).
		SetGasLimit(21_000).
		SetValue(big.NewInt(500)).
		SetV(27)

	assert.Equal(t, uint64(3), tx.GetNonce())
	assert.Equal(t, int64(20_000_000_000), tx.GetGasPrice().Int64())
	assert.Equal(t, uint64(21_000), tx.GetGasLimit())
	assert.Equal(t, int64(500), tx.GetValue().Int64())
	assert.Equal(t, uint64(27), tx.GetV())
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
