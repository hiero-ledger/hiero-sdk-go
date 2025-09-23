package hiero

// SPDX-License-Identifier: Apache-2.0

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

type _HbarTransfer struct {
	accountID                       *AccountID
	Amount                          Hbar
	IsApproved                      bool
	PreTransactionAllowanceHook     *HookCall
	PrePostTransactionAllowanceHook *HookCall
}

func _HbarTransferFromProtobuf(pb []*services.AccountAmount) []*_HbarTransfer {
	result := make([]*_HbarTransfer, 0)
	for _, acc := range pb {
		transfer := &_HbarTransfer{
			accountID:  _AccountIDFromProtobuf(acc.AccountID),
			Amount:     HbarFromTinybar(acc.Amount),
			IsApproved: acc.GetIsApproval(),
		}
		if acc.GetPreTxAllowanceHook() != nil {
			hookCall := hookCallFromProtobuf(acc.GetPreTxAllowanceHook())
			transfer.PreTransactionAllowanceHook = &hookCall
		}
		if acc.GetPrePostTxAllowanceHook() != nil {
			hookCall := hookCallFromProtobuf(acc.GetPrePostTxAllowanceHook())
			transfer.PrePostTransactionAllowanceHook = &hookCall
		}
		result = append(result, transfer)
	}

	return result
}

func (transfer *_HbarTransfer) _ToProtobuf() *services.AccountAmount { //nolint
	var account *services.AccountID
	if transfer.accountID != nil {
		account = transfer.accountID._ToProtobuf()
	}

	protoBody := &services.AccountAmount{
		AccountID:  account,
		Amount:     transfer.Amount.AsTinybar(),
		IsApproval: transfer.IsApproved,
	}

	if transfer.PreTransactionAllowanceHook != nil {
		protoBody.HookCall = &services.AccountAmount_PreTxAllowanceHook{
			PreTxAllowanceHook: transfer.PreTransactionAllowanceHook.toProtobuf(),
		}
	}
	if transfer.PrePostTransactionAllowanceHook != nil {
		protoBody.HookCall = &services.AccountAmount_PrePostTxAllowanceHook{
			PrePostTxAllowanceHook: transfer.PrePostTransactionAllowanceHook.toProtobuf(),
		}
	}

	return protoBody
}
