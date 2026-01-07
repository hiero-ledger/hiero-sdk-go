package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// ---- Struct to hold hiero.Client implementation and to implement the methods of the specification ----
type AccountService struct {
	sdkService *SDKService
}

// Variable to be set to `SetGrpcDeadline` for all transactions
var threeSecondsDuration = time.Second * 3

// SetSdkService We set object, which is holding our client param. Pass it by referance, because TCK is dynamically updating it
func (a *AccountService) SetSdkService(service *SDKService) {
	a.sdkService = service
}

// CreateAccount jRPC method for createAccount
func (a *AccountService) CreateAccount(_ context.Context, params param.CreateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set key
	if err := utils.SetKeyIfPresent(params.Key, transaction.SetKeyWithoutAlias); err != nil {
		return nil, err
	}
	if params.InitialBalance != nil {
		initialBalance, err := strconv.ParseInt(*params.InitialBalance, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetInitialBalance(hiero.HbarFromTinybar(initialBalance))
	}
	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}
	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}
	// Set staked account ID
	if err := utils.SetAccountIDIfPresent(params.StakedAccountId, transaction.SetStakedAccountID); err != nil {
		return nil, err
	}
	if params.StakedNodeId != nil {
		stakedNodeID, err := params.StakedNodeId.Int64()
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}
	if params.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*params.DeclineStakingReward)
	}
	if params.Memo != nil {
		transaction.SetAccountMemo(*params.Memo)
	}
	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}

		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}
	if params.Alias != nil {
		transaction.SetAlias(*params.Alias)
	}
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}
	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	var accId string
	if receipt.Status == hiero.StatusSuccess {
		accId = receipt.AccountID.String()
	}
	return &response.AccountResponse{AccountId: accId, Status: receipt.Status.String()}, nil
}

// buildCreateAccount builds an account create transaction without executing it (for scheduling)
func (a *AccountService) buildCreateAccount(params param.CreateAccountParams) (*hiero.AccountCreateTransaction, error) {
	transaction := hiero.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set key
	if err := utils.SetKeyIfPresent(params.Key, transaction.SetKeyWithoutAlias); err != nil {
		return nil, err
	}
	if params.InitialBalance != nil {
		initialBalance, err := strconv.ParseInt(*params.InitialBalance, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetInitialBalance(hiero.HbarFromTinybar(initialBalance))
	}
	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}
	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}
	// Set staked account ID
	if err := utils.SetAccountIDIfPresent(params.StakedAccountId, transaction.SetStakedAccountID); err != nil {
		return nil, err
	}
	if params.StakedNodeId != nil {
		stakedNodeID, err := params.StakedNodeId.Int64()
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}
	if params.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*params.DeclineStakingReward)
	}
	if params.Memo != nil {
		transaction.SetAccountMemo(*params.Memo)
	}
	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}

		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}
	if params.Alias != nil {
		transaction.SetAlias(*params.Alias)
	}
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

// UpdateAccount jRPC method for updateAccount
func (a *AccountService) UpdateAccount(_ context.Context, params param.UpdateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)
	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	// Set key
	if err := utils.SetKeyIfPresent(params.Key, transaction.SetKey); err != nil {
		return nil, err
	}

	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}

	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}

	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}

	// Set staked account ID
	if err := utils.SetAccountIDIfPresent(params.StakedAccountId, transaction.SetStakedAccountID); err != nil {
		return nil, err
	}

	if params.StakedNodeId != nil {
		stakedNodeID, err := params.StakedNodeId.Int64()
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}

	if params.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*params.DeclineStakingReward)
	}

	if params.Memo != nil {
		transaction.SetAccountMemo(*params.Memo)
	}

	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// DeleteAccount jRPC method for deleteAccount
func (a *AccountService) DeleteAccount(_ context.Context, params param.DeleteAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)
	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.DeleteAccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	// Set transfer account ID
	if err := utils.SetAccountIDIfPresent(params.TransferAccountId, transaction.SetTransferAccountID); err != nil {
		return nil, err
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// buildApproveAllowance builds an AccountAllowanceApproveTransaction from parameters
func (a *AccountService) buildApproveAllowance(params param.AccountAllowanceApproveParams) (*hiero.AccountAllowanceApproveTransaction, error) {
	transaction := hiero.NewAccountAllowanceApproveTransaction().SetGrpcDeadline(&threeSecondsDuration)

	allowances := *params.Allowances

	for _, allowance := range allowances {
		owner, err := hiero.AccountIDFromString(*allowance.OwnerAccountId)
		if err != nil {
			return nil, err
		}

		spender, err := hiero.AccountIDFromString(*allowance.SpenderAccountId)
		if err != nil {
			return nil, err
		}

		hbar := allowance.Hbar
		token := allowance.Token
		nft := allowance.Nft

		switch {
		case hbar != nil:
			// Process Hbar allowance
			hbarAmount, err := strconv.ParseInt(*hbar.Amount, 10, 64)
			if err != nil {
				return nil, err
			}
			transaction.ApproveHbarAllowance(owner, spender, hiero.HbarFromTinybar(hbarAmount))

		case token != nil:
			// Process Token allowance
			tokenID, err := hiero.TokenIDFromString(*token.TokenId)
			if err != nil {
				return nil, err
			}
			tokenAmount, err := strconv.ParseInt(*token.Amount, 10, 64)
			if err != nil {
				return nil, err
			}
			transaction.ApproveTokenAllowance(tokenID, owner, spender, tokenAmount)

		case nft != nil:
			// Process Nft allowance
			tokenID, err := hiero.TokenIDFromString(*nft.TokenId)
			if err != nil {
				return nil, err
			}

			switch {
			case nft.SerialNumbers != nil:
				for _, serialNumber := range *nft.SerialNumbers {
					serialNumberParsed, err := strconv.ParseInt(serialNumber, 10, 64)
					if err != nil {
						return nil, err
					}

					nftID := hiero.NftID{
						TokenID:      tokenID,
						SerialNumber: serialNumberParsed,
					}

					if nft.DelegateSpenderAccountId != nil {
						delegateSpenderAccountId, err := hiero.AccountIDFromString(*nft.DelegateSpenderAccountId)
						if err != nil {
							return nil, err
						}

						transaction.ApproveTokenNftAllowanceWithDelegatingSpender(
							nftID,
							owner,
							spender,
							delegateSpenderAccountId,
						)
					} else {
						transaction.ApproveTokenNftAllowance(
							nftID,
							owner,
							spender,
						)
					}
				}
			case nft.ApprovedForAll != nil && *nft.ApprovedForAll:
				transaction.ApproveTokenNftAllowanceAllSerials(
					tokenID,
					owner,
					spender,
				)
			default:
				transaction.DeleteTokenNftAllowanceAllSerials(tokenID, owner, spender)
			}

		default:
			return nil, errors.New("no valid allowance type provided")
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}
	return transaction, nil
}

// ApproveAllowance jRPC method for approveAllowance
func (a *AccountService) ApproveAllowance(_ context.Context, params param.AccountAllowanceApproveParams) (*response.AccountResponse, error) {
	transaction, err := a.buildApproveAllowance(params)
	if err != nil {
		return nil, err
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// DeleteAllowance jRPC method for deleteAllowance
func (a *AccountService) DeleteAllowance(_ context.Context, params param.AccountAllowanceDeleteParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountAllowanceDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)

	allowances := *params.Allowances

	// Loop through each allowance and process
	for _, allowance := range allowances {
		owner, err := hiero.AccountIDFromString(*allowance.OwnerAccountId)
		if err != nil {
			return nil, err
		}

		tokenID, err := hiero.TokenIDFromString(*allowance.TokenId)
		if err != nil {
			return nil, err
		}

		// Process NFT serial numbers if provided
		if allowance.SerialNumbers != nil {
			for _, serialNumber := range *allowance.SerialNumbers {
				serialNumberParsed, err := strconv.ParseInt(serialNumber, 10, 64)
				if err != nil {
					return nil, err
				}

				nftID := hiero.NftID{
					TokenID:      tokenID,
					SerialNumber: serialNumberParsed,
				}

				transaction.DeleteAllTokenNftAllowances(nftID, &owner)
			}
		} else {
			transaction.DeleteAllTokenNftAllowances(hiero.NftID{TokenID: tokenID}, &owner)
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// buildTransferCrypto builds a TransferTransaction from parameters
func (a *AccountService) buildTransferCrypto(params param.TransferCryptoParams) (*hiero.TransferTransaction, error) {
	transaction := hiero.NewTransferTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.Transfers == nil {
		return nil, response.NewInternalError("transferParams is required")
	}

	transferParams := *params.Transfers
	if len(transferParams) == 0 {
		return nil, response.NewInternalError("transferParams is required")
	}

	for _, transferParam := range transferParams {
		if err := utils.HandleTransferParam(transaction, transferParam); err != nil {
			return nil, err
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}
	return transaction, nil
}

// TransferCrypto jRPC method for transferCrypto
func (a *AccountService) TransferCrypto(_ context.Context, params param.TransferCryptoParams) (*response.AccountResponse, error) {
	transaction, err := a.buildTransferCrypto(params)
	if err != nil {
		return nil, err
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// GetAccountBalance jRPC method for getAccountBalance
func (a *AccountService) GetAccountBalance(_ context.Context, params param.GetAccountBalanceParams) (*response.AccountBalanceResponse, error) {
	query := hiero.NewAccountBalanceQuery().SetGrpcDeadline(&threeSecondsDuration)

	if params.AccountId != nil {
		accountID, err := hiero.AccountIDFromString(*params.AccountId)
		if err != nil {
			return nil, err
		}
		query.SetAccountID(accountID)
	}

	if params.ContractId != nil {
		contractID, err := hiero.ContractIDFromString(*params.ContractId)
		if err != nil {
			return nil, err
		}
		query.SetContractID(contractID)
	}

	balance, err := query.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	tokenBalances := make(map[string]uint64)
	for tokenID, amount := range balance.Tokens.GetAll() {
		tokenBalances[tokenID] = amount
	}

	tokenDecimals := make(map[string]uint64)
	for tokenID, decimals := range balance.TokenDecimals.GetAll() {
		tokenDecimals[tokenID] = decimals
	}

	return &response.AccountBalanceResponse{
		Hbar:          strconv.Itoa(int(balance.Hbars.AsTinybar())),
		TokenBalances: tokenBalances,
		TokenDecimals: tokenDecimals,
	}, nil
}
