package response

// SPDX-License-Identifier: Apache-2.0

type TopicResponse struct {
	TopicId string `json:"topicId"`
	Status  string `json:"status"`
}

type GetTopicInfoResponse struct {
	TopicId            string               `json:"topicId"`
	TopicMemo          string               `json:"topicMemo"`
	SequenceNumber     string               `json:"sequenceNumber"`
	RunningHash        string               `json:"runningHash"`
	AdminKey           *string              `json:"adminKey"`
	SubmitKey          *string              `json:"submitKey"`
	AutoRenewAccountId *string              `json:"autoRenewAccountId"`
	AutoRenewPeriod    string               `json:"autoRenewPeriod"`
	ExpirationTime     string               `json:"expirationTime"`
	FeeScheduleKey     *string              `json:"feeScheduleKey"`
	FeeExemptKeys      *[]string            `json:"feeExemptKeys"`
	CustomFees         *[]CustomFeeResponse `json:"customFees"`
	LedgerId           string               `json:"ledgerId"`
}

type CustomFeeResponse struct {
	FeeCollectorAccountId  *string          `json:"feeCollectorAccountId"`
	AllCollectorsAreExempt bool             `json:"allCollectorsAreExempt"`
	FixedFee               FixedFeeResponse `json:"fixedFee"`
}

type FixedFeeResponse struct {
	Amount              string  `json:"amount"`
	DenominatingTokenId *string `json:"denominatingTokenId"`
}
