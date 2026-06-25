package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// AccessListItem represents an entry in an EIP-2930 access list.
// Each entry specifies an Ethereum address and a list of storage keys
// that the transaction plans to access.
type AccessListItem struct {
	address     []byte
	storageKeys [][]byte
}

// NewAccessListItem creates a new AccessListItem with the given address and storage keys.
func NewAccessListItem(address []byte, storageKeys [][]byte) AccessListItem {
	return AccessListItem{
		address:     address,
		storageKeys: storageKeys,
	}
}

// GetAddress returns the Ethereum address for this access list entry.
func (a *AccessListItem) GetAddress() []byte { return a.address }

// SetAddress sets the Ethereum address for this access list entry.
func (a *AccessListItem) SetAddress(address []byte) *AccessListItem {
	a.address = address
	return a
}

// GetStorageKeys returns the storage keys for this access list entry.
func (a *AccessListItem) GetStorageKeys() [][]byte { return a.storageKeys }

// SetStorageKeys sets the storage keys for this access list entry.
func (a *AccessListItem) SetStorageKeys(storageKeys [][]byte) *AccessListItem {
	a.storageKeys = storageKeys
	return a
}

// AddStorageKey appends a storage key to this access list entry.
func (a *AccessListItem) AddStorageKey(storageKey []byte) *AccessListItem {
	a.storageKeys = append(a.storageKeys, storageKey)
	return a
}

// String returns a string representation of the AccessListItem.
func (a AccessListItem) String() string {
	var encodedKeys []string
	for _, key := range a.storageKeys {
		encodedKeys = append(encodedKeys, hex.EncodeToString(key))
	}
	return fmt.Sprintf("{Address: %s, StorageKeys: [%s]}",
		hex.EncodeToString(a.address),
		strings.Join(encodedKeys, ", "),
	)
}

// _accessListItemToBytes encodes an AccessListItem as RLP [address, [storageKeys...]].
func _accessListItemToBytes(a AccessListItem) []byte {
	item := NewRLPItem(LIST_TYPE)
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(a.address))

	storageKeysItem := NewRLPItem(LIST_TYPE)
	for _, key := range a.storageKeys {
		storageKeysItem.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(key))
	}
	item.PushBack(storageKeysItem)

	bytes, _ := item.Write()
	return bytes
}

// _accessListItemsToBytes converts a slice of AccessListItem into [][]byte.
func _accessListItemsToBytes(items []AccessListItem) [][]byte {
	result := make([][]byte, len(items))
	for i, item := range items {
		result[i] = _accessListItemToBytes(item)
	}
	return result
}

// _accessListItemFromBytes decodes a single RLP-encoded byte slice into an AccessListItem.
func _accessListItemFromBytes(data []byte) (AccessListItem, error) {
	item := NewRLPItem(LIST_TYPE)
	if err := item.Read(data); err != nil {
		return AccessListItem{}, fmt.Errorf("failed to decode access list entry: %w", err)
	}

	if item.itemType != LIST_TYPE || len(item.childItems) != 2 {
		return AccessListItem{}, fmt.Errorf("invalid access list entry: expected list of [address, storageKeys]")
	}

	address := item.childItems[0].itemValue

	var storageKeys [][]byte
	if item.childItems[1].itemType == LIST_TYPE {
		for _, keyItem := range item.childItems[1].childItems {
			storageKeys = append(storageKeys, keyItem.itemValue)
		}
	}

	return NewAccessListItem(address, storageKeys), nil
}

// _accessListItemsFromBytes converts [][]byte into a slice of AccessListItem.
// Entries that fail to decode are skipped.
func _accessListItemsFromBytes(data [][]byte) []AccessListItem {
	var result []AccessListItem
	for _, entry := range data {
		ali, err := _accessListItemFromBytes(entry)
		if err != nil {
			continue
		}
		result = append(result, ali)
	}
	return result
}
