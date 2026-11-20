package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprCloseChannelTransaction closes a CLPR Channel (HIP-1535). It requires
// CLPR admin authorization (Treasury or System Admin).
// An existing Channel drains and closes. When no Channel exists, the optional ownership commitment
// names an abandoned pending registration to sweep; it is ignored when the Channel exists.
type ClprCloseChannelTransaction struct {
	*Transaction[*ClprCloseChannelTransaction]
	fieldErrors         _ClprFieldErrors
	channelID           []byte
	ownershipCommitment []byte
}

// NewClprCloseChannelTransaction creates a ClprCloseChannelTransaction.
func NewClprCloseChannelTransaction() *ClprCloseChannelTransaction {
	tx := &ClprCloseChannelTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprCloseChannelTransactionFromProtobuf(tx Transaction[*ClprCloseChannelTransaction], pb *services.TransactionBody) ClprCloseChannelTransaction {
	body := pb.GetClprCloseChannel()

	clprCloseChannelTransaction := ClprCloseChannelTransaction{
		channelID:           _ClprBytesFromProtobuf(body.GetChannelId()),
		ownershipCommitment: _ClprBytesFromProtobuf(body.GetOwnershipCommitment()),
	}

	tx.childTransaction = &clprCloseChannelTransaction
	clprCloseChannelTransaction.Transaction = &tx
	return clprCloseChannelTransaction
}

// SetChannelID sets the 32-byte ID of the Channel to close.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprCloseChannelTransaction) SetChannelID(channelID []byte) *ClprCloseChannelTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("channelID", _ClprCheckLength("channelID", channelID)) {
		return tx
	}
	tx.channelID = _ClprCopyBytes(channelID)
	return tx
}

// GetChannelID returns a copy of the Channel ID, or nil if it is not set.
func (tx *ClprCloseChannelTransaction) GetChannelID() []byte {
	return _ClprCopyBytes(tx.channelID)
}

// SetOwnershipCommitment sets the optional 32-byte commitment of a pending registration to sweep.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprCloseChannelTransaction) SetOwnershipCommitment(commitment []byte) *ClprCloseChannelTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("ownershipCommitment", _ClprCheckLength("ownershipCommitment", commitment)) {
		return tx
	}
	tx.ownershipCommitment = _ClprCopyBytes(commitment)
	return tx
}

// GetOwnershipCommitment returns a copy of the ownership commitment, or nil if it is not set.
func (tx *ClprCloseChannelTransaction) GetOwnershipCommitment() []byte {
	return _ClprCopyBytes(tx.ownershipCommitment)
}

// ----------- Overridden functions ----------------

func (tx ClprCloseChannelTransaction) getName() string {
	return "ClprCloseChannelTransaction"
}

func (tx ClprCloseChannelTransaction) validateNetworkOnIDs(*Client) error {
	return nil
}

func (tx ClprCloseChannelTransaction) validateTransactionFields() error {
	if err := _ClprRequireID(tx.getName(), "channelID", tx.channelID, tx.fieldErrors); err != nil {
		return err
	}
	if err := tx.fieldErrors["ownershipCommitment"]; err != nil {
		return fmt.Errorf("%s: %w", tx.getName(), err)
	}
	if tx.ownershipCommitment != nil {
		if err := _ClprCheckLength("ownershipCommitment", tx.ownershipCommitment); err != nil {
			return fmt.Errorf("%s: %w", tx.getName(), err)
		}
	}

	return nil
}

func (tx ClprCloseChannelTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprCloseChannelTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprCloseChannel{
		ClprCloseChannel: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprCloseChannelTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprCloseChannelTransaction) buildProtoBody() *services.ClprCloseChannelTransactionBody {
	return &services.ClprCloseChannelTransactionBody{
		ChannelId:           _ClprCopyBytes(tx.channelID),
		OwnershipCommitment: _ClprCopyBytes(tx.ownershipCommitment),
	}
}

func (tx ClprCloseChannelTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().CloseChannel,
	}
}

func (tx ClprCloseChannelTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
