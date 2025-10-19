package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// _TokenNftTransfer is the information about a NFT transfer
type _TokenNftTransfer struct {
	SenderAccountID   AccountID
	ReceiverAccountID AccountID
	SerialNumber      int64
	IsApproved        bool
	SenderHookCall    *NftHookCall
	ReceiverHookCall  *NftHookCall
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

	senderHookCall := nftHookCallFromProtobuf(pb)
	receiverHookCall := nftHookCallFromProtobuf(pb)
	return _TokenNftTransfer{
		SenderAccountID:   senderAccountID,
		ReceiverAccountID: receiverAccountID,
		SerialNumber:      pb.SerialNumber,
		IsApproved:        pb.IsApproval,
		SenderHookCall:    senderHookCall,
		ReceiverHookCall:  receiverHookCall,
	}
}

func (transfer *_TokenNftTransfer) _ToProtobuf() *services.NftTransfer {
	pbBody := &services.NftTransfer{
		SenderAccountID:   transfer.SenderAccountID._ToProtobuf(),
		ReceiverAccountID: transfer.ReceiverAccountID._ToProtobuf(),
		SerialNumber:      transfer.SerialNumber,
		IsApproval:        transfer.IsApproved,
	}

	if transfer.SenderHookCall != nil {
		switch transfer.SenderHookCall.hookType {
		case PRE_HOOK_SENDER:
			pbBody.SenderAllowanceHookCall = &services.NftTransfer_PreTxSenderAllowanceHook{
				PreTxSenderAllowanceHook: transfer.SenderHookCall.toProtobuf(),
			}
		case PRE_POST_HOOK_SENDER:
			pbBody.SenderAllowanceHookCall = &services.NftTransfer_PrePostTxSenderAllowanceHook{
				PrePostTxSenderAllowanceHook: transfer.SenderHookCall.toProtobuf(),
			}
		}
	}

	if transfer.ReceiverHookCall != nil {
		switch transfer.ReceiverHookCall.hookType {
		case PRE_HOOK_RECEIVER:
			pbBody.ReceiverAllowanceHookCall = &services.NftTransfer_PreTxReceiverAllowanceHook{
				PreTxReceiverAllowanceHook: transfer.ReceiverHookCall.toProtobuf(),
			}
		case PRE_POST_HOOK_RECEIVER:
			pbBody.ReceiverAllowanceHookCall = &services.NftTransfer_PrePostTxReceiverAllowanceHook{
				PrePostTxReceiverAllowanceHook: transfer.ReceiverHookCall.toProtobuf(),
			}
		}
	}

	return pbBody
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
