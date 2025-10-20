//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitNftHookCallWithHookId(t *testing.T) {
	t.Parallel()

	hookId := int64(123)
	evmCall := *NewEvmHookCall().SetData([]byte{0x01, 0x02}).SetGasLimit(25000)

	nftHookCall := NewNftHookCallWithHookId(hookId, evmCall, PRE_HOOK_SENDER)

	require.NotNil(t, nftHookCall)
	require.Equal(t, hookId, nftHookCall.GetHookId())
	require.Equal(t, evmCall.GetData(), nftHookCall.GetEvmHookCall().GetData())
	require.Equal(t, evmCall.GetGasLimit(), nftHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_HOOK_SENDER, nftHookCall.hookType)
}

func TestUnitNftHookCallWithHookIdFull(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Contract: 456}
	entityId := NewHookEntityIdWithContractId(contractId)
	hookIdFull := NewHookId(*entityId, 789)
	evmCall := *NewEvmHookCall().SetData([]byte{0x03, 0x04}).SetGasLimit(30000)

	nftHookCall := NewNftHookCallWithHookIdFull(*hookIdFull, evmCall, PRE_POST_HOOK_SENDER)

	require.NotNil(t, nftHookCall)
	require.Equal(t, hookIdFull.GetHookId(), nftHookCall.GetHookIdFull().GetHookId())
	require.Equal(t, contractId.Contract, nftHookCall.GetHookIdFull().GetEntityId().GetContractId().Contract)
	require.Equal(t, evmCall.GetData(), nftHookCall.GetEvmHookCall().GetData())
	require.Equal(t, evmCall.GetGasLimit(), nftHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_POST_HOOK_SENDER, nftHookCall.hookType)
}

func TestUnitNftHookCallAllHookTypes(t *testing.T) {
	t.Parallel()

	hookId := int64(100)
	evmCall := *NewEvmHookCall().SetData([]byte{}).SetGasLimit(10000)

	testCases := []struct {
		name     string
		hookType NftHookType
	}{
		{"PRE_HOOK_SENDER", PRE_HOOK_SENDER},
		{"PRE_POST_HOOK_SENDER", PRE_POST_HOOK_SENDER},
		{"PRE_HOOK_RECEIVER", PRE_HOOK_RECEIVER},
		{"PRE_POST_HOOK_RECEIVER", PRE_POST_HOOK_RECEIVER},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hookCall := NewNftHookCallWithHookId(hookId, evmCall, tc.hookType)
			require.Equal(t, tc.hookType, hookCall.hookType)
		})
	}
}

func TestUnitNftHookCallToProtobuf(t *testing.T) {
	t.Parallel()

	hookId := int64(999)
	evmCall := *NewEvmHookCall().SetData([]byte{0xaa, 0xbb}).SetGasLimit(50000)

	nftHookCall := NewNftHookCallWithHookId(hookId, evmCall, PRE_HOOK_SENDER)

	proto := nftHookCall.toProtobuf()
	require.NotNil(t, proto)
	require.Equal(t, hookId, proto.GetHookId())
	require.NotNil(t, proto.GetEvmHookCall())
	require.Equal(t, evmCall.GetData(), proto.GetEvmHookCall().GetData())
	require.Equal(t, evmCall.GetGasLimit(), proto.GetEvmHookCall().GetGasLimit())
}

func TestUnitNftHookCallFromProtobufPreHookSender(t *testing.T) {
	t.Parallel()

	hookId := int64(111)
	data := []byte{0x11, 0x22, 0x33}
	gasLimit := uint64(15000)

	pbNftTransfer := &services.NftTransfer{
		SenderAllowanceHookCall: &services.NftTransfer_PreTxSenderAllowanceHook{
			PreTxSenderAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{
					HookId: hookId,
				},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     data,
						GasLimit: gasLimit,
					},
				},
			},
		},
	}

	nftHookCall := nftHookCallFromProtobuf(pbNftTransfer)
	require.NotNil(t, nftHookCall)
	require.Equal(t, hookId, nftHookCall.GetHookId())
	require.Equal(t, data, nftHookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, nftHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_HOOK_SENDER, nftHookCall.hookType)
}

func TestUnitNftHookCallFromProtobufPrePostHookSender(t *testing.T) {
	t.Parallel()

	hookId := int64(222)
	data := []byte{0x44, 0x55}
	gasLimit := uint64(20000)

	pbNftTransfer := &services.NftTransfer{
		SenderAllowanceHookCall: &services.NftTransfer_PrePostTxSenderAllowanceHook{
			PrePostTxSenderAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{
					HookId: hookId,
				},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     data,
						GasLimit: gasLimit,
					},
				},
			},
		},
	}

	nftHookCall := nftHookCallFromProtobuf(pbNftTransfer)
	require.NotNil(t, nftHookCall)
	require.Equal(t, hookId, nftHookCall.GetHookId())
	require.Equal(t, data, nftHookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, nftHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_POST_HOOK_SENDER, nftHookCall.hookType)
}

func TestUnitNftHookCallFromProtobufPreHookReceiver(t *testing.T) {
	t.Parallel()

	hookId := int64(333)
	data := []byte{0x66, 0x77, 0x88}
	gasLimit := uint64(35000)

	pbNftTransfer := &services.NftTransfer{
		ReceiverAllowanceHookCall: &services.NftTransfer_PreTxReceiverAllowanceHook{
			PreTxReceiverAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{
					HookId: hookId,
				},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     data,
						GasLimit: gasLimit,
					},
				},
			},
		},
	}

	nftHookCall := nftHookCallFromProtobuf(pbNftTransfer)
	require.NotNil(t, nftHookCall)
	require.Equal(t, hookId, nftHookCall.GetHookId())
	require.Equal(t, data, nftHookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, nftHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_HOOK_RECEIVER, nftHookCall.hookType)
}

func TestUnitNftHookCallFromProtobufPrePostHookReceiver(t *testing.T) {
	t.Parallel()

	hookId := int64(444)
	data := []byte{0x99, 0xaa}
	gasLimit := uint64(40000)

	pbNftTransfer := &services.NftTransfer{
		ReceiverAllowanceHookCall: &services.NftTransfer_PrePostTxReceiverAllowanceHook{
			PrePostTxReceiverAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{
					HookId: hookId,
				},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     data,
						GasLimit: gasLimit,
					},
				},
			},
		},
	}

	nftHookCall := nftHookCallFromProtobuf(pbNftTransfer)
	require.NotNil(t, nftHookCall)
	require.Equal(t, hookId, nftHookCall.GetHookId())
	require.Equal(t, data, nftHookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, nftHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_POST_HOOK_RECEIVER, nftHookCall.hookType)
}

func TestUnitNftHookCallFromProtobufNil(t *testing.T) {
	t.Parallel()

	nftHookCall := nftHookCallFromProtobuf(nil)
	require.Nil(t, nftHookCall)
}

func TestUnitNftHookCallFromProtobufNoHook(t *testing.T) {
	t.Parallel()

	pbNftTransfer := &services.NftTransfer{}

	nftHookCall := nftHookCallFromProtobuf(pbNftTransfer)
	require.Nil(t, nftHookCall)
}

// Tests for FungibleHookCall

func TestUnitFungibleHookCallWithHookId(t *testing.T) {
	t.Parallel()

	hookId := int64(555)
	evmCall := *NewEvmHookCall().SetData([]byte{0xcc, 0xdd}).SetGasLimit(45000)

	fungibleHookCall := NewFungibleHookCallWithHookId(hookId, evmCall, PRE_HOOK)

	require.NotNil(t, fungibleHookCall)
	require.Equal(t, hookId, fungibleHookCall.GetHookId())
	require.Equal(t, evmCall.GetData(), fungibleHookCall.GetEvmHookCall().GetData())
	require.Equal(t, evmCall.GetGasLimit(), fungibleHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_HOOK, fungibleHookCall.hookType)
}

func TestUnitFungibleHookCallWithHookIdFull(t *testing.T) {
	t.Parallel()

	contractId := ContractID{Contract: 777}
	entityId := NewHookEntityIdWithContractId(contractId)
	hookIdFull := NewHookId(*entityId, 888)
	evmCall := *NewEvmHookCall().SetData([]byte{0xee, 0xff}).SetGasLimit(55000)

	fungibleHookCall := NewFungibleHookCallWithHookIdFull(*hookIdFull, evmCall, PRE_POST_HOOK)

	require.NotNil(t, fungibleHookCall)
	require.Equal(t, hookIdFull.GetHookId(), fungibleHookCall.GetHookIdFull().GetHookId())
	require.Equal(t, contractId.Contract, fungibleHookCall.GetHookIdFull().GetEntityId().GetContractId().Contract)
	require.Equal(t, evmCall.GetData(), fungibleHookCall.GetEvmHookCall().GetData())
	require.Equal(t, evmCall.GetGasLimit(), fungibleHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_POST_HOOK, fungibleHookCall.hookType)
}

func TestUnitFungibleHookCallAllHookTypes(t *testing.T) {
	t.Parallel()

	hookId := int64(200)
	evmCall := *NewEvmHookCall().SetData([]byte{}).SetGasLimit(12000)

	testCases := []struct {
		name     string
		hookType FungibleHookType
	}{
		{"PRE_HOOK", PRE_HOOK},
		{"PRE_POST_HOOK", PRE_POST_HOOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hookCall := NewFungibleHookCallWithHookId(hookId, evmCall, tc.hookType)
			require.Equal(t, tc.hookType, hookCall.hookType)
		})
	}
}

func TestUnitFungibleHookCallToProtobuf(t *testing.T) {
	t.Parallel()

	hookId := int64(666)
	evmCall := *NewEvmHookCall().SetData([]byte{0x12, 0x34}).SetGasLimit(60000)

	fungibleHookCall := NewFungibleHookCallWithHookId(hookId, evmCall, PRE_HOOK)

	proto := fungibleHookCall.toProtobuf()
	require.NotNil(t, proto)
	require.Equal(t, hookId, proto.GetHookId())
	require.NotNil(t, proto.GetEvmHookCall())
	require.Equal(t, evmCall.GetData(), proto.GetEvmHookCall().GetData())
	require.Equal(t, evmCall.GetGasLimit(), proto.GetEvmHookCall().GetGasLimit())
}

func TestUnitFungibleHookCallFromProtobufPreHook(t *testing.T) {
	t.Parallel()

	hookId := int64(777)
	data := []byte{0xab, 0xcd, 0xef}
	gasLimit := uint64(70000)

	pbAccountAmount := &services.AccountAmount{
		HookCall: &services.AccountAmount_PreTxAllowanceHook{
			PreTxAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{
					HookId: hookId,
				},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     data,
						GasLimit: gasLimit,
					},
				},
			},
		},
	}

	fungibleHookCall := fungibleHookCallFromProtobuf(pbAccountAmount)
	require.NotNil(t, fungibleHookCall)
	require.Equal(t, hookId, fungibleHookCall.GetHookId())
	require.Equal(t, data, fungibleHookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, fungibleHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_HOOK, fungibleHookCall.hookType)
}

func TestUnitFungibleHookCallFromProtobufPrePostHook(t *testing.T) {
	t.Parallel()

	hookId := int64(888)
	data := []byte{0xfe, 0xdc}
	gasLimit := uint64(80000)

	pbAccountAmount := &services.AccountAmount{
		HookCall: &services.AccountAmount_PrePostTxAllowanceHook{
			PrePostTxAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{
					HookId: hookId,
				},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     data,
						GasLimit: gasLimit,
					},
				},
			},
		},
	}

	fungibleHookCall := fungibleHookCallFromProtobuf(pbAccountAmount)
	require.NotNil(t, fungibleHookCall)
	require.Equal(t, hookId, fungibleHookCall.GetHookId())
	require.Equal(t, data, fungibleHookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, fungibleHookCall.GetEvmHookCall().GetGasLimit())
	require.Equal(t, PRE_POST_HOOK, fungibleHookCall.hookType)
}

func TestUnitFungibleHookCallFromProtobufNilHookCall(t *testing.T) {
	t.Parallel()

	pbAccountAmount := &services.AccountAmount{
		HookCall: nil,
	}

	fungibleHookCall := fungibleHookCallFromProtobuf(pbAccountAmount)
	require.Nil(t, fungibleHookCall)
}

// Tests for base hookCall

func TestUnitHookCallFromProtobufNil(t *testing.T) {
	t.Parallel()

	hookCall := hookCallFromProtobuf(nil)
	require.Nil(t, hookCall)
}

func TestUnitHookCallFromProtobufWithHookId(t *testing.T) {
	t.Parallel()

	hookId := int64(1234)
	data := []byte{0x01, 0x02, 0x03}
	gasLimit := uint64(90000)

	pb := &services.HookCall{
		Id: &services.HookCall_HookId{
			HookId: hookId,
		},
		CallSpec: &services.HookCall_EvmHookCall{
			EvmHookCall: &services.EvmHookCall{
				Data:     data,
				GasLimit: gasLimit,
			},
		},
	}

	hookCall := hookCallFromProtobuf(pb)
	require.NotNil(t, hookCall)
	require.Equal(t, hookId, hookCall.GetHookId())
	require.Equal(t, data, hookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, hookCall.GetEvmHookCall().GetGasLimit())
}

func TestUnitHookCallFromProtobufWithFullHookId(t *testing.T) {
	t.Parallel()

	contractNum := int64(5678)
	hookIdNum := int64(9012)
	data := []byte{0x04, 0x05, 0x06}
	gasLimit := uint64(95000)

	pb := &services.HookCall{
		Id: &services.HookCall_FullHookId{
			FullHookId: &services.HookId{
				EntityId: &services.HookEntityId{
					EntityId: &services.HookEntityId_ContractId{
						ContractId: &services.ContractID{
							Contract: &services.ContractID_ContractNum{
								ContractNum: contractNum,
							},
						},
					},
				},
				HookId: hookIdNum,
			},
		},
		CallSpec: &services.HookCall_EvmHookCall{
			EvmHookCall: &services.EvmHookCall{
				Data:     data,
				GasLimit: gasLimit,
			},
		},
	}

	hookCall := hookCallFromProtobuf(pb)
	require.NotNil(t, hookCall)
	require.Equal(t, hookIdNum, hookCall.GetHookIdFull().GetHookId())
	require.Equal(t, uint64(contractNum), hookCall.GetHookIdFull().GetEntityId().GetContractId().Contract)
	require.Equal(t, data, hookCall.GetEvmHookCall().GetData())
	require.Equal(t, gasLimit, hookCall.GetEvmHookCall().GetGasLimit())
}

func TestUnitHookCallGetHookIdWhenNil(t *testing.T) {
	t.Parallel()

	hc := hookCall{
		hookId:      nil,
		hookIdFull:  nil,
		evmHookCall: *NewEvmHookCall(),
	}

	require.Equal(t, int64(0), hc.GetHookId())
}

func TestUnitHookCallGetHookIdFullWhenNil(t *testing.T) {
	t.Parallel()

	hc := hookCall{
		hookId:      nil,
		hookIdFull:  nil,
		evmHookCall: *NewEvmHookCall(),
	}

	hookIdFull := hc.GetHookIdFull()
	require.Equal(t, HookId{}, hookIdFull)
}

func TestUnitNftHookCallRoundTripProtobufWithHookId(t *testing.T) {
	t.Parallel()

	hookId := int64(11111)
	evmCall := *NewEvmHookCall().SetData([]byte{0xaa, 0xbb, 0xcc, 0xdd}).SetGasLimit(100000)

	original := NewNftHookCallWithHookId(hookId, evmCall, PRE_HOOK_SENDER)

	// Convert to protobuf
	proto := original.toProtobuf()
	require.NotNil(t, proto)

	// Create NFT transfer protobuf wrapper for round-trip
	pbNftTransfer := &services.NftTransfer{
		SenderAllowanceHookCall: &services.NftTransfer_PreTxSenderAllowanceHook{
			PreTxSenderAllowanceHook: proto,
		},
	}

	// Convert back from protobuf
	reconstructed := nftHookCallFromProtobuf(pbNftTransfer)

	// Verify round-trip
	require.Equal(t, original.GetHookId(), reconstructed.GetHookId())
	require.Equal(t, original.GetEvmHookCall().GetData(), reconstructed.GetEvmHookCall().GetData())
	require.Equal(t, original.GetEvmHookCall().GetGasLimit(), reconstructed.GetEvmHookCall().GetGasLimit())
	require.Equal(t, original.hookType, reconstructed.hookType)
}

func TestUnitFungibleHookCallRoundTripProtobufWithHookId(t *testing.T) {
	t.Parallel()

	hookId := int64(22222)
	evmCall := *NewEvmHookCall().SetData([]byte{0x11, 0x22, 0x33}).SetGasLimit(110000)

	original := NewFungibleHookCallWithHookId(hookId, evmCall, PRE_POST_HOOK)

	// Convert to protobuf
	proto := original.toProtobuf()
	require.NotNil(t, proto)

	// Create AccountAmount protobuf wrapper for round-trip
	pbAccountAmount := &services.AccountAmount{
		HookCall: &services.AccountAmount_PrePostTxAllowanceHook{
			PrePostTxAllowanceHook: proto,
		},
	}

	// Convert back from protobuf
	reconstructed := fungibleHookCallFromProtobuf(pbAccountAmount)

	// Verify round-trip
	require.Equal(t, original.GetHookId(), reconstructed.GetHookId())
	require.Equal(t, original.GetEvmHookCall().GetData(), reconstructed.GetEvmHookCall().GetData())
	require.Equal(t, original.GetEvmHookCall().GetGasLimit(), reconstructed.GetEvmHookCall().GetGasLimit())
	require.Equal(t, original.hookType, reconstructed.hookType)
}
