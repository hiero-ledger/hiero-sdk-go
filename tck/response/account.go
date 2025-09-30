package response

// SPDX-License-Identifier: Apache-2.0

type AccountResponse struct {
	AccountId string `json:"accountId"`
	Status    string `json:"status"`
}

type AccountBalanceResponse struct {
	HbarBalance string            `json:"hbarBalance"`
	Tokens      map[string]uint64 `json:"tokens,omitempty"`
}
