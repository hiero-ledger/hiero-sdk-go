package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprSubmitBundleTransaction submits an opaque CLPR message bundle for a Channel (HIP-1535). Only
// endpoint nodes in the address book may submit it.
type ClprSubmitBundleTransaction struct {
	*Transaction[*ClprSubmitBundleTransaction]
	fieldErrors   _ClprFieldErrors
	channelID     []byte
	bundlePayload []byte
}

// NewClprSubmitBundleTransaction creates a ClprSubmitBundleTransaction.
func NewClprSubmitBundleTransaction() *ClprSubmitBundleTransaction {
	tx := &ClprSubmitBundleTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprSubmitBundleTransactionFromProtobuf(tx Transaction[*ClprSubmitBundleTransaction], pb *services.TransactionBody) ClprSubmitBundleTransaction {
	body := pb.GetClprSubmitBundle()

	clprSubmitBundleTransaction := ClprSubmitBundleTransaction{
		channelID:     _ClprBytesFromProtobuf(body.GetChannelId()),
		bundlePayload: _ClprCopyBytes(body.GetBundlePayload()),
	}

	tx.childTransaction = &clprSubmitBundleTransaction
	clprSubmitBundleTransaction.Transaction = &tx
	return clprSubmitBundleTransaction
}

// SetChannelID sets the 32-byte ID of the Channel the bundle belongs to.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprSubmitBundleTransaction) SetChannelID(channelID []byte) *ClprSubmitBundleTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("channelID", _ClprCheckLength("channelID", channelID)) {
		return tx
	}
	tx.channelID = _ClprCopyBytes(channelID)
	return tx
}

// GetChannelID returns a copy of the Channel ID, or nil if it is not set.
func (tx *ClprSubmitBundleTransaction) GetChannelID() []byte {
	return _ClprCopyBytes(tx.channelID)
}

// SetBundlePayload sets the opaque bundle bytes. Any length is accepted, including empty.
func (tx *ClprSubmitBundleTransaction) SetBundlePayload(bundlePayload []byte) *ClprSubmitBundleTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx
	}
	tx.bundlePayload = _ClprCopyBytes(bundlePayload)
	return tx
}

// GetBundlePayload returns a copy of the bundle payload.
func (tx *ClprSubmitBundleTransaction) GetBundlePayload() []byte {
	return _ClprCopyBytes(tx.bundlePayload)
}

// ----------- Overridden functions ----------------

func (tx ClprSubmitBundleTransaction) getName() string {
	return "ClprSubmitBundleTransaction"
}

func (tx ClprSubmitBundleTransaction) validateNetworkOnIDs(*Client) error {
	return nil
}

func (tx ClprSubmitBundleTransaction) validateTransactionFields() error {
	return _ClprRequireID(tx.getName(), "channelID", tx.channelID, tx.fieldErrors)
}

func (tx ClprSubmitBundleTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprSubmitBundleTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprSubmitBundle{
		ClprSubmitBundle: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprSubmitBundleTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprSubmitBundleTransaction) buildProtoBody() *services.ClprSubmitBundleTransactionBody {
	return &services.ClprSubmitBundleTransactionBody{
		ChannelId:     _ClprCopyBytes(tx.channelID),
		BundlePayload: _ClprCopyBytes(tx.bundlePayload),
	}
}

func (tx ClprSubmitBundleTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().SubmitBundle,
	}
}

func (tx ClprSubmitBundleTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
