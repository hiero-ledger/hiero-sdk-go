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
	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TopicResponse{
		TopicId: receipt.TopicID.String(),
		Status:  receipt.Status.String(),
	}, nil
}

// buildCreateTopic builds a topic create transaction without executing it (for scheduling)
func (t *TopicService) buildCreateTopic(params param.CreateTopicParams) (*hiero.TopicCreateTransaction, error) {
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
		}
		transaction.SetCustomFees(topicCustomFees)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

// UpdateTopic jRPC method for updateTopic
func (t *TopicService) UpdateTopic(_ context.Context, params param.UpdateTopicParams) (*response.TopicResponse, error) {
	transaction := hiero.NewTopicUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TopicId != nil {
		topicId, err := hiero.TopicIDFromString(*params.TopicId)
		if err != nil {
			return nil, err
		}
		transaction.SetTopicID(topicId)
	}

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

	if params.FeeExemptKeys != nil {
		if len(*params.FeeExemptKeys) == 0 {
			transaction.ClearFeeExemptKeys()
		} else {
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

	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}

	if params.CustomFees != nil {
		if len(*params.CustomFees) == 0 {
			transaction.ClearCustomFees()
		} else {
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
	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TopicResponse{
		Status: receipt.Status.String(),
	}, nil
}

// DeleteTopic jRPC method for deleteTopic
func (t *TopicService) DeleteTopic(_ context.Context, params param.DeleteTopicParams) (*response.TopicResponse, error) {
	transaction := hiero.NewTopicDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TopicId != nil {
		topicId, err := hiero.TopicIDFromString(*params.TopicId)
		if err != nil {
			return nil, err
		}
		transaction.SetTopicID(topicId)
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
	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TopicResponse{
		Status: receipt.Status.String(),
	}, nil
}

// buildSubmitTopicMessage builds a TopicMessageSubmitTransaction from parameters
func (t *TopicService) buildSubmitTopicMessage(params param.SubmitTopicMessageParams) (*hiero.TopicMessageSubmitTransaction, error) {
	transaction := hiero.NewTopicMessageSubmitTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.TopicId != nil {
		topicId, err := hiero.TopicIDFromString(*params.TopicId)
		if err != nil {
			return nil, err
		}
		transaction.SetTopicID(topicId)
	}

	if params.Message != nil {
		transaction.SetMessage(*params.Message)
	}

	if params.MaxChunks != nil {
		if *params.MaxChunks < 0 {
			return nil, response.NewInternalError("maxChunks must be greater than 0")
		}
		transaction.SetMaxChunks(uint64(*params.MaxChunks))
	}

	if params.ChunkSize != nil {
		if *params.ChunkSize < 0 {
			return nil, response.NewInternalError("chunkSize must be greater than 0")
		}
		transaction.SetChunkSize(uint64(*params.ChunkSize))
	}

	if params.CustomFeeLimits != nil {
		for _, customFeeLimit := range *params.CustomFeeLimits {
			var sdkCustomFeeLimit hiero.CustomFeeLimit
			if customFeeLimit.PayerId != nil {
				payerId, err := hiero.AccountIDFromString(*customFeeLimit.PayerId)
				if err != nil {
					return nil, err
				}
				sdkCustomFeeLimit.SetPayerId(payerId)

				for _, fixedFee := range *customFeeLimit.FixedFees {
					sdkFixedFee := hiero.NewCustomFixedFee()
					if fixedFee.DenominatingTokenId != nil {
						tokenId, err := hiero.TokenIDFromString(*fixedFee.DenominatingTokenId)
						if err != nil {
							return nil, err
						}
						sdkFixedFee.SetDenominatingTokenID(tokenId)
					}
					parsed, _ := strconv.ParseInt(fixedFee.Amount, 10, 64)
					sdkFixedFee.SetAmount(parsed)
					sdkCustomFeeLimit.AddCustomFee(sdkFixedFee)
				}
			}
			transaction.AddCustomFeeLimit(&sdkCustomFeeLimit)
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

// SubmitTopicMessage jRPC method for submitTopicMessage
func (t *TopicService) SubmitTopicMessage(_ context.Context, params param.SubmitTopicMessageParams) (*response.TopicResponse, error) {
	transaction, err := t.buildSubmitTopicMessage(params)
	if err != nil {
		return nil, err
	}

	txResponse, err := transaction.Execute(t.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(t.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.TopicResponse{
		Status: receipt.Status.String(),
	}, nil
}
