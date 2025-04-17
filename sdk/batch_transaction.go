package hiero

import (
	"reflect"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// SPDX-License-Identifier: Apache-2.0

// BatchTransaction executes multiple transactions in a single consensus event.
// This allows for atomic execution of multiple transactions, where they either
// all succeed or all fail together.
//
// Requirements:
// - All inner transactions must be frozen before being added to the batch
// - All inner transactions must have a batch key set
// - All inner transactions must be signed as required for each individual transaction
// - The BatchTransaction must be signed by all batch keys of the inner transactions
// - Certain transaction types (FreezeTransaction, BatchTransaction) are not allowed in a batch
//
// Important notes:
// - Fees are assessed for each inner transaction separately
// - The maximum number of inner transactions in a batch is limited to 25
// - Inner transactions cannot be scheduled transactions
//
// Example usage:
//
//	batchKey := PrivateKeyGenerateECDSA()
//
//	// Create and prepare inner transaction
//	transaction := NewTransferTransaction().
//		AddHbarTransfer(sender, amount.Negated()).
//		AddHbarTransfer(receiver, amount).
//		Batchify(client, batchKey)
//
//	// Create and execute batch transaction
//	response, err := NewBatchTransaction().
//		AddInnerTransaction(transaction).
//		FreezeWith(client).
//		Sign(batchKey).
//		Execute(client)
type BatchTransaction struct {
	*Transaction[*BatchTransaction]
	innerTransactions []TransactionInterface
}

// blacklistedTransactions is a list of transaction types that are not allowed in a batch transaction.
// These transactions are prohibited due to their special nature or network-level implications.
var blacklistedTransactions = []reflect.Type{
	reflect.TypeOf(&FreezeTransaction{}),
	reflect.TypeOf(&BatchTransaction{}),
}

// NewBatchTransaction creates a new empty BatchTransaction.
func NewBatchTransaction() *BatchTransaction {
	tx := &BatchTransaction{
		innerTransactions: make([]TransactionInterface, 0),
	}
	tx.Transaction = _NewTransaction(tx)
	return tx
}

func _BatchTransactionFromProtobuf(tx Transaction[*BatchTransaction], pb *services.TransactionBody) BatchTransaction {
	var innerTransactions []TransactionInterface
	for _, innerTransaction := range pb.GetAtomicBatch().Transactions {
		var transaction services.Transaction
		transaction.SignedTransactionBytes = innerTransaction
		transactionBytes, _ := protobuf.Marshal(&transaction)
		transactionFromBytes, _ := TransactionFromBytes(transactionBytes)
		innerTransactions = append(innerTransactions, transactionFromBytes)
	}

	batchTransaction := BatchTransaction{
		innerTransactions: innerTransactions,
	}
	tx.childTransaction = &batchTransaction
	batchTransaction.Transaction = &tx
	return batchTransaction
}

// validateInnerTransaction validates if a transaction is allowed in a batch transaction.
// A transaction is valid if:
// - It is not a blacklisted type (FreezeTransaction or BatchTransaction)
// - It is frozen
// - It has a batch key set
func (tx *BatchTransaction) validateInnerTransaction(innerTransaction TransactionInterface) error {
	// Check for nil transaction
	if innerTransaction == nil {
		return errInnerTransactionNil
	}

	// Validate transaction type is not blacklisted
	txType := reflect.TypeOf(innerTransaction)
	for _, blacklistedType := range blacklistedTransactions {
		if txType == blacklistedType {
			return errTransactionTypeNotAllowed
		}
	}

	innerBaseTransaction := innerTransaction.getBaseTransaction()
	if !innerBaseTransaction.IsFrozen() {
		return errInnerTransactionShouldBeFrozen
	}

	if innerBaseTransaction.GetBatchKey() == nil {
		return errBatchKeyNotSet
	}

	return nil
}

// SetInnerTransactions sets the list of transactions to be executed as part of this BatchTransaction.
//
// Requirements for each inner transaction:
// - Must be frozen (use Transaction.Freeze() or Transaction.FreezeWith(client))
// - Must have a batch key set (use Transaction.SetBatchKey(key) or Transaction.Batchify(client, key))
// - Must not be a blacklisted transaction type
//
// Returns:
// - *BatchTransaction: The BatchTransaction instance for fluent API calls
// - error: If any validation fails
func (tx *BatchTransaction) SetInnerTransactions(transactions []TransactionInterface) *BatchTransaction {
	tx._RequireNotFrozen()

	// Validate all transactions before setting
	for _, transaction := range transactions {
		if err := tx.validateInnerTransaction(transaction); err != nil {
			tx.freezeError = err
			return tx
		}
	}

	tx.innerTransactions = make([]TransactionInterface, len(transactions))
	copy(tx.innerTransactions, transactions)
	return tx
}

// AddInnerTransaction appends a transaction to the list of transactions this BatchTransaction will execute.
//
// Requirements for the inner transaction:
// - Must be frozen (use Transaction.Freeze() or Transaction.FreezeWith(client))
// - Must have a batch key set (use Transaction.SetBatchKey(key) or Transaction.Batchify(client, key))
// - Must not be a blacklisted transaction type
//
// Returns:
// - *BatchTransaction: The BatchTransaction instance for fluent API calls
func (tx *BatchTransaction) AddInnerTransaction(innerTransaction TransactionInterface) *BatchTransaction {
	tx._RequireNotFrozen()

	if err := tx.validateInnerTransaction(innerTransaction); err != nil {
		tx.freezeError = err
		return tx
	}

	tx.innerTransactions = append(tx.innerTransactions, innerTransaction)
	return tx
}

// GetInnerTransactions returns the list of transactions this BatchTransaction is currently configured to execute.
func (tx *BatchTransaction) GetInnerTransactions() []TransactionInterface {
	return tx.innerTransactions
}

// GetInnerTransactionIDs returns the list of transaction IDs of each inner transaction of this BatchTransaction.
//
// This method is particularly useful after execution to:
// - Track individual transaction results
// - Query receipts for specific inner transactions
// - Monitor the status of each transaction in the batch
//
// NOTE: Transaction IDs will only be meaningful after the batch transaction has been
// executed or the IDs have been explicitly set on the inner transactions.
func (tx *BatchTransaction) GetInnerTransactionIDs() []TransactionID {
	var ids []TransactionID
	for _, tx := range tx.innerTransactions {
		ids = append(ids, tx.getBaseTransaction().GetTransactionID())
	}
	return ids
}

// ----------- Overridden functions ----------------

func (tx BatchTransaction) getName() string {
	return "BatchTransaction"
}

func (tx BatchTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, innerTx := range tx.innerTransactions {
		if err := innerTx.validateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx BatchTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionID:            tx.transactionID._ToProtobuf(),
		TransactionFee:           tx.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		Memo:                     tx.Transaction.memo,
		Data: &services.TransactionBody_AtomicBatch{
			AtomicBatch: tx.buildProtoBody(),
		},
	}
}

func (tx BatchTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `BatchTransaction`")
}

func (tx BatchTransaction) buildProtoBody() *services.AtomicBatchTransactionBody {
	body := &services.AtomicBatchTransactionBody{}
	for _, innerTransaction := range tx.innerTransactions {
		request := innerTransaction.makeRequest()
		switch request := request.(type) {
		case *services.Transaction:
			body.Transactions = append(body.Transactions, request.GetSignedTransactionBytes())
		default:
			// do nothing
			return nil
		}
	}

	return body
}

func (tx BatchTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetUtil().AtomicBatch,
	}
}

func (tx BatchTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx BatchTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
