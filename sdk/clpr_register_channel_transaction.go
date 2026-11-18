package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprRegisterChannelTransaction registers the ownership commitment of a CLPR Channel, the commit
// phase of the Channel commit-reveal (HIP-1535). Any payer may submit it. Compute the commitment
// with ComputeClprChannelOwnershipCommitment, then reveal it with ClprCompleteChannelTransaction.
type ClprRegisterChannelTransaction struct {
	*Transaction[*ClprRegisterChannelTransaction]
	fieldErrors         _ClprFieldErrors
	ownershipCommitment []byte
}

// NewClprRegisterChannelTransaction creates a ClprRegisterChannelTransaction.
func NewClprRegisterChannelTransaction() *ClprRegisterChannelTransaction {
	tx := &ClprRegisterChannelTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprRegisterChannelTransactionFromProtobuf(tx Transaction[*ClprRegisterChannelTransaction], pb *services.TransactionBody) ClprRegisterChannelTransaction {
	clprRegisterChannelTransaction := ClprRegisterChannelTransaction{
		ownershipCommitment: _ClprBytesFromProtobuf(pb.GetClprRegisterChannel().GetOwnershipCommitment()),
	}

	tx.childTransaction = &clprRegisterChannelTransaction
	clprRegisterChannelTransaction.Transaction = &tx
	return clprRegisterChannelTransaction
}

// SetOwnershipCommitment sets the 32-byte keccak256(channelID || ownership public key) commitment.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprRegisterChannelTransaction) SetOwnershipCommitment(commitment []byte) *ClprRegisterChannelTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("ownershipCommitment", _ClprCheckLength("ownershipCommitment", commitment)) {
		return tx
	}
	tx.ownershipCommitment = _ClprCopyBytes(commitment)
	return tx
}

// GetOwnershipCommitment returns a copy of the ownership commitment, or nil if it is not set.
func (tx *ClprRegisterChannelTransaction) GetOwnershipCommitment() []byte {
	return _ClprCopyBytes(tx.ownershipCommitment)
}

// ----------- Overridden functions ----------------

func (tx ClprRegisterChannelTransaction) getName() string {
	return "ClprRegisterChannelTransaction"
}

func (tx ClprRegisterChannelTransaction) validateNetworkOnIDs(*Client) error {
	return nil
}

func (tx ClprRegisterChannelTransaction) validateTransactionFields() error {
	return _ClprRequireID(tx.getName(), "ownershipCommitment", tx.ownershipCommitment, tx.fieldErrors)
}

func (tx ClprRegisterChannelTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprRegisterChannelTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprRegisterChannel{
		ClprRegisterChannel: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprRegisterChannelTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprRegisterChannelTransaction) buildProtoBody() *services.ClprRegisterChannelTransactionBody {
	return &services.ClprRegisterChannelTransactionBody{
		OwnershipCommitment: _ClprCopyBytes(tx.ownershipCommitment),
	}
}

func (tx ClprRegisterChannelTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().RegisterChannel,
	}
}

func (tx ClprRegisterChannelTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
