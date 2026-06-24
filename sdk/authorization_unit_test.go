//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitAuthorizationGettersSetters(t *testing.T) {
	t.Parallel()

	chainId := []byte{0x01, 0x2a}
	addr := []byte{0xAA, 0xBB}
	r := []byte{0x11}
	s := []byte{0x22}

	auth := NewAuthorization(chainId, addr, 5, 1, r, s)

	assert.Equal(t, chainId, auth.GetChainId())
	assert.Equal(t, addr, auth.GetAddress())
	assert.Equal(t, uint64(5), auth.GetNonce())
	assert.Equal(t, uint32(1), auth.GetYParity())
	assert.Equal(t, r, auth.GetR())
	assert.Equal(t, s, auth.GetS())

	auth.SetChainId([]byte{0x09}).
		SetAddress([]byte{0xCC}).
		SetNonce(7).
		SetYParity(0).
		SetR([]byte{0x33}).
		SetS([]byte{0x44})

	assert.Equal(t, []byte{0x09}, auth.GetChainId())
	assert.Equal(t, []byte{0xCC}, auth.GetAddress())
	assert.Equal(t, uint64(7), auth.GetNonce())
	assert.Equal(t, uint32(0), auth.GetYParity())
	assert.Equal(t, []byte{0x33}, auth.GetR())
	assert.Equal(t, []byte{0x44}, auth.GetS())
}

func TestUnitAuthorizationString(t *testing.T) {
	t.Parallel()

	auth := NewAuthorization([]byte{0xab, 0xcd}, []byte{0xde}, 3, 1, []byte{0x11}, []byte{0x22})
	str := auth.String()
	assert.Contains(t, str, "ChainId:")
	assert.Contains(t, str, "abcd")
	assert.Contains(t, str, "Nonce: 3")
	assert.Contains(t, str, "YParity: 1")
}

func TestUnitAuthorizationRLPRoundTrip(t *testing.T) {
	t.Parallel()

	auth := NewAuthorization([]byte{0x01, 0x2a}, make([]byte, 20), 5, 1, make([]byte, 32), make([]byte, 32))

	encoded, err := auth._toRLPItem().Write()
	require.NoError(t, err)

	item := NewRLPItem(LIST_TYPE)
	require.NoError(t, item.Read(encoded))

	decoded, err := _authorizationFromRLPItem(item)
	require.NoError(t, err)
	assert.Equal(t, auth, decoded)
}

func TestUnitAuthorizationFromRLPItemErrors(t *testing.T) {
	t.Parallel()

	// Not a list.
	_, err := _authorizationFromRLPItem(NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01}))
	assert.Error(t, err)

	// A list, but not 6 elements.
	short := NewRLPItem(LIST_TYPE)
	short.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01}))
	_, err = _authorizationFromRLPItem(short)
	assert.Error(t, err)
}
