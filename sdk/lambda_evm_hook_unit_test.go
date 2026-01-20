//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitLambdaEvmHookContractId(t *testing.T) {
	t.Parallel()

	leh := NewEvmHook()

	// Test default value
	assert.Nil(t, leh.GetContractId())

	// Test setting contract ID
	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)
	leh.SetContractId(&contractID)
	assert.Equal(t, &contractID, leh.GetContractId())

	// Test setting nil contract ID
	leh.SetContractId(nil)
	assert.Nil(t, leh.GetContractId())
}

func TestUnitLambdaEvmHookStorageUpdates(t *testing.T) {
	t.Parallel()

	leh := NewEvmHook()

	// Test default value
	assert.Equal(t, []EvmHookStorageUpdate(nil), leh.GetStorageUpdates())

	// Test setting storage updates
	storageSlot1 := NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1"))
	storageSlot2 := NewEvmHookStorageSlot().SetKey([]byte("key2")).SetValue([]byte("value2"))
	storageUpdates := []EvmHookStorageUpdate{storageSlot1, storageSlot2}

	leh.SetStorageUpdates(storageUpdates)
	assert.Equal(t, storageUpdates, leh.GetStorageUpdates())

	// Test setting empty storage updates
	leh.SetStorageUpdates([]EvmHookStorageUpdate{})
	assert.Equal(t, []EvmHookStorageUpdate{}, leh.GetStorageUpdates())

	// Test setting nil storage updates
	leh.SetStorageUpdates(nil)
	assert.Equal(t, []EvmHookStorageUpdate(nil), leh.GetStorageUpdates())
	pb := leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)
}

func TestUnitLambdaEvmHookAddStorageUpdate(t *testing.T) {
	t.Parallel()

	leh := NewEvmHook()

	// Test adding to empty slice
	storageSlot1 := NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1"))
	leh.AddStorageUpdate(storageSlot1)
	assert.Len(t, leh.GetStorageUpdates(), 1)
	assert.Equal(t, storageSlot1, leh.GetStorageUpdates()[0])

	// Test adding to existing slice
	storageSlot2 := NewEvmHookStorageSlot().SetKey([]byte("key2")).SetValue([]byte("value2"))
	leh.AddStorageUpdate(storageSlot2)
	assert.Len(t, leh.GetStorageUpdates(), 2)
	assert.Equal(t, storageSlot1, leh.GetStorageUpdates()[0])
	assert.Equal(t, storageSlot2, leh.GetStorageUpdates()[1])
}

func TestUnitLambdaEvmHookMethodChaining(t *testing.T) {
	t.Parallel()

	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	storageSlot1 := NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1"))
	storageSlot2 := NewEvmHookStorageSlot().SetKey([]byte("key2")).SetValue([]byte("value2"))

	// Test method chaining
	leh := NewEvmHook().
		SetContractId(&contractID).
		SetStorageUpdates([]EvmHookStorageUpdate{storageSlot1}).
		AddStorageUpdate(storageSlot2)

	assert.Equal(t, &contractID, leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 2)
	assert.Equal(t, storageSlot1, leh.GetStorageUpdates()[0])
	assert.Equal(t, storageSlot2, leh.GetStorageUpdates()[1])
}

func TestUnitLambdaEvmHookToProtobuf(t *testing.T) {
	t.Parallel()

	// Test with all fields set
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	storageSlot1 := NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1"))
	storageSlot2 := NewEvmHookStorageSlot().SetKey([]byte("key2")).SetValue([]byte("value2"))

	leh := NewEvmHook().
		SetContractId(&contractID).
		SetStorageUpdates([]EvmHookStorageUpdate{storageSlot1, storageSlot2})

	pb := leh.toProtobuf()
	require.NotNil(t, pb)
	assert.NotNil(t, pb.Spec)
	assert.NotNil(t, pb.Spec.GetContractId())
	assert.Len(t, pb.StorageUpdates, 2)
}

func TestUnitLambdaEvmHookToProtobufWithNilContractId(t *testing.T) {
	t.Parallel()

	// Test with nil contract ID - this should cause a panic
	leh := NewEvmHook()
	storageSlot := NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1"))
	leh.SetStorageUpdates([]EvmHookStorageUpdate{storageSlot})

	pb := leh.toProtobuf()
	assert.Nil(t, pb.GetSpec().GetContractId())
	lambdaEvmHookFromProtobuf(pb)
}

func TestUnitLambdaEvmHookToProtobufWithEmptyStorageUpdates(t *testing.T) {
	t.Parallel()

	// Test with empty storage updates
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	leh := NewEvmHook().SetContractId(&contractID)

	pb := leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)
	require.NotNil(t, pb)
	assert.NotNil(t, pb.Spec)
	assert.NotNil(t, pb.Spec.GetContractId())
	assert.Len(t, pb.StorageUpdates, 0)
}

func TestUnitLambdaEvmHookFromProtobuf(t *testing.T) {
	t.Parallel()

	// Create a protobuf message
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	storageSlot1 := NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1"))
	storageSlot2 := NewEvmHookStorageSlot().SetKey([]byte("key2")).SetValue([]byte("value2"))

	pb := &services.LambdaEvmHook{
		Spec: &services.EvmHookSpec{
			BytecodeSource: &services.EvmHookSpec_ContractId{
				ContractId: contractID._ToProtobuf(),
			},
		},
		StorageUpdates: []*services.LambdaStorageUpdate{
			storageSlot1.toProtobuf(),
			storageSlot2.toProtobuf(),
		},
	}

	leh := lambdaEvmHookFromProtobuf(pb)
	assert.Equal(t, contractID, *leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 2)
}

func TestUnitLambdaEvmHookFromProtobufWithEmptyStorageUpdates(t *testing.T) {
	t.Parallel()

	// Create a protobuf message with empty storage updates
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	pb := &services.LambdaEvmHook{
		Spec: &services.EvmHookSpec{
			BytecodeSource: &services.EvmHookSpec_ContractId{
				ContractId: contractID._ToProtobuf(),
			},
		},
		StorageUpdates: []*services.LambdaStorageUpdate{},
	}

	leh := lambdaEvmHookFromProtobuf(pb)
	assert.Equal(t, contractID, *leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 0)
}

func TestUnitLambdaEvmHookRoundTrip(t *testing.T) {
	t.Parallel()

	// Test round trip conversion: struct -> protobuf -> struct
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	storageSlot1 := NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1"))
	storageSlot2 := NewEvmHookStorageSlot().SetKey([]byte("key2")).SetValue([]byte("value2"))

	original := NewEvmHook().
		SetContractId(&contractID).
		SetStorageUpdates([]EvmHookStorageUpdate{storageSlot1, storageSlot2})

	// Convert to protobuf and back
	pb := original.toProtobuf()
	converted := lambdaEvmHookFromProtobuf(pb)

	// Compare original and converted
	assert.Equal(t, original.GetContractId(), converted.GetContractId())
	assert.Len(t, converted.GetStorageUpdates(), 2)
}

func TestUnitLambdaEvmHookWithDifferentContractIds(t *testing.T) {
	t.Parallel()

	// Test with different contract IDs
	contractID1, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)

	contractID2, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	leh1 := NewEvmHook().SetContractId(&contractID1)
	leh2 := NewEvmHook().SetContractId(&contractID2)

	assert.Equal(t, &contractID1, leh1.GetContractId())
	assert.Equal(t, &contractID2, leh2.GetContractId())

	// Test protobuf conversion
	pb1 := leh1.toProtobuf()
	pb2 := leh2.toProtobuf()

	converted1 := lambdaEvmHookFromProtobuf(pb1)
	converted2 := lambdaEvmHookFromProtobuf(pb2)

	assert.Equal(t, &contractID1, converted1.GetContractId())
	assert.Equal(t, &contractID2, converted2.GetContractId())
}

func TestUnitLambdaEvmHookEdgeCases(t *testing.T) {
	t.Parallel()

	// Test with zero values
	leh := NewEvmHook()
	assert.Nil(t, leh.GetContractId())
	assert.Equal(t, []EvmHookStorageUpdate(nil), leh.GetStorageUpdates())
	pb := leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)

	// Test with empty storage slot
	emptyStorageSlot := NewEvmHookStorageSlot()
	leh.AddStorageUpdate(emptyStorageSlot)
	assert.Len(t, leh.GetStorageUpdates(), 1)
	assert.Equal(t, emptyStorageSlot, leh.GetStorageUpdates()[0])
	pb = leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)

	// Test with large storage updates
	largeStorageUpdates := make([]EvmHookStorageUpdate, 100)
	for i := 0; i < 100; i++ {
		storageSlot := NewEvmHookStorageSlot().
			SetKey([]byte("key" + string(rune(i)))).
			SetValue([]byte("value" + string(rune(i))))
		largeStorageUpdates[i] = storageSlot
	}

	leh.SetStorageUpdates(largeStorageUpdates)
	assert.Len(t, leh.GetStorageUpdates(), 100)

	// Test with storage slot with nil key/value
	storageSlot := NewEvmHookStorageSlot().SetKey(nil).SetValue(nil)
	leh.AddStorageUpdate(storageSlot)
	assert.Len(t, leh.GetStorageUpdates(), 101)
	assert.Equal(t, storageSlot, leh.GetStorageUpdates()[100])
	pb = leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)
}

func TestUnitLambdaEvmHookStorageSlotIntegration(t *testing.T) {
	t.Parallel()

	// Test integration with LambdaStorageSlot
	contractID, err := ContractIDFromString("0.0.789")
	require.NoError(t, err)

	// Create storage slots with different key/value combinations
	storageSlots := []*EvmHookStorageSlot{
		NewEvmHookStorageSlot().SetKey([]byte("")).SetValue([]byte("")), // Empty key/value
		NewEvmHookStorageSlot().SetKey([]byte("key1")).SetValue([]byte("value1")),
		NewEvmHookStorageSlot().SetKey([]byte("key2")).SetValue([]byte("value2")),
		NewEvmHookStorageSlot().SetKey([]byte("very_long_key_123456789")).SetValue([]byte("very_long_value_987654321")),
	}

	leh := NewEvmHook().SetContractId(&contractID)

	// Add each storage slot
	for _, slot := range storageSlots {
		leh.AddStorageUpdate(slot)
	}

	assert.Equal(t, &contractID, leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), len(storageSlots))

	// Test protobuf conversion
	pb := leh.toProtobuf()
	converted := lambdaEvmHookFromProtobuf(pb)

	assert.Equal(t, &contractID, converted.GetContractId())
	assert.Len(t, converted.GetStorageUpdates(), len(storageSlots))
}

func TestUnitLambdaEvmHookNilHandling(t *testing.T) {
	t.Parallel()

	// Test setting nil contract ID
	leh := NewEvmHook()
	leh.SetContractId(nil)
	assert.Nil(t, leh.GetContractId())

	// Test setting nil storage updates
	leh.SetStorageUpdates(nil)
	assert.Equal(t, []EvmHookStorageUpdate(nil), leh.GetStorageUpdates())
	pb := leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)

	// Test adding nil storage update (this should work but the slice will contain nil)
	leh.AddStorageUpdate(nil)
	assert.Len(t, leh.GetStorageUpdates(), 1)
	assert.Nil(t, leh.GetStorageUpdates()[0])

	pb = leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)
}

func TestUnitLambdaEvmHookWithMappingEntries(t *testing.T) {
	t.Parallel()

	// Test with LambdaMappingEntries as storage update
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)

	// Create mapping entries
	mappingEntry1 := NewEvmHookMappingEntryWithKey([]byte("key1"), []byte("value1"))
	mappingEntry2 := NewEvmHookMappingEntryWithPreImage([]byte("preimage2"), []byte("value2"))
	mappingEntry3 := NewEvmHookMappingEntryWithKey([]byte("key3"), []byte("value3"))

	mappingEntries := NewLambdaMappingEntries().
		SetMappingSlot([]byte("mapping_slot_123")).
		AddMappingEntry(*mappingEntry1).
		AddMappingEntry(*mappingEntry2).
		AddMappingEntry(*mappingEntry3)

	leh := NewEvmHook().
		SetContractId(&contractID).
		AddStorageUpdate(mappingEntries)

	assert.Equal(t, &contractID, leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 1)
	assert.Equal(t, mappingEntries, leh.GetStorageUpdates()[0])

	// Test protobuf conversion
	pb := leh.toProtobuf()
	require.NotNil(t, pb)
	assert.NotNil(t, pb.Spec)
	assert.Len(t, pb.StorageUpdates, 1)

	converted := lambdaEvmHookFromProtobuf(pb)
	assert.Equal(t, &contractID, converted.GetContractId())
	assert.Len(t, converted.GetStorageUpdates(), 1)
}

func TestUnitLambdaEvmHookWithMixedStorageUpdates(t *testing.T) {
	t.Parallel()

	// Test with both LambdaStorageSlot and LambdaMappingEntries
	contractID, err := ContractIDFromString("0.0.789")
	require.NoError(t, err)

	// Create storage slot
	storageSlot := NewEvmHookStorageSlot().SetKey([]byte("slot_key")).SetValue([]byte("slot_value"))

	// Create mapping entries
	mappingEntry := NewEvmHookMappingEntryWithKey([]byte("mapping_key"), []byte("mapping_value"))
	mappingEntries := NewLambdaMappingEntries().
		SetMappingSlot([]byte("mapping_slot")).
		AddMappingEntry(*mappingEntry)

	leh := NewEvmHook().
		SetContractId(&contractID).
		AddStorageUpdate(storageSlot).
		AddStorageUpdate(mappingEntries)

	assert.Equal(t, &contractID, leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 2)
	assert.Equal(t, storageSlot, leh.GetStorageUpdates()[0])
	assert.Equal(t, mappingEntries, leh.GetStorageUpdates()[1])

	// Test protobuf conversion
	pb := leh.toProtobuf()
	require.NotNil(t, pb)
	assert.Len(t, pb.StorageUpdates, 2)

	converted := lambdaEvmHookFromProtobuf(pb)
	assert.Equal(t, &contractID, converted.GetContractId())
	assert.Len(t, converted.GetStorageUpdates(), 2)
}

func TestUnitLambdaEvmHookWithEmptyMappingEntries(t *testing.T) {
	t.Parallel()

	// Test with empty mapping entries
	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)

	emptyMappingEntries := NewLambdaMappingEntries().SetMappingSlot([]byte("empty_slot"))

	leh := NewEvmHook().
		SetContractId(&contractID).
		AddStorageUpdate(emptyMappingEntries)

	assert.Equal(t, &contractID, leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 1)
	assert.Equal(t, emptyMappingEntries, leh.GetStorageUpdates()[0])

	// Test protobuf conversion
	pb := leh.toProtobuf()
	require.NotNil(t, pb)
	assert.Len(t, pb.StorageUpdates, 1)

	converted := lambdaEvmHookFromProtobuf(pb)
	assert.Equal(t, &contractID, converted.GetContractId())
	assert.Len(t, converted.GetStorageUpdates(), 1)
}

func TestUnitLambdaEvmHookWithComplexMappingEntries(t *testing.T) {
	t.Parallel()

	// Test with complex mapping entries containing multiple entries
	contractID, err := ContractIDFromString("0.0.999")
	require.NoError(t, err)

	// Create multiple mapping entries with different key types
	mappingEntries := NewLambdaMappingEntries().SetMappingSlot([]byte("complex_mapping_slot"))

	// Add entries with keys
	for i := 0; i < 5; i++ {
		entry := NewEvmHookMappingEntryWithKey([]byte("key"+string(rune(i))), []byte("value"+string(rune(i))))
		mappingEntries.AddMappingEntry(*entry)
	}

	// Add entries with preimages
	for i := 0; i < 3; i++ {
		entry := NewEvmHookMappingEntryWithPreImage([]byte("preimage"+string(rune(i))), []byte("prevalue"+string(rune(i))))
		mappingEntries.AddMappingEntry(*entry)
	}

	leh := NewEvmHook().
		SetContractId(&contractID).
		AddStorageUpdate(mappingEntries)

	assert.Equal(t, &contractID, leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 1)
	assert.Len(t, mappingEntries.GetMappingEntries(), 8) // 5 + 3 entries

	// Test protobuf conversion
	pb := leh.toProtobuf()
	require.NotNil(t, pb)
	assert.Len(t, pb.StorageUpdates, 1)

	converted := lambdaEvmHookFromProtobuf(pb)
	assert.Equal(t, &contractID, converted.GetContractId())
	assert.Len(t, converted.GetStorageUpdates(), 1)
}

func TestUnitLambdaEvmHookMappingEntriesRoundTrip(t *testing.T) {
	t.Parallel()

	// Test round trip conversion with mapping entries
	contractID, err := ContractIDFromString("0.0.555")
	require.NoError(t, err)

	// Create mapping entries
	mappingEntry1 := NewEvmHookMappingEntryWithKey([]byte("roundtrip_key1"), []byte("roundtrip_value1"))
	mappingEntry2 := NewEvmHookMappingEntryWithPreImage([]byte("roundtrip_preimage"), []byte("roundtrip_value2"))

	mappingEntries := NewLambdaMappingEntries().
		SetMappingSlot([]byte("roundtrip_slot")).
		AddMappingEntry(*mappingEntry1).
		AddMappingEntry(*mappingEntry2)

	original := NewEvmHook().
		SetContractId(&contractID).
		AddStorageUpdate(mappingEntries)

	// Convert to protobuf and back
	pb := original.toProtobuf()
	converted := lambdaEvmHookFromProtobuf(pb)

	// Compare original and converted
	assert.Equal(t, original.GetContractId(), converted.GetContractId())
	assert.Len(t, converted.GetStorageUpdates(), 1)

	// Verify the mapping entries are preserved
	convertedMappingEntries := converted.GetStorageUpdates()[0].(EvmHookMappingEntries)
	assert.Equal(t, mappingEntries.GetMappingSlot(), convertedMappingEntries.GetMappingSlot())
	assert.Len(t, convertedMappingEntries.GetMappingEntries(), 2)
}

func TestUnitLambdaEvmHookMappingEntriesEdgeCases(t *testing.T) {
	t.Parallel()

	contractID, err := ContractIDFromString("0.0.555")
	require.NoError(t, err)

	// Test with mapping entries with nil key/value
	mappingEntry := NewEvmHookMappingEntryWithKey(nil, nil)
	mappingEntries := NewLambdaMappingEntries().
		SetMappingSlot(nil).
		AddMappingEntry(*mappingEntry)

	leh := NewEvmHook().
		SetContractId(&contractID).
		AddStorageUpdate(mappingEntries)

	assert.Equal(t, &contractID, leh.GetContractId())
	assert.Len(t, leh.GetStorageUpdates(), 1)
	assert.Equal(t, mappingEntries, leh.GetStorageUpdates()[0])
	pb := leh.toProtobuf()
	lambdaEvmHookFromProtobuf(pb)
}
