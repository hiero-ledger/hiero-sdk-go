//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitHookEntityIdWithAccountId(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Account: 123}
	entityId := NewHookEntityIdWithAccountId(accountId)

	require.NotNil(t, entityId)
	require.Equal(t, accountId, entityId.GetAccountId())
	require.Equal(t, ContractID{}, entityId.GetContractId())
}

func TestUnitHookEntityIdWithContractId(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Contract: 456}
	entityId := NewHookEntityIdWithContractId(contractId)

	require.NotNil(t, entityId)
	require.Equal(t, contractId, entityId.GetContractId())
	require.Equal(t, AccountID{}, entityId.GetAccountId())
}

func TestUnitHookEntityIdGetAccountIdWhenNil(t *testing.T) {
	t.Parallel()

	entityId := HookEntityId{}

	result := entityId.GetAccountId()
	require.Equal(t, AccountID{}, result)
}

func TestUnitHookEntityIdGetContractIdWhenNil(t *testing.T) {
	t.Parallel()

	entityId := HookEntityId{}

	result := entityId.GetContractId()
	require.Equal(t, ContractID{}, result)
}

func TestUnitHookEntityIdToProtobufWithAccountId(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Shard: 1, Realm: 2, Account: 3}
	entityId := NewHookEntityIdWithAccountId(accountId)

	proto := entityId.toProtobuf()
	require.NotNil(t, proto)
	require.NotNil(t, proto.GetAccountId())
	require.Equal(t, int64(1), proto.GetAccountId().GetShardNum())
	require.Equal(t, int64(2), proto.GetAccountId().GetRealmNum())
	require.Equal(t, int64(3), proto.GetAccountId().GetAccountNum())
}

func TestUnitHookEntityIdToProtobufWithContractId(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Shard: 4, Realm: 5, Contract: 6}
	entityId := NewHookEntityIdWithContractId(contractId)

	proto := entityId.toProtobuf()
	require.NotNil(t, proto)
	require.NotNil(t, proto.GetContractId())
	require.Equal(t, int64(4), proto.GetContractId().GetShardNum())
	require.Equal(t, int64(5), proto.GetContractId().GetRealmNum())
	require.Equal(t, int64(6), proto.GetContractId().GetContractNum())
}

func TestUnitHookEntityIdToProtobufWhenEmpty(t *testing.T) {
	t.Parallel()

	entityId := HookEntityId{}

	proto := entityId.toProtobuf()
	require.Nil(t, proto)
}

func TestUnitHookEntityIdFromProtobufWithAccountId(t *testing.T) {
	t.Parallel()

	proto := &services.HookEntityId{
		EntityId: &services.HookEntityId_AccountId{
			AccountId: &services.AccountID{
				ShardNum: 1,
				RealmNum: 2,
				Account:  &services.AccountID_AccountNum{AccountNum: 3},
			},
		},
	}

	entityId := hookEntityIdFromProtobuf(proto)
	require.Equal(t, uint64(1), entityId.GetAccountId().Shard)
	require.Equal(t, uint64(2), entityId.GetAccountId().Realm)
	require.Equal(t, uint64(3), entityId.GetAccountId().Account)
}

func TestUnitHookEntityIdFromProtobufWithContractId(t *testing.T) {
	t.Parallel()

	proto := &services.HookEntityId{
		EntityId: &services.HookEntityId_ContractId{
			ContractId: &services.ContractID{
				ShardNum: 4,
				RealmNum: 5,
				Contract: &services.ContractID_ContractNum{ContractNum: 6},
			},
		},
	}

	entityId := hookEntityIdFromProtobuf(proto)
	require.Equal(t, uint64(4), entityId.GetContractId().Shard)
	require.Equal(t, uint64(5), entityId.GetContractId().Realm)
	require.Equal(t, uint64(6), entityId.GetContractId().Contract)
}

func TestUnitHookEntityIdRoundTripWithAccountId(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Shard: 10, Realm: 20, Account: 30}
	original := NewHookEntityIdWithAccountId(accountId)

	// Convert to protobuf
	proto := original.toProtobuf()
	require.NotNil(t, proto)

	// Convert back from protobuf
	reconstructed := hookEntityIdFromProtobuf(proto)

	// Verify round-trip
	require.Equal(t, original.GetAccountId(), reconstructed.GetAccountId())
}

func TestUnitHookEntityIdRoundTripWithContractId(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Shard: 11, Realm: 22, Contract: 33}
	original := NewHookEntityIdWithContractId(contractId)

	// Convert to protobuf
	proto := original.toProtobuf()
	require.NotNil(t, proto)

	// Convert back from protobuf
	reconstructed := hookEntityIdFromProtobuf(proto)

	// Verify round-trip
	require.Equal(t, original.GetContractId(), reconstructed.GetContractId())
}

func TestUnitHookIdNew(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Account: 100}
	entityId := NewHookEntityIdWithAccountId(accountId)
	hookIdNum := int64(200)

	hookId := NewHookId(*entityId, hookIdNum)

	require.NotNil(t, hookId)
	require.Equal(t, hookIdNum, hookId.GetHookId())
	require.Equal(t, accountId, hookId.GetEntityId().GetAccountId())
}

func TestUnitHookIdWithAccountEntity(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Shard: 1, Realm: 2, Account: 300}
	entityId := NewHookEntityIdWithAccountId(accountId)
	hookIdNum := int64(400)

	hookId := NewHookId(*entityId, hookIdNum)

	require.NotNil(t, hookId)
	require.Equal(t, hookIdNum, hookId.GetHookId())
	require.Equal(t, accountId.Shard, hookId.GetEntityId().GetAccountId().Shard)
	require.Equal(t, accountId.Realm, hookId.GetEntityId().GetAccountId().Realm)
	require.Equal(t, accountId.Account, hookId.GetEntityId().GetAccountId().Account)
}

func TestUnitHookIdWithContractEntity(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Shard: 3, Realm: 4, Contract: 500}
	entityId := NewHookEntityIdWithContractId(contractId)
	hookIdNum := int64(600)

	hookId := NewHookId(*entityId, hookIdNum)

	require.NotNil(t, hookId)
	require.Equal(t, hookIdNum, hookId.GetHookId())
	require.Equal(t, contractId.Shard, hookId.GetEntityId().GetContractId().Shard)
	require.Equal(t, contractId.Realm, hookId.GetEntityId().GetContractId().Realm)
	require.Equal(t, contractId.Contract, hookId.GetEntityId().GetContractId().Contract)
}

func TestUnitHookIdToProtobufWithAccountEntity(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Shard: 5, Realm: 6, Account: 700}
	entityId := NewHookEntityIdWithAccountId(accountId)
	hookIdNum := int64(800)

	hookId := NewHookId(*entityId, hookIdNum)

	proto := hookId.toProtobuf()
	require.NotNil(t, proto)
	require.Equal(t, hookIdNum, proto.GetHookId())
	require.NotNil(t, proto.GetEntityId())
	require.NotNil(t, proto.GetEntityId().GetAccountId())
	require.Equal(t, int64(5), proto.GetEntityId().GetAccountId().GetShardNum())
	require.Equal(t, int64(6), proto.GetEntityId().GetAccountId().GetRealmNum())
	require.Equal(t, int64(700), proto.GetEntityId().GetAccountId().GetAccountNum())
}

func TestUnitHookIdToProtobufWithContractEntity(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Shard: 7, Realm: 8, Contract: 900}
	entityId := NewHookEntityIdWithContractId(contractId)
	hookIdNum := int64(1000)

	hookId := NewHookId(*entityId, hookIdNum)

	proto := hookId.toProtobuf()
	require.NotNil(t, proto)
	require.Equal(t, hookIdNum, proto.GetHookId())
	require.NotNil(t, proto.GetEntityId())
	require.NotNil(t, proto.GetEntityId().GetContractId())
	require.Equal(t, int64(7), proto.GetEntityId().GetContractId().GetShardNum())
	require.Equal(t, int64(8), proto.GetEntityId().GetContractId().GetRealmNum())
	require.Equal(t, int64(900), proto.GetEntityId().GetContractId().GetContractNum())
}

func TestUnitHookIdFromProtobufWithAccountEntity(t *testing.T) {
	t.Parallel()

	proto := &services.HookId{
		EntityId: &services.HookEntityId{
			EntityId: &services.HookEntityId_AccountId{
				AccountId: &services.AccountID{
					ShardNum: 9,
					RealmNum: 10,
					Account:  &services.AccountID_AccountNum{AccountNum: 1100},
				},
			},
		},
		HookId: 1200,
	}

	hookId := hookIdFromProtobuf(proto)
	require.Equal(t, int64(1200), hookId.GetHookId())
	require.Equal(t, uint64(9), hookId.GetEntityId().GetAccountId().Shard)
	require.Equal(t, uint64(10), hookId.GetEntityId().GetAccountId().Realm)
	require.Equal(t, uint64(1100), hookId.GetEntityId().GetAccountId().Account)
}

func TestUnitHookIdFromProtobufWithContractEntity(t *testing.T) {
	t.Parallel()

	proto := &services.HookId{
		EntityId: &services.HookEntityId{
			EntityId: &services.HookEntityId_ContractId{
				ContractId: &services.ContractID{
					ShardNum: 11,
					RealmNum: 12,
					Contract: &services.ContractID_ContractNum{ContractNum: 1300},
				},
			},
		},
		HookId: 1400,
	}

	hookId := hookIdFromProtobuf(proto)
	require.Equal(t, int64(1400), hookId.GetHookId())
	require.Equal(t, uint64(11), hookId.GetEntityId().GetContractId().Shard)
	require.Equal(t, uint64(12), hookId.GetEntityId().GetContractId().Realm)
	require.Equal(t, uint64(1300), hookId.GetEntityId().GetContractId().Contract)
}

func TestUnitHookIdRoundTripWithAccountEntity(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Shard: 13, Realm: 14, Account: 1500}
	entityId := NewHookEntityIdWithAccountId(accountId)
	hookIdNum := int64(1600)

	original := NewHookId(*entityId, hookIdNum)

	// Convert to protobuf
	proto := original.toProtobuf()
	require.NotNil(t, proto)

	// Convert back from protobuf
	reconstructed := hookIdFromProtobuf(proto)

	// Verify round-trip
	require.Equal(t, original.GetHookId(), reconstructed.GetHookId())
	require.Equal(t, original.GetEntityId().GetAccountId().Shard, reconstructed.GetEntityId().GetAccountId().Shard)
	require.Equal(t, original.GetEntityId().GetAccountId().Realm, reconstructed.GetEntityId().GetAccountId().Realm)
	require.Equal(t, original.GetEntityId().GetAccountId().Account, reconstructed.GetEntityId().GetAccountId().Account)
}

func TestUnitHookIdRoundTripWithContractEntity(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Shard: 15, Realm: 16, Contract: 1700}
	entityId := NewHookEntityIdWithContractId(contractId)
	hookIdNum := int64(1800)

	original := NewHookId(*entityId, hookIdNum)

	// Convert to protobuf
	proto := original.toProtobuf()
	require.NotNil(t, proto)

	// Convert back from protobuf
	reconstructed := hookIdFromProtobuf(proto)

	// Verify round-trip
	require.Equal(t, original.GetHookId(), reconstructed.GetHookId())
	require.Equal(t, original.GetEntityId().GetContractId().Shard, reconstructed.GetEntityId().GetContractId().Shard)
	require.Equal(t, original.GetEntityId().GetContractId().Realm, reconstructed.GetEntityId().GetContractId().Realm)
	require.Equal(t, original.GetEntityId().GetContractId().Contract, reconstructed.GetEntityId().GetContractId().Contract)
}

func TestUnitHookIdZeroHookId(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Account: 100}
	entityId := NewHookEntityIdWithAccountId(accountId)
	hookIdNum := int64(0)

	hookId := NewHookId(*entityId, hookIdNum)

	require.NotNil(t, hookId)
	require.Equal(t, int64(0), hookId.GetHookId())
}

func TestUnitHookIdNegativeHookId(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Account: 100}
	entityId := NewHookEntityIdWithAccountId(accountId)
	hookIdNum := int64(-1)

	hookId := NewHookId(*entityId, hookIdNum)

	require.NotNil(t, hookId)
	require.Equal(t, int64(-1), hookId.GetHookId())

	// Test round-trip with negative value
	proto := hookId.toProtobuf()
	reconstructed := hookIdFromProtobuf(proto)
	require.Equal(t, int64(-1), reconstructed.GetHookId())
}

func TestUnitHookIdValidateChecksumWithAccountId(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	accountId, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithAccountId(accountId)
	hookId := NewHookId(*entityId, 100)

	err = hookId.validateChecksum(client)
	require.NoError(t, err)
}

func TestUnitHookIdValidateChecksumWithContractId(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetAutoValidateChecksums(true)

	contractId, err := ContractIDFromString("0.0.456-rmkykd")
	require.NoError(t, err)

	entityId := NewHookEntityIdWithContractId(contractId)
	hookId := NewHookId(*entityId, 200)

	err = hookId.validateChecksum(client)
	require.Error(t, err)
}

func TestUnitHookIdMultipleRoundTrips(t *testing.T) {
	t.Parallel()

	accountId := AccountID{Shard: 66, Realm: 55, Account: 44}
	entityId := NewHookEntityIdWithAccountId(accountId)
	hookIdNum := int64(33)

	original := NewHookId(*entityId, hookIdNum)

	// First round trip
	proto1 := original.toProtobuf()
	reconstructed1 := hookIdFromProtobuf(proto1)

	// Second round trip
	proto2 := reconstructed1.toProtobuf()
	reconstructed2 := hookIdFromProtobuf(proto2)

	// Third round trip
	proto3 := reconstructed2.toProtobuf()
	reconstructed3 := hookIdFromProtobuf(proto3)

	// All should be equal
	require.Equal(t, original.GetHookId(), reconstructed1.GetHookId())
	require.Equal(t, original.GetHookId(), reconstructed2.GetHookId())
	require.Equal(t, original.GetHookId(), reconstructed3.GetHookId())
	require.Equal(t, original.GetEntityId().GetAccountId(), reconstructed1.GetEntityId().GetAccountId())
	require.Equal(t, original.GetEntityId().GetAccountId(), reconstructed2.GetEntityId().GetAccountId())
	require.Equal(t, original.GetEntityId().GetAccountId(), reconstructed3.GetEntityId().GetAccountId())
}
