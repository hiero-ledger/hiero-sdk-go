package utils

import (
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func SetAccountIDIfPresent[T any](accountIDStr *string, setter func(hiero.AccountID) T) error {
	if accountIDStr != nil {
		accountID, err := hiero.AccountIDFromString(*accountIDStr)
		if err != nil {
			return err
		}
		setter(accountID)
	}
	return nil
}
