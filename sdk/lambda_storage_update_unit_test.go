//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitLambdaStorageSlotKey(t *testing.T) {
	t.Parallel()

	ls := NewEvmHookStorageSlot()

	// Test default value
	assert.Nil(t, ls.GetKey())

	// Test setting key
	key := []byte("test_key")
	ls.SetKey(key)
	assert.Equal(t, key, ls.GetKey())

	// Test setting nil key
	ls.SetKey(nil)
	assert.Nil(t, ls.GetKey())

	// Test setting empty key
	ls.SetKey([]byte{})
	assert.Equal(t, []byte{}, ls.GetKey())
}

func TestUnitLambdaStorageSlotValue(t *testing.T) {
	t.Parallel()

	ls := NewEvmHookStorageSlot()

	// Test default value
	assert.Nil(t, ls.GetValue())

	// Test setting value
	value := []byte("test_value")
	ls.SetValue(value)
	assert.Equal(t, value, ls.GetValue())

	// Test setting nil value
	ls.SetValue(nil)
	assert.Nil(t, ls.GetValue())

	// Test setting empty value
	ls.SetValue([]byte{})
	assert.Equal(t, []byte{}, ls.GetValue())
}

func TestUnitLambdaStorageSlotMethodChaining(t *testing.T) {
	t.Parallel()

	key := []byte("chained_key")
	value := []byte("chained_value")

	// Test method chaining
	ls := NewEvmHookStorageSlot().SetKey(key).SetValue(value)

	assert.Equal(t, key, ls.GetKey())
	assert.Equal(t, value, ls.GetValue())
}

func TestUnitLambdaStorageSlotToProtobuf(t *testing.T) {
	t.Parallel()

	// Test with all fields set
	key := []byte("protobuf_key")
	value := []byte("protobuf_value")

	ls := NewEvmHookStorageSlot().SetKey(key).SetValue(value)
	pb := ls.toProtobuf()

	require.NotNil(t, pb)
	assert.NotNil(t, pb.Update)
	assert.NotNil(t, pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot)
	assert.Equal(t, key, pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot.Key)
	assert.Equal(t, value, pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot.Value)
}

func TestUnitLambdaStorageSlotToProtobufWithNilValues(t *testing.T) {
	t.Parallel()

	// Test with nil values
	ls := NewEvmHookStorageSlot()
	pb := ls.toProtobuf()

	require.NotNil(t, pb)
	assert.NotNil(t, pb.Update)
	assert.NotNil(t, pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot)
	assert.Nil(t, pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot.Key)
	assert.Nil(t, pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot.Value)
}

func TestUnitLambdaStorageSlotFromProtobuf(t *testing.T) {
	t.Parallel()

	// Test conversion from protobuf
	key := []byte("from_protobuf_key")
	value := []byte("from_protobuf_value")

	pb := &services.LambdaStorageSlot{
		Key:   key,
		Value: value,
	}

	ls := lambdaStorageSlotFromProtobuf(pb)
	assert.Equal(t, key, ls.GetKey())
	assert.Equal(t, value, ls.GetValue())
}

func TestUnitLambdaStorageSlotFromProtobufWithNilValues(t *testing.T) {
	t.Parallel()

	// Test conversion from protobuf with nil values
	pb := &services.LambdaStorageSlot{
		Key:   nil,
		Value: nil,
	}

	ls := lambdaStorageSlotFromProtobuf(pb)
	assert.Nil(t, ls.GetKey())
	assert.Nil(t, ls.GetValue())
}

func TestUnitLambdaStorageSlotRoundTrip(t *testing.T) {
	t.Parallel()

	// Test round trip conversion
	key := []byte("roundtrip_key")
	value := []byte("roundtrip_value")

	original := NewEvmHookStorageSlot().SetKey(key).SetValue(value)

	// Convert to protobuf and back
	pb := original.toProtobuf()
	converted := lambdaStorageSlotFromProtobuf(pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot)

	assert.Equal(t, original.GetKey(), converted.GetKey())
	assert.Equal(t, original.GetValue(), converted.GetValue())
}

func TestUnitLambdaMappingEntriesNew(t *testing.T) {
	t.Parallel()

	lme := NewLambdaMappingEntries()
	require.NotNil(t, lme)
	assert.Nil(t, lme.GetMappingSlot())
	assert.Equal(t, []EvmHookMappingEntry(nil), lme.GetMappingEntries())
}

func TestUnitLambdaMappingEntriesMappingSlot(t *testing.T) {
	t.Parallel()

	lme := NewLambdaMappingEntries()

	// Test default value
	assert.Nil(t, lme.GetMappingSlot())

	// Test setting mapping slot
	mappingSlot := []byte("mapping_slot_123")
	lme.SetMappingSlot(mappingSlot)
	assert.Equal(t, mappingSlot, lme.GetMappingSlot())

	// Test setting nil mapping slot
	lme.SetMappingSlot(nil)
	assert.Nil(t, lme.GetMappingSlot())

	// Test setting empty mapping slot
	lme.SetMappingSlot([]byte{})
	assert.Equal(t, []byte{}, lme.GetMappingSlot())

	// Test method chaining
	result := lme.SetMappingSlot(mappingSlot)
	assert.Equal(t, lme, result)
}

func TestUnitLambdaMappingEntriesMappingEntries(t *testing.T) {
	t.Parallel()

	lme := NewLambdaMappingEntries()

	// Test default value
	assert.Equal(t, []EvmHookMappingEntry(nil), lme.GetMappingEntries())

	// Test setting mapping entries
	entry1 := NewEvmHookMappingEntryWithKey([]byte("key1"), []byte("value1"))
	entry2 := NewEvmHookMappingEntryWithPreImage([]byte("preimage2"), []byte("value2"))
	entries := []EvmHookMappingEntry{*entry1, *entry2}

	lme.SetMappingEntries(entries)
	assert.Equal(t, entries, lme.GetMappingEntries())

	// Test setting empty mapping entries
	lme.SetMappingEntries([]EvmHookMappingEntry{})
	assert.Equal(t, []EvmHookMappingEntry{}, lme.GetMappingEntries())

	// Test setting nil mapping entries
	lme.SetMappingEntries(nil)
	assert.Equal(t, []EvmHookMappingEntry(nil), lme.GetMappingEntries())

	// Test method chaining
	result := lme.SetMappingEntries(entries)
	assert.Equal(t, lme, result)
}

func TestUnitLambdaMappingEntriesAddMappingEntry(t *testing.T) {
	t.Parallel()

	lme := NewLambdaMappingEntries()

	// Test adding to empty slice
	entry1 := NewEvmHookMappingEntryWithKey([]byte("key1"), []byte("value1"))
	lme.AddMappingEntry(*entry1)
	assert.Len(t, lme.GetMappingEntries(), 1)
	assert.Equal(t, *entry1, lme.GetMappingEntries()[0])

	// Test adding to existing slice
	entry2 := NewEvmHookMappingEntryWithPreImage([]byte("preimage2"), []byte("value2"))
	lme.AddMappingEntry(*entry2)
	assert.Len(t, lme.GetMappingEntries(), 2)
	assert.Equal(t, *entry1, lme.GetMappingEntries()[0])
	assert.Equal(t, *entry2, lme.GetMappingEntries()[1])

	// Test method chaining
	entry3 := NewEvmHookMappingEntryWithKey([]byte("key3"), []byte("value3"))
	result := lme.AddMappingEntry(*entry3)
	assert.Equal(t, lme, result)
	assert.Len(t, lme.GetMappingEntries(), 3)
}

func TestUnitLambdaMappingEntriesMethodChaining(t *testing.T) {
	t.Parallel()

	mappingSlot := []byte("chained_slot")
	entry1 := NewEvmHookMappingEntryWithKey([]byte("key1"), []byte("value1"))
	entry2 := NewEvmHookMappingEntryWithPreImage([]byte("preimage2"), []byte("value2"))

	// Test method chaining
	lme := NewLambdaMappingEntries().
		SetMappingSlot(mappingSlot).
		AddMappingEntry(*entry1).
		AddMappingEntry(*entry2)

	assert.Equal(t, mappingSlot, lme.GetMappingSlot())
	assert.Len(t, lme.GetMappingEntries(), 2)
	assert.Equal(t, *entry1, lme.GetMappingEntries()[0])
	assert.Equal(t, *entry2, lme.GetMappingEntries()[1])
}

func TestUnitLambdaMappingEntriesToProtobuf(t *testing.T) {
	t.Parallel()

	// Test with all fields set
	mappingSlot := []byte("protobuf_slot")
	entry1 := NewEvmHookMappingEntryWithKey([]byte("key1"), []byte("value1"))
	entry2 := NewEvmHookMappingEntryWithPreImage([]byte("preimage2"), []byte("value2"))

	lme := NewLambdaMappingEntries().
		SetMappingSlot(mappingSlot).
		AddMappingEntry(*entry1).
		AddMappingEntry(*entry2)

	pb := lme.toProtobuf()

	require.NotNil(t, pb)
	assert.NotNil(t, pb.Update)
	assert.NotNil(t, pb.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries)
	assert.Equal(t, mappingSlot, pb.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries.MappingSlot)
	assert.Len(t, pb.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries.Entries, 2)
}

func TestUnitLambdaMappingEntriesToProtobufWithEmptyEntries(t *testing.T) {
	t.Parallel()

	// Test with empty entries
	mappingSlot := []byte("empty_slot")

	lme := NewLambdaMappingEntries().SetMappingSlot(mappingSlot)
	pb := lme.toProtobuf()

	require.NotNil(t, pb)
	assert.NotNil(t, pb.Update)
	assert.NotNil(t, pb.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries)
	assert.Equal(t, mappingSlot, pb.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries.MappingSlot)
	assert.Len(t, pb.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries.Entries, 0)
}

func TestUnitLambdaMappingEntriesFromProtobuf(t *testing.T) {
	t.Parallel()

	// Test conversion from protobuf
	mappingSlot := []byte("from_protobuf_slot")
	entry1 := &services.LambdaMappingEntry{
		EntryKey: &services.LambdaMappingEntry_Key{Key: []byte("key1")},
		Value:    []byte("value1"),
	}
	entry2 := &services.LambdaMappingEntry{
		EntryKey: &services.LambdaMappingEntry_Preimage{Preimage: []byte("preimage2")},
		Value:    []byte("value2"),
	}

	pb := &services.LambdaMappingEntries{
		MappingSlot: mappingSlot,
		Entries:     []*services.LambdaMappingEntry{entry1, entry2},
	}

	lme := lambdaMappingEntriesFromProtobuf(pb)
	assert.Equal(t, mappingSlot, lme.GetMappingSlot())
	assert.Len(t, lme.GetMappingEntries(), 2)
}

func TestUnitLambdaMappingEntriesFromProtobufWithEmptyEntries(t *testing.T) {
	t.Parallel()

	// Test conversion from protobuf with empty entries
	mappingSlot := []byte("empty_from_protobuf_slot")

	pb := &services.LambdaMappingEntries{
		MappingSlot: mappingSlot,
		Entries:     []*services.LambdaMappingEntry{},
	}

	lme := lambdaMappingEntriesFromProtobuf(pb)
	assert.Equal(t, mappingSlot, lme.GetMappingSlot())
	assert.Len(t, lme.GetMappingEntries(), 0)
}

func TestUnitLambdaMappingEntriesRoundTrip(t *testing.T) {
	t.Parallel()

	// Test round trip conversion
	mappingSlot := []byte("roundtrip_slot")
	entry1 := NewEvmHookMappingEntryWithKey([]byte("key1"), []byte("value1"))
	entry2 := NewEvmHookMappingEntryWithPreImage([]byte("preimage2"), []byte("value2"))

	original := NewLambdaMappingEntries().
		SetMappingSlot(mappingSlot).
		AddMappingEntry(*entry1).
		AddMappingEntry(*entry2)

	// Convert to protobuf and back
	pb := original.toProtobuf()
	converted := lambdaMappingEntriesFromProtobuf(pb.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries)

	assert.Equal(t, original.GetMappingSlot(), converted.GetMappingSlot())
	assert.Len(t, converted.GetMappingEntries(), 2)
}

func TestUnitLambdaMappingEntryWithKey(t *testing.T) {
	t.Parallel()

	key := []byte("test_key")
	value := []byte("test_value")

	entry := NewEvmHookMappingEntryWithKey(key, value)
	require.NotNil(t, entry)
	assert.Equal(t, key, entry.GetKey())
	assert.Equal(t, value, entry.GetValue())
	assert.Nil(t, entry.GetPreImage())
}

func TestUnitLambdaMappingEntryWithPreImage(t *testing.T) {
	t.Parallel()

	preImage := []byte("test_preimage")
	value := []byte("test_value")

	entry := NewEvmHookMappingEntryWithPreImage(preImage, value)
	require.NotNil(t, entry)
	assert.Equal(t, preImage, entry.GetPreImage())
	assert.Equal(t, value, entry.GetValue())
	assert.Nil(t, entry.GetKey())
}

func TestUnitLambdaMappingEntryKey(t *testing.T) {
	t.Parallel()

	entry := &EvmHookMappingEntry{}

	// Test default value
	assert.Nil(t, entry.GetKey())

	// Test setting key
	key := []byte("test_key")
	entry.SetKey(key)
	assert.Equal(t, key, entry.GetKey())
	assert.Nil(t, entry.GetPreImage()) // Should clear preimage

	// Test setting nil key
	entry.SetKey(nil)
	assert.Nil(t, entry.GetKey())

	// Test setting empty key
	entry.SetKey([]byte{})
	assert.Equal(t, []byte{}, entry.GetKey())

	// Test method chaining
	result := entry.SetKey(key)
	assert.Equal(t, entry, result)
}

func TestUnitLambdaMappingEntryPreImage(t *testing.T) {
	t.Parallel()

	entry := &EvmHookMappingEntry{}

	// Test default value
	assert.Nil(t, entry.GetPreImage())

	// Test setting preimage
	preImage := []byte("test_preimage")
	entry.SetPreImage(preImage)
	assert.Equal(t, preImage, entry.GetPreImage())
	assert.Nil(t, entry.GetKey()) // Should clear key

	// Test setting nil preimage
	entry.SetPreImage(nil)
	assert.Nil(t, entry.GetPreImage())

	// Test setting empty preimage
	entry.SetPreImage([]byte{})
	assert.Equal(t, []byte{}, entry.GetPreImage())

	// Test method chaining
	result := entry.SetPreImage(preImage)
	assert.Equal(t, entry, result)
}

func TestUnitLambdaMappingEntryValue(t *testing.T) {
	t.Parallel()

	entry := &EvmHookMappingEntry{}

	// Test default value
	assert.Nil(t, entry.GetValue())

	// Test setting value
	value := []byte("test_value")
	entry.SetValue(value)
	assert.Equal(t, value, entry.GetValue())

	// Test setting nil value
	entry.SetValue(nil)
	assert.Nil(t, entry.GetValue())

	// Test setting empty value
	entry.SetValue([]byte{})
	assert.Equal(t, []byte{}, entry.GetValue())

	// Test method chaining
	result := entry.SetValue(value)
	assert.Equal(t, entry, result)
}

func TestUnitLambdaMappingEntryKeyPreImageInteraction(t *testing.T) {
	t.Parallel()

	entry := &EvmHookMappingEntry{}

	// Test that setting key clears preimage
	entry.SetPreImage([]byte("preimage"))
	assert.Equal(t, []byte("preimage"), entry.GetPreImage())
	assert.Nil(t, entry.GetKey())

	entry.SetKey([]byte("key"))
	assert.Equal(t, []byte("key"), entry.GetKey())
	assert.Nil(t, entry.GetPreImage())

	// Test that setting preimage clears key
	entry.SetKey([]byte("key"))
	assert.Equal(t, []byte("key"), entry.GetKey())
	assert.Nil(t, entry.GetPreImage())

	entry.SetPreImage([]byte("preimage"))
	assert.Equal(t, []byte("preimage"), entry.GetPreImage())
	assert.Nil(t, entry.GetKey())
}

func TestUnitLambdaMappingEntryToProtobufWithKey(t *testing.T) {
	t.Parallel()

	// Test with key
	key := []byte("protobuf_key")
	value := []byte("protobuf_value")

	entry := NewEvmHookMappingEntryWithKey(key, value)
	pb := entry.toProtobuf()

	require.NotNil(t, pb)
	assert.Equal(t, value, pb.Value)
	assert.NotNil(t, pb.EntryKey)
	assert.NotNil(t, pb.EntryKey.(*services.LambdaMappingEntry_Key).Key)
	assert.Equal(t, key, pb.EntryKey.(*services.LambdaMappingEntry_Key).Key)
}

func TestUnitLambdaMappingEntryToProtobufWithPreImage(t *testing.T) {
	t.Parallel()

	// Test with preimage
	preImage := []byte("protobuf_preimage")
	value := []byte("protobuf_value")

	entry := NewEvmHookMappingEntryWithPreImage(preImage, value)
	pb := entry.toProtobuf()

	require.NotNil(t, pb)
	assert.Equal(t, value, pb.Value)
	assert.NotNil(t, pb.EntryKey)
	assert.NotNil(t, pb.EntryKey.(*services.LambdaMappingEntry_Preimage).Preimage)
	assert.Equal(t, preImage, pb.EntryKey.(*services.LambdaMappingEntry_Preimage).Preimage)
}

func TestUnitLambdaMappingEntryToProtobufWithNilValues(t *testing.T) {
	t.Parallel()

	// Test with nil values
	entry := &EvmHookMappingEntry{}
	pb := entry.toProtobuf()

	require.NotNil(t, pb)
	assert.Nil(t, pb.Value)
	assert.Nil(t, pb.EntryKey)
}

func TestUnitLambdaMappingEntryFromProtobufWithKey(t *testing.T) {
	t.Parallel()

	// Test conversion from protobuf with key
	key := []byte("from_protobuf_key")
	value := []byte("from_protobuf_value")

	pb := &services.LambdaMappingEntry{
		EntryKey: &services.LambdaMappingEntry_Key{Key: key},
		Value:    value,
	}

	entry := lambdaMappingEntryFromProtobuf(pb)
	assert.Equal(t, key, entry.GetKey())
	assert.Equal(t, value, entry.GetValue())
	assert.Nil(t, entry.GetPreImage())
}

func TestUnitLambdaMappingEntryFromProtobufWithPreImage(t *testing.T) {
	t.Parallel()

	// Test conversion from protobuf with preimage
	preImage := []byte("from_protobuf_preimage")
	value := []byte("from_protobuf_value")

	pb := &services.LambdaMappingEntry{
		EntryKey: &services.LambdaMappingEntry_Preimage{Preimage: preImage},
		Value:    value,
	}

	entry := lambdaMappingEntryFromProtobuf(pb)
	assert.Equal(t, preImage, entry.GetPreImage())
	assert.Equal(t, value, entry.GetValue())
	assert.Nil(t, entry.GetKey())
}

func TestUnitLambdaMappingEntryFromProtobufWithNilValues(t *testing.T) {
	t.Parallel()

	// Test conversion from protobuf with nil values
	pb := &services.LambdaMappingEntry{
		EntryKey: nil,
		Value:    nil,
	}

	entry := lambdaMappingEntryFromProtobuf(pb)
	assert.Nil(t, entry.GetKey())
	assert.Nil(t, entry.GetPreImage())
	assert.Nil(t, entry.GetValue())
}

func TestUnitLambdaMappingEntryRoundTrip(t *testing.T) {
	t.Parallel()

	// Test round trip conversion with key
	key := []byte("roundtrip_key")
	value := []byte("roundtrip_value")

	original := NewEvmHookMappingEntryWithKey(key, value)
	pb := original.toProtobuf()
	converted := lambdaMappingEntryFromProtobuf(pb)

	assert.Equal(t, original.GetKey(), converted.GetKey())
	assert.Equal(t, original.GetValue(), converted.GetValue())
	assert.Equal(t, original.GetPreImage(), converted.GetPreImage())

	// Test round trip conversion with preimage
	preImage := []byte("roundtrip_preimage")
	value2 := []byte("roundtrip_value2")

	original2 := NewEvmHookMappingEntryWithPreImage(preImage, value2)
	pb2 := original2.toProtobuf()
	converted2 := lambdaMappingEntryFromProtobuf(pb2)

	assert.Equal(t, original2.GetKey(), converted2.GetKey())
	assert.Equal(t, original2.GetValue(), converted2.GetValue())
	assert.Equal(t, original2.GetPreImage(), converted2.GetPreImage())
}

func TestUnitLambdaStorageUpdateFromProtobufWithStorageSlot(t *testing.T) {
	t.Parallel()

	// Test with storage slot
	key := []byte("storage_slot_key")
	value := []byte("storage_slot_value")

	pb := &services.LambdaStorageUpdate{
		Update: &services.LambdaStorageUpdate_StorageSlot{
			StorageSlot: &services.LambdaStorageSlot{
				Key:   key,
				Value: value,
			},
		},
	}

	update := lambdaStorageUpdateFromProtobuf(pb)
	require.NotNil(t, update)
	storageSlot := update.(EvmHookStorageSlot)
	assert.Equal(t, key, storageSlot.GetKey())
	assert.Equal(t, value, storageSlot.GetValue())
}

func TestUnitLambdaStorageUpdateFromProtobufWithMappingEntries(t *testing.T) {
	t.Parallel()

	// Test with mapping entries
	mappingSlot := []byte("mapping_slot")
	entry := &services.LambdaMappingEntry{
		EntryKey: &services.LambdaMappingEntry_Key{Key: []byte("key")},
		Value:    []byte("value"),
	}

	pb := &services.LambdaStorageUpdate{
		Update: &services.LambdaStorageUpdate_MappingEntries{
			MappingEntries: &services.LambdaMappingEntries{
				MappingSlot: mappingSlot,
				Entries:     []*services.LambdaMappingEntry{entry},
			},
		},
	}

	update := lambdaStorageUpdateFromProtobuf(pb)
	require.NotNil(t, update)
	mappingEntries := update.(EvmHookMappingEntries)
	assert.Equal(t, mappingSlot, mappingEntries.GetMappingSlot())
	assert.Len(t, mappingEntries.GetMappingEntries(), 1)
}

func TestUnitLambdaStorageUpdateFromProtobufWithNil(t *testing.T) {
	t.Parallel()

	// Test with nil update
	pb := &services.LambdaStorageUpdate{
		Update: nil,
	}

	update := lambdaStorageUpdateFromProtobuf(pb)
	assert.Nil(t, update)
}

func TestUnitLambdaStorageUpdateEdgeCases(t *testing.T) {
	t.Parallel()

	// Test storage slot with empty values
	ls := NewEvmHookStorageSlot().SetKey([]byte{}).SetValue([]byte{})
	assert.Equal(t, []byte{}, ls.GetKey())
	assert.Equal(t, []byte{}, ls.GetValue())

	pb := ls.toProtobuf()
	require.NotNil(t, pb)
	converted := lambdaStorageSlotFromProtobuf(pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot)
	assert.Equal(t, []byte{}, converted.GetKey())
	assert.Equal(t, []byte{}, converted.GetValue())

	// Test mapping entries with empty values
	entry := NewEvmHookMappingEntryWithKey([]byte{}, []byte{})
	mappingEntries := NewLambdaMappingEntries().
		SetMappingSlot([]byte{}).
		AddMappingEntry(*entry)

	assert.Equal(t, []byte{}, mappingEntries.GetMappingSlot())
	assert.Len(t, mappingEntries.GetMappingEntries(), 1)

	pb2 := mappingEntries.toProtobuf()
	require.NotNil(t, pb2)
	converted2 := lambdaMappingEntriesFromProtobuf(pb2.Update.(*services.LambdaStorageUpdate_MappingEntries).MappingEntries)
	assert.Equal(t, []byte{}, converted2.GetMappingSlot())
	assert.Len(t, converted2.GetMappingEntries(), 1)
}

func TestUnitLambdaStorageUpdateLargeData(t *testing.T) {
	t.Parallel()

	// Test with large data
	largeKey := make([]byte, 1000)
	largeValue := make([]byte, 2000)
	for i := range largeKey {
		largeKey[i] = byte(i % 256)
	}
	for i := range largeValue {
		largeValue[i] = byte(i % 256)
	}

	// Test storage slot with large data
	ls := NewEvmHookStorageSlot().SetKey(largeKey).SetValue(largeValue)
	pb := ls.toProtobuf()
	converted := lambdaStorageSlotFromProtobuf(pb.Update.(*services.LambdaStorageUpdate_StorageSlot).StorageSlot)
	assert.Equal(t, largeKey, converted.GetKey())
	assert.Equal(t, largeValue, converted.GetValue())

	// Test mapping entry with large data
	entry := NewEvmHookMappingEntryWithKey(largeKey, largeValue)
	pb2 := entry.toProtobuf()
	converted2 := lambdaMappingEntryFromProtobuf(pb2)
	assert.Equal(t, largeKey, converted2.GetKey())
	assert.Equal(t, largeValue, converted2.GetValue())
}
