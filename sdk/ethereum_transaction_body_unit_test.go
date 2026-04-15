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

func newTestEcdsaKey(t *testing.T) PrivateKey {
	t.Helper()
	key, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	return key
}

func TestUnitEthereumEIP1559BuilderAndSign(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")

	tx := (&EthereumEIP1559Transaction{}).
		SetChainId(298).
		SetNonce(0).
		SetMaxPriorityGas(big.NewInt(0)).
		SetMaxGas(big.NewInt(900_000_000_000)).
		SetGasLimit(150_000).
		SetTo(to).
		SetValue(big.NewInt(0)).
		SetCallData([]byte{0xDE, 0xAD, 0xBE, 0xEF}).
		AddAccessListItem(NewAccessListItem(to, [][]byte{{0x01}}))

	signed, err := tx.Sign(key)
	require.NoError(t, err)
	require.NotEmpty(t, signed)
	assert.Equal(t, byte(0x02), signed[0])

	assert.Equal(t, 32, len(tx.GetR()))
	assert.Equal(t, 32, len(tx.GetS()))
	rid := tx.GetRecoveryId()
	assert.True(t, rid == 0 || rid == 1, "recovery id must be 0 or 1, got %d", rid)

	// Round-trip through EthereumTransactionDataFromBytes.
	roundTrip, err := EthereumTransactionDataFromBytes(signed)
	require.NoError(t, err)
	body := roundTrip.GetTransaction()
	require.NotNil(t, body)

	concrete, ok := body.(*EthereumEIP1559Transaction)
	require.True(t, ok)
	assert.Equal(t, uint64(298), concrete.GetChainId())
	assert.Equal(t, tx.GetR(), concrete.GetR())
	assert.Equal(t, tx.GetS(), concrete.GetS())
}

func TestUnitEthereumEIP2930BuilderAndSign(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")

	tx := (&EthereumEIP2930Transaction{}).
		SetChainId(298).
		SetNonce(5).
		SetGasPrice(big.NewInt(0xa54f4c3c00)).
		SetGasLimit(0x2dc6c0).
		SetTo(to).
		SetValue(big.NewInt(0x02540be400)).
		SetCallData(nil)

	signed, err := tx.Sign(key)
	require.NoError(t, err)
	assert.Equal(t, byte(0x01), signed[0])

	roundTrip, err := EthereumTransactionDataFromBytes(signed)
	require.NoError(t, err)
	body, ok := roundTrip.GetTransaction().(*EthereumEIP2930Transaction)
	require.True(t, ok)
	assert.Equal(t, uint64(298), body.GetChainId())
	assert.Equal(t, to, body.GetTo())
}

func TestUnitEthereumEIP7702BuilderAndSign(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")
	authTuple := AuthorizationTuple{
		{0x01, 0x2a},
		to,
		{0x00},
		{0x00},
		make([]byte, 32),
		make([]byte, 32),
	}

	tx := (&EthereumEIP7702Transaction{}).
		SetChainId(298).
		SetNonce(1).
		SetMaxPriorityGas(big.NewInt(0)).
		SetMaxGas(big.NewInt(5)).
		SetGasLimit(150_000).
		SetTo(to).
		SetValue(big.NewInt(0)).
		SetCallData([]byte{0x01, 0x02}).
		AddAuthorization(authTuple)

	signed, err := tx.Sign(key)
	require.NoError(t, err)
	assert.Equal(t, byte(0x04), signed[0])

	roundTrip, err := EthereumTransactionDataFromBytes(signed)
	require.NoError(t, err)
	body, ok := roundTrip.GetTransaction().(*EthereumEIP7702Transaction)
	require.True(t, ok)
	require.Equal(t, 1, len(body.GetAuthorizationList()))
}

func TestUnitEthereumLegacyBuilderAndSign(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)

	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")

	tx := (&EthereumLegacyTransaction{}).
		SetNonce(1).
		SetGasPrice(big.NewInt(0x04a817c800)).
		SetGasLimit(21_000).
		SetTo(to).
		SetValue(big.NewInt(0)).
		SetCallData(nil)

	signed, err := tx.Sign(key)
	require.NoError(t, err)
	require.NotEmpty(t, signed)

	assert.Equal(t, 32, len(tx.GetR()))
	assert.Equal(t, 32, len(tx.GetS()))
	assert.True(t, tx.GetV() == 27 || tx.GetV() == 28)

	roundTrip, err := EthereumTransactionDataFromBytes(signed)
	require.NoError(t, err)
	body, ok := roundTrip.GetTransaction().(*EthereumLegacyTransaction)
	require.True(t, ok)
	assert.Equal(t, to, body.GetTo())
}

func TestUnitEthereumTransactionDataWrapAndSign(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)

	inner := (&EthereumEIP1559Transaction{}).
		SetChainId(298).
		SetNonce(0).
		SetMaxPriorityGas(big.NewInt(0)).
		SetMaxGas(big.NewInt(5)).
		SetGasLimit(150_000).
		SetTo(make([]byte, 20)).
		SetValue(big.NewInt(0)).
		SetCallData([]byte{0xCA, 0xFE})

	data := NewEthereumTransactionData(inner)
	require.NotNil(t, data)

	// GetTransaction returns the same pointer we passed in.
	assert.Same(t, inner, data.GetTransaction())

	signed, err := data.Sign(key)
	require.NoError(t, err)
	assert.Equal(t, byte(0x02), signed[0])

	// Inner was mutated with R/S/RecoveryId.
	assert.NotEmpty(t, inner.GetR())
	assert.NotEmpty(t, inner.GetS())
	assert.NotEmpty(t, inner.GetRecoveryIdBytes())

	// ToBytes() on the data matches the signed payload.
	dataBytes, err := data.ToBytes()
	require.NoError(t, err)
	assert.Equal(t, signed, dataBytes)
}

func TestUnitEthereumTransactionDataUnsupportedBody(t *testing.T) {
	t.Parallel()
	assert.Nil(t, NewEthereumTransactionData(nil))
}

func TestUnitEthereumEIP1559SignWithNonEcdsaKeyFails(t *testing.T) {
	t.Parallel()
	edKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx := (&EthereumEIP1559Transaction{}).
		SetChainId(1).
		SetNonce(0).
		SetGasLimit(21_000).
		SetTo(make([]byte, 20)).
		SetValue(big.NewInt(0))

	_, err = tx.Sign(edKey)
	assert.Error(t, err)
}
