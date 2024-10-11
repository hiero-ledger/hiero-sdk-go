package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/generated/services"
)

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

/**
 * A transaction to create a new node in the network address book.
 * The transaction, once complete, enables a new consensus node
 * to join the network, and requires governing council authorization.
 * <p>
 * This transaction body SHALL be considered a "privileged transaction".
 * <p>
 *
 * - MUST be signed by the governing council.
 * - MUST be signed by the `Key` assigned to the
 *   `admin_key` field.
 * - The newly created node information SHALL be added to the network address
 *   book information in the network state.
 * - The new entry SHALL be created in "state" but SHALL NOT participate in
 *   network consensus and SHALL NOT be present in network "configuration"
 *   until the next "upgrade" transaction (as noted below).
 * - All new address book entries SHALL be added to the active network
 *   configuration during the next `freeze` transaction with the field
 *   `freeze_type` set to `PREPARE_UPGRADE`.
 *
 * ### Record Stream Effects
 * Upon completion the newly assigned `node_id` SHALL be in the transaction
 * receipt.
 */
type NodeCreateTransaction struct {
	Transaction
	accountID           *AccountID
	description         string
	gossipEndpoints     []Endpoint
	serviceEndpoints    []Endpoint
	gossipCaCertificate []byte
	grpcCertificateHash []byte
	adminKey            Key
}

func NewNodeCreateTransaction() *NodeCreateTransaction {
	tx := &NodeCreateTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _NodeCreateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *NodeCreateTransaction {
	adminKey, err := _KeyFromProtobuf(pb.GetNodeCreate().GetAdminKey())
	if err != nil {
		return &NodeCreateTransaction{}
	}

	accountID := _AccountIDFromProtobuf(pb.GetNodeCreate().GetAccountId())
	gossipEndpoints := make([]Endpoint, 0)
	for _, endpoint := range pb.GetNodeCreate().GetGossipEndpoint() {
		gossipEndpoints = append(gossipEndpoints, EndpointFromProtobuf(endpoint))
	}
	serviceEndpoints := make([]Endpoint, 0)
	for _, endpoint := range pb.GetNodeCreate().GetServiceEndpoint() {
		serviceEndpoints = append(serviceEndpoints, EndpointFromProtobuf(endpoint))
	}

	return &NodeCreateTransaction{
		Transaction:         transaction,
		accountID:           accountID,
		description:         pb.GetNodeCreate().GetDescription(),
		gossipEndpoints:     gossipEndpoints,
		serviceEndpoints:    serviceEndpoints,
		gossipCaCertificate: pb.GetNodeCreate().GetGossipCaCertificate(),
		grpcCertificateHash: pb.GetNodeCreate().GetGrpcCertificateHash(),
		adminKey:            adminKey,
	}
}

// GetAccountID AccountID of the node
func (tx *NodeCreateTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// SetAccountID get the AccountID of the node
func (tx *NodeCreateTransaction) SetAccountID(accountID AccountID) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetDescription get the description of the node
func (tx *NodeCreateTransaction) GetDescription() string {
	return tx.description
}

// SetDescription set the description of the node
func (tx *NodeCreateTransaction) SetDescription(description string) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.description = description
	return tx
}

// GetServiceEndpoints the list of service endpoints for gossip.
func (tx *NodeCreateTransaction) GetGossipEndpoints() []Endpoint {
	return tx.gossipEndpoints
}

// SetServiceEndpoints the list of service endpoints for gossip.
func (tx *NodeCreateTransaction) SetGossipEndpoints(gossipEndpoints []Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.gossipEndpoints = gossipEndpoints
	return tx
}

// AddGossipEndpoint add an endpoint for gossip to the list of service endpoints for gossip.
func (tx *NodeCreateTransaction) AddGossipEndpoint(endpoint Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.gossipEndpoints = append(tx.gossipEndpoints, endpoint)
	return tx
}

// GetServiceEndpoints the list of service endpoints for gRPC calls.
func (tx *NodeCreateTransaction) GetServiceEndpoints() []Endpoint {
	return tx.serviceEndpoints
}

// SetServiceEndpoints the list of service endpoints for gRPC calls.
func (tx *NodeCreateTransaction) SetServiceEndpoints(serviceEndpoints []Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = serviceEndpoints
	return tx
}

// AddServiceEndpoint the list of service endpoints for gRPC calls.
func (tx *NodeCreateTransaction) AddServiceEndpoint(endpoint Endpoint) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.serviceEndpoints = append(tx.serviceEndpoints, endpoint)
	return tx
}

// GetGossipCaCertificate the certificate used to sign gossip events.
func (tx *NodeCreateTransaction) GetGossipCaCertificate() []byte {
	return tx.gossipCaCertificate
}

// SetGossipCaCertificate the certificate used to sign gossip events.
// This value MUST be the DER encoding of the certificate presented.
func (tx *NodeCreateTransaction) SetGossipCaCertificate(gossipCaCertificate []byte) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.gossipCaCertificate = gossipCaCertificate
	return tx
}

// GetGrpcCertificateHash the hash of the node gRPC TLS certificate.
func (tx *NodeCreateTransaction) GetGrpcCertificateHash() []byte {
	return tx.grpcCertificateHash
}

// SetGrpcCertificateHash the hash of the node gRPC TLS certificate.
// This value MUST be a SHA-384 hash.
func (tx *NodeCreateTransaction) SetGrpcCertificateHash(grpcCertificateHash []byte) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.grpcCertificateHash = grpcCertificateHash
	return tx
}

// GetAdminKey an administrative key controlled by the node operator.
func (tx *NodeCreateTransaction) GetAdminKey() Key {
	return tx.adminKey
}

// SetAdminKey an administrative key controlled by the node operator.
func (tx *NodeCreateTransaction) SetAdminKey(adminKey Key) *NodeCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = adminKey
	return tx
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *NodeCreateTransaction) Sign(privateKey PrivateKey) *NodeCreateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *NodeCreateTransaction) SignWithOperator(client *Client) (*NodeCreateTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *NodeCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NodeCreateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *NodeCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *NodeCreateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *NodeCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *NodeCreateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *NodeCreateTransaction) Freeze() (*NodeCreateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *NodeCreateTransaction) FreezeWith(client *Client) (*NodeCreateTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this NodeCreateTransaction.
func (tx *NodeCreateTransaction) SetMaxTransactionFee(fee Hbar) *NodeCreateTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *NodeCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *NodeCreateTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this NodeCreateTransaction.
func (tx *NodeCreateTransaction) SetTransactionMemo(memo string) *NodeCreateTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this NodeCreateTransaction.
func (tx *NodeCreateTransaction) SetTransactionValidDuration(duration time.Duration) *NodeCreateTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *NodeCreateTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this NodeCreateTransaction.
func (tx *NodeCreateTransaction) SetTransactionID(transactionID TransactionID) *NodeCreateTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this NodeCreateTransaction.
func (tx *NodeCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *NodeCreateTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *NodeCreateTransaction) SetMaxRetry(count int) *NodeCreateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *NodeCreateTransaction) SetMaxBackoff(max time.Duration) *NodeCreateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *NodeCreateTransaction) SetMinBackoff(min time.Duration) *NodeCreateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *NodeCreateTransaction) SetLogLevel(level LogLevel) *NodeCreateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *NodeCreateTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *NodeCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *NodeCreateTransaction) getName() string {
	return "NodeCreateTransaction"
}

func (tx *NodeCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *NodeCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_NodeCreate{
			NodeCreate: tx.buildProtoBody(),
		},
	}
}

func (tx *NodeCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_NodeCreate{
			NodeCreate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *NodeCreateTransaction) buildProtoBody() *services.NodeCreateTransactionBody {
	body := &services.NodeCreateTransactionBody{
		Description: tx.description,
	}

	if tx.accountID != nil {
		body.AccountId = tx.accountID._ToProtobuf()
	}

	for _, endpoint := range tx.gossipEndpoints {
		body.GossipEndpoint = append(body.GossipEndpoint, endpoint._ToProtobuf())
	}

	for _, endpoint := range tx.serviceEndpoints {
		body.ServiceEndpoint = append(body.ServiceEndpoint, endpoint._ToProtobuf())
	}

	if tx.gossipCaCertificate != nil {
		body.GossipCaCertificate = tx.gossipCaCertificate
	}

	if tx.grpcCertificateHash != nil {
		body.GrpcCertificateHash = tx.grpcCertificateHash
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	return body
}

func (tx *NodeCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetAddressBook().CreateNode,
	}
}

func (tx *NodeCreateTransaction) preFreezeWith(client *Client) {
	// No special actions needed.
}

func (tx *NodeCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
