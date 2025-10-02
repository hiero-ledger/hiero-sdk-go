package response

// SPDX-License-Identifier: Apache-2.0

type AccountResponse struct {
	AccountId string `json:"accountId"`
	Status    string `json:"status"`
}

type AccountBalanceResponse struct {
	Hbar          string            `json:"hbars"`
	TokenBalances map[string]uint64 `json:"tokenBalances,omitempty"`
	TokenDecimals map[string]uint64 `json:"tokenDecimals,omitempty"`
}
