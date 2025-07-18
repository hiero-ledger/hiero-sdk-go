package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"strconv"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type TopicService struct {
	sdkService *SDKService
}

func (t *TopicService) SetSdkService(service *SDKService) {
	t.sdkService = service
}

// CreateTopic jRPC method for createTopic
func (t *TopicService) CreateTopic(_ context.Context, params param.CreateTopicParams) (*response.TopicResponse, error) {
	transaction := hiero.NewTopicCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.Memo != nil {
		transaction.SetTopicMemo(*params.Memo)
	}

	if err := utils.SetKeyIfPresent(params.AdminKey, transaction.SetAdminKey); err != nil {
		return nil, err
	}

	if err := utils.SetKeyIfPresent(params.SubmitKey, transaction.SetSubmitKey); err != nil {
		return nil, err
	}

	if err := utils.SetKeyIfPresent(params.FeeScheduleKey, transaction.SetFeeScheduleKey); err != nil {
		return nil, err
	}

	if params.FeeExemptKeys != nil && len(*params.FeeExemptKeys) > 0 {
		var keys []hiero.Key
		for _, keyStr := range *params.FeeExemptKeys {
			key, err := utils.GetKeyFromString(keyStr)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		transaction.SetFeeExemptKeys(keys)
	}

	if params.AutoRenewPeriod != nil {
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}

	if err := utils.SetAccountIDIfPresent(params.AutoRenewAccountId, transaction.SetAutoRenewAccountID); err != nil {
		return nil, err
	}

	if params.CustomFees != nil {
		customFees, err := utils.ParseCustomFees(*params.CustomFees)
		if err != nil {
			return nil, err
		}

		var topicCustomFees []*hiero.CustomFixedFee
		for _, fee := range customFees {
			if fixedFee, ok := fee.(*hiero.CustomFixedFee); ok {
				topicCustomFees = append(topicCustomFees, fixedFee)
			}
			// Ignore fractional and royalty fees as topics don't support them
		}
		transaction.SetCustomFees(topicCustomFees)
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

	return &response.TopicResponse{
		TopicId: receipt.TopicID.String(),
		Status:  receipt.Status.String(),
	}, nil
}
