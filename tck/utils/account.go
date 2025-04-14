package utils

import (
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func SetAccountIDIfPresent[T any](accountIDStr *string, setter func(hiero.AccountID) T) error {
	if accountIDStr != nil {
		accountID, err := hiero.AccountIDFromString(*accountIDStr)
		if err != nil {
			return response.NewInternalError(err.Error())
		}
		setter(accountID)
	}
	return nil
}
