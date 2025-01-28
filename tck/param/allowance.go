package param

// SPDX-License-Identifier: Apache-2.0

type AllowanceParams struct {
	OwnerAccountId            *string               `json:"ownerAccountId,omitempty"`
	SpenderAccountId          *string               `json:"spenderAccountId,omitempty"`
	ReceiverSignatureRequired *bool                 `json:"receiverSignatureRequired,omitempty"`
	Hbar                      *HbarAllowanceParams  `json:"hbar,omitempty"`
	Token                     *TokenAllowanceParams `json:"token,omitempty"`
	Nft                       *NftAllowanceParams   `json:"nft,omitempty"`
}

type HbarAllowanceParams struct {
	Amount *string `json:"amount,omitempty"`
}

type TokenAllowanceParams struct {
	TokenId *string `json:"tokenId,omitempty"`
	Amount  *string `json:"amount,omitempty"`
}

type NftAllowanceParams struct {
	TokenId                  *string   `json:"tokenId,omitempty"`
	SerialNumbers            *[]string `json:"serialNumbers,omitempty"`
	ApprovedForAll           *bool     `json:"approvedForAll,omitempty"`
	DelegateSpenderAccountId *string   `json:"delegateSpenderAccountId,omitempty"`
}
