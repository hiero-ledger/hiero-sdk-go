package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

type CustomFeeLimit struct {
	PayerId    *AccountID
	CustomFees []*CustomFixedFee
}

func NewCustomFeeLimit() *CustomFeeLimit {
	return &CustomFeeLimit{}
}

func customFeeLimitFromProtobuf(customFeeLimit *services.CustomFeeLimit) CustomFeeLimit {
	if customFeeLimit == nil {
		return CustomFeeLimit{}
	}

	var payerId *AccountID
	if customFeeLimit.AccountId != nil {
		payerId = _AccountIDFromProtobuf(customFeeLimit.AccountId)
	}

	var customFees []*CustomFixedFee
	for _, customFee := range customFeeLimit.Fees {
		customFixedFee := _CustomFixedFeeFromProtobuf(customFee, CustomFee{})
		customFees = append(customFees, &customFixedFee)
	}

	return CustomFeeLimit{
		PayerId:    payerId,
		CustomFees: customFees,
	}
}

func (feeLimit *CustomFeeLimit) SetPayerId(payerId AccountID) *CustomFeeLimit {
	feeLimit.PayerId = &payerId
	return feeLimit
}

func (feeLimit *CustomFeeLimit) GetPayerId() AccountID {
	return *feeLimit.PayerId
}

func (feeLimit *CustomFeeLimit) SetCustomFees(customFees []*CustomFixedFee) *CustomFeeLimit {
	feeLimit.CustomFees = customFees
	return feeLimit
}

func (feeLimit *CustomFeeLimit) AddCustomFee(customFee *CustomFixedFee) *CustomFeeLimit {
	feeLimit.CustomFees = append(feeLimit.CustomFees, customFee)
	return feeLimit
}

func (feeLimit *CustomFeeLimit) GetCustomFees() []*CustomFixedFee {
	return feeLimit.CustomFees
}

func (feeLimit *CustomFeeLimit) toProtobuf() *services.CustomFeeLimit {
	var fees []*services.FixedFee
	for _, customFee := range feeLimit.CustomFees {
		fees = append(fees, customFee._ToProtobuf().GetFixedFee())
	}

	return &services.CustomFeeLimit{
		AccountId: feeLimit.PayerId._ToProtobuf(),
		Fees:      fees,
	}
}

func (feeLimit *CustomFeeLimit) String() string {
	customFeesStr := "["
	for i, fee := range feeLimit.CustomFees {
		if i > 0 {
			customFeesStr += ", "
		}
		customFeesStr += fee.String()
	}
	customFeesStr += "]"
	return "CustomFeeLimit{PayerId: " + feeLimit.PayerId.String() + ", CustomFees: " + customFeesStr + "}"
}
