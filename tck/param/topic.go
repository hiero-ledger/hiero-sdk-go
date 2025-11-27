package param

// SPDX-License-Identifier: Apache-2.0
type CreateTopicParams struct {
	BaseParams
	Memo               *string      `json:"memo"`
	AdminKey           *string      `json:"adminKey"`
	SubmitKey          *string      `json:"submitKey"`
	AutoRenewPeriod    *string      `json:"autoRenewPeriod"`
	AutoRenewAccountId *string      `json:"autoRenewAccountId"`
	FeeScheduleKey     *string      `json:"feeScheduleKey"`
	FeeExemptKeys      *[]string    `json:"feeExemptKeys"`
	CustomFees         *[]CustomFee `json:"customFees"`
}

type UpdateTopicParams struct {
	BaseParams
	TopicId            *string      `json:"topicId"`
	Memo               *string      `json:"memo"`
	AdminKey           *string      `json:"adminKey"`
	SubmitKey          *string      `json:"submitKey"`
	FeeScheduleKey     *string      `json:"feeScheduleKey"`
	FeeExemptKeys      *[]string    `json:"feeExemptKeys"`
	CustomFees         *[]CustomFee `json:"customFees"`
	AutoRenewPeriod    *string      `json:"autoRenewPeriod"`
	AutoRenewAccountId *string      `json:"autoRenewAccountId"`
	ExpirationTime     *string      `json:"expirationTime"`
}

type DeleteTopicParams struct {
	BaseParams
	TopicId *string `json:"topicId"`
}

type SubmitTopicMessageParams struct {
	BaseParams
	TopicId         *string           `json:"topicId"`
	Message         *string           `json:"message"`
	MaxChunks       *int64            `json:"maxChunks"`
	ChunkSize       *int64            `json:"chunkSize"`
	CustomFeeLimits *[]CustomFeeLimit `json:"customFeeLimits"`
}

type CustomFeeLimit struct {
	PayerId   *string     `json:"payerId"`
	FixedFees *[]FixedFee `json:"fixedFees"`
}
