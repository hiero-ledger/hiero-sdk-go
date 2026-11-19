package response

// SPDX-License-Identifier: Apache-2.0

type TopicResponse struct {
	TopicId string `json:"topicId"`
	Status  string `json:"status"`
}

type TopicInfoResponse struct {
	TopicId            string               `json:"topicId"`
	TopicMemo          string               `json:"topicMemo"`
	SequenceNumber     string               `json:"sequenceNumber"`
	RunningHash        string               `json:"runningHash"`
	AdminKey           *string              `json:"adminKey,omitempty"`
	SubmitKey          *string              `json:"submitKey,omitempty"`
	AutoRenewAccountId *string              `json:"autoRenewAccountId,omitempty"`
	AutoRenewPeriod    string               `json:"autoRenewPeriod"`
	ExpirationTime     string               `json:"expirationTime"`
	FeeScheduleKey     *string              `json:"feeScheduleKey,omitempty"`
	FeeExemptKeys      *[]string            `json:"feeExemptKeys,omitempty"`
	CustomFees         *[]CustomFeeResponse `json:"customFees,omitempty"`
	LedgerId           string               `json:"ledgerId"`
}

type CustomFeeResponse struct {
	FeeCollectorAccountId  *string          `json:"feeCollectorAccountId,omitempty"`
	AllCollectorsAreExempt bool             `json:"allCollectorsAreExempt"`
	FixedFee               FixedFeeResponse `json:"fixedFee"`
}

type FixedFeeResponse struct {
	Amount              string  `json:"amount"`
	DenominatingTokenId *string `json:"denominatingTokenId,omitempty"`
}
