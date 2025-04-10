package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
)

// SPDX-License-Identifier: Apache-2.0

type BatchTransaction struct {
	*Transaction[*BatchTransaction]
	innerTransactions   []TransactionInterface
	innerTransactionIDs []TransactionID
}

func NewBatchTransaction() *BatchTransaction {
	tx := &BatchTransaction{
		innerTransactions: make([]TransactionInterface, 0),
	}
	tx.Transaction = _NewTransaction(tx)
	return tx
}

func (tx *BatchTransaction) AddInnerTransaction(innerTransaction TransactionInterface) *BatchTransaction {
	tx._RequireNotFrozen()
	innerBaseTransaction := innerTransaction.getBaseTransaction()
	if !innerBaseTransaction.IsFrozen() {
		tx.freezeError = errInnerTransactionShouldBeFrozen
	}
	tx.innerTransactions = append(tx.innerTransactions, innerTransaction)
	tx.innerTransactionIDs = append(tx.innerTransactionIDs, innerBaseTransaction.GetTransactionID())
	return tx
}

func (tx *BatchTransaction) GetInnerTransactions() []TransactionInterface {
	return tx.innerTransactions
}

func (tx *BatchTransaction) GetInnerTransactionIDs() []TransactionID {
	return tx.innerTransactionIDs
}

// ----------- Overridden functions ----------------

func (tx BatchTransaction) getName() string {
	return "BatchTransaction"
}

func (tx BatchTransaction) validateNetworkOnIDs(client *Client) error {
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
