package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
)

// SPDX-License-Identifier: Apache-2.0

type LambdaSStoreTransaction struct {
	*Transaction[*LambdaSStoreTransaction]
	hookId         HookId
	storageUpdates []LambdaStorageUpdate
}

// Adds or removes key/value pairs in the storage of a lambda.
func NewLambdaSStoreTransaction() *LambdaSStoreTransaction {
	tx := &LambdaSStoreTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func lambdaSStoreTransactionFromProtobuf(tx Transaction[*LambdaSStoreTransaction], pb *services.TransactionBody) LambdaSStoreTransaction {
	protoBody := pb.GetLambdaSstore()
	storageUpdates := make([]LambdaStorageUpdate, 0)
	for _, storageUpdate := range protoBody.GetStorageUpdates() {
		storageUpdates = append(storageUpdates, lambdaStorageUpdateFromProtobuf(storageUpdate))
	}
	lambdaSStoreTransaction := LambdaSStoreTransaction{
		hookId:         hookIdFromProtobuf(protoBody.GetHookId()),
		storageUpdates: storageUpdates,
	}

	tx.childTransaction = &lambdaSStoreTransaction
	lambdaSStoreTransaction.Transaction = &tx
	return lambdaSStoreTransaction
}

// SetHookId sets the hook ID for the LambdaSStoreTransaction.
func (tx *LambdaSStoreTransaction) SetHookId(hookId HookId) *LambdaSStoreTransaction {
	tx._RequireNotFrozen()
	tx.hookId = hookId
	return tx
}

// GetHookId returns the hook ID for the LambdaSStoreTransaction.
func (tx LambdaSStoreTransaction) GetHookId() HookId {
	return tx.hookId
}

// AddStorageUpdate adds a storage update to the LambdaSStoreTransaction.
func (tx *LambdaSStoreTransaction) AddStorageUpdate(storageUpdate LambdaStorageUpdate) *LambdaSStoreTransaction {
	tx._RequireNotFrozen()
	tx.storageUpdates = append(tx.storageUpdates, storageUpdate)
	return tx
}

// SetStorageUpdates sets the storage updates for the LambdaSStoreTransaction.
func (tx *LambdaSStoreTransaction) SetStorageUpdates(storageUpdates []LambdaStorageUpdate) *LambdaSStoreTransaction {
	tx._RequireNotFrozen()
	tx.storageUpdates = storageUpdates
	return tx
}

// GetStorageUpdates returns the storage updates for the LambdaSStoreTransaction.
func (tx LambdaSStoreTransaction) GetStorageUpdates() []LambdaStorageUpdate {
	return tx.storageUpdates
}

// ----------- Overridden functions ----------------

func (tx LambdaSStoreTransaction) getName() string {
	return "LambdaSStoreTransaction"
}

func (tx LambdaSStoreTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := tx.hookId.validateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (tx LambdaSStoreTransaction) build() *services.TransactionBody {
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

func (tx LambdaSStoreTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LambdaSStoreTransaction`")
}

func (tx LambdaSStoreTransaction) buildProtoBody() *services.LambdaSStoreTransactionBody {
	body := &services.LambdaSStoreTransactionBody{
		HookId: tx.hookId.toProtobuf(),
	}

	storageUpdates := make([]*services.LambdaStorageUpdate, 0)
	for _, storageUpdate := range tx.storageUpdates {
		storageUpdates = append(storageUpdates, storageUpdate.toProtobuf())
	}
	body.StorageUpdates = storageUpdates

	return body
}

func (tx LambdaSStoreTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().LambdaSStore,
	}
}

func (tx LambdaSStoreTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx LambdaSStoreTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
