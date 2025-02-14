package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type TokenService struct {
	sdkService *SDKService
}

func (t *TokenService) SetSdkService(service *SDKService) {
	t.sdkService = service
}

//nolint:gocyclo,gocritic // CreateToken jRPC method for createToken
func (t *TokenService) CreateToken(_ context.Context, params param.CreateTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AdminKey != nil {
		key, err := utils.GetKeyFromString(*params.AdminKey)
		if err != nil {
			return nil, err
		}
		transaction.SetAdminKey(key)
	}

	if params.KycKey != nil {
		key, err := utils.GetKeyFromString(*params.KycKey)
		if err != nil {
			return nil, err
		}
		transaction.SetKycKey(key)
	}

	if params.FreezeKey != nil {
		key, err := utils.GetKeyFromString(*params.FreezeKey)
		if err != nil {
			return nil, err
		}
		transaction.SetFreezeKey(key)
	}

	if params.WipeKey != nil {
		key, err := utils.GetKeyFromString(*params.WipeKey)
		if err != nil {
			return nil, err
		}
		transaction.SetWipeKey(key)
	}

	if params.PauseKey != nil {
		key, err := utils.GetKeyFromString(*params.PauseKey)
		if err != nil {
			return nil, err
		}
		transaction.SetPauseKey(key)
	}

	if params.MetadataKey != nil {
		key, err := utils.GetKeyFromString(*params.MetadataKey)
		if err != nil {
			return nil, err
		}
		transaction.SetMetadataKey(key)
	}

	if params.SupplyKey != nil {
		key, err := utils.GetKeyFromString(*params.SupplyKey)
		if err != nil {
			return nil, err
		}
		transaction.SetSupplyKey(key)
	}

	if params.FeeScheduleKey != nil {
		key, err := utils.GetKeyFromString(*params.FeeScheduleKey)
		if err != nil {
			return nil, err
		}
		transaction.SetFeeScheduleKey(key)
	}

	if params.Name != nil {
		transaction.SetTokenName(*params.Name)
	}
	if params.Symbol != nil {
		transaction.SetTokenSymbol(*params.Symbol)
	}
	if params.Decimals != nil {
		transaction.SetDecimals(uint(*params.Decimals))
	}
	if params.Memo != nil {
		transaction.SetTokenMemo(*params.Memo)
	}
	if params.TokenType != nil {
		if *params.TokenType == "ft" {
			transaction.SetTokenType(hiero.TokenTypeFungibleCommon)
		} else if *params.TokenType == "nft" {
			transaction.SetTokenType(hiero.TokenTypeNonFungibleUnique)
		} else {
			return nil, response.InvalidParams.WithData("Invalid token type")
		}
	}
	if params.SupplyType != nil {
		if *params.SupplyType == "finite" {
			transaction.SetSupplyType(hiero.TokenSupplyTypeFinite)
		} else if *params.SupplyType == "infinite" {
			transaction.SetSupplyType(hiero.TokenSupplyTypeInfinite)
		} else {
			return nil, response.InvalidParams.WithData("Invalid supply type")
		}
	}
	if params.MaxSupply != nil {
		maxSupply, err := strconv.ParseInt(*params.MaxSupply, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetMaxSupply(maxSupply)
	}
	if params.InitialSupply != nil {
		initialSupply, err := strconv.ParseInt(*params.InitialSupply, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetInitialSupply(uint64(initialSupply))
	}
	if params.TreasuryAccountId != nil {
		accountID, err := hiero.AccountIDFromString(*params.TreasuryAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetTreasuryAccountID(accountID)
	}
	if params.FreezeDefault != nil {
		transaction.SetFreezeDefault(*params.FreezeDefault)
	}
	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}
	if params.AutoRenewAccountId != nil {
		autoRenewAccountId, err := hiero.AccountIDFromString(*params.AutoRenewAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewAccount(autoRenewAccountId)
	}
	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}

	if params.Metadata != nil {
		transaction.SetTokenMetadata([]byte(*params.Metadata))
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	if params.CustomFees != nil {
		customFees, err := utils.ParseCustomFees(*params.CustomFees)
		if err != nil {
			return nil, err
		}
		transaction.SetCustomFees(customFees)
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{TokenId: receipt.TokenID.String(), Status: receipt.Status.String()}, nil
}

//nolint:gocyclo // UpdateToken jRPC method for updateToken
func (t *TokenService) UpdateToken(_ context.Context, params param.UpdateTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)
		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}
	if params.AdminKey != nil {
		key, err := utils.GetKeyFromString(*params.AdminKey)
		if err != nil {
			return nil, err
		}
		transaction.SetAdminKey(key)
	}

	if params.KycKey != nil {
		key, err := utils.GetKeyFromString(*params.KycKey)
		if err != nil {
			return nil, err
		}
		transaction.SetKycKey(key)
	}

	if params.FreezeKey != nil {
		key, err := utils.GetKeyFromString(*params.FreezeKey)
		if err != nil {
			return nil, err
		}
		transaction.SetFreezeKey(key)
	}

	if params.WipeKey != nil {
		key, err := utils.GetKeyFromString(*params.WipeKey)
		if err != nil {
			return nil, err
		}
		transaction.SetWipeKey(key)
	}

	if params.PauseKey != nil {
		key, err := utils.GetKeyFromString(*params.PauseKey)
		if err != nil {
			return nil, err
		}
		transaction.SetPauseKey(key)
	}

	if params.MetadataKey != nil {
		key, err := utils.GetKeyFromString(*params.MetadataKey)
		if err != nil {
			return nil, err
		}
		transaction.SetMetadataKey(key)
	}

	if params.SupplyKey != nil {
		key, err := utils.GetKeyFromString(*params.SupplyKey)
		if err != nil {
			return nil, err
		}
		transaction.SetSupplyKey(key)
	}

	if params.FeeScheduleKey != nil {
		key, err := utils.GetKeyFromString(*params.FeeScheduleKey)
		if err != nil {
			return nil, err
		}
		transaction.SetFeeScheduleKey(key)
	}

	if params.Name != nil {
		transaction.SetTokenName(*params.Name)
	}
	if params.Symbol != nil {
		transaction.SetTokenSymbol(*params.Symbol)
	}
	if params.Memo != nil {
		transaction.SetTokenMemo(*params.Memo)
	}
	if params.TreasuryAccountId != nil {
		accountID, err := hiero.AccountIDFromString(*params.TreasuryAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetTreasuryAccountID(accountID)
	}
	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}
	if params.AutoRenewAccountId != nil {
		autoRenewAccountId, err := hiero.AccountIDFromString(*params.AutoRenewAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewAccount(autoRenewAccountId)
	}
	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}

	if params.Metadata != nil {
		transaction.SetTokenMetadata([]byte(*params.Metadata))
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// DeleteToken jRPC method for deleteToken
func (t *TokenService) DeleteToken(_ context.Context, params param.DeleteTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)
	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)
		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}
	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// UpdateTokenFeeSchedule jRPC method for updateTokenFeeSchedule
func (t *TokenService) UpdateTokenFeeSchedule(_ context.Context, params param.UpdateTokenFeeScheduleParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenFeeScheduleUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)
		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CustomFees != nil {
		customFees, err := utils.ParseCustomFees(*params.CustomFees)
		if err != nil {
			return nil, err
		}
		transaction.SetCustomFees(customFees)
	}
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// AssociateToken jRPC method for associateToken
func (t *TokenService) AssociateToken(_ context.Context, params param.AssociateDissociatesTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenAssociateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.AccountId)

		if err != nil {
			return nil, err
		}
		transaction.SetAccountID(accountId)
	}

	if params.TokenIds != nil {
		// Dereference the pointer to access the slice
		tokenIds := *params.TokenIds

		// Create a slice to hold the parsed Token IDs
		var parsedTokenIds []hiero.TokenID

		// Iterate and parse each Token ID
		for _, tokenIDStr := range tokenIds {
			parsedTokenID, err := hiero.TokenIDFromString(tokenIDStr)

			if err != nil {
				return nil, err
			}

			parsedTokenIds = append(parsedTokenIds, parsedTokenID)
		}

		// Set the parsed Token IDs in the transaction
		transaction.SetTokenIDs(parsedTokenIds...)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// DisassociateToken jRPC method for dissociateToken
func (t *TokenService) DissociatesToken(_ context.Context, params param.AssociateDissociatesTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenDissociateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.AccountId)

		if err != nil {
			return nil, err
		}
		transaction.SetAccountID(accountId)
	}

	if params.TokenIds != nil {
		tokenIds := *params.TokenIds
		var parsedTokenIds []hiero.TokenID

		for _, tokenIDStr := range tokenIds {
			parsedTokenID, err := hiero.TokenIDFromString(tokenIDStr)
			if err != nil {
				return nil, err
			}

			parsedTokenIds = append(parsedTokenIds, parsedTokenID)
		}

		// Set the parsed Token IDs in the transaction
		transaction.SetTokenIDs(parsedTokenIds...)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// PauseToken jRPC method for pauseToken
func (t *TokenService) PauseToken(_ context.Context, params param.PauseUnPauseTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenPauseTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// UnpauseToken jRPC method for unpauseToken
func (t *TokenService) UnpauseToken(_ context.Context, params param.PauseUnPauseTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenUnpauseTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// FreezeToken jRPC method for freezeToken
func (t *TokenService) FreezeToken(_ context.Context, params param.FreezeUnFreezeTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenFreezeTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.AccountId)

		if err != nil {
			return nil, err
		}
		transaction.SetAccountID(accountId)
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// UnfreezeToken jRPC method for unfreezeToken
func (t *TokenService) UnfreezeToken(_ context.Context, params param.FreezeUnFreezeTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenUnfreezeTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.AccountId)

		if err != nil {
			return nil, err
		}
		transaction.SetAccountID(accountId)
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// GrantTokenKyc jRPC method for grantTokenKyc
func (t *TokenService) GrantTokenKyc(_ context.Context, params param.GrantRevokeTokenKycParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenGrantKycTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.AccountId)

		if err != nil {
			return nil, err
		}
		transaction.SetAccountID(accountId)
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// RevokeTokenKyc jRPC method for revokeTokenKyc
func (t *TokenService) RevokeTokenKyc(_ context.Context, params param.GrantRevokeTokenKycParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenRevokeKycTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.AccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.AccountId)

		if err != nil {
			return nil, err
		}
		transaction.SetAccountID(accountId)
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// MintToken jRPC method for mintToken
func (t *TokenService) MintToken(_ context.Context, params param.MintTokenParams) (*response.TokenMintResponse, error) {
	transaction := hiero.NewTokenMintTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.Metadata != nil {
		var allMetadata [][]byte
		for _, metadataValue := range *params.Metadata {
			decodedMetadata, err := hex.DecodeString(metadataValue)
			if err != nil {
				return nil, fmt.Errorf("failed to decode metadata: %w", err)
			}
			allMetadata = append(allMetadata, decodedMetadata)
		}

		// Set the separate metadata slices on the transaction
		transaction.SetMetadatas(allMetadata)
	}

	if params.Amount != nil {
		amount, err := strconv.ParseUint(*params.Amount, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAmount(amount)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	// Construct the response
	status := receipt.Status.String()
	newTotalSupply := strconv.FormatUint(receipt.TotalSupply, 10)
	serialNumbers := utils.MapSerialNumbersToString(receipt.SerialNumbers)

	return &response.TokenMintResponse{
		TokenId:        params.TokenId,
		NewTotalSupply: &newTotalSupply,
		SerialNumbers:  &serialNumbers,
		Status:         &status,
	}, nil
}

// BurnToken jRPC method for burnToken
func (t *TokenService) BurnToken(_ context.Context, params param.BurnTokenParams) (*response.TokenBurnResponse, error) {
	transaction := hiero.NewTokenBurnTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.Amount != nil {
		amount, err := strconv.ParseUint(*params.Amount, 10, 64)
		if err != nil {
			return nil, err
		}

		transaction.SetAmount(amount)
	}

	if params.SerialNumbers != nil {
		var allSerialNumbers []int64

		for _, serialNumber := range *params.SerialNumbers {
			serial, err := strconv.ParseInt(serialNumber, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse serial number: %w", err)
			}
			allSerialNumbers = append(allSerialNumbers, serial)
		}
		transaction.SetSerialNumbers(allSerialNumbers)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	// Construct the response
	status := receipt.Status.String()
	newTotalSupply := strconv.FormatUint(receipt.TotalSupply, 10)

	return &response.TokenBurnResponse{
		TokenId:        params.TokenId,
		NewTotalSupply: &newTotalSupply,
		Status:         &status,
	}, nil
}
