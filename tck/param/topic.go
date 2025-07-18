package param

// SPDX-License-Identifier: Apache-2.0
type CreateTopicParams struct {
	Memo                    *string                  `json:"memo"`
	AdminKey                *string                  `json:"adminKey"`
	SubmitKey               *string                  `json:"submitKey"`
	AutoRenewPeriod         *string                  `json:"autoRenewPeriod"`
	AutoRenewAccountId      *string                  `json:"autoRenewAccountId"`
	FeeScheduleKey          *string                  `json:"feeScheduleKey"`
	FeeExemptKeys           *[]string                `json:"feeExemptKeys"`
	CustomFees              *[]CustomFee             `json:"customFees"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}
