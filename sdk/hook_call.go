package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

type NftHookType int32

const (
	PRE_HOOK_SENDER NftHookType = iota
	PRE_POST_HOOK_SENDER
	PRE_HOOK_RECEIVER
	PRE_POST_HOOK_RECEIVER
)

type FungibleHookType uint32

const (
	PRE_HOOK FungibleHookType = iota
	PRE_POST_HOOK
)

type hookCall struct {
	hookId      *int64
	hookIdFull  *HookId
	evmHookCall EvmHookCall
}

// GetHookId returns the hook ID
func (hc hookCall) GetHookId() int64 {
	if hc.hookId == nil {
		return 0
	}
	return *hc.hookId
}

func (hc hookCall) GetHookIdFull() HookId {
	if hc.hookIdFull == nil {
		return HookId{}
	}
	return *hc.hookIdFull
}

// GetEvmHookCall returns the EVM hook call details
func (hc hookCall) GetEvmHookCall() EvmHookCall {
	return hc.evmHookCall
}

func (hc hookCall) toProtobuf() *services.HookCall {
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

	protoBody.CallSpec = &services.HookCall_EvmHookCall{
		EvmHookCall: hc.evmHookCall.toProtobuf(),
	}

	return protoBody
}

func (hc hookCall) validateChecksum(client *Client) error {
	if hc.hookIdFull != nil {
		return hc.hookIdFull.validateChecksum(client)
	}
	return nil
}

func hookCallFromProtobuf(pb *services.HookCall) *hookCall {
	if pb == nil {
		return nil
	}
	hookCall := &hookCall{
		evmHookCall: evmHookCallFromProtobuf(pb.GetEvmHookCall()),
	}

	if pb.GetFullHookId() != nil {
		hookIdFull := hookIdFromProtobuf(pb.GetFullHookId())
		hookCall.hookIdFull = &hookIdFull
	} else {
		hookId := pb.GetHookId()
		hookCall.hookId = &hookId
	}

	return hookCall
}

type NftHookCall struct {
	hookCall
	hookType NftHookType
}

func NewNftHookCallWithHookId(hookId int64, evmHookCall EvmHookCall, hookType NftHookType) *NftHookCall {
	return &NftHookCall{
		hookCall: hookCall{
			hookId:      &hookId,
			hookIdFull:  nil,
			evmHookCall: evmHookCall,
		},
		hookType: hookType,
	}
}

func NewNftHookCallWithHookIdFull(hookIdFull HookId, evmHookCall EvmHookCall, hookType NftHookType) *NftHookCall {
	return &NftHookCall{
		hookCall: hookCall{
			hookId:      nil,
			hookIdFull:  &hookIdFull,
			evmHookCall: evmHookCall,
		},
		hookType: hookType,
	}
}

func nftHookCallFromProtobuf(pb *services.NftTransfer) *NftHookCall {
	if pb == nil {
		return nil
	}
	if pb.GetPreTxSenderAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPreTxSenderAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				hookIdFull:  base.hookIdFull,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_HOOK_SENDER,
		}
	} else if pb.GetPrePostTxSenderAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPrePostTxSenderAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				hookIdFull:  base.hookIdFull,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_POST_HOOK_SENDER,
		}
	}
	if pb.GetPreTxReceiverAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPreTxReceiverAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				hookIdFull:  base.hookIdFull,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_HOOK_RECEIVER,
		}
	} else if pb.GetPrePostTxReceiverAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPrePostTxReceiverAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				hookIdFull:  base.hookIdFull,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_POST_HOOK_RECEIVER,
		}
	}
	return nil
}

type FungibleHookCall struct {
	hookCall
	hookType FungibleHookType
}

func NewFungibleHookCallWithHookId(hookId int64, evmHookCall EvmHookCall, hookType FungibleHookType) *FungibleHookCall {
	return &FungibleHookCall{
		hookCall: hookCall{
			hookId:      &hookId,
			hookIdFull:  nil,
			evmHookCall: evmHookCall,
		},
		hookType: hookType,
	}
}

func NewFungibleHookCallWithHookIdFull(hookIdFull HookId, evmHookCall EvmHookCall, hookType FungibleHookType) *FungibleHookCall {
	return &FungibleHookCall{
		hookCall: hookCall{
			hookId:      nil,
			hookIdFull:  &hookIdFull,
			evmHookCall: evmHookCall,
		},
		hookType: hookType,
	}
}

func fungibleHookCallFromProtobuf(pb *services.AccountAmount) *FungibleHookCall {
	if pb.HookCall == nil {
		return nil
	}

	if pb.GetPreTxAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPreTxAllowanceHook())
		return &FungibleHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				hookIdFull:  base.hookIdFull,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_HOOK,
		}

	} else if pb.GetPrePostTxAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPrePostTxAllowanceHook())
		return &FungibleHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				hookIdFull:  base.hookIdFull,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_POST_HOOK,
		}
	}

	return nil
}
