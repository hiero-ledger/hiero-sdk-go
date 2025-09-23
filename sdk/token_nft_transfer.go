package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// _TokenNftTransfer is the information about a NFT transfer
type _TokenNftTransfer struct {
	SenderAccountID                         AccountID
	ReceiverAccountID                       AccountID
	SerialNumber                            int64
	IsApproved                              bool
	PreTransactionSenderAllowanceHook       *HookCall
	PrePostTransactionSenderAllowanceHook   *HookCall
	PreTransactionReceiverAllowanceHook     *HookCall
	PrePostTransactionReceiverAllowanceHook *HookCall
}

func _NftTransferFromProtobuf(pb *services.NftTransfer) _TokenNftTransfer {
	if pb == nil {
		return _TokenNftTransfer{}
	}

	senderAccountID := AccountID{}
	if pb.SenderAccountID != nil {
		senderAccountID = *_AccountIDFromProtobuf(pb.SenderAccountID)
	}

	receiverAccountID := AccountID{}
	if pb.ReceiverAccountID != nil {
		receiverAccountID = *_AccountIDFromProtobuf(pb.ReceiverAccountID)
	}

	var preTransactionSenderAllowanceHook *HookCall
	if pb.GetPreTxSenderAllowanceHook() != nil {
		preTransactionSenderAllowanceHookValue := hookCallFromProtobuf(pb.GetPreTxSenderAllowanceHook())
		preTransactionSenderAllowanceHook = &preTransactionSenderAllowanceHookValue
	}

	var prePostTransactionSenderAllowanceHook *HookCall
	if pb.GetPrePostTxSenderAllowanceHook() != nil {
		prePostTransactionSenderAllowanceHookValue := hookCallFromProtobuf(pb.GetPrePostTxSenderAllowanceHook())
		prePostTransactionSenderAllowanceHook = &prePostTransactionSenderAllowanceHookValue
	}

	var preTransactionReceiverAllowanceHook *HookCall
	if pb.GetPreTxReceiverAllowanceHook() != nil {
		preTransactionReceiverAllowanceHookValue := hookCallFromProtobuf(pb.GetPreTxReceiverAllowanceHook())
		preTransactionReceiverAllowanceHook = &preTransactionReceiverAllowanceHookValue
	}

	var prePostTransactionReceiverAllowanceHook *HookCall
	if pb.GetPrePostTxReceiverAllowanceHook() != nil {
		prePostTransactionReceiverAllowanceHookValue := hookCallFromProtobuf(pb.GetPrePostTxReceiverAllowanceHook())
		prePostTransactionReceiverAllowanceHook = &prePostTransactionReceiverAllowanceHookValue
	}

	return _TokenNftTransfer{
		SenderAccountID:                         senderAccountID,
		ReceiverAccountID:                       receiverAccountID,
		SerialNumber:                            pb.SerialNumber,
		IsApproved:                              pb.IsApproval,
		PreTransactionSenderAllowanceHook:       preTransactionSenderAllowanceHook,
		PrePostTransactionSenderAllowanceHook:   prePostTransactionSenderAllowanceHook,
		PreTransactionReceiverAllowanceHook:     preTransactionReceiverAllowanceHook,
		PrePostTransactionReceiverAllowanceHook: prePostTransactionReceiverAllowanceHook,
	}
}

func (transfer *_TokenNftTransfer) _ToProtobuf() *services.NftTransfer {
	protoBody := &services.NftTransfer{
		SenderAccountID:   transfer.SenderAccountID._ToProtobuf(),
		ReceiverAccountID: transfer.ReceiverAccountID._ToProtobuf(),
		SerialNumber:      transfer.SerialNumber,
		IsApproval:        transfer.IsApproved,
	}

	if transfer.PreTransactionSenderAllowanceHook != nil {
		protoBody.SenderAllowanceHookCall = &services.NftTransfer_PreTxSenderAllowanceHook{
			PreTxSenderAllowanceHook: transfer.PreTransactionSenderAllowanceHook.toProtobuf(),
		}
	}
	if transfer.PrePostTransactionSenderAllowanceHook != nil {
		protoBody.SenderAllowanceHookCall = &services.NftTransfer_PrePostTxSenderAllowanceHook{
			PrePostTxSenderAllowanceHook: transfer.PrePostTransactionSenderAllowanceHook.toProtobuf(),
		}
	}

	if transfer.PreTransactionReceiverAllowanceHook != nil {
		protoBody.ReceiverAllowanceHookCall = &services.NftTransfer_PreTxReceiverAllowanceHook{
			PreTxReceiverAllowanceHook: transfer.PreTransactionReceiverAllowanceHook.toProtobuf(),
		}
	}
	if transfer.PrePostTransactionReceiverAllowanceHook != nil {
		protoBody.ReceiverAllowanceHookCall = &services.NftTransfer_PrePostTxReceiverAllowanceHook{
			PrePostTxReceiverAllowanceHook: transfer.PrePostTransactionReceiverAllowanceHook.toProtobuf(),
		}
	}

	return protoBody
}

// ToBytes returns the byte representation of the TokenNftTransfer
func (transfer _TokenNftTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenNftTransfersFromBytes returns the TokenNftTransfer from a raw protobuf bytes representation
func NftTransferFromBytes(data []byte) (_TokenNftTransfer, error) {
	if data == nil {
		return _TokenNftTransfer{}, errByteArrayNull
	}
	pb := services.NftTransfer{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return _TokenNftTransfer{}, err
	}

	return _NftTransferFromProtobuf(&pb), nil
}
