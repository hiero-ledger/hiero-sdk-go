package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/hex"
	"strconv"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// ---- Struct to hold hiero.Client implementation and to implement the methods of the specification ----
type ContractService struct {
	sdkService *SDKService
}

// SetSdkService We set object, which is holding our client param. Pass it by reference, because TCK is dynamically updating it
func (c *ContractService) SetSdkService(service *SDKService) {
	c.sdkService = service
}

// CreateContract jRPC method for contractCreate
func (c *ContractService) CreateContract(_ context.Context, params param.ContractCreateTransactionParams) (*response.ContractResponse, error) {
	transaction := hiero.NewContractCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.BytecodeFileId != nil {
		fileID, err := hiero.FileIDFromString(*params.BytecodeFileId)
		if err != nil {
			return nil, err
		}
		transaction.SetBytecodeFileID(fileID)
	}

	if err := utils.SetKeyIfPresent(params.AdminKey, transaction.SetAdminKey); err != nil {
		return nil, err
	}

	if params.Gas != nil {
		gas, err := strconv.ParseUint(*params.Gas, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetGas(gas)
	}

	if params.InitialBalance != nil {
		initialBalance, err := strconv.ParseInt(*params.InitialBalance, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetInitialBalance(hiero.HbarFromTinybar(initialBalance))
	}

	if params.ConstructorParameters != nil {
		constructorParams, err := hex.DecodeString(*params.ConstructorParameters)
		if err != nil {
			return nil, err
		}
		transaction.SetConstructorParametersRaw(constructorParams)
	}

	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriodInt(autoRenewPeriodSeconds)
	}

	if params.AutoRenewAccountId != nil {
		autoRenewAccountID, err := hiero.AccountIDFromString(*params.AutoRenewAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewAccountID(autoRenewAccountID)
	}

	if params.Memo != nil {
		transaction.SetContractMemo(*params.Memo)
	}

	if params.StakedAccountId != nil {
		stakedAccountID, err := hiero.AccountIDFromString(*params.StakedAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedAccountID(stakedAccountID)
	}

	if params.StakedNodeId != nil {
		stakedNodeIDVal, err := strconv.ParseInt(string(*params.StakedNodeId), 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedNodeID(stakedNodeIDVal)
	}

	if params.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*params.DeclineStakingReward)
	}

	if params.MaxAutomaticTokenAssociation != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociation)
	}

	if params.Initcode != nil {
		initcode, err := hex.DecodeString(*params.Initcode)
		if err != nil {
			return nil, err
		}
		transaction.SetBytecode(initcode)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, c.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(c.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(c.sdkService.Client)
	if err != nil {
		return nil, err
	}

	var contractId string
	if receipt.Status == hiero.StatusSuccess {
		contractId = receipt.ContractID.String()
	}

	return &response.ContractResponse{ContractId: contractId, Status: receipt.Status.String()}, nil
}
