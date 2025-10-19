package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

type HookId struct {
	entityId HookEntityId
	hookId   int64
}

// NewHookId creates a new HookId
func NewHookId(entityId HookEntityId, hookId int64) *HookId {
	return &HookId{
		entityId: entityId,
		hookId:   hookId,
	}
}

// GetEntityId returns the entity ID
func (id HookId) GetEntityId() HookEntityId {
	return id.entityId
}

// GetHookId returns the hook ID
func (id HookId) GetHookId() int64 {
	return id.hookId
}

func (id HookId) validateChecksum(client *Client) error {
	return id.entityId.validateChecksum(client)
}

func (id HookId) toProtobuf() *services.HookId {
	return &services.HookId{
		EntityId: id.entityId.toProtobuf(),
		HookId:   id.hookId,
	}
}

func hookIdFromProtobuf(pb *services.HookId) HookId {
	return HookId{
		entityId: hookEntityIdFromProtobuf(pb.GetEntityId()),
		hookId:   pb.GetHookId(),
	}
}

// hook entity id
type HookEntityId struct {
	accountId  *AccountID
	contractId *ContractID
}

// NewHookEntityId creates a new HookEntityId
func NewHookEntityIdWithAccountId(accountId AccountID) *HookEntityId {
	return &HookEntityId{
		accountId: &accountId,
	}
}
func NewHookEntityIdWithContractId(contractId ContractID) *HookEntityId {
	return &HookEntityId{
		contractId: &contractId,
	}
}

// GetAccountId returns the account ID
func (id HookEntityId) GetAccountId() AccountID {
	if id.accountId == nil {
		return AccountID{}
	}
	return *id.accountId
}

func (id HookEntityId) GetContractId() ContractID {
	if id.contractId == nil {
		return ContractID{}
	}
	return *id.contractId
}

func (id HookEntityId) validateChecksum(client *Client) error {
	if id.accountId != nil {
		err := id.accountId.ValidateChecksum(client)
		if err != nil {
			return err
		}
	}
	if id.contractId != nil {
		err := id.contractId.ValidateChecksum(client)
		if err != nil {
			return err
		}
	}
	return nil
}

func (id HookEntityId) toProtobuf() *services.HookEntityId {
	// TODO
	return &services.HookEntityId{
		EntityId: &services.HookEntityId_AccountId{
			AccountId: id.accountId._ToProtobuf(),
		},
	}
}

func hookEntityIdFromProtobuf(pb *services.HookEntityId) HookEntityId {
	// TODO
	accountId := _AccountIDFromProtobuf(pb.GetAccountId())
	if accountId != nil {
		return HookEntityId{
			accountId: accountId,
		}
	}
	return HookEntityId{}
}
