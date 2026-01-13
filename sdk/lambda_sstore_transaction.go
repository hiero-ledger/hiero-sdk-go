package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
)

// SPDX-License-Identifier: Apache-2.0

type HookStoreTransaction struct {
	*Transaction[*HookStoreTransaction]
	hookId         HookId
	storageUpdates []EvmHookStorageUpdate
}

// Adds or removes key/value pairs in the storage of a lambda.
func NewHookStoreTransaction() *HookStoreTransaction {
	tx := &HookStoreTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func lambdaSStoreTransactionFromProtobuf(tx Transaction[*HookStoreTransaction], pb *services.TransactionBody) HookStoreTransaction {
	protoBody := pb.GetLambdaSstore()
	storageUpdates := make([]EvmHookStorageUpdate, 0)
	for _, storageUpdate := range protoBody.GetStorageUpdates() {
		storageUpdates = append(storageUpdates, lambdaStorageUpdateFromProtobuf(storageUpdate))
	}
	lambdaSStoreTransaction := HookStoreTransaction{
		hookId:         hookIdFromProtobuf(protoBody.GetHookId()),
		storageUpdates: storageUpdates,
	}

	tx.childTransaction = &lambdaSStoreTransaction
	lambdaSStoreTransaction.Transaction = &tx
	return lambdaSStoreTransaction
}

// SetHookId sets the hook ID for the LambdaSStoreTransaction.
func (tx *HookStoreTransaction) SetHookId(hookId HookId) *HookStoreTransaction {
	tx._RequireNotFrozen()
	tx.hookId = hookId
	return tx
}

// GetHookId returns the hook ID for the LambdaSStoreTransaction.
func (tx HookStoreTransaction) GetHookId() HookId {
	return tx.hookId
}

// AddStorageUpdate adds a storage update to the LambdaSStoreTransaction.
func (tx *HookStoreTransaction) AddStorageUpdate(storageUpdate EvmHookStorageUpdate) *HookStoreTransaction {
	tx._RequireNotFrozen()
	tx.storageUpdates = append(tx.storageUpdates, storageUpdate)
	return tx
}

// SetStorageUpdates sets the storage updates for the LambdaSStoreTransaction.
func (tx *HookStoreTransaction) SetStorageUpdates(storageUpdates []EvmHookStorageUpdate) *HookStoreTransaction {
	tx._RequireNotFrozen()
	tx.storageUpdates = storageUpdates
	return tx
}

// GetStorageUpdates returns the storage updates for the LambdaSStoreTransaction.
func (tx HookStoreTransaction) GetStorageUpdates() []EvmHookStorageUpdate {
	return tx.storageUpdates
}

// ----------- Overridden functions ----------------

func (tx HookStoreTransaction) getName() string {
	return "LambdaSStoreTransaction"
}

func (tx HookStoreTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := tx.hookId.validateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (tx HookStoreTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_LambdaSstore{
			LambdaSstore: tx.buildProtoBody(),
		},
	}
}

func (tx HookStoreTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LambdaSStoreTransaction`")
}

func (tx HookStoreTransaction) buildProtoBody() *services.LambdaSStoreTransactionBody {
	body := &services.LambdaSStoreTransactionBody{
		HookId: tx.hookId.toProtobuf(),
	}

	for _, storageUpdate := range tx.storageUpdates {
		body.StorageUpdates = append(body.StorageUpdates, storageUpdate.toProtobuf())
	}

	return body
}

func (tx HookStoreTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().LambdaSStore,
	}
}

func (tx HookStoreTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx HookStoreTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
