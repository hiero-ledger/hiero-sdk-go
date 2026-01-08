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

// CreateToken jRPC method for createToken
func (t *TokenService) CreateToken(_ context.Context, params param.CreateTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set admin key
	if err := utils.SetKeyIfPresent(params.AdminKey, transaction.SetAdminKey); err != nil {
		return nil, err
	}
	// Set kyc key
	if err := utils.SetKeyIfPresent(params.KycKey, transaction.SetKycKey); err != nil {
		return nil, err
	}
	// Set freeze key
	if err := utils.SetKeyIfPresent(params.FreezeKey, transaction.SetFreezeKey); err != nil {
		return nil, err
	}
	// Set wipe key
	if err := utils.SetKeyIfPresent(params.WipeKey, transaction.SetWipeKey); err != nil {
		return nil, err
	}
	// Set pause key
	if err := utils.SetKeyIfPresent(params.PauseKey, transaction.SetPauseKey); err != nil {
		return nil, err
	}
	// Set metadata key
	if err := utils.SetKeyIfPresent(params.MetadataKey, transaction.SetMetadataKey); err != nil {
		return nil, err
	}
	// Set supply key
	if err := utils.SetKeyIfPresent(params.SupplyKey, transaction.SetSupplyKey); err != nil {
		return nil, err
	}
	// Set fee schedule key
	if err := utils.SetKeyIfPresent(params.FeeScheduleKey, transaction.SetFeeScheduleKey); err != nil {
		return nil, err
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

	// Set token types
	if err := utils.SetTokenTypes(transaction, params); err != nil {
		return nil, err
	}

	// Set token supply params
	if err := utils.SetTokenSupplyParams(transaction, params); err != nil {
		return nil, err
	}

	// Set treasury account ID
	if err := utils.SetAccountIDIfPresent(params.TreasuryAccountId, transaction.SetTreasuryAccountID); err != nil {
		return nil, err
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

	// Set auto renew account ID
	if err := utils.SetAccountIDIfPresent(params.AutoRenewAccountId, transaction.SetAutoRenewAccount); err != nil {
		return nil, err
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
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

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{TokenId: receipt.TokenID.String(), Status: receipt.Status.String()}, nil
}

// UpdateToken jRPC method for updateToken
func (t *TokenService) UpdateToken(_ context.Context, params param.UpdateTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)
		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	// Set admin key
	if err := utils.SetKeyIfPresent(params.AdminKey, transaction.SetAdminKey); err != nil {
		return nil, err
	}
	// Set kyc key
	if err := utils.SetKeyIfPresent(params.KycKey, transaction.SetKycKey); err != nil {
		return nil, err
	}
	// Set freeze key
	if err := utils.SetKeyIfPresent(params.FreezeKey, transaction.SetFreezeKey); err != nil {
		return nil, err
	}
	// Set wipe key
	if err := utils.SetKeyIfPresent(params.WipeKey, transaction.SetWipeKey); err != nil {
		return nil, err
	}
	// Set pause key
	if err := utils.SetKeyIfPresent(params.PauseKey, transaction.SetPauseKey); err != nil {
		return nil, err
	}
	// Set metadata key
	if err := utils.SetKeyIfPresent(params.MetadataKey, transaction.SetMetadataKey); err != nil {
		return nil, err
	}
	// Set supply key
	if err := utils.SetKeyIfPresent(params.SupplyKey, transaction.SetSupplyKey); err != nil {
		return nil, err
	}
	// Set fee schedule key
	if err := utils.SetKeyIfPresent(params.FeeScheduleKey, transaction.SetFeeScheduleKey); err != nil {
		return nil, err
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
	// Set treasury account ID
	if err := utils.SetAccountIDIfPresent(params.TreasuryAccountId, transaction.SetTreasuryAccountID); err != nil {
		return nil, err
	}

	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}
	// Set auto renew account ID
	if err := utils.SetAccountIDIfPresent(params.AutoRenewAccountId, transaction.SetAutoRenewAccount); err != nil {
		return nil, err
	}
	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriodSeconds(autoRenewPeriodSeconds)
	}

	if params.Metadata != nil {
		transaction.SetTokenMetadata([]byte(*params.Metadata))
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}
	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// AssociateToken jRPC method for associateToken
func (t *TokenService) AssociateToken(_ context.Context, params param.AssociateDissociatesTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenAssociateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// DisassociateToken jRPC method for dissociateToken
func (t *TokenService) DisassociateToken(_ context.Context, params param.AssociateDissociatesTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenDissociateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// FreezeToken jRPC method for freezeToken
func (t *TokenService) FreezeToken(_ context.Context, params param.FreezeUnFreezeTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenFreezeTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// UnfreezeToken jRPC method for unfreezeToken
func (t *TokenService) UnfreezeToken(_ context.Context, params param.FreezeUnFreezeTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenUnfreezeTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// GrantTokenKyc jRPC method for grantTokenKyc
func (t *TokenService) GrantTokenKyc(_ context.Context, params param.GrantRevokeTokenKycParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenGrantKycTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// RevokeTokenKyc jRPC method for revokeTokenKyc
func (t *TokenService) RevokeTokenKyc(_ context.Context, params param.GrantRevokeTokenKycParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenRevokeKycTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set accountId
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}
		transaction.SetTokenID(tokenId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// buildMintToken builds a TokenMintTransaction from parameters
func (t *TokenService) buildMintToken(params param.MintTokenParams) (*hiero.TokenMintTransaction, error) {
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

// MintToken jRPC method for mintToken
func (t *TokenService) MintToken(_ context.Context, params param.MintTokenParams) (*response.TokenMintResponse, error) {
	transaction, err := t.buildMintToken(params)
	if err != nil {
		return nil, err
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
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

// buildBurnToken builds a TokenBurnTransaction from parameters
func (t *TokenService) buildBurnToken(params param.BurnTokenParams) (*hiero.TokenBurnTransaction, error) {
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

// BurnToken jRPC method for burnToken
func (t *TokenService) BurnToken(_ context.Context, params param.BurnTokenParams) (*response.TokenBurnResponse, error) {
	transaction, err := t.buildBurnToken(params)
	if err != nil {
		return nil, err
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
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

// jRPC method for wipeToken
func (t *TokenService) WipeToken(_ context.Context, params param.WipeTokenParams) (*response.TokenWipeResponse, error) {
	transaction := hiero.NewTokenWipeTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId != nil {
		tokenId, err := hiero.TokenIDFromString(*params.TokenId)

		if err != nil {
			return nil, err
		}

		transaction.SetTokenID(tokenId)
	}

	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
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
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	// Construct the response
	status := receipt.Status.String()
	newTotalSupply := strconv.FormatUint(receipt.TotalSupply, 10)

	return &response.TokenWipeResponse{
		TokenId:        params.TokenId,
		NewTotalSupply: &newTotalSupply,
		Status:         &status,
	}, nil
}

// AirdropToken jRPC method for airdropToken
func (t *TokenService) AirdropToken(_ context.Context, params param.AirdropParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenAirdropTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenTransfers == nil {
		return nil, response.NewInternalError("transferParams is required")
	}

	transferParams := *params.TokenTransfers

	for _, transferParam := range transferParams {
		if err := utils.HandleAirdropParam(transaction, transferParam); err != nil {
			return nil, err
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// ClaimToken jRPC method for claimToken
func (t *TokenService) ClaimToken(_ context.Context, params param.ClaimTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenClaimAirdropTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TokenId == nil || params.SenderAccountId == nil || params.ReceiverAccountId == nil {
		return nil, response.NewInternalError("tokenId, senderAccountId, and receiverAccountId are required")
	}

	tokenID, err := hiero.TokenIDFromString(*params.TokenId)
	if err != nil {
		return nil, err
	}

	senderID, err := hiero.AccountIDFromString(*params.SenderAccountId)
	if err != nil {
		return nil, err
	}

	receiverID, err := hiero.AccountIDFromString(*params.ReceiverAccountId)
	if err != nil {
		return nil, err
	}

	// NFT token claiming
	if params.SerialNumbers != nil && len(*params.SerialNumbers) > 0 {
		for _, serialNumber := range *params.SerialNumbers {
			serial, err := strconv.ParseInt(serialNumber, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse serial number: %w", err)
			}

			nftID := hiero.NftID{
				TokenID:      tokenID,
				SerialNumber: serial,
			}

			pendingAirdropID := &hiero.PendingAirdropId{}
			pendingAirdropID.SetSender(senderID)
			pendingAirdropID.SetReceiver(receiverID)
			pendingAirdropID.SetNftID(nftID)

			transaction.AddPendingAirdropId(*pendingAirdropID)
		}
	} else {
		// Fungible token claiming
		pendingAirdropID := &hiero.PendingAirdropId{}
		pendingAirdropID.SetSender(senderID)
		pendingAirdropID.SetReceiver(receiverID)
		pendingAirdropID.SetTokenID(tokenID)

		transaction.AddPendingAirdropId(*pendingAirdropID)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{Status: receipt.Status.String()}, nil
}

// CancelAirdrop jRPC method for cancelAirdrop
func (t *TokenService) CancelAirdrop(_ context.Context, params param.AirdropCancelTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenCancelAirdropTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.PendingAirdrops == nil {
		return nil, response.NewInternalError("pendingAirdrops is required")
	}

	for _, pendingAirdrop := range *params.PendingAirdrops {
		if pendingAirdrop.SenderAccountId == nil || pendingAirdrop.ReceiverAccountId == nil || pendingAirdrop.TokenId == nil {
			return nil, response.NewInternalError("senderAccountId, receiverAccountId, and tokenId are required for each pending airdrop")
		}

		senderID, err := hiero.AccountIDFromString(*pendingAirdrop.SenderAccountId)
		if err != nil {
			return nil, err
		}

		receiverID, err := hiero.AccountIDFromString(*pendingAirdrop.ReceiverAccountId)
		if err != nil {
			return nil, err
		}

		tokenID, err := hiero.TokenIDFromString(*pendingAirdrop.TokenId)
		if err != nil {
			return nil, err
		}

		// NFT token canceling
		if pendingAirdrop.SerialNumbers != nil && len(*pendingAirdrop.SerialNumbers) > 0 {
			for _, serialNumberStr := range *pendingAirdrop.SerialNumbers {
				serialNumber, err := strconv.ParseInt(serialNumberStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse serial number: %w", err)
				}

				nftID := hiero.NftID{
					TokenID:      tokenID,
					SerialNumber: serialNumber,
				}

				pendingAirdropID := &hiero.PendingAirdropId{}
				pendingAirdropID.SetSender(senderID)
				pendingAirdropID.SetReceiver(receiverID)
				pendingAirdropID.SetNftID(nftID)

				transaction.AddPendingAirdropId(*pendingAirdropID)
			}
		} else {
			// Fungible token canceling
			pendingAirdropID := &hiero.PendingAirdropId{}
			pendingAirdropID.SetSender(senderID)
			pendingAirdropID.SetReceiver(receiverID)
			pendingAirdropID.SetTokenID(tokenID)

			transaction.AddPendingAirdropId(*pendingAirdropID)
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		Status: receipt.Status.String(),
	}, nil
}

// RejectToken jRPC method for rejectToken
func (t *TokenService) RejectToken(_ context.Context, params param.RejectTokenParams) (*response.TokenResponse, error) {
	transaction := hiero.NewTokenRejectTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set owner ID if present
	if err := utils.SetAccountIDIfPresent(params.OwnerId, transaction.SetOwnerID); err != nil {
		return nil, err
	}

	// Handle token IDs and serial numbers
	if params.TokenIds != nil && len(*params.TokenIds) > 0 {
		tokenIds := *params.TokenIds

		// If no serial numbers, add all token IDs as fungible tokens
		if params.SerialNumbers == nil || len(*params.SerialNumbers) == 0 {
			for _, tokenIdStr := range tokenIds {
				tokenId, err := hiero.TokenIDFromString(tokenIdStr)
				if err != nil {
					return nil, err
				}
				transaction.AddTokenID(tokenId)
			}
		} else {
			// NFT token rejecting - add specific NFTs
			serialNumbers := *params.SerialNumbers
			for _, tokenIdStr := range tokenIds {
				tokenId, err := hiero.TokenIDFromString(tokenIdStr)
				if err != nil {
					return nil, err
				}

				for _, serialNumberStr := range serialNumbers {
					serialNumber, err := strconv.ParseInt(serialNumberStr, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("failed to parse serial number: %w", err)
					}

					nftID := hiero.NftID{
						TokenID:      tokenId,
						SerialNumber: serialNumber,
					}
					transaction.AddNftID(nftID)
				}
			}
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		Status: receipt.Status.String(),
	}, nil
}
