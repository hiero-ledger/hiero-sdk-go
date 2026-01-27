package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

// Specifies a key/value pair in the storage of a lambda, either by the explicit storage
// slot contents; or by a combination of a Solidity mapping's slot key and the key into
// that mapping.
type EvmHookStorageUpdate interface {
	toProtobuf() *services.EvmHookStorageUpdate
}

func lambdaStorageUpdateFromProtobuf(pb *services.EvmHookStorageUpdate) EvmHookStorageUpdate {
	if pb.GetStorageSlot() != nil {
		return lambdaStorageSlotFromProtobuf(pb.GetStorageSlot())
	}
	if pb.GetMappingEntries() != nil {
		return lambdaMappingEntriesFromProtobuf(pb.GetMappingEntries())
	}
	return nil
}

/**
 * A slot in the storage of a lambda EVM hook.
 */
type EvmHookStorageSlot struct {
	key   []byte
	value []byte
}

// NewEvmHookStorageSlot creates a new LambdaStorageSlot
func NewEvmHookStorageSlot() *EvmHookStorageSlot {
	return &EvmHookStorageSlot{}
}

// GetKey returns the storage slot key
func (ls EvmHookStorageSlot) GetKey() []byte {
	return ls.key
}

// SetKey sets the storage slot key
func (ls *EvmHookStorageSlot) SetKey(key []byte) *EvmHookStorageSlot {
	ls.key = key
	return ls
}

// GetValue returns the storage slot value
func (ls EvmHookStorageSlot) GetValue() []byte {
	return ls.value
}

// SetValue sets the storage slot value
func (ls *EvmHookStorageSlot) SetValue(value []byte) *EvmHookStorageSlot {
	ls.value = value
	return ls
}

func (ls EvmHookStorageSlot) toProtobuf() *services.EvmHookStorageUpdate {
	return &services.EvmHookStorageUpdate{
		Update: &services.EvmHookStorageUpdate_StorageSlot{
			StorageSlot: &services.EvmHookStorageSlot{
				Key:   ls.key,
				Value: ls.value,
			},
		},
	}
}

func lambdaStorageSlotFromProtobuf(pb *services.EvmHookStorageSlot) EvmHookStorageSlot {
	return EvmHookStorageSlot{
		key:   pb.GetKey(),
		value: pb.GetValue(),
	}
}

// Specifies storage slot updates via indirection into a Solidity mapping.
// <p>
// Concretely, if the Solidity mapping is itself at slot `mapping_slot`, then
// the * storage slot for key `key` in the mapping is defined by the relationship
// `key_storage_slot = keccak256(abi.encodePacked(mapping_slot, key))`.
// <p>
// This message lets a metaprotocol be specified in terms of changes to a
// Solidity mapping's entries. If only raw slots could be updated, then a block
// stream consumer following the metaprotocol would have to invert the Keccak256
// hash to determine which mapping entry was being updated, which is not possible.
type EvmHookMappingEntries struct {
	mappingSlot    []byte
	mappingEntries []EvmHookMappingEntry
}

// NewLambdaMappingEntries creates a new LambdaMappingEntries
func NewLambdaMappingEntries() *EvmHookMappingEntries {
	return &EvmHookMappingEntries{}
}

// GetMappingSlot returns the mapping slot
func (le EvmHookMappingEntries) GetMappingSlot() []byte {
	return le.mappingSlot
}

// SetMappingSlot sets the mapping slot
func (le *EvmHookMappingEntries) SetMappingSlot(mappingSlot []byte) *EvmHookMappingEntries {
	le.mappingSlot = mappingSlot
	return le
}

// GetMappingEntries returns the mapping entries slice
func (le EvmHookMappingEntries) GetMappingEntries() []EvmHookMappingEntry {
	return le.mappingEntries
}

// SetMappingEntries sets the mapping entries slice
func (le *EvmHookMappingEntries) SetMappingEntries(mappingEntries []EvmHookMappingEntry) *EvmHookMappingEntries {
	le.mappingEntries = mappingEntries
	return le
}

// AddMappingEntry adds a mapping entry to the slice
func (le *EvmHookMappingEntries) AddMappingEntry(entry EvmHookMappingEntry) *EvmHookMappingEntries {
	le.mappingEntries = append(le.mappingEntries, entry)
	return le
}

func (le EvmHookMappingEntries) toProtobuf() *services.EvmHookStorageUpdate {
	mappingEntries := &services.EvmHookMappingEntries{
		MappingSlot: le.mappingSlot,
	}

	for _, entry := range le.mappingEntries {
		mappingEntries.Entries = append(mappingEntries.Entries, entry.toProtobuf())
	}

	return &services.EvmHookStorageUpdate{
		Update: &services.EvmHookStorageUpdate_MappingEntries{
			MappingEntries: mappingEntries,
		},
	}
}

func lambdaMappingEntriesFromProtobuf(pb *services.EvmHookMappingEntries) EvmHookMappingEntries {
	mappingEntries := EvmHookMappingEntries{
		mappingSlot: pb.GetMappingSlot(),
	}

	for _, entry := range pb.GetEntries() {
		mappingEntries.mappingEntries = append(mappingEntries.mappingEntries, lambdaMappingEntryFromProtobuf(entry))
	}

	return mappingEntries
}

// An entry in a Solidity mapping. Very helpful for protocols that apply
// `LambdaSStore` to manage the entries of a hook contract's mapping instead
// its raw storage slots.
// <p>
// This is especially attractive when the mapping value itself fits in a single
// word; for more complicated value storage layouts it becomes necessary to
// combine the mapping update with additional `LambdaStorageSlot` updates that
// specify the complete storage slots of the value type.
type EvmHookMappingEntry struct {
	key      []byte
	preImage []byte
	value    []byte
}

// NewEvmHookMappingEntryWithKey creates a new LambdaMappingEntry with key
func NewEvmHookMappingEntryWithKey(key []byte, value []byte) *EvmHookMappingEntry {
	return &EvmHookMappingEntry{
		key:   key,
		value: value,
	}
}

// NewEvmHookMappingEntryWithPreImage creates a new LambdaMappingEntry with preimage
func NewEvmHookMappingEntryWithPreImage(preImage []byte, value []byte) *EvmHookMappingEntry {
	return &EvmHookMappingEntry{
		preImage: preImage,
		value:    value,
	}
}

// GetKey returns the mapping entry key
func (le EvmHookMappingEntry) GetKey() []byte {
	return le.key
}

// SetKey sets the mapping entry key and clears preimage
func (le *EvmHookMappingEntry) SetKey(key []byte) *EvmHookMappingEntry {
	le.key = key
	le.preImage = nil
	return le
}

// GetPreImage returns the mapping entry preimage
func (le EvmHookMappingEntry) GetPreImage() []byte {
	return le.preImage
}

// SetPreImage sets the mapping entry preimage and clears key
func (le *EvmHookMappingEntry) SetPreImage(preImage []byte) *EvmHookMappingEntry {
	le.preImage = preImage
	le.key = nil
	return le
}

// GetValue returns the mapping entry value
func (le EvmHookMappingEntry) GetValue() []byte {
	return le.value
}

// SetValue sets the mapping entry value
func (le *EvmHookMappingEntry) SetValue(value []byte) *EvmHookMappingEntry {
	le.value = value
	return le
}

func (le EvmHookMappingEntry) toProtobuf() *services.EvmHookMappingEntry {
	pbBody := &services.EvmHookMappingEntry{
		Value: le.value,
	}

	if len(le.key) > 0 {
		pbBody.EntryKey = &services.EvmHookMappingEntry_Key{
			Key: le.key,
		}
	}
	if len(le.preImage) > 0 {
		pbBody.EntryKey = &services.EvmHookMappingEntry_Preimage{
			Preimage: le.preImage,
		}
	}

	return pbBody
}

func lambdaMappingEntryFromProtobuf(pb *services.EvmHookMappingEntry) EvmHookMappingEntry {
	return EvmHookMappingEntry{
		key:      pb.GetKey(),
		preImage: pb.GetPreimage(),
		value:    pb.GetValue(),
	}
}
