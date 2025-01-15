package response

// SPDX-License-Identifier: Apache-2.0

type TokenResponse struct {
	TokenId string `json:"tokenId"`
	Status  string `json:"status"`
}

type TokenMintResponse struct {
	TokenId        *string   `json:"tokenId,omitempty"`
	NewTotalSupply *string   `json:"newTotalSupply,omitempty"`
	SerialNumbers  *[]string `json:"serialNumbers,omitempty"`
	Status         *string   `json:"status,omitempty"`
}

type TokenBurnResponse struct {
	TokenId        *string `json:"tokenId,omitempty"`
	NewTotalSupply *string `json:"newTotalSupply,omitempty"`
	Status         *string `json:"status,omitempty"`
}
