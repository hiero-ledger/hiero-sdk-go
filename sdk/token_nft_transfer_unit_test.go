//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitNftTransferFromProtobufWithSenderHookOnly(t *testing.T) {
	t.Parallel()

	pb := &services.NftTransfer{
		SenderAccountID:   &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 100}},
		ReceiverAccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 200}},
		SerialNumber:      1,
		IsApproval:        false,
		SenderAllowanceHookCall: &services.NftTransfer_PreTxSenderAllowanceHook{
			PreTxSenderAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{HookId: 2},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     []byte{0x01, 0x02},
						GasLimit: 25000,
					},
				},
			},
		},
	}

	result := _NftTransferFromProtobuf(pb)

	assert.Equal(t, AccountID{Account: 100}, result.SenderAccountID)
	assert.Equal(t, AccountID{Account: 200}, result.ReceiverAccountID)
	assert.Equal(t, int64(1), result.SerialNumber)
	assert.False(t, result.IsApproved)

	// Should have sender hook but no receiver hook
	require.NotNil(t, result.SenderHookCall)
	assert.Equal(t, PRE_HOOK_SENDER, result.SenderHookCall.hookType)
	assert.Equal(t, int64(2), result.SenderHookCall.GetHookId())
	assert.Equal(t, []byte{0x01, 0x02}, result.SenderHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, uint64(25000), result.SenderHookCall.GetEvmHookCall().GetGasLimit())

	assert.Nil(t, result.ReceiverHookCall)
}

func TestUnitNftTransferFromProtobufWithReceiverHookOnly(t *testing.T) {
	t.Parallel()

	pb := &services.NftTransfer{
		SenderAccountID:   &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 100}},
		ReceiverAccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 200}},
		SerialNumber:      1,
		IsApproval:        false,
		ReceiverAllowanceHookCall: &services.NftTransfer_PreTxReceiverAllowanceHook{
			PreTxReceiverAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{HookId: 3},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     []byte{0x03, 0x04},
						GasLimit: 30000,
					},
				},
			},
		},
	}

	result := _NftTransferFromProtobuf(pb)

	assert.Equal(t, AccountID{Account: 100}, result.SenderAccountID)
	assert.Equal(t, AccountID{Account: 200}, result.ReceiverAccountID)
	assert.Equal(t, int64(1), result.SerialNumber)
	assert.False(t, result.IsApproved)

	// Should have receiver hook but no sender hook
	assert.Nil(t, result.SenderHookCall)

	require.NotNil(t, result.ReceiverHookCall)
	assert.Equal(t, PRE_HOOK_RECEIVER, result.ReceiverHookCall.hookType)
	assert.Equal(t, int64(3), result.ReceiverHookCall.GetHookId())
	assert.Equal(t, []byte{0x03, 0x04}, result.ReceiverHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, uint64(30000), result.ReceiverHookCall.GetEvmHookCall().GetGasLimit())
}

func TestUnitNftTransferFromProtobufWithBothSenderAndReceiverHooks(t *testing.T) {
	t.Parallel()

	pb := &services.NftTransfer{
		SenderAccountID:   &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 100}},
		ReceiverAccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 200}},
		SerialNumber:      1,
		IsApproval:        false,
		SenderAllowanceHookCall: &services.NftTransfer_PreTxSenderAllowanceHook{
			PreTxSenderAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{HookId: 2},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     []byte{0x01, 0x02},
						GasLimit: 25000,
					},
				},
			},
		},
		ReceiverAllowanceHookCall: &services.NftTransfer_PreTxReceiverAllowanceHook{
			PreTxReceiverAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{HookId: 3},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     []byte{0x03, 0x04},
						GasLimit: 30000,
					},
				},
			},
		},
	}

	result := _NftTransferFromProtobuf(pb)

	assert.Equal(t, AccountID{Account: 100}, result.SenderAccountID)
	assert.Equal(t, AccountID{Account: 200}, result.ReceiverAccountID)
	assert.Equal(t, int64(1), result.SerialNumber)
	assert.False(t, result.IsApproved)

	// Should have both sender and receiver hooks with different values
	require.NotNil(t, result.SenderHookCall)
	assert.Equal(t, PRE_HOOK_SENDER, result.SenderHookCall.hookType)
	assert.Equal(t, int64(2), result.SenderHookCall.GetHookId())
	assert.Equal(t, []byte{0x01, 0x02}, result.SenderHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, uint64(25000), result.SenderHookCall.GetEvmHookCall().GetGasLimit())

	require.NotNil(t, result.ReceiverHookCall)
	assert.Equal(t, PRE_HOOK_RECEIVER, result.ReceiverHookCall.hookType)
	assert.Equal(t, int64(3), result.ReceiverHookCall.GetHookId())
	assert.Equal(t, []byte{0x03, 0x04}, result.ReceiverHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, uint64(30000), result.ReceiverHookCall.GetEvmHookCall().GetGasLimit())

	// Verify they are different objects
	assert.NotEqual(t, result.SenderHookCall.GetHookId(), result.ReceiverHookCall.GetHookId())
}

func TestUnitNftTransferFromProtobufWithPrePostHooks(t *testing.T) {
	t.Parallel()

	pb := &services.NftTransfer{
		SenderAccountID:   &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 100}},
		ReceiverAccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 200}},
		SerialNumber:      1,
		IsApproval:        false,
		SenderAllowanceHookCall: &services.NftTransfer_PrePostTxSenderAllowanceHook{
			PrePostTxSenderAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{HookId: 4},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     []byte{0x05, 0x06},
						GasLimit: 35000,
					},
				},
			},
		},
		ReceiverAllowanceHookCall: &services.NftTransfer_PrePostTxReceiverAllowanceHook{
			PrePostTxReceiverAllowanceHook: &services.HookCall{
				Id: &services.HookCall_HookId{HookId: 5},
				CallSpec: &services.HookCall_EvmHookCall{
					EvmHookCall: &services.EvmHookCall{
						Data:     []byte{0x07, 0x08},
						GasLimit: 40000,
					},
				},
			},
		},
	}

	result := _NftTransferFromProtobuf(pb)

	assert.Equal(t, AccountID{Account: 100}, result.SenderAccountID)
	assert.Equal(t, AccountID{Account: 200}, result.ReceiverAccountID)
	assert.Equal(t, int64(1), result.SerialNumber)
	assert.False(t, result.IsApproved)

	// Should have both sender and receiver hooks with PRE_POST types
	require.NotNil(t, result.SenderHookCall)
	assert.Equal(t, PRE_POST_HOOK_SENDER, result.SenderHookCall.hookType)
	assert.Equal(t, int64(4), result.SenderHookCall.GetHookId())
	assert.Equal(t, []byte{0x05, 0x06}, result.SenderHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, uint64(35000), result.SenderHookCall.GetEvmHookCall().GetGasLimit())

	require.NotNil(t, result.ReceiverHookCall)
	assert.Equal(t, PRE_POST_HOOK_RECEIVER, result.ReceiverHookCall.hookType)
	assert.Equal(t, int64(5), result.ReceiverHookCall.GetHookId())
	assert.Equal(t, []byte{0x07, 0x08}, result.ReceiverHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, uint64(40000), result.ReceiverHookCall.GetEvmHookCall().GetGasLimit())

	// Verify they are different objects
	assert.NotEqual(t, result.SenderHookCall.GetHookId(), result.ReceiverHookCall.GetHookId())
}

func TestUnitNftTransferFromProtobufWithNoHooks(t *testing.T) {
	t.Parallel()

	pb := &services.NftTransfer{
		SenderAccountID:   &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 100}},
		ReceiverAccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 200}},
		SerialNumber:      1,
		IsApproval:        false,
	}

	result := _NftTransferFromProtobuf(pb)

	assert.Equal(t, AccountID{Account: 100}, result.SenderAccountID)
	assert.Equal(t, AccountID{Account: 200}, result.ReceiverAccountID)
	assert.Equal(t, int64(1), result.SerialNumber)
	assert.False(t, result.IsApproved)

	// Should have no hooks
	assert.Nil(t, result.SenderHookCall)
	assert.Nil(t, result.ReceiverHookCall)
}

func TestUnitNftTransferFromProtobufNilInput(t *testing.T) {
	t.Parallel()

	result := _NftTransferFromProtobuf(nil)

	assert.Equal(t, _TokenNftTransfer{}, result)
}

func TestUnitNftTransferRoundTripWithBothHooks(t *testing.T) {
	t.Parallel()

	// Create original transfer with both sender and receiver hooks
	original := _TokenNftTransfer{
		SenderAccountID:   AccountID{Account: 100},
		ReceiverAccountID: AccountID{Account: 200},
		SerialNumber:      1,
		IsApproved:        false,
		SenderHookCall: &NftHookCall{
			hookCall: hookCall{
				hookId:      2,
				evmHookCall: EvmHookCall{data: []byte{0x01, 0x02}, gasLimit: 25000},
			},
			hookType: PRE_HOOK_SENDER,
		},
		ReceiverHookCall: &NftHookCall{
			hookCall: hookCall{
				hookId:      3,
				evmHookCall: EvmHookCall{data: []byte{0x03, 0x04}, gasLimit: 30000},
			},
			hookType: PRE_HOOK_RECEIVER,
		},
	}

	// Convert to protobuf and back
	pb := original._ToProtobuf()
	result := _NftTransferFromProtobuf(pb)

	// Verify round-trip preserves all data
	assert.Equal(t, original.SenderAccountID, result.SenderAccountID)
	assert.Equal(t, original.ReceiverAccountID, result.ReceiverAccountID)
	assert.Equal(t, original.SerialNumber, result.SerialNumber)
	assert.Equal(t, original.IsApproved, result.IsApproved)

	// Verify sender hook
	require.NotNil(t, result.SenderHookCall)
	assert.Equal(t, original.SenderHookCall.hookType, result.SenderHookCall.hookType)
	assert.Equal(t, original.SenderHookCall.GetHookId(), result.SenderHookCall.GetHookId())
	assert.Equal(t, original.SenderHookCall.GetEvmHookCall().GetData(), result.SenderHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, original.SenderHookCall.GetEvmHookCall().GetGasLimit(), result.SenderHookCall.GetEvmHookCall().GetGasLimit())

	// Verify receiver hook
	require.NotNil(t, result.ReceiverHookCall)
	assert.Equal(t, original.ReceiverHookCall.hookType, result.ReceiverHookCall.hookType)
	assert.Equal(t, original.ReceiverHookCall.GetHookId(), result.ReceiverHookCall.GetHookId())
	assert.Equal(t, original.ReceiverHookCall.GetEvmHookCall().GetData(), result.ReceiverHookCall.GetEvmHookCall().GetData())
	assert.Equal(t, original.ReceiverHookCall.GetEvmHookCall().GetGasLimit(), result.ReceiverHookCall.GetEvmHookCall().GetGasLimit())

	// Verify they are different
	assert.NotEqual(t, result.SenderHookCall.GetHookId(), result.ReceiverHookCall.GetHookId())
}
