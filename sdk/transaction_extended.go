package hiero

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// TransactionFromBytes converts Transaction bytes to a related *Transaction.
func CreateTransferTransactionFromBytes(data []byte) (*TransferTransaction, error) { // nolint
	tx := _NewTransaction[*TransferTransaction](nil)
	fmt.Println("Parse Single transaction...")

	var first *services.TransactionBody = nil
	var err error

	var signedTransaction services.SignedTransaction
	if err := protobuf.Unmarshal(data, &signedTransaction); err != nil {
		fmt.Printf("Error in unmarshal: %v", err)
		return nil, errors.Wrap(err, "error deserializing SignedTransactionBytes in TransactionFromBytes")
	}

	fmt.Println("have a signed tx")
	tx.signedTransactions = tx.signedTransactions._Push(&signedTransaction)

	for _, sigPair := range signedTransaction.GetSigMap().GetSigPair() {
		key, err := PublicKeyFromBytes(sigPair.GetPubKeyPrefix())
		if err != nil {
			return nil, err
		}

		tx.publicKeys = append(tx.publicKeys, key)
		tx.transactionSigners = append(tx.transactionSigners, nil)
	}

	var body services.TransactionBody
	if err := protobuf.Unmarshal(signedTransaction.GetBodyBytes(), &body); err != nil {
		return nil, errors.Wrap(err, "error deserializing BodyBytes in TransactionFromBytes")
	}

	if first == nil {
		first = &body
	}
	var transactionID TransactionID
	var nodeAccountID AccountID
	if body.GetTransactionID() != nil {
		transactionID = _TransactionIDFromProtobuf(body.GetTransactionID())
	}

	if body.GetNodeAccountID() != nil {
		nodeAccountID = *_AccountIDFromProtobuf(body.GetNodeAccountID())
	}

	found := false

	for _, value := range tx.transactionIDs.slice {
		id := value.(TransactionID)
		if id.AccountID != nil && transactionID.AccountID != nil &&
			id.AccountID._Equals(*transactionID.AccountID) &&
			id.ValidStart != nil && transactionID.ValidStart != nil &&
			id.ValidStart.Equal(*transactionID.ValidStart) {
			found = true
			break
		}
	}

	if !found {
		tx.transactionIDs = tx.transactionIDs._Push(transactionID)
	}

	for _, id := range tx.GetNodeAccountIDs() {
		if id._Equals(nodeAccountID) {
			found = true
			break
		}
	}

	if !found {
		tx.nodeAccountIDs = tx.nodeAccountIDs._Push(nodeAccountID)
	}

	tx.memo = body.Memo
	if body.TransactionFee != 0 {
		tx.transactionFee = body.TransactionFee
	}

	if tx.transactionIDs._Length() > 0 {
		tx.transactionIDs.locked = true
	}

	if tx.nodeAccountIDs._Length() > 0 {
		tx.nodeAccountIDs.locked = true
	}

	if first == nil {
		return nil, errNoTransactionInBytes
	}

	ret := _TransferTransactionFromProtobuf(*tx, first)
	return &ret, nil
}
