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

func TestUnitEthereumEIP1559TypedAccessors(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")

	tx := (&EthereumEIP1559Transaction{}).
		SetChainId(298).
		SetNonce(5).
		SetMaxPriorityGas(big.NewInt(2_000_000_000)).
		SetMaxGas(big.NewInt(20_000_000_000)).
		SetGasLimit(150_000).
		SetTo(to).
		SetValue(big.NewInt(1_000_000)).
		SetCallData([]byte{0xDE, 0xAD}).
		SetRecoveryId(1).
		SetR([]byte{0xAA}).
		SetS([]byte{0xBB})

	assert.Equal(t, uint64(298), tx.GetChainId())
	assert.Equal(t, uint64(5), tx.GetNonce())
	assert.Equal(t, int64(2_000_000_000), tx.GetMaxPriorityGas().Int64())
	assert.Equal(t, int64(20_000_000_000), tx.GetMaxGas().Int64())
	assert.Equal(t, uint64(150_000), tx.GetGasLimit())
	assert.Equal(t, to, tx.GetTo())
	assert.Equal(t, int64(1_000_000), tx.GetValue().Int64())
	assert.Equal(t, []byte{0xDE, 0xAD}, tx.GetCallData())
	assert.Equal(t, 1, tx.GetRecoveryId())
	assert.Equal(t, []byte{0xAA}, tx.GetR())
	assert.Equal(t, []byte{0xBB}, tx.GetS())

	assert.Equal(t, []byte{0x01, 0x2a}, tx.GetChainIdBytes())
	assert.Equal(t, []byte{0x05}, tx.GetNonceBytes())
	assert.Equal(t, []byte{0x77, 0x35, 0x94, 0x00}, tx.GetMaxPriorityGasBytes())
	assert.Equal(t, []byte{0x04, 0xa8, 0x17, 0xc8, 0x00}, tx.GetMaxGasBytes())
	assert.Equal(t, []byte{0x02, 0x49, 0xf0}, tx.GetGasLimitBytes())
	assert.Equal(t, []byte{0x0f, 0x42, 0x40}, tx.GetValueBytes())
	assert.Equal(t, []byte{0x01}, tx.GetRecoveryIdBytes())
}

func TestUnitEthereumEIP1559BytesAccessors(t *testing.T) {
	t.Parallel()

	tx := (&EthereumEIP1559Transaction{}).
		SetChainIdBytes([]byte{0x01, 0x2a}).
		SetNonceBytes([]byte{0x05}).
		SetMaxPriorityGasBytes([]byte{0x02}).
		SetMaxGasBytes([]byte{0x03}).
		SetGasLimitBytes([]byte{0x04}).
		SetValueBytes([]byte{0x05}).
		SetRecoveryIdBytes([]byte{0x01})

	assert.Equal(t, uint64(298), tx.GetChainId())
	assert.Equal(t, uint64(5), tx.GetNonce())
	assert.Equal(t, int64(2), tx.GetMaxPriorityGas().Int64())
	assert.Equal(t, int64(3), tx.GetMaxGas().Int64())
	assert.Equal(t, uint64(4), tx.GetGasLimit())
	assert.Equal(t, int64(5), tx.GetValue().Int64())
	assert.Equal(t, 1, tx.GetRecoveryId())
}

func TestUnitEthereumEIP1559String(t *testing.T) {
	t.Parallel()

	tx := (&EthereumEIP1559Transaction{}).SetChainId(1).SetMaxGas(big.NewInt(5))
	s := tx.String()
	assert.Contains(t, s, "ChainId:")
	assert.Contains(t, s, "MaxGas:")
	assert.Contains(t, s, "AccessList:")
}

func TestUnitEthereumEIP1559FromBytesErrors(t *testing.T) {
	t.Parallel()

	_, err := EthereumEIP1559TransactionFromBytes(nil)
	assert.Error(t, err)

	_, err = EthereumEIP1559TransactionFromBytes([]byte{0x01, 0x00})
	assert.Error(t, err, "wrong type prefix should be rejected")

	list := NewRLPItem(LIST_TYPE)
	list.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01}))
	body, err := list.Write()
	require.NoError(t, err)
	_, err = EthereumEIP1559TransactionFromBytes(append([]byte{0x02}, body...))
	assert.Error(t, err)
}
