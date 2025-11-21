package hiero

// SPDX-License-Identifier: Apache-2.0

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

type _HbarTransfer struct {
	accountID  *AccountID
	amount     Hbar
	isApproved bool
	hookCall   *FungibleHookCall
}

func _HbarTransferFromProtobuf(pb []*services.AccountAmount) []*_HbarTransfer {
	result := make([]*_HbarTransfer, 0)
	for _, acc := range pb {
		result = append(result, &_HbarTransfer{
			accountID:  _AccountIDFromProtobuf(acc.AccountID),
			amount:     HbarFromTinybar(acc.Amount),
			isApproved: acc.GetIsApproval(),
			hookCall:   fungibleHookCallFromProtobuf(acc),
		})
	}

	return result
}

func (transfer *_HbarTransfer) _ToProtobuf() *services.AccountAmount { //nolint
	var account *services.AccountID
	if transfer.accountID != nil {
		account = transfer.accountID._ToProtobuf()
	}

	pbBody := &services.AccountAmount{
		AccountID:  account,
		Amount:     transfer.amount.AsTinybar(),
		IsApproval: transfer.isApproved,
	}

	if transfer.hookCall != nil {
		switch transfer.hookCall.hookType {
		case PRE_HOOK:
			pbBody.HookCall = &services.AccountAmount_PreTxAllowanceHook{
				PreTxAllowanceHook: transfer.hookCall.toProtobuf(),
			}
		case PRE_POST_HOOK:
			pbBody.HookCall = &services.AccountAmount_PrePostTxAllowanceHook{
				PrePostTxAllowanceHook: transfer.hookCall.toProtobuf(),
			}
		}
	}

	return pbBody
}
