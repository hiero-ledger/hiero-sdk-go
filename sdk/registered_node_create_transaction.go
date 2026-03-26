package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

/**
 * A transaction to create a new registered node in the network
 * address book.
 *
 * This transaction, once complete, SHALL add a new registered node to the
 * network state.
 * The new registered node SHALL be visible and discoverable upon
 * completion of this transaction.
 *
 */
type RegisteredNodeCreateTransaction struct {
	*Transaction[*RegisteredNodeCreateTransaction]
	adminKey         Key
	description      string
	serviceEndpoints []RegisteredServiceEndpoint
}

// NewRegisteredNodeCreateTransaction creates a new RegisteredNodeCreateTransaction.
func NewRegisteredNodeCreateTransaction() *RegisteredNodeCreateTransaction {
	tx := &RegisteredNodeCreateTransaction{}
	tx.Transaction = _NewTransaction(tx)
	return tx
}

func _RegisteredNodeCreateTransactionFromProtobuf(tx Transaction[*RegisteredNodeCreateTransaction], pb *services.TransactionBody) RegisteredNodeCreateTransaction {
	var adminKey Key
	if pb.GetRegisteredNodeCreate().GetAdminKey() != nil {
		adminKey, _ = _KeyFromProtobuf(pb.GetRegisteredNodeCreate().GetAdminKey())
	}

	serviceEndpoints := make([]RegisteredServiceEndpoint, 0)
	for _, endpoint := range pb.GetRegisteredNodeCreate().GetServiceEndpoint() {
		serviceEndpoints = append(serviceEndpoints, _RegisteredServiceEndpointFromProtobuf(endpoint))
	}

	registeredNodeCreateTransaction := RegisteredNodeCreateTransaction{
		adminKey:         adminKey,
		description:      pb.GetRegisteredNodeCreate().GetDescription(),
		serviceEndpoints: serviceEndpoints,
	}

	tx.childTransaction = &registeredNodeCreateTransaction
	registeredNodeCreateTransaction.Transaction = &tx
	return registeredNodeCreateTransaction
}

// SetAdminKey sets the administrative key for the registered node.
func (tx *RegisteredNodeCreateTransaction) SetAdminKey(key Key) *RegisteredNodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = key
	return tx
}

// GetAdminKey returns the administrative key for the registered node.
func (tx *RegisteredNodeCreateTransaction) GetAdminKey() Key {
	return tx.adminKey
}

// SetDescription sets the description for the registered node.
func (tx *RegisteredNodeCreateTransaction) SetDescription(description string) *RegisteredNodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.description = description
	return tx
}

// GetDescription returns the description for the registered node.
func (tx *RegisteredNodeCreateTransaction) GetDescription() string {
	return tx.description
}

// SetServiceEndpoints sets the service endpoints for the registered node.
func (tx *RegisteredNodeCreateTransaction) SetServiceEndpoints(endpoints []RegisteredServiceEndpoint) *RegisteredNodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = endpoints
	return tx
}

// GetServiceEndpoints returns the service endpoints for the registered node.
func (tx *RegisteredNodeCreateTransaction) GetServiceEndpoints() []RegisteredServiceEndpoint {
	return tx.serviceEndpoints
}

// AddServiceEndpoint adds a service endpoint to the registered node.
func (tx *RegisteredNodeCreateTransaction) AddServiceEndpoint(endpoint RegisteredServiceEndpoint) *RegisteredNodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = append(tx.serviceEndpoints, endpoint)
	return tx
}

// ----------- Overridden functions ----------------

func (tx RegisteredNodeCreateTransaction) getName() string {
	return "RegisteredNodeCreateTransaction"
}

func (tx RegisteredNodeCreateTransaction) validateNetworkOnIDs(_ *Client) error {
	return nil
}

func (tx RegisteredNodeCreateTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_RegisteredNodeCreate{
		RegisteredNodeCreate: tx.buildProtoBody(),
	}

	return body
}

func (tx RegisteredNodeCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	body := tx.buildSchedulableTransactionBody()
	body.Data = &services.SchedulableTransactionBody_RegisteredNodeCreate{
		RegisteredNodeCreate: tx.buildProtoBody(),
	}

	return body, nil
}

func (tx RegisteredNodeCreateTransaction) buildProtoBody() *services.RegisteredNodeCreateTransactionBody {
	body := &services.RegisteredNodeCreateTransactionBody{
		Description: tx.description,
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	for _, endpoint := range tx.serviceEndpoints {
		body.ServiceEndpoint = append(body.ServiceEndpoint, endpoint._ToProtobuf())
	}

	return body
}

func (tx RegisteredNodeCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().CreateRegisteredNode,
	}
}

func (tx RegisteredNodeCreateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx RegisteredNodeCreateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
