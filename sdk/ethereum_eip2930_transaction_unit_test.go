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

func TestUnitEthereumEIP2930TypedAccessors(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")

	tx := (&EthereumEIP2930Transaction{}).
		SetChainId(298).
		SetNonce(5).
		SetGasPrice(big.NewInt(20_000_000_000)).
		SetGasLimit(0x2dc6c0).
		SetTo(to).
		SetValue(big.NewInt(0x02540be400)).
		SetCallData([]byte{0xDE, 0xAD}).
		SetRecoveryId(1).
		SetR([]byte{0xAA}).
		SetS([]byte{0xBB})

	assert.Equal(t, uint64(298), tx.GetChainId())
	assert.Equal(t, uint64(5), tx.GetNonce())
	assert.Equal(t, int64(20_000_000_000), tx.GetGasPrice().Int64())
	assert.Equal(t, uint64(0x2dc6c0), tx.GetGasLimit())
	assert.Equal(t, to, tx.GetTo())
	assert.Equal(t, int64(0x02540be400), tx.GetValue().Int64())
	assert.Equal(t, []byte{0xDE, 0xAD}, tx.GetCallData())
	assert.Equal(t, 1, tx.GetRecoveryId())
	assert.Equal(t, []byte{0xAA}, tx.GetR())
	assert.Equal(t, []byte{0xBB}, tx.GetS())

	assert.Equal(t, []byte{0x01, 0x2a}, tx.GetChainIdBytes())
	assert.Equal(t, []byte{0x05}, tx.GetNonceBytes())
	assert.Equal(t, []byte{0x04, 0xa8, 0x17, 0xc8, 0x00}, tx.GetGasPriceBytes())
	assert.Equal(t, []byte{0x2d, 0xc6, 0xc0}, tx.GetGasLimitBytes())
	assert.Equal(t, []byte{0x02, 0x54, 0x0b, 0xe4, 0x00}, tx.GetValueBytes())
	assert.Equal(t, []byte{0x01}, tx.GetRecoveryIdBytes())
}

func TestUnitEthereumEIP2930BytesAccessors(t *testing.T) {
	t.Parallel()

	tx := (&EthereumEIP2930Transaction{}).
		SetChainIdBytes([]byte{0x01, 0x2a}).
		SetNonceBytes([]byte{0x05}).
		SetGasPriceBytes([]byte{0x02}).
		SetGasLimitBytes([]byte{0x03}).
		SetValueBytes([]byte{0x04}).
		SetRecoveryIdBytes([]byte{0x01})

	assert.Equal(t, uint64(298), tx.GetChainId())
	assert.Equal(t, uint64(5), tx.GetNonce())
	assert.Equal(t, int64(2), tx.GetGasPrice().Int64())
	assert.Equal(t, uint64(3), tx.GetGasLimit())
	assert.Equal(t, int64(4), tx.GetValue().Int64())
	assert.Equal(t, 1, tx.GetRecoveryId())

	// RecoveryId of 0 is stored as empty bytes, and reads back as 0.
	tx.SetRecoveryId(0)
	assert.Equal(t, []byte{}, tx.GetRecoveryIdBytes())
	assert.Equal(t, 0, tx.GetRecoveryId())
}

func TestUnitEthereumEIP2930SetAccessListItems(t *testing.T) {
	t.Parallel()

	tx := &EthereumEIP2930Transaction{}
	tx.SetAccessListItems([]AccessListItem{
		NewAccessListItem([]byte{0xAA}, [][]byte{{0x01}}),
		NewAccessListItem([]byte{0xBB}, nil),
	})

	got := tx.GetAccessListItems()
	require.Equal(t, 2, len(got))
	assert.Equal(t, []byte{0xAA}, got[0].GetAddress())
}

func TestUnitEthereumEIP2930String(t *testing.T) {
	t.Parallel()

	tx := (&EthereumEIP2930Transaction{}).
		SetChainId(1).
		AddAccessListItem(NewAccessListItem([]byte{0xAA}, [][]byte{{0x01}}))
	s := tx.String()
	assert.Contains(t, s, "ChainId:")
	assert.Contains(t, s, "AccessList:")
	assert.Contains(t, s, "RecoveryId:")
}

func TestUnitEthereumEIP2930FromBytesErrors(t *testing.T) {
	t.Parallel()

	_, err := EthereumEIP2930TransactionFromBytes(nil)
	assert.Error(t, err)

	_, err = EthereumEIP2930TransactionFromBytes([]byte{0x02, 0x00})
	assert.Error(t, err, "wrong type prefix should be rejected")

	// Correct 0x01 prefix but wrong element count.
	list := NewRLPItem(LIST_TYPE)
	list.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01}))
	body, err := list.Write()
	require.NoError(t, err)
	_, err = EthereumEIP2930TransactionFromBytes(append([]byte{0x01}, body...))
	assert.Error(t, err)
}
