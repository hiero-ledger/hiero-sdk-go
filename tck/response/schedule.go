package response

// SPDX-License-Identifier: Apache-2.0

// ScheduleResponse represents the response from schedule-related operations
type ScheduleResponse struct {
	ScheduleId    string `json:"scheduleId,omitempty"`
	TransactionId string `json:"transactionId,omitempty"`
	Status        string `json:"status"`
}

// ScheduleInfoResponse represent the response from scheduleInfo operation
type ScheduleInfoResponse struct {
	ScheduleId             string   `json:"scheduleId"`
	CreatorAccountId       string   `json:"creatorAccountId"`
	PayerAccountId         string   `json:"payerAccountId"`
	AdminKey               string   `json:"adminKey,omitempty"`
	Signers                []string `json:"signers"`
	ScheduleMemo           string   `json:"scheduleMemo"`
	ExpirationTime         string   `json:"expirationTime"`
	ExecutedAt             string   `json:"executedAt,omitempty"`
	DeletedAt              string   `json:"deletedAt,omitempty"`
	ScheduledTransactionId string   `json:"scheduledTransactionId"`
	WaitForExpiry          bool     `json:"waitForExpiry"`
	Cost                   string   `json:"cost,omitempty"`
}
