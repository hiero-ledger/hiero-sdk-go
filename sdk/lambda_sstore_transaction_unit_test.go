//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitLambdaSStoreTransactionNew(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	assert.NotNil(t, tx)
	assert.NotNil(t, tx.Transaction)
	assert.Equal(t, HookId{}, tx.hookId)
	assert.Empty(t, tx.storageUpdates)
}

func TestUnitLambdaSStoreTransactionSetHookId(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	tx.SetHookId(*hookId)

	assert.Equal(t, *hookId, tx.GetHookId())
}

func TestUnitLambdaSStoreTransactionSetHookIdChaining(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId1 := NewHookId(*entityId, 456)
	hookId2 := NewHookId(*entityId, 789)

	result := tx.SetHookId(*hookId1).SetHookId(*hookId2)

	assert.Equal(t, tx, result)
	assert.Equal(t, *hookId2, tx.GetHookId())
}

func TestUnitLambdaSStoreTransactionAddStorageUpdate(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	storageSlot := NewLambdaStorageSlot().
		SetKey([]byte{0x01, 0x02}).
		SetValue([]byte{0x03, 0x04})

	tx.AddStorageUpdate(storageSlot)

	updates := tx.GetStorageUpdates()
	require.Len(t, updates, 1)
	assert.Equal(t, storageSlot, updates[0])
}

func TestUnitLambdaSStoreTransactionAddMultipleStorageUpdates(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	storageSlot1 := NewLambdaStorageSlot().
		SetKey([]byte{0x01, 0x02}).
		SetValue([]byte{0x03, 0x04})

	storageSlot2 := NewLambdaStorageSlot().
		SetKey([]byte{0x05, 0x06}).
		SetValue([]byte{0x07, 0x08})

	tx.AddStorageUpdate(storageSlot1).AddStorageUpdate(storageSlot2)

	updates := tx.GetStorageUpdates()
	require.Len(t, updates, 2)
	assert.Equal(t, storageSlot1, updates[0])
	assert.Equal(t, storageSlot2, updates[1])
}

func TestUnitLambdaSStoreTransactionSetStorageUpdates(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	storageSlot1 := NewLambdaStorageSlot().
		SetKey([]byte{0x01, 0x02}).
		SetValue([]byte{0x03, 0x04})

	storageSlot2 := NewLambdaStorageSlot().
		SetKey([]byte{0x05, 0x06}).
		SetValue([]byte{0x07, 0x08})

	updates := []LambdaStorageUpdate{storageSlot1, storageSlot2}
	tx.SetStorageUpdates(updates)

	result := tx.GetStorageUpdates()
	require.Len(t, result, 2)
	assert.Equal(t, storageSlot1, result[0])
	assert.Equal(t, storageSlot2, result[1])
}

func TestUnitLambdaSStoreTransactionSetStorageUpdatesEmpty(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	storageSlot := NewLambdaStorageSlot().
		SetKey([]byte{0x01, 0x02}).
		SetValue([]byte{0x03, 0x04})
	tx.AddStorageUpdate(storageSlot)
	tx.SetStorageUpdates([]LambdaStorageUpdate{})

	result := tx.GetStorageUpdates()
	assert.Empty(t, result)
}

func TestUnitLambdaSStoreTransactionAddMappingEntries(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	mappingEntry := NewLambdaMappingEntryWithKey([]byte{0x01}, []byte{0x02})
	mappingEntries := NewLambdaMappingEntries().
		SetMappingSlot([]byte{0x03}).
		AddMappingEntry(*mappingEntry)

	tx.AddStorageUpdate(mappingEntries)

	updates := tx.GetStorageUpdates()
	require.Len(t, updates, 1)
	assert.Equal(t, mappingEntries, updates[0])
}

func TestUnitLambdaSStoreTransactionGetName(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	assert.Equal(t, "LambdaSStoreTransaction", tx.getName())
}

func TestUnitLambdaSStoreTransactionValidateNetworkOnIDs(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	contractID, err := ContractIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	tx := NewLambdaSStoreTransaction().SetHookId(*hookId)

	err = tx.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitLambdaSStoreTransactionValidateNetworkOnIDsWrongChecksum(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	contractID, err := ContractIDFromString("0.0.123-wrong")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	tx := NewLambdaSStoreTransaction().SetHookId(*hookId)

	err = tx.validateNetworkOnIDs(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network mismatch or wrong checksum given")
}

func TestUnitLambdaSStoreTransactionValidateNetworkOnIDsNilClient(t *testing.T) {
	t.Parallel()

	contractID, err := ContractIDFromString("0.0.123-quros")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	tx := NewLambdaSStoreTransaction().SetHookId(*hookId)

	err = tx.validateNetworkOnIDs(nil)
	require.NoError(t, err)
}

func TestUnitLambdaSStoreTransactionValidateNetworkOnIDsDisabledChecksumValidation(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(false)

	contractID, err := ContractIDFromString("0.0.123-wrong")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	tx := NewLambdaSStoreTransaction().SetHookId(*hookId)

	err = tx.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitLambdaSStoreTransactionBuildProtoBody(t *testing.T) {
	t.Parallel()

	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	storageSlot := NewLambdaStorageSlot().
		SetKey([]byte{0x01, 0x02}).
		SetValue([]byte{0x03, 0x04})

	tx := NewLambdaSStoreTransaction().
		SetHookId(*hookId).
		AddStorageUpdate(storageSlot)

	protoBody := tx.buildProtoBody()

	assert.NotNil(t, protoBody)
	assert.NotNil(t, protoBody.HookId)
	assert.Len(t, protoBody.StorageUpdates, 1)

	assert.Equal(t, int64(456), protoBody.HookId.HookId)
	assert.Equal(t, int64(123), protoBody.HookId.EntityId.GetContractId().GetContractNum())
	require.NotNil(t, protoBody.StorageUpdates[0].GetStorageSlot())
	assert.Equal(t, []byte{0x01, 0x02}, protoBody.StorageUpdates[0].GetStorageSlot().Key)
	assert.Equal(t, []byte{0x03, 0x04}, protoBody.StorageUpdates[0].GetStorageSlot().Value)
}

func TestUnitLambdaSStoreTransactionBuildProtoBodyWithMappingEntries(t *testing.T) {
	t.Parallel()

	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	mappingEntry := NewLambdaMappingEntryWithKey([]byte{0x01}, []byte{0x02})
	mappingEntries := NewLambdaMappingEntries().
		SetMappingSlot([]byte{0x03}).
		AddMappingEntry(*mappingEntry)

	tx := NewLambdaSStoreTransaction().
		SetHookId(*hookId).
		AddStorageUpdate(mappingEntries)

	protoBody := tx.buildProtoBody()

	assert.NotNil(t, protoBody)
	assert.NotNil(t, protoBody.HookId)
	assert.Len(t, protoBody.StorageUpdates, 1)

	assert.Equal(t, int64(456), protoBody.HookId.HookId)
	assert.Equal(t, int64(123), protoBody.HookId.EntityId.GetContractId().GetContractNum())
	require.NotNil(t, protoBody.StorageUpdates[0].GetMappingEntries())
	assert.Equal(t, []byte{0x03}, protoBody.StorageUpdates[0].GetMappingEntries().MappingSlot)
	assert.Len(t, protoBody.StorageUpdates[0].GetMappingEntries().Entries, 1)
	assert.Equal(t, []byte{0x01}, protoBody.StorageUpdates[0].GetMappingEntries().Entries[0].GetKey())
	assert.Equal(t, []byte{0x02}, protoBody.StorageUpdates[0].GetMappingEntries().Entries[0].Value)
}

func TestUnitLambdaSStoreTransactionBuild(t *testing.T) {
	t.Parallel()

	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractID)
	hookId := NewHookId(*entityId, 456)

	storageSlot := NewLambdaStorageSlot().
		SetKey([]byte{0x01, 0x02}).
		SetValue([]byte{0x03, 0x04})

	tx := NewLambdaSStoreTransaction().
		SetHookId(*hookId).
		AddStorageUpdate(storageSlot).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionMemo("test memo").
		SetTransactionValidDuration(24 * time.Hour)

	protoBody := tx.build()

	assert.NotNil(t, protoBody)
	assert.Equal(t, uint64(100000000), protoBody.TransactionFee)
	assert.Equal(t, "test memo", protoBody.Memo)
	assert.NotNil(t, protoBody.TransactionValidDuration)
	assert.NotNil(t, protoBody.TransactionID)

	require.NotNil(t, protoBody.GetLambdaSstore())
	lambdaSstore := protoBody.GetLambdaSstore()
	assert.NotNil(t, lambdaSstore.HookId)
	assert.Len(t, lambdaSstore.StorageUpdates, 1)
}

func TestUnitLambdaSStoreTransactionBuildScheduled(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	scheduled, err := tx.buildScheduled()

	assert.Nil(t, scheduled)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot schedule `LambdaSStoreTransaction`")
}

func TestUnitLambdaSStoreTransactionConstructScheduleProtobuf(t *testing.T) {
	t.Parallel()

	tx := NewLambdaSStoreTransaction()

	scheduled, err := tx.constructScheduleProtobuf()

	assert.Nil(t, scheduled)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot schedule `LambdaSStoreTransaction`")
}

func TestUnitLambdaSStoreTransactionFromProtobuf(t *testing.T) {
	t.Parallel()

	contractID := &services.ContractID{
		Contract: &services.ContractID_ContractNum{ContractNum: 123},
	}
	hookId := &services.HookId{
		HookId: 456,
		EntityId: &services.HookEntityId{
			EntityId: &services.HookEntityId_ContractId{ContractId: contractID},
		},
	}

	storageSlot := &services.LambdaStorageSlot{
		Key:   []byte{0x01, 0x02},
		Value: []byte{0x03, 0x04},
	}

	storageUpdate := &services.LambdaStorageUpdate{
		Update: &services.LambdaStorageUpdate_StorageSlot{
			StorageSlot: storageSlot,
		},
	}

	lambdaSstore := &services.LambdaSStoreTransactionBody{
		HookId:         hookId,
		StorageUpdates: []*services.LambdaStorageUpdate{storageUpdate},
	}

	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_LambdaSstore{
			LambdaSstore: lambdaSstore,
		},
	}

	baseTx := Transaction[*LambdaSStoreTransaction]{}
	result := lambdaSStoreTransactionFromProtobuf(baseTx, pbBody)

	assert.Equal(t, int64(456), result.hookId.GetHookId())
	assert.Equal(t, uint64(123), result.hookId.GetEntityId().GetContractId().Contract)
	require.Len(t, result.storageUpdates, 1)
	storageSlotResult, ok := result.storageUpdates[0].(LambdaStorageSlot)
	require.True(t, ok)
	assert.Equal(t, []byte{0x01, 0x02}, storageSlotResult.GetKey())
	assert.Equal(t, []byte{0x03, 0x04}, storageSlotResult.GetValue())
}

func TestUnitLambdaSStoreTransactionFromProtobufWithMappingEntries(t *testing.T) {
	t.Parallel()

	contractID := &services.ContractID{
		Contract: &services.ContractID_ContractNum{ContractNum: 123},
	}
	hookId := &services.HookId{
		HookId: 456,
		EntityId: &services.HookEntityId{
			EntityId: &services.HookEntityId_ContractId{ContractId: contractID},
		},
	}

	mappingEntry := &services.LambdaMappingEntry{
		EntryKey: &services.LambdaMappingEntry_Key{Key: []byte{0x01}},
		Value:    []byte{0x02},
	}

	mappingEntries := &services.LambdaMappingEntries{
		MappingSlot: []byte{0x03},
		Entries:     []*services.LambdaMappingEntry{mappingEntry},
	}

	storageUpdate := &services.LambdaStorageUpdate{
		Update: &services.LambdaStorageUpdate_MappingEntries{
			MappingEntries: mappingEntries,
		},
	}

	lambdaSstore := &services.LambdaSStoreTransactionBody{
		HookId:         hookId,
		StorageUpdates: []*services.LambdaStorageUpdate{storageUpdate},
	}

	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_LambdaSstore{
			LambdaSstore: lambdaSstore,
		},
	}

	baseTx := Transaction[*LambdaSStoreTransaction]{}
	result := lambdaSStoreTransactionFromProtobuf(baseTx, pbBody)

	assert.Equal(t, int64(456), result.hookId.GetHookId())
	assert.Equal(t, uint64(123), result.hookId.GetEntityId().GetContractId().Contract)
	require.Len(t, result.storageUpdates, 1)
	mappingEntriesResult, ok := result.storageUpdates[0].(LambdaMappingEntries)
	require.True(t, ok)
	assert.Equal(t, []byte{0x03}, mappingEntriesResult.GetMappingSlot())
	assert.Len(t, mappingEntriesResult.GetMappingEntries(), 1)
	assert.Equal(t, []byte{0x01}, mappingEntriesResult.GetMappingEntries()[0].GetKey())
	assert.Equal(t, []byte{0x02}, mappingEntriesResult.GetMappingEntries()[0].GetValue())
}

func TestUnitLambdaSStoreTransactionFromProtobufEmpty(t *testing.T) {
	t.Parallel()

	lambdaSstore := &services.LambdaSStoreTransactionBody{
		HookId:         &services.HookId{},
		StorageUpdates: []*services.LambdaStorageUpdate{},
	}

	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_LambdaSstore{
			LambdaSstore: lambdaSstore,
		},
	}

	baseTx := Transaction[*LambdaSStoreTransaction]{}
	result := lambdaSStoreTransactionFromProtobuf(baseTx, pbBody)

	assert.Equal(t, HookId{}, result.hookId)
	assert.Empty(t, result.storageUpdates)
}
