package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/hex"
	"strconv"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type EthereumService struct {
	sdkService *SDKService
}

func (e *EthereumService) SetSdkService(service *SDKService) {
	e.sdkService = service
}

func (e *EthereumService) CreateEthereumTransaction(ctx context.Context, params param.EthereumTransactionParams) (*response.EthereumTransactionResponse, error) {
	transaction := hiero.NewEthereumTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.EthereumData != nil {
		ethereumData, err := hex.DecodeString(*params.EthereumData)
		if err != nil {
			return nil, err
		}
		transaction.SetEthereumData(ethereumData)
	}

	if params.CallDataFileID != nil {
		parsedFileID, err := hiero.FileIDFromString(*params.CallDataFileID)
		if err != nil {
			return nil, err
		}
		transaction.SetCallDataFileID(parsedFileID)
	}

	if params.MaxGasAllowed != nil {
		maxGasAllowed, err := strconv.ParseInt(*params.MaxGasAllowed, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetMaxGasAllowed(maxGasAllowed)
	}

	if params.CommonTransactionParams != nil {
		params.CommonTransactionParams.FillOutTransaction(transaction, e.sdkService.Client)
	}

	txResponse, err := transaction.Execute(e.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(e.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.EthereumTransactionResponse{
		ContractId: receipt.ContractID.String(),
		Status:     receipt.Status.String(),
	}, nil
}
