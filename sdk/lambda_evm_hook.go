package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

// Definition of a lambda EVM hook.
type EvmHook struct {
	evmHookSpec
	storageUpdates []EvmHookStorageUpdate
}

// NewEvmHook creates a new EvmHook
func NewEvmHook() *EvmHook {
	return &EvmHook{}
}

// GetContractId returns the contract ID
func (leh EvmHook) GetContractId() *ContractID {
	return leh.evmHookSpec.contractId
}

// SetContractId sets the contract ID
func (leh *EvmHook) SetContractId(contractId *ContractID) *EvmHook {
	leh.evmHookSpec.contractId = contractId
	return leh
}

// GetStorageUpdates returns the storage updates
func (leh EvmHook) GetStorageUpdates() []EvmHookStorageUpdate {
	return leh.storageUpdates
}

// SetStorageUpdates sets the storage updates
func (leh *EvmHook) SetStorageUpdates(storageUpdates []EvmHookStorageUpdate) *EvmHook {
	leh.storageUpdates = storageUpdates
	return leh
}

// AddStorageUpdate adds a storage update to the slice
func (leh *EvmHook) AddStorageUpdate(storageUpdate EvmHookStorageUpdate) *EvmHook {
	leh.storageUpdates = append(leh.storageUpdates, storageUpdate)
	return leh
}

func lambdaEvmHookFromProtobuf(pb *services.LambdaEvmHook) EvmHook {
	body := EvmHook{
		evmHookSpec: evmHookSpecFromProtobuf(pb.GetSpec()),
	}

	var storageUpdates []EvmHookStorageUpdate
	for _, storageUpdate := range pb.GetStorageUpdates() {
		storageUpdates = append(storageUpdates, lambdaStorageUpdateFromProtobuf(storageUpdate))
	}
	body.storageUpdates = storageUpdates
	return body
}

func (leh EvmHook) toProtobuf() *services.LambdaEvmHook {
	protoBody := &services.LambdaEvmHook{
		Spec: leh.evmHookSpec.toProtobuf(),
	}

	for _, storageUpdate := range leh.storageUpdates {
		if storageUpdate != nil {
			protoBody.StorageUpdates = append(protoBody.StorageUpdates, storageUpdate.toProtobuf())
		}
	}

	return protoBody
}

// Shared specifications for an EVM hook. May be used for any extension point.
type evmHookSpec struct {
	contractId *ContractID
}

func evmHookSpecFromProtobuf(pb *services.EvmHookSpec) evmHookSpec {
	return evmHookSpec{
		contractId: _ContractIDFromProtobuf(pb.GetContractId()),
	}
}

func (eh evmHookSpec) toProtobuf() *services.EvmHookSpec {
	pbBody := &services.EvmHookSpec{}
	if eh.contractId != nil {
		pbBody.BytecodeSource = &services.EvmHookSpec_ContractId{
			ContractId: eh.contractId._ToProtobuf(),
		}
	}

	return pbBody
}
