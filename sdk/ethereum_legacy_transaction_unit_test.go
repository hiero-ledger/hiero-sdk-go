//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitEthereumLegacyTypedAccessors(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")

	tx := (&EthereumLegacyTransaction{}).
		SetNonce(7).
		SetGasPrice(big.NewInt(20_000_000_000)).
		SetGasLimit(21_000).
		SetTo(to).
		SetValue(big.NewInt(1_000_000)).
		SetCallData([]byte{0xDE, 0xAD}).
		SetV(28).
		SetR([]byte{0xAA}).
		SetS([]byte{0xBB})

	assert.Equal(t, uint64(7), tx.GetNonce())
	assert.Equal(t, int64(20_000_000_000), tx.GetGasPrice().Int64())
	assert.Equal(t, uint64(21_000), tx.GetGasLimit())
	assert.Equal(t, to, tx.GetTo())
	assert.Equal(t, int64(1_000_000), tx.GetValue().Int64())
	assert.Equal(t, []byte{0xDE, 0xAD}, tx.GetCallData())
	assert.Equal(t, uint64(28), tx.GetV())
	assert.Equal(t, []byte{0xAA}, tx.GetR())
	assert.Equal(t, []byte{0xBB}, tx.GetS())

	// Typed setters store canonical big-endian bytes.
	assert.Equal(t, []byte{0x07}, tx.GetNonceBytes())
	assert.Equal(t, []byte{0x04, 0xa8, 0x17, 0xc8, 0x00}, tx.GetGasPriceBytes())
	assert.Equal(t, []byte{0x52, 0x08}, tx.GetGasLimitBytes())
	assert.Equal(t, []byte{0x0f, 0x42, 0x40}, tx.GetValueBytes())
	assert.Equal(t, []byte{0x1c}, tx.GetVBytes())
}

func TestUnitEthereumLegacyBytesAccessors(t *testing.T) {
	t.Parallel()

	tx := (&EthereumLegacyTransaction{}).
		SetNonceBytes([]byte{0x01}).
		SetGasPriceBytes([]byte{0x02}).
		SetGasLimitBytes([]byte{0x03}).
		SetValueBytes([]byte{0x04}).
		SetVBytes([]byte{0x1b})

	assert.Equal(t, uint64(1), tx.GetNonce())
	assert.Equal(t, int64(2), tx.GetGasPrice().Int64())
	assert.Equal(t, uint64(3), tx.GetGasLimit())
	assert.Equal(t, int64(4), tx.GetValue().Int64())
	assert.Equal(t, uint64(27), tx.GetV())
}

func TestUnitEthereumLegacyConstructorAndToBytesRoundTrip(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")
	tx := NewEthereumLegacyTransaction(
		[]byte{0x01},       // nonce
		[]byte{0x04, 0xa8}, // gasPrice
		[]byte{0x52, 0x08}, // gasLimit
		to,                 // to
		[]byte{},           // value
		[]byte{0xCA, 0xFE}, // callData
		[]byte{0x1c},       // v
		[]byte{0x11},       // r
		[]byte{0x22},       // s
	)

	encoded, err := tx.ToBytes()
	require.NoError(t, err)
	require.NotEmpty(t, encoded)

	decoded, err := EthereumLegacyTransactionFromBytes(encoded)
	require.NoError(t, err)
	assert.Equal(t, tx.Nonce, decoded.Nonce)
	assert.Equal(t, tx.GasPrice, decoded.GasPrice)
	assert.Equal(t, tx.GasLimit, decoded.GasLimit)
	assert.Equal(t, tx.To, decoded.To)
	assert.Equal(t, tx.Value, decoded.Value)
	assert.Equal(t, tx.CallData, decoded.CallData)
	assert.Equal(t, tx.V, decoded.V)
	assert.Equal(t, tx.R, decoded.R)
	assert.Equal(t, tx.S, decoded.S)
}

func TestUnitEthereumLegacyString(t *testing.T) {
	t.Parallel()

	tx := (&EthereumLegacyTransaction{}).SetNonce(1).SetGasLimit(21_000)
	s := tx.String()
	assert.Contains(t, s, "Nonce:")
	assert.Contains(t, s, "GasLimit:")
	assert.Contains(t, s, "V:")
}

func TestUnitEthereumLegacyFromBytesErrors(t *testing.T) {
	t.Parallel()

	// Not an RLP list (a bare value).
	value := NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01})
	notAList, err := value.Write()
	require.NoError(t, err)
	_, err = EthereumLegacyTransactionFromBytes(notAList)
	assert.Error(t, err)

	// A list, but not 9 elements.
	short := NewRLPItem(LIST_TYPE)
	short.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01}))
	shortBytes, err := short.Write()
	require.NoError(t, err)
	_, err = EthereumLegacyTransactionFromBytes(shortBytes)
	assert.Error(t, err)
}
