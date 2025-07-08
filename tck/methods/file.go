package methods

// SPDX-License-Identifier: Apache-2.0
import (
	"context"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type FileService struct {
	sdkService *SDKService
}

func (t *FileService) SetSdkService(service *SDKService) {
	t.sdkService = service
}

// createFile jRPC method for createFile
func (t *FileService) CreateFile(_ context.Context, params param.CreateFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.Keys != nil && len(*params.Keys) > 0 {
		var keys []hiero.Key
		for _, keyStr := range *params.Keys {
			key, err := utils.GetKeyFromString(keyStr)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		transaction.SetKeys(keys...)
	}

	if params.Contents != nil {
		transaction.SetContents([]byte(*params.Contents))
	}

	if params.ExpirationTime != nil {
		expirationTime, err := time.Parse(time.RFC3339, *params.ExpirationTime)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(expirationTime)
	}

	if params.Memo != nil {
		transaction.SetMemo(*params.Memo)
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

	return &response.FileResponse{
		FileId: receipt.FileID.String(),
		Status: receipt.Status.String(),
	}, nil
}
