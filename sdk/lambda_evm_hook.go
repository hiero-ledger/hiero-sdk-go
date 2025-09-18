package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

type LambdaEvmHook struct {
	EvmHookSpec
	storageUpdates []LambdaStorageUpdate
}

// NewLambdaEvmHook creates a new LambdaEvmHook
func NewLambdaEvmHook() *LambdaEvmHook {
	return &LambdaEvmHook{}
}

// GetEvmHookSpec returns the EVM hook specification
func (leh LambdaEvmHook) GetEvmHookSpec() EvmHookSpec {
	return leh.EvmHookSpec
}

// SetEvmHookSpec sets the EVM hook specification
func (leh *LambdaEvmHook) SetEvmHookSpec(evmHookSpec EvmHookSpec) *LambdaEvmHook {
	leh.EvmHookSpec = evmHookSpec
	return leh
}

// GetStorageUpdates returns the storage updates
func (leh LambdaEvmHook) GetStorageUpdates() []LambdaStorageUpdate {
	return leh.storageUpdates
}

// SetStorageUpdates sets the storage updates
func (leh *LambdaEvmHook) SetStorageUpdates(storageUpdates []LambdaStorageUpdate) *LambdaEvmHook {
	leh.storageUpdates = storageUpdates
	return leh
}

// AddStorageUpdate adds a storage update to the slice
func (leh *LambdaEvmHook) AddStorageUpdate(storageUpdate LambdaStorageUpdate) *LambdaEvmHook {
	leh.storageUpdates = append(leh.storageUpdates, storageUpdate)
	return leh
}

func _LambdaEvmHookFromProtobuf(pb *services.LambdaEvmHook) LambdaEvmHook {
	body := LambdaEvmHook{
		EvmHookSpec: evmHookSpecFromProtobuf(pb.GetSpec()),
	}

	storageUpdates := make([]LambdaStorageUpdate, 0)
	for _, storageUpdate := range pb.GetStorageUpdates() {
		storageUpdates = append(storageUpdates, lambdaStorageUpdateFromProtobuf(storageUpdate))
	}
	body.storageUpdates = storageUpdates
	return body
}

func (leh LambdaEvmHook) toProtobuf() *services.LambdaEvmHook {
	protoBody := &services.LambdaEvmHook{
		Spec: leh.EvmHookSpec.toProtobuf(),
	}

	for _, storageUpdate := range leh.storageUpdates {
		protoBody.StorageUpdates = append(protoBody.StorageUpdates, storageUpdate.toProtobuf())
	}

	return protoBody
}

type EvmHookSpec struct {
	contractId *ContractID
}

// NewEvmHookSpec creates a new EvmHookSpec
func NewEvmHookSpec() *EvmHookSpec {
	return &EvmHookSpec{}
}

// GetContractId returns the contract ID
func (eh EvmHookSpec) GetContractId() ContractID {
	return *eh.contractId
}

// SetContractId sets the contract ID
func (eh *EvmHookSpec) SetContractId(contractId ContractID) *EvmHookSpec {
	eh.contractId = &contractId
	return eh
}

func evmHookSpecFromProtobuf(pb *services.EvmHookSpec) EvmHookSpec {
	return EvmHookSpec{
		contractId: _ContractIDFromProtobuf(pb.GetContractId()),
	}
}

func (eh EvmHookSpec) toProtobuf() *services.EvmHookSpec {
	protoBody := &services.EvmHookSpec{}
	if eh.contractId != nil {
		protoBody.BytecodeSource = &services.EvmHookSpec_ContractId{
			ContractId: eh.contractId._ToProtobuf(),
		}
	}
	return protoBody
}
