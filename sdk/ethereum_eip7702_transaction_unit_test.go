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

func TestUnitEthereumEIP7702TypedAccessors(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")

	tx := (&EthereumEIP7702Transaction{}).
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

func TestUnitEthereumEIP7702BytesAccessors(t *testing.T) {
	t.Parallel()

	tx := (&EthereumEIP7702Transaction{}).
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

func TestUnitEthereumEIP7702AuthorizationList(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")
	auth1 := NewAuthorization([]byte{0x01, 0x2a}, to, 0, 0, make([]byte, 32), make([]byte, 32))
	auth2 := NewAuthorization([]byte{0x01, 0x2a}, to, 1, 1, make([]byte, 32), make([]byte, 32))

	tx := (&EthereumEIP7702Transaction{}).SetAuthorizationList([]Authorization{auth1})
	require.Equal(t, 1, len(tx.GetAuthorizationList()))

	tx.AddAuthorization(auth2)
	require.Equal(t, 2, len(tx.GetAuthorizationList()))
	assert.Equal(t, auth2, tx.GetAuthorizationList()[1])
}

func TestUnitEthereumEIP7702AddAccessListItem(t *testing.T) {
	t.Parallel()

	tx := &EthereumEIP7702Transaction{}
	tx.AddAccessListItem(NewAccessListItem([]byte{0xAA}, [][]byte{{0x01}}))
	require.Equal(t, 1, len(tx.GetAccessListItems()))
	assert.Equal(t, []byte{0xAA}, tx.GetAccessListItems()[0].GetAddress())
}

func TestUnitEthereumEIP7702ToBytesRoundTripWithAuthList(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")
	// Canonical (non-zero / minimal) field values so the RLP round-trip is exact.
	auth := NewAuthorization([]byte{0x01, 0x2a}, to, 5, 1, make([]byte, 32), make([]byte, 32))

	tx := (&EthereumEIP7702Transaction{}).
		SetChainId(298).
		SetNonce(1).
		SetMaxGas(big.NewInt(5)).
		SetGasLimit(150_000).
		SetTo(to).
		SetValue(big.NewInt(0)).
		SetRecoveryId(0).
		SetR([]byte{0x11}).
		SetS([]byte{0x22}).
		AddAuthorization(auth)

	encoded, err := tx.ToBytes()
	require.NoError(t, err)
	require.Equal(t, byte(0x04), encoded[0])

	decoded, err := EthereumEIP7702TransactionFromBytes(encoded)
	require.NoError(t, err)
	require.Equal(t, 1, len(decoded.GetAuthorizationList()))
	assert.Equal(t, auth, decoded.GetAuthorizationList()[0])
}

func TestUnitEthereumEIP7702String(t *testing.T) {
	t.Parallel()

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")
	tx := (&EthereumEIP7702Transaction{}).
		SetChainId(1).
		AddAuthorization(NewAuthorization([]byte{0x01}, to, 0, 0, []byte{0x11}, []byte{0x22}))
	s := tx.String()
	assert.Contains(t, s, "ChainId:")
	assert.Contains(t, s, "AuthorizationList:")
}

func TestUnitEthereumEIP7702FromBytesErrors(t *testing.T) {
	t.Parallel()

	_, err := EthereumEIP7702TransactionFromBytes(nil)
	assert.Error(t, err)

	_, err = EthereumEIP7702TransactionFromBytes([]byte{0x02, 0x00})
	assert.Error(t, err, "wrong type prefix should be rejected")

	// Correct prefix, but an authorization tuple without 6 elements.
	list := NewRLPItem(LIST_TYPE)
	for i := 0; i < 8; i++ {
		list.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{}))
	}
	list.PushBack(NewRLPItem(LIST_TYPE)) // empty access list (index 8)

	authList := NewRLPItem(LIST_TYPE)
	badTuple := NewRLPItem(LIST_TYPE)
	badTuple.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01}))
	authList.PushBack(badTuple)
	list.PushBack(authList) // authorization list (index 9)

	list.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{})) // recoveryId
	list.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{})) // r
	list.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{})) // s

	body, err := list.Write()
	require.NoError(t, err)
	_, err = EthereumEIP7702TransactionFromBytes(append([]byte{0x04}, body...))
	assert.Error(t, err)
}
