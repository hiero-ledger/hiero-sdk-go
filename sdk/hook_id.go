package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

// hook id
type HookId struct {
	entityId HookEntityId
	hookId   int64
}

// NewHookId creates a new HookId
func NewHookId() *HookId {
	return &HookId{}
}

// GetEntityId returns the entity ID
func (id HookId) GetEntityId() HookEntityId {
	return id.entityId
}

// SetEntityId sets the entity ID
func (id *HookId) SetEntityId(entityId HookEntityId) *HookId {
	id.entityId = entityId
	return id
}

// GetHookId returns the hook ID
func (id HookId) GetHookId() int64 {
	return id.hookId
}

// SetHookId sets the hook ID
func (id *HookId) SetHookId(hookId int64) *HookId {
	id.hookId = hookId
	return id
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
	accountId AccountID
}

// NewHookEntityId creates a new HookEntityId
func NewHookEntityId() *HookEntityId {
	return &HookEntityId{}
}

// GetAccountId returns the account ID
func (id HookEntityId) GetAccountId() AccountID {
	return id.accountId
}

// SetAccountId sets the account ID
func (id *HookEntityId) SetAccountId(accountId AccountID) *HookEntityId {
	id.accountId = accountId
	return id
}

func (id HookEntityId) validateChecksum(client *Client) error {
	return id.accountId.ValidateChecksum(client)
}

func (id HookEntityId) toProtobuf() *services.HookEntityId {
	return &services.HookEntityId{
		EntityId: &services.HookEntityId_AccountId{
			AccountId: id.accountId._ToProtobuf(),
		},
	}
}

func hookEntityIdFromProtobuf(pb *services.HookEntityId) HookEntityId {
	accountId := _AccountIDFromProtobuf(pb.GetAccountId())
	if accountId != nil {
		return HookEntityId{
			accountId: *accountId,
		}
	}
	return HookEntityId{}
}
