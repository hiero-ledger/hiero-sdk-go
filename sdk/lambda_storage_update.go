package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

// Specifies a key/value pair in the storage of a lambda, either by the explicit storage slot contents; or by a combination of a Solidity mapping's slot key and the key into that mapping.
type LambdaStorageUpdate struct {
	// oneof (mutually exclusive):

	// An explicit storage slot update.
	storageSlot *LambdaStorageSlot

	// An implicit storage slot update specified as a Solidity mapping entry.
	mappingEntries *LambdaMappingEntries
}

// NewLambdaStorageUpdateWithStorageSlot creates a new LambdaStorageUpdate with storage slot
func NewLambdaStorageUpdateWithStorageSlot(storageSlot LambdaStorageSlot) LambdaStorageUpdate {
	return LambdaStorageUpdate{
		storageSlot: &storageSlot,
	}
}

// NewLambdaStorageUpdateWithMappingEntries creates a new LambdaStorageUpdate with mapping entries
func NewLambdaStorageUpdateWithMappingEntries(mappingEntries LambdaMappingEntries) LambdaStorageUpdate {
	return LambdaStorageUpdate{
		mappingEntries: &mappingEntries,
	}
}

// GetStorageSlot returns the storage slot
func (lu LambdaStorageUpdate) GetStorageSlot() *LambdaStorageSlot {
	return lu.storageSlot
}

// SetStorageSlot sets the storage slot and clears mapping entries
func (lu *LambdaStorageUpdate) SetStorageSlot(storageSlot LambdaStorageSlot) *LambdaStorageUpdate {
	lu.storageSlot = &storageSlot
	lu.mappingEntries = nil
	return lu
}

// GetMappingEntries returns the mapping entries
func (lu LambdaStorageUpdate) GetMappingEntries() *LambdaMappingEntries {
	return lu.mappingEntries
}

// SetMappingEntries sets the mapping entries and clears storage slot
func (lu *LambdaStorageUpdate) SetMappingEntries(mappingEntries LambdaMappingEntries) *LambdaStorageUpdate {
	lu.mappingEntries = &mappingEntries
	lu.storageSlot = nil
	return lu
}

func lambdaStorageUpdateFromProtobuf(pb *services.LambdaStorageUpdate) LambdaStorageUpdate {
	lu := LambdaStorageUpdate{}
	if pb.GetStorageSlot() != nil {
		storageSlot := pb.GetStorageSlot()
		lu.storageSlot = &LambdaStorageSlot{
			key:   storageSlot.GetKey(),
			value: storageSlot.GetValue(),
		}
	}
	if pb.GetMappingEntries() != nil {
		mappingEntries := pb.GetMappingEntries()
		entries := make([]LambdaMappingEntry, 0)
		for _, entry := range mappingEntries.GetEntries() {
			entries = append(entries, LambdaMappingEntry{
				key:      entry.GetKey(),
				value:    entry.GetValue(),
				preImage: entry.GetPreimage(),
			})
		}
		lu.mappingEntries = &LambdaMappingEntries{
			mappingSlot:    mappingEntries.GetMappingSlot(),
			mappingEntries: entries,
		}
	}
	return lu
}

func (lu LambdaStorageUpdate) toProtobuf() *services.LambdaStorageUpdate {
	if lu.storageSlot != nil {
		return &services.LambdaStorageUpdate{
			Update: &services.LambdaStorageUpdate_StorageSlot{
				StorageSlot: lu.storageSlot.toProtobuf(),
			},
		}
	}
	if lu.mappingEntries != nil {
		return &services.LambdaStorageUpdate{
			Update: &services.LambdaStorageUpdate_MappingEntries{
				MappingEntries: lu.mappingEntries.toProtobuf(),
			},
		}
	}
	return nil
}

type LambdaStorageSlot struct {
	// 32-byte storage slot key
	key []byte

	// 32-byte storage slot value (empty to delete)
	value []byte
}

// NewLambdaStorageSlot creates a new LambdaStorageSlot
func NewLambdaStorageSlot(key []byte, value []byte) LambdaStorageSlot {
	return LambdaStorageSlot{
		key:   key,
		value: value,
	}
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

func (ls LambdaStorageSlot) toProtobuf() *services.LambdaStorageSlot {
	return &services.LambdaStorageSlot{
		Key:   ls.key,
		Value: ls.value,
	}
}

type LambdaMappingEntries struct {
	mappingSlot    []byte
	mappingEntries []LambdaMappingEntry
}

// NewLambdaMappingEntries creates a new LambdaMappingEntries
func NewLambdaMappingEntries(mappingSlot []byte, mappingEntries []LambdaMappingEntry) LambdaMappingEntries {
	return LambdaMappingEntries{
		mappingSlot:    mappingSlot,
		mappingEntries: mappingEntries,
	}
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

func (le LambdaMappingEntries) toProtobuf() *services.LambdaMappingEntries {
	pbBody := &services.LambdaMappingEntries{
		MappingSlot: le.mappingSlot,
	}

	for _, entry := range le.mappingEntries {
		pbBody.Entries = append(pbBody.Entries, entry.toProtobuf())
	}

	return pbBody
}

type LambdaMappingEntry struct {
	// one of:

	// 32-byte key of mapping entry
	key []byte
	// The bytes that are the preimage of the Keccak256 hash that forms the mapping key.
	preImage []byte

	// 32-byte value of mapping entry (empty to delete)
	value []byte
}

// NewLambdaMappingEntryWithKey creates a new LambdaMappingEntry with key
func NewLambdaMappingEntryWithKey(key []byte, value []byte) LambdaMappingEntry {
	return LambdaMappingEntry{
		key:   key,
		value: value,
	}
}

// NewLambdaMappingEntryWithPreImage creates a new LambdaMappingEntry with preimage
func NewLambdaMappingEntryWithPreImage(preImage []byte, value []byte) LambdaMappingEntry {
	return LambdaMappingEntry{
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
