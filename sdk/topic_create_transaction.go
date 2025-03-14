package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// A TopicCreateTransaction is for creating a new Topic on HCS.
type TopicCreateTransaction struct {
	*Transaction[*TopicCreateTransaction]
	autoRenewAccountID *AccountID
	adminKey           Key
	submitKey          Key
	feeScheduleKey     Key
	feeExemptKeys      []Key
	customFees         []*CustomFixedFee
	memo               string
	autoRenewPeriod    *time.Duration
}

// NewTopicCreateTransaction creates a TopicCreateTransaction transaction which can be
// used to construct and execute a  Create Topic Transaction.
func NewTopicCreateTransaction() *TopicCreateTransaction {
	tx := &TopicCreateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx.SetAutoRenewPeriod(7890000 * time.Second)
	tx.SetMaxTransactionFee(NewHbar(25))

	return tx
}

func _TopicCreateTransactionFromProtobuf(tx Transaction[*TopicCreateTransaction], pb *services.TransactionBody) TopicCreateTransaction {
	var adminKey Key
	if pb.GetConsensusCreateTopic().GetAdminKey() != nil {
		adminKey, _ = _KeyFromProtobuf(pb.GetConsensusCreateTopic().GetAdminKey())
	}
	var submitKey Key
	if pb.GetConsensusCreateTopic().GetSubmitKey() != nil {
		submitKey, _ = _KeyFromProtobuf(pb.GetConsensusCreateTopic().GetSubmitKey())
	}
	var feeScheduleKey Key
	if pb.GetConsensusCreateTopic().GetFeeScheduleKey() != nil {
		feeScheduleKey, _ = _KeyFromProtobuf(pb.GetConsensusCreateTopic().GetFeeScheduleKey())
	}
	var feeExemptKeys []Key = nil
	if pb.GetConsensusCreateTopic().GetFeeExemptKeyList() != nil {
		protobufKeysList := pb.GetConsensusCreateTopic().GetFeeExemptKeyList()
		for _, key := range protobufKeysList {
			key, _ := _KeyFromProtobuf(key)
			feeExemptKeys = append(feeExemptKeys, key)
		}
	}
	var customFixedFees []*CustomFixedFee = nil
	if pb.GetConsensusCreateTopic().GetCustomFees() != nil {
		protobufCustomFixedFees := pb.GetConsensusCreateTopic().GetCustomFees()
		for _, customFixedFee := range protobufCustomFixedFees {
			customFee := CustomFee{FeeCollectorAccountID: _AccountIDFromProtobuf(customFixedFee.FeeCollectorAccountId)}
			customFixedFees = append(customFixedFees, _CustomFixedFeeFromProtobuf(customFixedFee.FixedFee, customFee))
		}
	}
	var autoRenew *time.Duration
	if pb.GetConsensusCreateTopic().GetAutoRenewPeriod() != nil {
		autoRenewVal := _DurationFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewPeriod())
		autoRenew = &autoRenewVal
	}
	topicCreateTransaction := TopicCreateTransaction{
		autoRenewAccountID: _AccountIDFromProtobuf(pb.GetConsensusCreateTopic().GetAutoRenewAccount()),
		adminKey:           adminKey,
		submitKey:          submitKey,
		feeScheduleKey:     feeScheduleKey,
		feeExemptKeys:      feeExemptKeys,
		customFees:         customFixedFees,
		memo:               pb.GetConsensusCreateTopic().GetMemo(),
		autoRenewPeriod:    autoRenew,
	}

	tx.childTransaction = &topicCreateTransaction
	topicCreateTransaction.Transaction = &tx
	return topicCreateTransaction
}

// SetAdminKey sets the key required to update or delete the topic. If unspecified, anyone can increase the topic's
// expirationTime.
func (tx *TopicCreateTransaction) SetAdminKey(publicKey Key) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = publicKey
	return tx
}

// GetAdminKey returns the key required to update or delete the topic
func (tx *TopicCreateTransaction) GetAdminKey() (Key, error) {
	return tx.adminKey, nil
}

// SetSubmitKey sets the key required for submitting messages to the topic. If unspecified, all submissions are allowed.
func (tx *TopicCreateTransaction) SetSubmitKey(publicKey Key) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.submitKey = publicKey
	return tx
}

// GetSubmitKey returns the key required for submitting messages to the topic
func (tx *TopicCreateTransaction) GetSubmitKey() (Key, error) {
	return tx.submitKey, nil
}

// SetFeeScheduleKey sets the key which allows updates to the new topic’s fees.
func (tx *TopicCreateTransaction) SetFeeScheduleKey(publicKey Key) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.feeScheduleKey = publicKey
	return tx
}

// GetFeeScheduleKey returns the key which allows updates to the new topic’s fees.
func (tx *TopicCreateTransaction) GetFeeScheduleKey() Key {
	return tx.feeScheduleKey
}

// SetFeeExemptKeys sets the keys that will be exempt from paying fees.
func (tx *TopicCreateTransaction) SetFeeExemptKeys(keys []Key) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.feeExemptKeys = keys
	return tx
}

// AddFeeExemptKey adds a key that will be exempt from paying fees.
func (tx *TopicCreateTransaction) AddFeeExemptKey(key Key) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.feeExemptKeys = append(tx.feeExemptKeys, key)
	return tx
}

// ClearFeeExemptKeys removes all keys that will be exempt from paying fees.
func (tx *TopicCreateTransaction) ClearFeeExemptKeys() *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.feeExemptKeys = []Key{}
	return tx
}

// GetFeeExemptKeys returns the keys that will be exempt from paying fees.
func (tx *TopicCreateTransaction) GetFeeExemptKeys() []Key {
	return tx.feeExemptKeys
}

// SetCustomFees Sets the fixed fees to assess when a message is submitted to the new topic.
func (tx *TopicCreateTransaction) SetCustomFees(fees []*CustomFixedFee) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.customFees = fees
	return tx
}

// AddCustomFee adds a fixed fee to assess when a message is submitted to the new topic.
func (tx *TopicCreateTransaction) AddCustomFee(fee *CustomFixedFee) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.customFees = append(tx.customFees, fee)
	return tx
}

// ClearCustomFees removes all custom fees to assess when a message is submitted to the new topic.
func (tx *TopicCreateTransaction) ClearCustomFees() *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.customFees = []*CustomFixedFee{}
	return tx
}

// GetCustomFees returns the fixed fees to assess when a message is submitted to the new topic.
func (tx *TopicCreateTransaction) GetCustomFees() []*CustomFixedFee {
	return tx.customFees
}

// SetTopicMemo sets a short publicly visible memo about the topic. No guarantee of uniqueness.
func (tx *TopicCreateTransaction) SetTopicMemo(memo string) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// GetTopicMemo returns the memo for this topic
func (tx *TopicCreateTransaction) GetTopicMemo() string {
	return tx.memo
}

// SetAutoRenewPeriod sets the initial lifetime of the topic and the amount of time to extend the topic's lifetime
// automatically at expirationTime if the autoRenewAccount is configured and has sufficient funds.
//
// Required. Limited to a maximum of 90 days (server-sIDe configuration which may change).
func (tx *TopicCreateTransaction) SetAutoRenewPeriod(period time.Duration) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &period
	return tx
}

// GetAutoRenewPeriod returns the auto renew period for this topic
func (tx *TopicCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return *tx.autoRenewPeriod
	}

	return time.Duration(0)
}

// SetAutoRenewAccountID sets an optional account to be used at the topic's expirationTime to extend the life of the
// topic. The topic lifetime will be extended up to a maximum of the autoRenewPeriod or however long the topic can be
// extended using all funds on the account (whichever is the smaller duration/amount).
//
// If specified, there must be an adminKey and the autoRenewAccount must sign this transaction.
func (tx *TopicCreateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) *TopicCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &autoRenewAccountID
	return tx
}

// GetAutoRenewAccountID returns the auto renew account ID for this topic
func (tx *TopicCreateTransaction) GetAutoRenewAccountID() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// ----------- Overridden functions ----------------

func (tx TopicCreateTransaction) getName() string {
	return "TopicCreateTransaction"
}

func (tx TopicCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.autoRenewAccountID != nil {
		if err := tx.autoRenewAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TopicCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: tx.buildProtoBody(),
		},
	}
}

func (tx TopicCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TopicCreateTransaction) buildProtoBody() *services.ConsensusCreateTopicTransactionBody {
	body := &services.ConsensusCreateTopicTransactionBody{
		Memo: tx.memo,
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.autoRenewAccountID != nil {
		body.AutoRenewAccount = tx.autoRenewAccountID._ToProtobuf()
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.submitKey != nil {
		body.SubmitKey = tx.submitKey._ToProtoKey()
	}

	if tx.feeScheduleKey != nil {
		body.FeeScheduleKey = tx.feeScheduleKey._ToProtoKey()
	}

	if tx.feeExemptKeys != nil {
		protobufKeysList := make([]*services.Key, 0)
		for _, key := range tx.feeExemptKeys {
			protobufKeysList = append(protobufKeysList, key._ToProtoKey())
		}
		body.FeeExemptKeyList = protobufKeysList
	}

	if tx.customFees != nil {
		protobufCustomFixedFees := make([]*services.FixedCustomFee, 0)
		for _, customFixedFee := range tx.customFees {
			protobufCustomFixedFees = append(protobufCustomFixedFees, customFixedFee._ToTopicFeeProtobuf())
		}
		body.CustomFees = protobufCustomFixedFees
	}

	return body
}

func (tx TopicCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetTopic().CreateTopic,
	}
}

// TODO	 Temporarily disabled due to issues with consensus node version 0.60.
// This will be reintroduced once all networks (previewnet, testnet, mainnet)
// are on version 0.60.
// func (tx TopicCreateTransaction) preFreezeWith(client *Client, self TransactionInterface) {
// 	if selfTopicCreate, ok := self.(*TopicCreateTransaction); ok {
// 		if selfTopicCreate.GetAutoRenewAccountID()._IsZero() && tx.Transaction.transactionIDs != nil && !tx.Transaction.transactionIDs._IsEmpty() {
// 			selfTopicCreate.SetAutoRenewAccountID(*tx.Transaction.GetTransactionID().AccountID)
// 		}

// 		if selfTopicCreate.GetAutoRenewAccountID()._IsZero() && client != nil && selfTopicCreate.adminKey != nil {
// 			selfTopicCreate.SetAutoRenewAccountID(client.GetOperatorAccountID())
// 		}
// 	}
// }

func (tx TopicCreateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TopicCreateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
