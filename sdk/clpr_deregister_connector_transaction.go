package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprDeregisterConnectorTransaction removes a CLPR Connector and returns its locked stake to the
// named recipient (HIP-1535). The Connector admin key and the stake recipient must both sign. The
// commitment stays registered, so registering again needs a new salt and Connector ID.
type ClprDeregisterConnectorTransaction struct {
	*Transaction[*ClprDeregisterConnectorTransaction]
	fieldErrors    _ClprFieldErrors
	channelID      []byte
	connectorID    []byte
	stakeRecipient *AccountID
}

// NewClprDeregisterConnectorTransaction creates a ClprDeregisterConnectorTransaction.
func NewClprDeregisterConnectorTransaction() *ClprDeregisterConnectorTransaction {
	tx := &ClprDeregisterConnectorTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprDeregisterConnectorTransactionFromProtobuf(tx Transaction[*ClprDeregisterConnectorTransaction], pb *services.TransactionBody) ClprDeregisterConnectorTransaction {
	body := pb.GetClprDeregisterConnector()

	clprDeregisterConnectorTransaction := ClprDeregisterConnectorTransaction{
		channelID:   _ClprBytesFromProtobuf(body.GetChannelId()),
		connectorID: _ClprBytesFromProtobuf(body.GetConnectorId()),
	}
	if body.GetStakeRecipient() != nil {
		clprDeregisterConnectorTransaction.stakeRecipient = _AccountIDFromProtobuf(body.GetStakeRecipient())
	}

	tx.childTransaction = &clprDeregisterConnectorTransaction
	clprDeregisterConnectorTransaction.Transaction = &tx
	return clprDeregisterConnectorTransaction
}

// SetChannelID sets the 32-byte Channel ID.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprDeregisterConnectorTransaction) SetChannelID(channelID []byte) *ClprDeregisterConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("channelID", _ClprCheckLength("channelID", channelID)) {
		return tx
	}
	tx.channelID = _ClprCopyBytes(channelID)
	return tx
}

// GetChannelID returns a copy of the Channel ID, or nil if it is not set.
func (tx *ClprDeregisterConnectorTransaction) GetChannelID() []byte {
	return _ClprCopyBytes(tx.channelID)
}

// SetConnectorID sets the 32-byte Connector ID.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprDeregisterConnectorTransaction) SetConnectorID(connectorID []byte) *ClprDeregisterConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("connectorID", _ClprCheckLength("connectorID", connectorID)) {
		return tx
	}
	tx.connectorID = _ClprCopyBytes(connectorID)
	return tx
}

// GetConnectorID returns a copy of the Connector ID, or nil if it is not set.
func (tx *ClprDeregisterConnectorTransaction) GetConnectorID() []byte {
	return _ClprCopyBytes(tx.connectorID)
}

// SetStakeRecipient sets the account that receives the locked stake. It must sign the transaction.
func (tx *ClprDeregisterConnectorTransaction) SetStakeRecipient(accountID AccountID) *ClprDeregisterConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx
	}
	tx.stakeRecipient = _ClprCopyAccountID(&accountID)
	return tx
}

// GetStakeRecipient returns the stake recipient, or an empty AccountID if it is not set.
func (tx *ClprDeregisterConnectorTransaction) GetStakeRecipient() AccountID {
	if tx.stakeRecipient == nil {
		return AccountID{}
	}

	return *_ClprCopyAccountID(tx.stakeRecipient)
}

// ----------- Overridden functions ----------------

func (tx ClprDeregisterConnectorTransaction) getName() string {
	return "ClprDeregisterConnectorTransaction"
}

func (tx ClprDeregisterConnectorTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.stakeRecipient != nil {
		if err := tx.stakeRecipient.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx ClprDeregisterConnectorTransaction) validateTransactionFields() error {
	if err := _ClprRequireID(tx.getName(), "channelID", tx.channelID, tx.fieldErrors); err != nil {
		return err
	}
	if err := _ClprRequireID(tx.getName(), "connectorID", tx.connectorID, tx.fieldErrors); err != nil {
		return err
	}

	return _ClprRequireSet(tx.getName(), "stakeRecipient", tx.stakeRecipient != nil, tx.fieldErrors)
}

func (tx ClprDeregisterConnectorTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprDeregisterConnectorTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprDeregisterConnector{
		ClprDeregisterConnector: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprDeregisterConnectorTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprDeregisterConnectorTransaction) buildProtoBody() *services.ClprDeregisterConnectorTransactionBody {
	body := &services.ClprDeregisterConnectorTransactionBody{
		ChannelId:   _ClprCopyBytes(tx.channelID),
		ConnectorId: _ClprCopyBytes(tx.connectorID),
	}
	if tx.stakeRecipient != nil {
		body.StakeRecipient = tx.stakeRecipient._ToProtobuf()
	}

	return body
}

func (tx ClprDeregisterConnectorTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().DeregisterConnector,
	}
}

func (tx ClprDeregisterConnectorTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
