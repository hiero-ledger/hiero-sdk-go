package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

// Definition of a lambda EVM hook.
type LambdaEvmHook struct {
	evmHookSpec
	storageUpdates []LambdaStorageUpdate
}

// NewLambdaEvmHook creates a new LambdaEvmHook
func NewLambdaEvmHook() *LambdaEvmHook {
	return &LambdaEvmHook{}
}

// GetEvmHookSpec returns the EVM hook specification
func (leh LambdaEvmHook) GetEvmHookSpec() evmHookSpec {
	return leh.evmHookSpec
}

// SetEvmHookSpec sets the EVM hook specification
func (leh *LambdaEvmHook) SetEvmHookSpec(evmHookSpec evmHookSpec) *LambdaEvmHook {
	leh.evmHookSpec = evmHookSpec
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

func lambdaEvmHookFromProtobuf(pb *services.LambdaEvmHook) LambdaEvmHook {
	body := LambdaEvmHook{
		evmHookSpec: evmHookSpecFromProtobuf(pb.GetSpec()),
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
		Spec: leh.evmHookSpec.toProtobuf(),
	}

	for _, storageUpdate := range leh.storageUpdates {
		protoBody.StorageUpdates = append(protoBody.StorageUpdates, storageUpdate.toProtobuf())
	}

	return protoBody
}

type evmHookSpec struct {
	contractId *ContractID
}

// GetContractId returns the contract ID
func (eh evmHookSpec) GetContractId() *ContractID {
	return eh.contractId
}

// SetContractId sets the contract ID
func (eh *evmHookSpec) SetContractId(contractId *ContractID) *evmHookSpec {
	eh.contractId = contractId
	return eh
}

func evmHookSpecFromProtobuf(pb *services.EvmHookSpec) evmHookSpec {
	return evmHookSpec{
		contractId: _ContractIDFromProtobuf(pb.GetContractId()),
	}
}

func (eh evmHookSpec) toProtobuf() *services.EvmHookSpec {
	return &services.EvmHookSpec{
		BytecodeSource: &services.EvmHookSpec_ContractId{
			ContractId: eh.contractId._ToProtobuf(),
		},
	}
}
