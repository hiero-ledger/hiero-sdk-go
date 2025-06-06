package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// TopicInfo is the information about a topic
type TopicInfo struct {
	TopicMemo          string
	RunningHash        []byte
	SequenceNumber     uint64
	ExpirationTime     time.Time
	AdminKey           Key
	SubmitKey          Key
	FeeScheduleKey     Key
	FeeExemptKeys      []Key
	CustomFees         []*CustomFixedFee
	AutoRenewPeriod    time.Duration
	AutoRenewAccountID *AccountID
	LedgerID           LedgerID
}

func _TopicInfoFromProtobuf(topicInfo *services.ConsensusTopicInfo) (TopicInfo, error) {
	if topicInfo == nil {
		return TopicInfo{}, errParameterNull
	}
	var err error
	tempTopicInfo := TopicInfo{
		TopicMemo:      topicInfo.Memo,
		RunningHash:    topicInfo.RunningHash,
		SequenceNumber: topicInfo.SequenceNumber,
		LedgerID:       LedgerID{topicInfo.LedgerId},
	}

	if autoRenewPeriod := topicInfo.AutoRenewPeriod; autoRenewPeriod != nil {
		tempTopicInfo.AutoRenewPeriod = _DurationFromProtobuf(topicInfo.AutoRenewPeriod)
	}

	if expirationTime := topicInfo.ExpirationTime; expirationTime != nil {
		tempTopicInfo.ExpirationTime = _TimeFromProtobuf(expirationTime)
	}

	if adminKey := topicInfo.AdminKey; adminKey != nil {
		tempTopicInfo.AdminKey, err = _KeyFromProtobuf(adminKey)
	}

	if submitKey := topicInfo.SubmitKey; submitKey != nil {
		tempTopicInfo.SubmitKey, err = _KeyFromProtobuf(submitKey)
	}

	if feeScheduleKey := topicInfo.FeeScheduleKey; feeScheduleKey != nil {
		tempTopicInfo.FeeScheduleKey, err = _KeyFromProtobuf(feeScheduleKey)
	}

	if len(topicInfo.FeeExemptKeyList) > 0 {
		for _, protoFeeExemptKey := range topicInfo.FeeExemptKeyList {
			feeExemptKey, _ := _KeyFromProtobuf(protoFeeExemptKey)
			tempTopicInfo.FeeExemptKeys = append(tempTopicInfo.FeeExemptKeys, feeExemptKey)
		}
	}

	if len(topicInfo.CustomFees) > 0 {
		for _, protoCustomFee := range topicInfo.CustomFees {
			customFee := CustomFee{FeeCollectorAccountID: _AccountIDFromProtobuf(protoCustomFee.FeeCollectorAccountId)}
			customFixedFee := _CustomFixedFeeFromProtobuf(protoCustomFee.FixedFee, customFee)
			tempTopicInfo.CustomFees = append(tempTopicInfo.CustomFees, customFixedFee)
		}
	}

	if autoRenewAccount := topicInfo.AutoRenewAccount; autoRenewAccount != nil {
		tempTopicInfo.AutoRenewAccountID = _AccountIDFromProtobuf(autoRenewAccount)
	}

	return tempTopicInfo, err
}

func (topicInfo *TopicInfo) _ToProtobuf() *services.ConsensusTopicInfo {
	txBody := &services.ConsensusTopicInfo{
		Memo:           topicInfo.TopicMemo,
		RunningHash:    topicInfo.RunningHash,
		SequenceNumber: topicInfo.SequenceNumber,
		ExpirationTime: &services.Timestamp{
			Seconds: int64(topicInfo.ExpirationTime.Second()),
			Nanos:   int32(topicInfo.ExpirationTime.Nanosecond()),
		},
		AdminKey:         topicInfo.AdminKey._ToProtoKey(),
		SubmitKey:        topicInfo.SubmitKey._ToProtoKey(),
		FeeScheduleKey:   topicInfo.FeeScheduleKey._ToProtoKey(),
		AutoRenewPeriod:  _DurationToProtobuf(topicInfo.AutoRenewPeriod),
		AutoRenewAccount: topicInfo.AutoRenewAccountID._ToProtobuf(),
		LedgerId:         topicInfo.LedgerID.ToBytes(),
	}

	if len(topicInfo.FeeExemptKeys) > 0 {
		for _, feeExemptKey := range topicInfo.FeeExemptKeys {
			txBody.FeeExemptKeyList = append(txBody.FeeExemptKeyList, feeExemptKey._ToProtoKey())
		}
	}

	if len(topicInfo.CustomFees) > 0 {
		for _, customFee := range topicInfo.CustomFees {
			protoCustomFee := &services.FixedCustomFee{
				FeeCollectorAccountId: customFee.FeeCollectorAccountID._ToProtobuf(),
				FixedFee:              customFee._ToProtobuf().GetFixedFee(),
			}
			txBody.CustomFees = append(txBody.CustomFees, protoCustomFee)
		}
	}

	return txBody
}

// ToBytes returns a byte array representation of the TopicInfo object
func (topicInfo TopicInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(topicInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TopicInfoFromBytes returns a TopicInfo object from a byte array
func TopicInfoFromBytes(data []byte) (TopicInfo, error) {
	if data == nil {
		return TopicInfo{}, errByteArrayNull
	}
	pb := services.ConsensusTopicInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TopicInfo{}, err
	}

	info, err := _TopicInfoFromProtobuf(&pb)
	if err != nil {
		return TopicInfo{}, err
	}

	return info, nil
}
