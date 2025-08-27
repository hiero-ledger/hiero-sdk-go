package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type ScheduleService struct {
	sdkService *SDKService
}

func (s *ScheduleService) SetSdkService(service *SDKService) {
	s.sdkService = service
}

// CreateSchedule jRPC method for createSchedule
func (s *ScheduleService) CreateSchedule(_ context.Context, params param.ScheduleCreateParams) (*response.ScheduleResponse, error) {
	transaction := hiero.NewScheduleCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.ScheduledTransaction != nil {
		scheduledTx, err := s.buildScheduledTransaction(params.ScheduledTransaction)
		if err != nil {
			return nil, fmt.Errorf("failed to build scheduled transaction: %w", err)
		}
		_, err = transaction.SetScheduledTransaction(scheduledTx)
		if err != nil {
			return nil, fmt.Errorf("failed to set scheduled transaction: %w", err)
		}
	}

	if params.Memo != nil {
		transaction.SetScheduleMemo(*params.Memo)
	}

	if err := utils.SetKeyIfPresent(params.AdminKey, transaction.SetAdminKey); err != nil {
		return nil, err
	}

	if params.PayerAccountId != nil {
		payerAccountID, err := hiero.AccountIDFromString(*params.PayerAccountId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payer account ID: %w", err)
		}
		transaction.SetPayerAccountID(payerAccountID)
	}

	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expiration time: %w", err)
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}

	if params.WaitForExpiry != nil {
		transaction.SetWaitForExpiry(*params.WaitForExpiry)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, s.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(s.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(s.sdkService.Client)
	if err != nil {
		return nil, err
	}

	var scheduleId string
	if receipt.Status == hiero.StatusSuccess {
		scheduleId = receipt.ScheduleID.String()
	}

	return &response.ScheduleResponse{
		ScheduleId: scheduleId,
		Status:     receipt.Status.String(),
	}, nil
}

// buildScheduledTransaction creates the appropriate transaction based on method name
func (s *ScheduleService) buildScheduledTransaction(scheduledTx *param.ScheduledTransaction) (hiero.TransactionInterface, error) {
	switch scheduledTx.Method {
	case "createAccount":
		var params param.CreateAccountParams
		if err := json.Unmarshal(scheduledTx.Params, &params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal createAccount params: %w", err)
		}
		accountService := &AccountService{sdkService: s.sdkService}
		return accountService.buildCreateAccount(params)

	case "updateAccount":
		var params param.UpdateAccountParams
		if err := json.Unmarshal(scheduledTx.Params, &params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal updateAccount params: %w", err)
		}
		accountService := &AccountService{sdkService: s.sdkService}
		return accountService.buildUpdateAccount(params)

	case "deleteAccount":
		var params param.DeleteAccountParams
		if err := json.Unmarshal(scheduledTx.Params, &params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal deleteAccount params: %w", err)
		}
		accountService := &AccountService{sdkService: s.sdkService}
		return accountService.buildDeleteAccount(params)

	case "approveAllowance":
		var params param.AccountAllowanceApproveParams
		if err := json.Unmarshal(scheduledTx.Params, &params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal approveAllowance params: %w", err)
		}
		accountService := &AccountService{sdkService: s.sdkService}
		return accountService.buildApproveAllowance(params)

	case "transferCrypto":
		var params param.TransferCryptoParams
		if err := json.Unmarshal(scheduledTx.Params, &params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transferCrypto params: %w", err)
		}
		accountService := &AccountService{sdkService: s.sdkService}
		return accountService.buildTransferCrypto(params)

	case "createToken":
		var params param.CreateTokenParams
		if err := json.Unmarshal(scheduledTx.Params, &params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal createToken params: %w", err)
		}
		tokenService := &TokenService{sdkService: s.sdkService}
		return tokenService.buildCreateToken(params)

	default:
		return nil, fmt.Errorf("unsupported scheduled transaction method: %s", scheduledTx.Method)
	}
}
