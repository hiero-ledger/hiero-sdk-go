package param

// SPDX-License-Identifier: Apache-2.0

import "encoding/json"

// ScheduledTransaction represents the transaction to be scheduled
type ScheduledTransaction struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// ScheduleCreateParams represents the parameters for creating a schedule
type ScheduleCreateParams struct {
	ScheduledTransaction    *ScheduledTransaction    `json:"scheduledTransaction,omitempty"`
	Memo                    *string                  `json:"memo,omitempty"`
	AdminKey                *string                  `json:"adminKey,omitempty"`
	PayerAccountId          *string                  `json:"payerAccountId,omitempty"`
	ExpirationTime          *string                  `json:"expirationTime,omitempty"`
	WaitForExpiry           *bool                    `json:"waitForExpiry,omitempty"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams,omitempty"`
}

// ScheduleSignParams represents the parameters for signing a schedule
type ScheduleSignParams struct {
	ScheduleId              *string                  `json:"scheduleId,omitempty"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams,omitempty"`
}
