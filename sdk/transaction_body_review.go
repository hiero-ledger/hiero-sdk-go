package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// Shared by building and hash validation; this projection never changes stored bytes.
func (tx *Transaction[T]) _updatedTransactionBody(original *services.TransactionBody, id *services.TransactionID, node *services.AccountID) *services.TransactionBody {
	updated := protobuf.Clone(original).(*services.TransactionBody)
	if updated.NodeAccountID == nil {
		updated.NodeAccountID = node
	}
	if !_transactionIDsEqual(updated.TransactionID, id) {
		updated.TransactionID = id
	}
	updated.Memo = tx.memo
	updated.HighVolume = tx.highVolume
	if tx.transactionFee != 0 {
		updated.TransactionFee = tx.transactionFee
	} else {
		updated.TransactionFee = tx.defaultMaxTransactionFee
	}
	tx._populateCommonFeeLimitsAndBatchKey(updated)
	return updated
}

// Creation and late building share configured fields, but have different scalar fee defaults.
func (tx *Transaction[T]) _populateCommonFeeLimitsAndBatchKey(body *services.TransactionBody) {
	if tx.customFeeLimits != nil {
		limits := make([]*services.CustomFeeLimit, 0, len(tx.customFeeLimits))
		for i, limit := range tx.customFeeLimits {
			proposed := limit.toProtobuf()
			if i < len(body.MaxCustomFees) && protobuf.Equal(customFeeLimitFromProtobuf(body.MaxCustomFees[i]).toProtobuf(), proposed) {
				// Compare editable values, retaining unknown nested bytes when unchanged.
				proposed = body.MaxCustomFees[i]
			}
			limits = append(limits, proposed)
		}
		body.MaxCustomFees = limits
	}
	if tx.batchKey != nil {
		proposed := tx.batchKey._ToProtoKey()
		originalKey, err := _KeyFromProtobuf(body.BatchKey)
		if err != nil || !protobuf.Equal(originalKey._ToProtoKey(), proposed) {
			body.BatchKey = proposed
		}
	}
}

func (tx *Transaction[T]) _hasPendingCommonBodyChanges(body *services.TransactionBody) bool {
	updated := tx._updatedTransactionBody(body, body.TransactionID, body.NodeAccountID)
	return !protobuf.Equal(body, updated) ||
		body.GetTransactionValidDuration().GetSeconds() != _DurationToProtobuf(tx.GetTransactionValidDuration()).GetSeconds()
}

func (tx *Transaction[T]) _validateReviewedBodies() error {
	if !tx.bodyReviewed {
		return nil
	}
	for _, value := range tx.signedTransactions.slice {
		signed, ok := value.(*services.SignedTransaction)
		if !ok || signed == nil {
			return fmt.Errorf("invalid prepared signed transaction")
		}
		var body services.TransactionBody
		if err := protobuf.Unmarshal(signed.BodyBytes, &body); err != nil {
			return err
		}
		if tx._hasPendingCommonBodyChanges(&body) {
			return errTransactionIsFrozen
		}
	}
	return nil
}

func (tx *Transaction[T]) _rejectBodyReviewEdit() bool {
	if tx.bodyReviewed {
		tx.freezeError = errTransactionIsFrozen
		return true
	}
	return false
}

// Existing free setters keep their pre-review error behavior.
func (tx *Transaction[T]) _bodyReviewMutationError() error {
	if tx.bodyReviewed {
		return tx.freezeError
	}
	return nil
}

func (tx *Transaction[T]) _preparedTransactionIndex(id *services.TransactionID, node *services.AccountID) (int, error) {
	for index, value := range tx.signedTransactions.slice {
		var body services.TransactionBody
		if err := protobuf.Unmarshal(value.(*services.SignedTransaction).BodyBytes, &body); err != nil {
			return 0, err
		}
		if _transactionIDsEqual(body.TransactionID, id) && _accountIDValueKey(body.NodeAccountID) == _accountIDValueKey(node) {
			return index, nil
		}
	}
	return 0, fmt.Errorf("no prepared transaction for selected transaction ID and node ID")
}
