package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

type EvmHookCall struct {
	data     []byte
	gasLimit uint64
}

// NewEvmHookCall creates a new EvmHookCall instance
func NewEvmHookCall() *EvmHookCall {
	return &EvmHookCall{}
}

// GetData returns the call data
func (ehc EvmHookCall) GetData() []byte {
	return ehc.data
}

// SetData sets the call data
func (ehc *EvmHookCall) SetData(data []byte) *EvmHookCall {
	ehc.data = data
	return ehc
}

// GetGasLimit returns the gas limit
func (ehc EvmHookCall) GetGasLimit() uint64 {
	return ehc.gasLimit
}

// SetGasLimit sets the gas limit
func (ehc *EvmHookCall) SetGasLimit(gasLimit uint64) *EvmHookCall {
	ehc.gasLimit = gasLimit
	return ehc
}

func (ehc EvmHookCall) toProtobuf() *services.EvmHookCall {
	return &services.EvmHookCall{
		Data:     ehc.data,
		GasLimit: ehc.gasLimit,
	}
}

func evmHookCallFromProtobuf(pb *services.EvmHookCall) EvmHookCall {
	return EvmHookCall{
		data:     pb.GetData(),
		gasLimit: pb.GetGasLimit(),
	}
}
