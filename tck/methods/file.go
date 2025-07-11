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
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
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

// updateFile jRPC method for updateFile
func (t *FileService) UpdateFile(_ context.Context, params param.UpdateFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.FileId != nil {
		fileId, err := hiero.FileIDFromString(*params.FileId)
		if err != nil {
			return nil, err
		}
		transaction.SetFileID(fileId)
	}

	if params.Keys != nil {
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
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}

	if params.Memo != nil {
		transaction.SetFileMemo(*params.Memo)
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
		Status: receipt.Status.String(),
	}, nil
}

// deleteFile jRPC method for deleteFile
func (t *FileService) DeleteFile(_ context.Context, params param.UpdateFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.FileId != nil {
		fileId, err := hiero.FileIDFromString(*params.FileId)
		if err != nil {
			return nil, err
		}
		transaction.SetFileID(fileId)
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
		Status: receipt.Status.String(),
	}, nil
}

// appendFile jRPC method for appendFile
func (t *FileService) AppendFile(_ context.Context, params param.AppendFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileAppendTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.FileId != nil {
		fileId, err := hiero.FileIDFromString(*params.FileId)
		if err != nil {
			return nil, err
		}
		transaction.SetFileID(fileId)
	}

	if params.Contents != nil {
		transaction.SetContents([]byte(*params.Contents))
	}

	if params.ChunkSize != nil {
		if *params.ChunkSize <= 0 {
			return nil, errors.New("internal error")
		}
		transaction.SetMaxChunkSize(*params.ChunkSize)
	}

	if params.MaxChunks != nil {
		if *params.MaxChunks <= 0 {
			return nil, errors.New("internal error")
		}
		transaction.SetMaxChunks(uint64(*params.MaxChunks))
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
		Status: receipt.Status.String(),
	}, nil
}
