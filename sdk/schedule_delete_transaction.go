package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ScheduleDeleteTransaction Marks a schedule in the network's action queue as deleted. Must be signed by the admin key of the
// target schedule.  A deleted schedule cannot receive any additional signing keys, nor will it be
// executed.
type ScheduleDeleteTransaction struct {
	*Transaction[*ScheduleDeleteTransaction]
	scheduleID *ScheduleID
}

// NewScheduleDeleteTransaction creates ScheduleDeleteTransaction which marks a schedule in the network's action queue as deleted.
// Must be signed by the admin key of the target schedule.
// A deleted schedule cannot receive any additional signing keys, nor will it be executed.
func NewScheduleDeleteTransaction() *ScheduleDeleteTransaction {
	tx := &ScheduleDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _ScheduleDeleteTransactionFromProtobuf(tx Transaction[*ScheduleDeleteTransaction], pb *services.TransactionBody) ScheduleDeleteTransaction {
	scheduleDeleteTransaction := ScheduleDeleteTransaction{
		scheduleID: _ScheduleIDFromProtobuf(pb.GetScheduleDelete().GetScheduleID()),
	}

	tx.childTransaction = &scheduleDeleteTransaction
	scheduleDeleteTransaction.Transaction = &tx
	return scheduleDeleteTransaction
}

// SetScheduleID Sets the ScheduleID of the scheduled transaction to be deleted
func (tx *ScheduleDeleteTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleDeleteTransaction {
	tx._RequireNotFrozen()
	tx.scheduleID = &scheduleID
	return tx
}

func (tx *ScheduleDeleteTransaction) GetScheduleID() ScheduleID {
	if tx.scheduleID == nil {
		return ScheduleID{}
	}

	return *tx.scheduleID
}

// ----------- Overridden functions ----------------

func (tx ScheduleDeleteTransaction) getName() string {
	return "ScheduleDeleteTransaction"
}

func (tx ScheduleDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.scheduleID != nil {
		if err := tx.scheduleID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx ScheduleDeleteTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ScheduleDelete{
		ScheduleDelete: tx.buildProtoBody(),
	}

	return body
}

func (tx ScheduleDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	body := tx.buildSchedulableTransactionBody()
	body.Data = &services.SchedulableTransactionBody_ScheduleDelete{
		ScheduleDelete: tx.buildProtoBody(),
	}

	return body, nil
}

func (tx ScheduleDeleteTransaction) buildProtoBody() *services.ScheduleDeleteTransactionBody {
	body := &services.ScheduleDeleteTransactionBody{}
	if tx.scheduleID != nil {
		body.ScheduleID = tx.scheduleID._ToProtobuf()
	}

	return body
}

func (tx ScheduleDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetSchedule().DeleteSchedule,
	}
}
func (tx ScheduleDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx ScheduleDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
