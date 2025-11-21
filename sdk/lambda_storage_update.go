package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

// Specifies a key/value pair in the storage of a lambda, either by the explicit storage
// slot contents; or by a combination of a Solidity mapping's slot key and the key into
// that mapping.
type LambdaStorageUpdate interface {
	toProtobuf() *services.LambdaStorageUpdate
}

func lambdaStorageUpdateFromProtobuf(pb *services.LambdaStorageUpdate) LambdaStorageUpdate {
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
type LambdaStorageSlot struct {
	key   []byte
	value []byte
}

// NewLambdaStorageSlot creates a new LambdaStorageSlot
func NewLambdaStorageSlot() *LambdaStorageSlot {
	return &LambdaStorageSlot{}
}

// GetKey returns the storage slot key
func (ls LambdaStorageSlot) GetKey() []byte {
	return ls.key
}

// SetKey sets the storage slot key
func (ls *LambdaStorageSlot) SetKey(key []byte) *LambdaStorageSlot {
	ls.key = key
	return ls
}

// GetValue returns the storage slot value
func (ls LambdaStorageSlot) GetValue() []byte {
	return ls.value
}

// SetValue sets the storage slot value
func (ls *LambdaStorageSlot) SetValue(value []byte) *LambdaStorageSlot {
	ls.value = value
	return ls
}

func (ls LambdaStorageSlot) toProtobuf() *services.LambdaStorageUpdate {
	return &services.LambdaStorageUpdate{
		Update: &services.LambdaStorageUpdate_StorageSlot{
			StorageSlot: &services.LambdaStorageSlot{
				Key:   ls.key,
				Value: ls.value,
			},
		},
	}
}

func lambdaStorageSlotFromProtobuf(pb *services.LambdaStorageSlot) LambdaStorageSlot {
	return LambdaStorageSlot{
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
type LambdaMappingEntries struct {
	mappingSlot    []byte
	mappingEntries []LambdaMappingEntry
}

// NewLambdaMappingEntries creates a new LambdaMappingEntries
func NewLambdaMappingEntries() *LambdaMappingEntries {
	return &LambdaMappingEntries{}
}

// GetMappingSlot returns the mapping slot
func (le LambdaMappingEntries) GetMappingSlot() []byte {
	return le.mappingSlot
}

// SetMappingSlot sets the mapping slot
func (le *LambdaMappingEntries) SetMappingSlot(mappingSlot []byte) *LambdaMappingEntries {
	le.mappingSlot = mappingSlot
	return le
}

// GetMappingEntries returns the mapping entries slice
func (le LambdaMappingEntries) GetMappingEntries() []LambdaMappingEntry {
	return le.mappingEntries
}

// SetMappingEntries sets the mapping entries slice
func (le *LambdaMappingEntries) SetMappingEntries(mappingEntries []LambdaMappingEntry) *LambdaMappingEntries {
	le.mappingEntries = mappingEntries
	return le
}

// AddMappingEntry adds a mapping entry to the slice
func (le *LambdaMappingEntries) AddMappingEntry(entry LambdaMappingEntry) *LambdaMappingEntries {
	le.mappingEntries = append(le.mappingEntries, entry)
	return le
}

func (le LambdaMappingEntries) toProtobuf() *services.LambdaStorageUpdate {
	mappingEntries := &services.LambdaMappingEntries{
		MappingSlot: le.mappingSlot,
	}

	for _, entry := range le.mappingEntries {
		mappingEntries.Entries = append(mappingEntries.Entries, entry.toProtobuf())
	}

	return &services.LambdaStorageUpdate{
		Update: &services.LambdaStorageUpdate_MappingEntries{
			MappingEntries: mappingEntries,
		},
	}
}

func lambdaMappingEntriesFromProtobuf(pb *services.LambdaMappingEntries) LambdaMappingEntries {
	mappingEntries := LambdaMappingEntries{
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
type LambdaMappingEntry struct {
	key      []byte
	preImage []byte
	value    []byte
}

// NewLambdaMappingEntryWithKey creates a new LambdaMappingEntry with key
func NewLambdaMappingEntryWithKey(key []byte, value []byte) *LambdaMappingEntry {
	return &LambdaMappingEntry{
		key:   key,
		value: value,
	}
}

// NewLambdaMappingEntryWithPreImage creates a new LambdaMappingEntry with preimage
func NewLambdaMappingEntryWithPreImage(preImage []byte, value []byte) *LambdaMappingEntry {
	return &LambdaMappingEntry{
		preImage: preImage,
		value:    value,
	}
}

// GetKey returns the mapping entry key
func (le LambdaMappingEntry) GetKey() []byte {
	return le.key
}

// SetKey sets the mapping entry key and clears preimage
func (le *LambdaMappingEntry) SetKey(key []byte) *LambdaMappingEntry {
	le.key = key
	le.preImage = nil
	return le
}

// GetPreImage returns the mapping entry preimage
func (le LambdaMappingEntry) GetPreImage() []byte {
	return le.preImage
}

// SetPreImage sets the mapping entry preimage and clears key
func (le *LambdaMappingEntry) SetPreImage(preImage []byte) *LambdaMappingEntry {
	le.preImage = preImage
	le.key = nil
	return le
}

// GetValue returns the mapping entry value
func (le LambdaMappingEntry) GetValue() []byte {
	return le.value
}

// SetValue sets the mapping entry value
func (le *LambdaMappingEntry) SetValue(value []byte) *LambdaMappingEntry {
	le.value = value
	return le
}

func (le LambdaMappingEntry) toProtobuf() *services.LambdaMappingEntry {
	pbBody := &services.LambdaMappingEntry{
		Value: le.value,
	}

	if len(le.key) > 0 {
		pbBody.EntryKey = &services.LambdaMappingEntry_Key{
			Key: le.key,
		}
	}
	if len(le.preImage) > 0 {
		pbBody.EntryKey = &services.LambdaMappingEntry_Preimage{
			Preimage: le.preImage,
		}
	}

	return pbBody
}

func lambdaMappingEntryFromProtobuf(pb *services.LambdaMappingEntry) LambdaMappingEntry {
	return LambdaMappingEntry{
		key:      pb.GetKey(),
		preImage: pb.GetPreimage(),
		value:    pb.GetValue(),
	}
}
