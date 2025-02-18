package hiero

import "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// SPDX-License-Identifier: Apache-2.0

// A maximum custom fee that the user is willing to pay.
// See HIP-991 for more information hips.hedera.com/hip/hip-991
type CustomFeeLimit struct {
	PayerId    *AccountID
	CustomFees []*CustomFixedFee
}

func NewCustomFeeLimit() *CustomFeeLimit {
	return &CustomFeeLimit{}
}

func customFeeLimitFromProtobuf(customFeeLimit *services.CustomFeeLimit) *CustomFeeLimit {
	if customFeeLimit == nil {
		return &CustomFeeLimit{}
	}

	var payerId *AccountID
	if customFeeLimit.AccountId != nil {
		payerId = _AccountIDFromProtobuf(customFeeLimit.AccountId)
	}

	var customFees []*CustomFixedFee
	for _, customFee := range customFeeLimit.Fees {
		customFixedFee := _CustomFixedFeeFromProtobuf(customFee, CustomFee{})
		customFees = append(customFees, customFixedFee)
	}

	return &CustomFeeLimit{
		PayerId:    payerId,
		CustomFees: customFees,
	}
}

// SetPayerId sets the account ID of the payer.
func (feeLimit *CustomFeeLimit) SetPayerId(payerId AccountID) *CustomFeeLimit {
	feeLimit.PayerId = &payerId
	return feeLimit
}

// GetPayerId returns the account ID of the payer.
func (feeLimit *CustomFeeLimit) GetPayerId() AccountID {
	return *feeLimit.PayerId
}

// SetCustomFees sets the custom fees.
func (feeLimit *CustomFeeLimit) SetCustomFees(customFees []*CustomFixedFee) *CustomFeeLimit {
	feeLimit.CustomFees = customFees
	return feeLimit
}

// AddCustomFee adds a custom fee.
func (feeLimit *CustomFeeLimit) AddCustomFee(customFee *CustomFixedFee) *CustomFeeLimit {
	feeLimit.CustomFees = append(feeLimit.CustomFees, customFee)
	return feeLimit
}

// GetCustomFees returns the custom fees.
func (feeLimit *CustomFeeLimit) GetCustomFees() []*CustomFixedFee {
	return feeLimit.CustomFees
}

func (feeLimit *CustomFeeLimit) toProtobuf() *services.CustomFeeLimit {
	var fees []*services.FixedFee
	for _, customFee := range feeLimit.CustomFees {
		fees = append(fees, customFee._ToProtobuf().GetFixedFee())
	}

	var payerId *services.AccountID
	if feeLimit.PayerId != nil {
		payerId = feeLimit.PayerId._ToProtobuf()
	}

	return &services.CustomFeeLimit{
		AccountId: payerId,
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
