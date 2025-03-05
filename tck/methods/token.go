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
