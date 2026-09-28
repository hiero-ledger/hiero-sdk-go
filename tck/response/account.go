package response

// SPDX-License-Identifier: Apache-2.0

type AccountResponse struct {
	AccountId string `json:"accountId"`
	Status    string `json:"status"`
}

type MirrorNodeAccountBalanceResponse struct {
	Hbar string `json:"hbars"`
}

type DeprecatedAccountBalanceQueryResponse struct {
	ConstructionWarning *string `json:"constructionWarning"`
	ExecutionError      *string `json:"executionError"`
}
