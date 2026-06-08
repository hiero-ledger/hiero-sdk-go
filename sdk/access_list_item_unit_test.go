//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitAccessListItemGettersSetters(t *testing.T) {
	t.Parallel()

	addr := []byte{0x01, 0x02, 0x03}
	k1 := []byte{0x10}
	k2 := []byte{0x20}

	item := NewAccessListItem(addr, [][]byte{k1, k2})

	assert.Equal(t, addr, item.GetAddress())
	assert.Equal(t, [][]byte{k1, k2}, item.GetStorageKeys())

	newAddr := []byte{0xAA, 0xBB}
	item.SetAddress(newAddr)
	assert.Equal(t, newAddr, item.GetAddress())

	newKeys := [][]byte{{0x30}, {0x40}}
	item.SetStorageKeys(newKeys)
	assert.Equal(t, newKeys, item.GetStorageKeys())

	item.AddStorageKey([]byte{0x50})
	assert.Equal(t, 3, len(item.GetStorageKeys()))
	assert.Equal(t, []byte{0x50}, item.GetStorageKeys()[2])
}

func TestUnitAccessListItemChaining(t *testing.T) {
	t.Parallel()

	addr := []byte{0xDE, 0xAD}
	item := NewAccessListItem(nil, nil)
	item.SetAddress(addr).
		SetStorageKeys([][]byte{{0x01}}).
		AddStorageKey([]byte{0x02})

	assert.Equal(t, addr, item.GetAddress())
	require.Equal(t, 2, len(item.GetStorageKeys()))
}

func TestUnitAccessListItemRoundTrip(t *testing.T) {
	t.Parallel()

	original := NewAccessListItem(
		[]byte{0x01, 0x02, 0x03, 0x04},
		[][]byte{
			{0xA1, 0xA2},
			{0xB1, 0xB2, 0xB3},
		},
	)

	encoded := _accessListItemToBytes(original)
	decoded, err := _accessListItemFromBytes(encoded)
	require.NoError(t, err)

	assert.True(t, bytes.Equal(original.GetAddress(), decoded.GetAddress()))
	require.Equal(t, len(original.GetStorageKeys()), len(decoded.GetStorageKeys()))
	for i, k := range original.GetStorageKeys() {
		assert.True(t, bytes.Equal(k, decoded.GetStorageKeys()[i]))
	}
}

func TestUnitAccessListItemsRoundTrip(t *testing.T) {
	t.Parallel()

	items := []AccessListItem{
		NewAccessListItem([]byte{0x11}, [][]byte{{0x01}, {0x02}}),
		NewAccessListItem([]byte{0x22}, nil),
	}

	encoded := _accessListItemsToBytes(items)
	require.Equal(t, 2, len(encoded))

	decoded := _accessListItemsFromBytes(encoded)
	require.Equal(t, 2, len(decoded))
	assert.Equal(t, items[0].GetAddress(), decoded[0].GetAddress())
	assert.Equal(t, items[1].GetAddress(), decoded[1].GetAddress())
}

func TestUnitEIP1559AccessListItems(t *testing.T) {
	t.Parallel()

	tx := &EthereumEIP1559Transaction{}
	items := []AccessListItem{
		NewAccessListItem([]byte{0xAA}, [][]byte{{0x01}}),
		NewAccessListItem([]byte{0xBB}, [][]byte{{0x02}, {0x03}}),
	}
	tx.SetAccessListItems(items)

	roundTrip := tx.GetAccessListItems()
	require.Equal(t, 2, len(roundTrip))
	assert.Equal(t, []byte{0xAA}, roundTrip[0].GetAddress())
	assert.Equal(t, []byte{0xBB}, roundTrip[1].GetAddress())

	tx.AddAccessListItem(NewAccessListItem([]byte{0xCC}, nil))
	assert.Equal(t, 3, len(tx.GetAccessListItems()))
}

func TestUnitEIP2930AccessListItems(t *testing.T) {
	t.Parallel()

	tx := &EthereumEIP2930Transaction{}
	tx.AddAccessListItem(NewAccessListItem([]byte{0xAA}, [][]byte{{0x01}}))
	tx.AddAccessListItem(NewAccessListItem([]byte{0xBB}, nil))

	got := tx.GetAccessListItems()
	require.Equal(t, 2, len(got))
	assert.Equal(t, []byte{0xAA}, got[0].GetAddress())
	assert.Equal(t, []byte{0xBB}, got[1].GetAddress())
}

func TestUnitEIP7702AccessListItems(t *testing.T) {
	t.Parallel()

	tx := &EthereumEIP7702Transaction{}
	items := []AccessListItem{
		NewAccessListItem([]byte{0x01}, [][]byte{{0xFE}}),
	}
	tx.SetAccessListItems(items)

	got := tx.GetAccessListItems()
	require.Equal(t, 1, len(got))
	assert.Equal(t, []byte{0x01}, got[0].GetAddress())
	assert.Equal(t, [][]byte{{0xFE}}, got[0].GetStorageKeys())
}

func TestUnitAccessListItemEmptyList(t *testing.T) {
	t.Parallel()

	tx := &EthereumEIP1559Transaction{}
	assert.Equal(t, 0, len(tx.GetAccessListItems()))
	tx.SetAccessListItems(nil)
	assert.Equal(t, 0, len(tx.GetAccessListItems()))
}
