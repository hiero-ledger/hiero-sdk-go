package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

type HookCall struct {
	hookId      *int64
	hookIdFull  *HookId
	evmHookCall *EvmHookCall
}

// NewHookCall creates a new HookCall instance
func NewHookCall() *HookCall {
	return &HookCall{}
}

// GetHookId returns the hook ID
func (hc HookCall) GetHookId() int64 {
	if hc.hookId == nil {
		return 0
	}
	return *hc.hookId
}

// SetHookId sets the hook ID
func (hc *HookCall) SetHookId(hookId int64) *HookCall {
	hc.hookId = &hookId
	return hc
}

func (hc HookCall) GetHookIdFull() HookId {
	if hc.hookIdFull == nil {
		return HookId{}
	}
	return *hc.hookIdFull
}

func (hc *HookCall) SetHookIdFull(hookIdFull HookId) *HookCall {
	hc.hookIdFull = &hookIdFull
	return hc
}

// GetEvmHookCall returns the EVM hook call details
func (hc HookCall) GetEvmHookCall() EvmHookCall {
	if hc.evmHookCall == nil {
		return EvmHookCall{}
	}
	return *hc.evmHookCall
}

// SetEvmHookCall sets the EVM hook call details
func (hc *HookCall) SetEvmHookCall(evmHookCall EvmHookCall) *HookCall {
	hc.evmHookCall = &evmHookCall
	return hc
}

func (hc HookCall) toProtobuf() *services.HookCall {
	protoBody := &services.HookCall{}

	if hc.hookIdFull != nil {
		protoBody.Id = &services.HookCall_FullHookId{
			FullHookId: hc.hookIdFull.toProtobuf(),
		}
	}

	if hc.hookId != nil {
		protoBody.Id = &services.HookCall_HookId{
			HookId: *hc.hookId,
		}
	}

	if hc.evmHookCall != nil {
		protoBody.CallSpec = &services.HookCall_EvmHookCall{
			EvmHookCall: hc.evmHookCall.toProtobuf(),
		}
	}

	return protoBody
}

func hookCallFromProtobuf(pb *services.HookCall) HookCall {
	hookId := pb.GetHookId()
	evmHookCall := evmHookCallFromProtobuf(pb.GetEvmHookCall())
	return HookCall{
		hookId:      &hookId,
		evmHookCall: &evmHookCall,
	}
}

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
