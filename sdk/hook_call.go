package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

// NftHookType represents the type of hook that can be executed for NFT transfers.
type NftHookType int32

const (
	// executes a hook on the sender's account before the NFT transfer.
	PRE_HOOK_SENDER NftHookType = iota

	//  xecutes hooks on the sender's account both before and after
	PRE_POST_HOOK_SENDER

	// executes a hook on the receiver's account before the NFT transfer.
	PRE_HOOK_RECEIVER

	// PRE_POST_HOOK_RECEIVER executes hooks on the receiver's account both before and after
	PRE_POST_HOOK_RECEIVER
)

// FungibleHookType represents the type of hook that can be executed for fungible token transfers.
type FungibleHookType uint32

const (
	// executes a hook before the fungible token transfer.
	PRE_HOOK FungibleHookType = iota

	// executes hooks both before and after the fungible token transfer.
	PRE_POST_HOOK
)

type hookCall struct {
	hookId      int64
	evmHookCall EvmHookCall
}

// GetHookId returns the hook ID
func (hc hookCall) GetHookId() int64 {
	return hc.hookId
}

// GetEvmHookCall returns the EVM hook call details
func (hc hookCall) GetEvmHookCall() EvmHookCall {
	return hc.evmHookCall
}

func (hc hookCall) toProtobuf() *services.HookCall {
	protoBody := &services.HookCall{}
	protoBody.Id = &services.HookCall_HookId{
		HookId: hc.hookId,
	}

	protoBody.CallSpec = &services.HookCall_EvmHookCall{
		EvmHookCall: hc.evmHookCall.toProtobuf(),
	}

	return protoBody
}

func hookCallFromProtobuf(pb *services.HookCall) *hookCall {
	if pb == nil {
		return nil
	}
	return &hookCall{
		evmHookCall: evmHookCallFromProtobuf(pb.GetEvmHookCall()),
		hookId:      pb.GetHookId(),
	}
}

type NftHookCall struct {
	hookCall
	hookType NftHookType
}

func NewNftHookCall(hookId int64, evmHookCall EvmHookCall, hookType NftHookType) *NftHookCall {
	return &NftHookCall{
		hookCall: hookCall{
			hookId:      hookId,
			evmHookCall: evmHookCall,
		},
		hookType: hookType,
	}
}

func nftSenderHookCallFromProtobuf(pb *services.NftTransfer) *NftHookCall {
	if pb == nil {
		return nil
	}
	if pb.GetPreTxSenderAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPreTxSenderAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_HOOK_SENDER,
		}
	} else if pb.GetPrePostTxSenderAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPrePostTxSenderAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_POST_HOOK_SENDER,
		}
	}
	return nil
}

func nftReceiverHookCallFromProtobuf(pb *services.NftTransfer) *NftHookCall {
	if pb == nil {
		return nil
	}
	if pb.GetPreTxReceiverAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPreTxReceiverAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_HOOK_RECEIVER,
		}
	} else if pb.GetPrePostTxReceiverAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPrePostTxReceiverAllowanceHook())
		return &NftHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
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

func NewFungibleHookCall(hookId int64, evmHookCall EvmHookCall, hookType FungibleHookType) *FungibleHookCall {
	return &FungibleHookCall{
		hookCall: hookCall{
			hookId:      hookId,
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
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_HOOK,
		}

	} else if pb.GetPrePostTxAllowanceHook() != nil {
		base := hookCallFromProtobuf(pb.GetPrePostTxAllowanceHook())
		return &FungibleHookCall{
			hookCall: hookCall{
				hookId:      base.hookId,
				evmHookCall: base.evmHookCall,
			},
			hookType: PRE_POST_HOOK,
		}
	}

	return nil
}
