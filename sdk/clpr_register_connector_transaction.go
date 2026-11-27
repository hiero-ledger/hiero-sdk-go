package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprRegisterConnectorTransaction registers the ownership commitment of a Connector on an existing
// CLPR Channel, the commit phase of the Connector commit-reveal (HIP-1535). Derive the Connector ID
// with ComputeClprConnectorID and the commitment with ComputeClprConnectorOwnershipCommitment, then
// reveal it with ClprCompleteConnectorTransaction.
type ClprRegisterConnectorTransaction struct {
	*Transaction[*ClprRegisterConnectorTransaction]
	fieldErrors _ClprFieldErrors
	commitment  []byte
	channelID   []byte
}

// NewClprRegisterConnectorTransaction creates a ClprRegisterConnectorTransaction.
func NewClprRegisterConnectorTransaction() *ClprRegisterConnectorTransaction {
	tx := &ClprRegisterConnectorTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprRegisterConnectorTransactionFromProtobuf(tx Transaction[*ClprRegisterConnectorTransaction], pb *services.TransactionBody) ClprRegisterConnectorTransaction {
	body := pb.GetClprRegisterConnector()

	clprRegisterConnectorTransaction := ClprRegisterConnectorTransaction{
		commitment: _ClprBytesFromProtobuf(body.GetCommitment()),
		channelID:  _ClprBytesFromProtobuf(body.GetChannelId()),
	}

	tx.childTransaction = &clprRegisterConnectorTransaction
	clprRegisterConnectorTransaction.Transaction = &tx
	return clprRegisterConnectorTransaction
}

// SetCommitment sets the 32-byte keccak256(connectorID || ownership public key) commitment.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprRegisterConnectorTransaction) SetCommitment(commitment []byte) *ClprRegisterConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("commitment", _ClprCheckLength("commitment", commitment)) {
		return tx
	}
	tx.commitment = _ClprCopyBytes(commitment)
	return tx
}

// GetCommitment returns a copy of the ownership commitment, or nil if it is not set.
func (tx *ClprRegisterConnectorTransaction) GetCommitment() []byte {
	return _ClprCopyBytes(tx.commitment)
}

// SetChannelID sets the 32-byte ID of the existing Channel the Connector registers on.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprRegisterConnectorTransaction) SetChannelID(channelID []byte) *ClprRegisterConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("channelID", _ClprCheckLength("channelID", channelID)) {
		return tx
	}
	tx.channelID = _ClprCopyBytes(channelID)
	return tx
}

// GetChannelID returns a copy of the Channel ID, or nil if it is not set.
func (tx *ClprRegisterConnectorTransaction) GetChannelID() []byte {
	return _ClprCopyBytes(tx.channelID)
}

// ----------- Overridden functions ----------------

func (tx ClprRegisterConnectorTransaction) getName() string {
	return "ClprRegisterConnectorTransaction"
}

func (tx ClprRegisterConnectorTransaction) validateNetworkOnIDs(*Client) error {
	return nil
}

func (tx ClprRegisterConnectorTransaction) validateTransactionFields() error {
	if err := _ClprRequireID(tx.getName(), "commitment", tx.commitment, tx.fieldErrors); err != nil {
		return err
	}

	return _ClprRequireID(tx.getName(), "channelID", tx.channelID, tx.fieldErrors)
}

func (tx ClprRegisterConnectorTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprRegisterConnectorTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprRegisterConnector{
		ClprRegisterConnector: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprRegisterConnectorTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprRegisterConnectorTransaction) buildProtoBody() *services.ClprRegisterConnectorTransactionBody {
	return &services.ClprRegisterConnectorTransactionBody{
		Commitment: _ClprCopyBytes(tx.commitment),
		ChannelId:  _ClprCopyBytes(tx.channelID),
	}
}

func (tx ClprRegisterConnectorTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().RegisterConnector,
	}
}

func (tx ClprRegisterConnectorTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
