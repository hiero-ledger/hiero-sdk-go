package param

// SPDX-License-Identifier: Apache-2.0

type CreateTokenParams struct {
	BaseParams
	Name               *string      `json:"name,omitempty"`
	Symbol             *string      `json:"symbol,omitempty"`
	Decimals           *int         `json:"decimals,omitempty"`
	InitialSupply      *string      `json:"initialSupply,omitempty"`
	TreasuryAccountId  *string      `json:"treasuryAccountId,omitempty"`
	AdminKey           *string      `json:"adminKey,omitempty"`
	KycKey             *string      `json:"kycKey,omitempty"`
	FreezeKey          *string      `json:"freezeKey,omitempty"`
	WipeKey            *string      `json:"wipeKey,omitempty"`
	SupplyKey          *string      `json:"supplyKey,omitempty"`
	FeeScheduleKey     *string      `json:"feeScheduleKey,omitempty"`
	PauseKey           *string      `json:"pauseKey,omitempty"`
	MetadataKey        *string      `json:"metadataKey,omitempty"`
	FreezeDefault      *bool        `json:"freezeDefault,omitempty"`
	ExpirationTime     *string      `json:"expirationTime,omitempty"`
	AutoRenewAccountId *string      `json:"autoRenewAccountId,omitempty"`
	AutoRenewPeriod    *string      `json:"autoRenewPeriod,omitempty"`
	Memo               *string      `json:"memo,omitempty"`
	TokenType          *string      `json:"tokenType,omitempty"`
	SupplyType         *string      `json:"supplyType,omitempty"`
	MaxSupply          *string      `json:"maxSupply,omitempty"`
	CustomFees         *[]CustomFee `json:"customFees,omitempty"`
	Metadata           *string      `json:"metadata,omitempty"`
}

type UpdateTokenParams struct {
	BaseParams
	TokenId            *string `json:"tokenId,omitempty"`
	Name               *string `json:"name,omitempty"`
	Symbol             *string `json:"symbol,omitempty"`
	TreasuryAccountId  *string `json:"treasuryAccountId,omitempty"`
	AdminKey           *string `json:"adminKey,omitempty"`
	KycKey             *string `json:"kycKey,omitempty"`
	FreezeKey          *string `json:"freezeKey,omitempty"`
	WipeKey            *string `json:"wipeKey,omitempty"`
	SupplyKey          *string `json:"supplyKey,omitempty"`
	FeeScheduleKey     *string `json:"feeScheduleKey,omitempty"`
	PauseKey           *string `json:"pauseKey,omitempty"`
	MetadataKey        *string `json:"metadataKey,omitempty"`
	ExpirationTime     *string `json:"expirationTime,omitempty"`
	AutoRenewAccountId *string `json:"autoRenewAccountId,omitempty"`
	AutoRenewPeriod    *string `json:"autoRenewPeriod,omitempty"`
	Memo               *string `json:"memo,omitempty"`
	Metadata           *string `json:"metadata,omitempty"`
}

type DeleteTokenParams struct {
	BaseParams
	TokenId *string `json:"tokenId,omitempty"`
}

type UpdateTokenFeeScheduleParams struct {
	BaseParams
	TokenId    *string      `json:"tokenId,omitempty"`
	CustomFees *[]CustomFee `json:"customFees,omitempty"`
}

type AssociateDissociatesTokenParams struct {
	BaseParams
	AccountId *string   `json:"accountId,omitempty"`
	TokenIds  *[]string `json:"tokenIds,omitempty"`
}

type PauseUnPauseTokenParams struct {
	BaseParams
	TokenId *string `json:"tokenId,omitempty"`
}

type FreezeUnFreezeTokenParams struct {
	BaseParams
	AccountId *string `json:"accountId,omitempty"`
	TokenId   *string `json:"tokenId,omitempty"`
}

type GrantRevokeTokenKycParams struct {
	BaseParams
	AccountId *string `json:"accountId,omitempty"`
	TokenId   *string `json:"tokenId,omitempty"`
}

type MintTokenParams struct {
	BaseParams
	TokenId  *string   `json:"tokenId,omitempty"`
	Amount   *string   `json:"amount,omitempty"`
	Metadata *[]string `json:"metadata,omitempty"`
}

type BurnTokenParams struct {
	BaseParams
	TokenId       *string   `json:"tokenId,omitempty"`
	Amount        *string   `json:"amount,omitempty"`
	SerialNumbers *[]string `json:"serialNumbers,omitempty"`
}

type WipeTokenParams struct {
	BaseParams
	TokenId       *string   `json:"tokenId,omitempty"`
	AccountId     *string   `json:"accountId,omitempty"`
	Amount        *string   `json:"amount,omitempty"`
	SerialNumbers *[]string `json:"serialNumbers,omitempty"`
}

type AirdropParams struct {
	BaseParams
	TokenTransfers *[]TransferParams `json:"tokenTransfers,omitempty"`
}

type ClaimTokenParams struct {
	BaseParams
	SenderAccountId   *string   `json:"senderAccountId,omitempty"`
	ReceiverAccountId *string   `json:"receiverAccountId,omitempty"`
	TokenId           *string   `json:"tokenId,omitempty"`
	SerialNumbers     *[]string `json:"serialNumbers,omitempty"`
}

type PendingAirdropParams struct {
	SenderAccountId   *string   `json:"senderAccountId,omitempty"`
	ReceiverAccountId *string   `json:"receiverAccountId,omitempty"`
	TokenId           *string   `json:"tokenId,omitempty"`
	SerialNumbers     *[]string `json:"serialNumbers,omitempty"`
}

type AirdropCancelTokenParams struct {
	BaseParams
	PendingAirdrops *[]PendingAirdropParams `json:"pendingAirdrops,omitempty"`
}

type RejectTokenParams struct {
	BaseParams
	OwnerId       *string   `json:"ownerId,omitempty"`
	TokenIds      *[]string `json:"tokenIds,omitempty"`
	SerialNumbers *[]string `json:"serialNumbers,omitempty"`
}
