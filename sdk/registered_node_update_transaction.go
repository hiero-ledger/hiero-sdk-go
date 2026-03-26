package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// RegisteredNodeUpdateTransaction updates an existing registered node in the network address book.
type RegisteredNodeUpdateTransaction struct {
	*Transaction[*RegisteredNodeUpdateTransaction]
	registeredNodeId *uint64
	adminKey         Key
	description      *string
	serviceEndpoints []RegisteredServiceEndpoint
}

// NewRegisteredNodeUpdateTransaction creates a new RegisteredNodeUpdateTransaction.
func NewRegisteredNodeUpdateTransaction() *RegisteredNodeUpdateTransaction {
	tx := &RegisteredNodeUpdateTransaction{}
	tx.Transaction = _NewTransaction(tx)
	return tx
}

func _RegisteredNodeUpdateTransactionFromProtobuf(tx Transaction[*RegisteredNodeUpdateTransaction], pb *services.TransactionBody) RegisteredNodeUpdateTransaction {
	var adminKey Key
	if pb.GetRegisteredNodeUpdate().GetAdminKey() != nil {
		adminKey, _ = _KeyFromProtobuf(pb.GetRegisteredNodeUpdate().GetAdminKey())
	}

	serviceEndpoints := make([]RegisteredServiceEndpoint, 0)
	for _, endpoint := range pb.GetRegisteredNodeUpdate().GetServiceEndpoint() {
		serviceEndpoints = append(serviceEndpoints, _RegisteredServiceEndpointFromProtobuf(endpoint))
	}

	nodeId := pb.GetRegisteredNodeUpdate().GetRegisteredNodeId()

	var description *string
	if pb.GetRegisteredNodeUpdate().GetDescription() != nil {
		description = &pb.GetRegisteredNodeUpdate().GetDescription().Value
	}

	registeredNodeUpdateTransaction := RegisteredNodeUpdateTransaction{
		registeredNodeId: &nodeId,
		adminKey:         adminKey,
		description:      description,
		serviceEndpoints: serviceEndpoints,
	}

	tx.childTransaction = &registeredNodeUpdateTransaction
	registeredNodeUpdateTransaction.Transaction = &tx
	return registeredNodeUpdateTransaction
}

// SetRegisteredNodeId sets the registered node ID to update.
func (tx *RegisteredNodeUpdateTransaction) SetRegisteredNodeId(id uint64) *RegisteredNodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.registeredNodeId = &id
	return tx
}

// GetRegisteredNodeId returns the registered node ID.
func (tx *RegisteredNodeUpdateTransaction) GetRegisteredNodeId() uint64 {
	if tx.registeredNodeId == nil {
		return 0
	}
	return *tx.registeredNodeId
}

// SetAdminKey sets the new admin key for the registered node.
func (tx *RegisteredNodeUpdateTransaction) SetAdminKey(key Key) *RegisteredNodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = key
	return tx
}

// GetAdminKey returns the admin key.
func (tx *RegisteredNodeUpdateTransaction) GetAdminKey() Key {
	return tx.adminKey
}

// SetDescription sets the description for the registered node.
func (tx *RegisteredNodeUpdateTransaction) SetDescription(description string) *RegisteredNodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.description = &description
	return tx
}

// GetDescription returns the description.
func (tx *RegisteredNodeUpdateTransaction) GetDescription() string {
	if tx.description != nil {
		return *tx.description
	}

	return ""
}

// SetServiceEndpoints sets the service endpoints, replacing any existing ones.
func (tx *RegisteredNodeUpdateTransaction) SetServiceEndpoints(endpoints []RegisteredServiceEndpoint) *RegisteredNodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = endpoints
	return tx
}

// GetServiceEndpoints returns the service endpoints.
func (tx *RegisteredNodeUpdateTransaction) GetServiceEndpoints() []RegisteredServiceEndpoint {
	return tx.serviceEndpoints
}

// AddServiceEndpoint adds a service endpoint.
func (tx *RegisteredNodeUpdateTransaction) AddServiceEndpoint(endpoint RegisteredServiceEndpoint) *RegisteredNodeUpdateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = append(tx.serviceEndpoints, endpoint)
	return tx
}

// ----------- Overridden functions ----------------

func (tx RegisteredNodeUpdateTransaction) getName() string {
	return "RegisteredNodeUpdateTransaction"
}

func (tx RegisteredNodeUpdateTransaction) validateNetworkOnIDs(_ *Client) error {
	return nil
}

func (tx RegisteredNodeUpdateTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_RegisteredNodeUpdate{
		RegisteredNodeUpdate: tx.buildProtoBody(),
	}

	return body
}

func (tx RegisteredNodeUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	body := tx.buildSchedulableTransactionBody()
	body.Data = &services.SchedulableTransactionBody_RegisteredNodeUpdate{
		RegisteredNodeUpdate: tx.buildProtoBody(),
	}

	return body, nil
}

func (tx RegisteredNodeUpdateTransaction) buildProtoBody() *services.RegisteredNodeUpdateTransactionBody {
	body := &services.RegisteredNodeUpdateTransactionBody{}

	if tx.registeredNodeId != nil {
		body.RegisteredNodeId = *tx.registeredNodeId
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.description != nil {
		body.Description = &wrapperspb.StringValue{Value: *tx.description}
	}

	for _, endpoint := range tx.serviceEndpoints {
		body.ServiceEndpoint = append(body.ServiceEndpoint, endpoint._ToProtobuf())
	}

	return body
}

func (tx RegisteredNodeUpdateTransaction) validateTransactionFields() error {
	return nil
}

func (tx RegisteredNodeUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().UpdateRegisteredNode,
	}
}

func (tx RegisteredNodeUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx RegisteredNodeUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
