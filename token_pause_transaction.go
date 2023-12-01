package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenPauseTransaction
// Pauses the Token from being involved in any kind of transaction until it is unpaused.
// Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as paused and will be not able to be a part of any transaction.
// The operation is idempotent - becomes a no-op if the Token is already Paused.
type TokenPauseTransaction struct {
	transaction
	tokenID *TokenID
}

// NewTokenPauseTransaction creates TokenPauseTransaction which
// pauses the Token from being involved in any kind of transaction until it is unpaused.
// Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as paused and will be not able to be a part of any transaction.
// The operation is idempotent - becomes a no-op if the Token is already Paused.
func NewTokenPauseTransaction() *TokenPauseTransaction {
	tx := TokenPauseTransaction{
		transaction: _NewTransaction(),
	}

	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenPauseTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *TokenPauseTransaction {
	resultTx := &TokenPauseTransaction{
		transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}
	resultTx.e = resultTx
	return resultTx
}

// SetTokenID Sets the token to be paused
func (tx *TokenPauseTransaction) SetTokenID(tokenID TokenID) *TokenPauseTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token to be paused
func (tx *TokenPauseTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenPauseTransaction) Sign(privateKey PrivateKey) *TokenPauseTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenPauseTransaction) SignWithOperator(client *Client) (*TokenPauseTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenPauseTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenPauseTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenPauseTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenPauseTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenPauseTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenPauseTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenPauseTransaction) Freeze() (*TokenPauseTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenPauseTransaction) FreezeWith(client *Client) (*TokenPauseTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenPauseTransaction.
func (tx *TokenPauseTransaction) SetMaxTransactionFee(fee Hbar) *TokenPauseTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenPauseTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenPauseTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenPauseTransaction.
func (tx *TokenPauseTransaction) SetTransactionMemo(memo string) *TokenPauseTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenPauseTransaction.
func (tx *TokenPauseTransaction) SetTransactionValidDuration(duration time.Duration) *TokenPauseTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenPauseTransaction.
func (tx *TokenPauseTransaction) SetTransactionID(transactionID TransactionID) *TokenPauseTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenPauseTransaction.
func (tx *TokenPauseTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenPauseTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenPauseTransaction) SetMaxRetry(count int) *TokenPauseTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenPauseTransaction) SetMaxBackoff(max time.Duration) *TokenPauseTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenPauseTransaction) SetMinBackoff(min time.Duration) *TokenPauseTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenPauseTransaction) SetLogLevel(level LogLevel) *TokenPauseTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *TokenPauseTransaction) getName() string {
	return "TokenPauseTransaction"
}

func (tx *TokenPauseTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenPauseTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenPause{
			TokenPause: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenPauseTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) { //nolint
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenPause{
			TokenPause: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenPauseTransaction) buildProtoBody() *services.TokenPauseTransactionBody {
	body := &services.TokenPauseTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}
	return body

}

func (tx *TokenPauseTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DeleteToken,
	}
}

func (tx *TokenPauseTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}