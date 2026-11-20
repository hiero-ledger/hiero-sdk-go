package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprUpdateLedgerConfigurationTransaction updates the CLPR ledger throttles (HIP-1535). It requires
// CLPR admin authorization (Treasury or System Admin).
// Send every limit: read them with ClprGetLedgerConfigurationQuery, change some, and submit.
// Only the throttles are sent; the other configuration fields are set by the network.
type ClprUpdateLedgerConfigurationTransaction struct {
	*Transaction[*ClprUpdateLedgerConfigurationTransaction]
	throttles *ClprThrottles
}

// NewClprUpdateLedgerConfigurationTransaction creates a ClprUpdateLedgerConfigurationTransaction.
func NewClprUpdateLedgerConfigurationTransaction() *ClprUpdateLedgerConfigurationTransaction {
	tx := &ClprUpdateLedgerConfigurationTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprUpdateLedgerConfigurationTransactionFromProtobuf(tx Transaction[*ClprUpdateLedgerConfigurationTransaction], pb *services.TransactionBody) ClprUpdateLedgerConfigurationTransaction {
	clprUpdateLedgerConfigurationTransaction := ClprUpdateLedgerConfigurationTransaction{}
	if pbThrottles := pb.GetClprUpdateLedgerConfiguration().GetConfiguration().GetThrottles(); pbThrottles != nil {
		throttles := _ClprThrottlesFromProtobuf(pbThrottles)
		clprUpdateLedgerConfigurationTransaction.throttles = &throttles
	}

	tx.childTransaction = &clprUpdateLedgerConfigurationTransaction
	clprUpdateLedgerConfigurationTransaction.Transaction = &tx
	return clprUpdateLedgerConfigurationTransaction
}

// SetThrottles sets the ledger throttles. They are validated at freeze.
func (tx *ClprUpdateLedgerConfigurationTransaction) SetThrottles(throttles ClprThrottles) *ClprUpdateLedgerConfigurationTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx
	}
	tx.throttles = &throttles
	return tx
}

// GetThrottles returns the throttles and whether they are set.
func (tx *ClprUpdateLedgerConfigurationTransaction) GetThrottles() (ClprThrottles, bool) {
	if tx.throttles == nil {
		return ClprThrottles{}, false
	}

	return *tx.throttles, true
}

// ----------- Overridden functions ----------------

func (tx ClprUpdateLedgerConfigurationTransaction) getName() string {
	return "ClprUpdateLedgerConfigurationTransaction"
}

func (tx ClprUpdateLedgerConfigurationTransaction) validateNetworkOnIDs(*Client) error {
	return nil
}

func (tx ClprUpdateLedgerConfigurationTransaction) validateTransactionFields() error {
	if err := _ClprRequireSet(tx.getName(), "throttles", tx.throttles != nil, nil); err != nil {
		return err
	}
	if err := tx.throttles.Validate(); err != nil {
		return fmt.Errorf("%s: %w", tx.getName(), err)
	}

	return nil
}

func (tx ClprUpdateLedgerConfigurationTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprUpdateLedgerConfigurationTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprUpdateLedgerConfiguration{
		ClprUpdateLedgerConfiguration: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprUpdateLedgerConfigurationTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprUpdateLedgerConfigurationTransaction) buildProtoBody() *services.ClprUpdateLedgerConfigurationTransactionBody {
	body := &services.ClprUpdateLedgerConfigurationTransactionBody{
		Configuration: &services.ClprLedgerConfiguration{},
	}
	if tx.throttles != nil {
		body.Configuration.Throttles = tx.throttles._ToProtobuf()
	}

	return body
}

func (tx ClprUpdateLedgerConfigurationTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().UpdateLedgerConfiguration,
	}
}

func (tx ClprUpdateLedgerConfigurationTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
