package response

// SPDX-License-Identifier: Apache-2.0

// ScheduleResponse represents the response from schedule-related operations
type ScheduleResponse struct {
	ScheduleId    string `json:"scheduleId,omitempty"`
	TransactionId string `json:"transactionId,omitempty"`
	Status        string `json:"status"`
}
