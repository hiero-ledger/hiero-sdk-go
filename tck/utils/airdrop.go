package utils

// SPDX-License-Identifier: Apache-2.0

import (
	"strconv"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// Handles a single transfer parameter and adds it to the transaction
func HandleAirdropParam(transaction *hiero.TokenAirdropTransaction, transferParam param.TransferParams) error {
	switch {
	case transferParam.Token != nil:
		return HandleAirdropTokenTransfer(transaction, transferParam)
	case transferParam.Nft != nil:
		return HandleAirdropNftTransfer(transaction, transferParam)
	default:
		return response.NewInternalError("Invalid transfer parameter")
	}
}

// Handles a Token transfer parameter
func HandleAirdropTokenTransfer(transaction *hiero.TokenAirdropTransaction, transferParam param.TransferParams) error {
	token := transferParam.Token

	accountId, err := hiero.AccountIDFromString(*token.AccountId)
	if err != nil {
		return err
	}

	tokenId, err := hiero.TokenIDFromString(*token.TokenId)
	if err != nil {
		return err
	}

	amount, err := strconv.ParseInt(*token.Amount, 10, 64)
	if err != nil {
		return err
	}

	if token.Decimals != nil {
		if transferParam.Approved != nil && *transferParam.Approved {
			transaction.AddApprovedTokenTransferWithDecimals(tokenId, accountId, amount, uint32(*token.Decimals), true)
		} else {
			transaction.AddTokenTransferWithDecimals(tokenId, accountId, amount, uint32(*token.Decimals))
		}
	} else {
		if transferParam.Approved != nil && *transferParam.Approved {
			transaction.AddApprovedTokenTransfer(tokenId, accountId, amount, true)
		} else {
			transaction.AddTokenTransfer(tokenId, accountId, amount)
		}
	}

	return nil
}

// Handles an NFT transfer parameter
func HandleAirdropNftTransfer(transaction *hiero.TokenAirdropTransaction, transferParam param.TransferParams) error {
	nft := transferParam.Nft

	senderAccountId, err := hiero.AccountIDFromString(*nft.SenderAccountId)
	if err != nil {
		return err
	}

	receiverAccountId, err := hiero.AccountIDFromString(*nft.ReceiverAccountId)
	if err != nil {
		return err
	}

	serialNumberParsed, err := strconv.ParseInt(*nft.SerialNumber, 10, 64)
	if err != nil {
		return err
	}

	tokenId, err := hiero.TokenIDFromString(*nft.TokenId)
	if err != nil {
		return err
	}

	nftId := hiero.NftID{
		TokenID:      tokenId,
		SerialNumber: serialNumberParsed,
	}

	if transferParam.Approved != nil && *transferParam.Approved {
		transaction.AddApprovedNftTransfer(nftId, senderAccountId, receiverAccountId, true)
	} else {
		transaction.AddNftTransfer(nftId, senderAccountId, receiverAccountId)
	}

	return nil
}
