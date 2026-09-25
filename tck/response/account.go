package response

// SPDX-License-Identifier: Apache-2.0

type AccountResponse struct {
	AccountId string `json:"accountId"`
	Status    string `json:"status"`
}

type AccountBalanceResponse struct {
	Hbar string `json:"hbars"`
}

// DeprecatedAccountBalanceQueryResponse fields are null when the SDK logged no warning or returned no error
type DeprecatedAccountBalanceQueryResponse struct {
	ConstructionWarning *string `json:"constructionWarning"`
	ExecutionError      *string `json:"executionError"`
}
