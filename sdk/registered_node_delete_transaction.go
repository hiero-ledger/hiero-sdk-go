package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

/**
 * A transaction to delete a registered node from the network
 * address book.
 *
 * This transaction, once complete, SHALL remove the identified registered
 * node from the network state.
 * This transaction MUST be signed by the existing entry `admin_key` or
 * authorized by the Hiero network governance structure.
 *
 */
type RegisteredNodeDeleteTransaction struct {
	*Transaction[*RegisteredNodeDeleteTransaction]
	registeredNodeId *uint64
}

// NewRegisteredNodeDeleteTransaction creates a new RegisteredNodeDeleteTransaction.
func NewRegisteredNodeDeleteTransaction() *RegisteredNodeDeleteTransaction {
	tx := &RegisteredNodeDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	return tx
}

func _RegisteredNodeDeleteTransactionFromProtobuf(tx Transaction[*RegisteredNodeDeleteTransaction], pb *services.TransactionBody) RegisteredNodeDeleteTransaction {
	nodeId := pb.GetRegisteredNodeDelete().GetRegisteredNodeId()

	registeredNodeDeleteTransaction := RegisteredNodeDeleteTransaction{
		registeredNodeId: &nodeId,
	}

	tx.childTransaction = &registeredNodeDeleteTransaction
	registeredNodeDeleteTransaction.Transaction = &tx
	return registeredNodeDeleteTransaction
}

// SetRegisteredNodeId sets the registered node ID to delete.
func (tx *RegisteredNodeDeleteTransaction) SetRegisteredNodeId(id uint64) *RegisteredNodeDeleteTransaction {
	tx._RequireNotFrozen()
	tx.registeredNodeId = &id
	return tx
}

// GetRegisteredNodeId returns the registered node ID.
func (tx *RegisteredNodeDeleteTransaction) GetRegisteredNodeId() uint64 {
	if tx.registeredNodeId == nil {
		return 0
	}
	return *tx.registeredNodeId
}

// ----------- Overridden functions ----------------

func (tx RegisteredNodeDeleteTransaction) getName() string {
	return "RegisteredNodeDeleteTransaction"
}

func (tx RegisteredNodeDeleteTransaction) validateNetworkOnIDs(_ *Client) error {
	return nil
}

func (tx RegisteredNodeDeleteTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_RegisteredNodeDelete{
		RegisteredNodeDelete: tx.buildProtoBody(),
	}

	return body
}

func (tx RegisteredNodeDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	body := tx.buildSchedulableTransactionBody()
	body.Data = &services.SchedulableTransactionBody_RegisteredNodeDelete{
		RegisteredNodeDelete: tx.buildProtoBody(),
	}

	return body, nil
}

func (tx RegisteredNodeDeleteTransaction) buildProtoBody() *services.RegisteredNodeDeleteTransactionBody {
	body := &services.RegisteredNodeDeleteTransactionBody{}

	if tx.registeredNodeId != nil {
		body.RegisteredNodeId = *tx.registeredNodeId
	}

	return body
}

func (tx RegisteredNodeDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().DeleteRegisteredNode,
	}
}

func (tx RegisteredNodeDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx RegisteredNodeDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
