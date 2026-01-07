package param

// SPDX-License-Identifier: Apache-2.0

type TransferCryptoParams struct {
	BaseParams
	Transfers *[]TransferParams `json:"transfers,omitempty"`
}

type TransferParams struct {
	Hbar     *HbarTransferParams  `json:"hbar,omitempty"`
	Token    *TokenTransferParams `json:"token,omitempty"`
	Nft      *NftTransferParams   `json:"nft,omitempty"`
	Approved *bool                `json:"approved,omitempty"`
}

type HbarTransferParams struct {
	AccountId  *string `json:"accountId,omitempty"`
	EvmAddress *string `json:"evmAddress,omitempty"`
	Amount     *string `json:"amount,omitempty"`
}

type NftTransferParams struct {
	SenderAccountId   *string `json:"senderAccountId,omitempty"`
	ReceiverAccountId *string `json:"receiverAccountId,omitempty"`
	TokenId           *string `json:"tokenId,omitempty"`
	SerialNumber      *string `json:"serialNumber,omitempty"`
}

type TokenTransferParams struct {
	AccountId *string `json:"accountId,omitempty"`
	TokenId   *string `json:"tokenId,omitempty"`
	Amount    *string `json:"amount,omitempty"`
	Decimals  *int    `json:"decimals,omitempty"`
}
