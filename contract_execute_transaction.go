package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use tx file except in compliance with the License.
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
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// ContractExecuteTransaction calls a function of the given smart contract instance, giving it ContractFuncionParams as
// its inputs. it can use the given amount of gas, and any unspent gas will be refunded to the paying account.
//
// If tx function stores information, it is charged gas to store it. There is a fee in hbars to maintain that storage
// until the expiration time, and that fee is added as part of the transaction fee.
//
// For a cheaper but more limited _Method to call functions, see ContractCallQuery.
type ContractExecuteTransaction struct {
	transaction
	contractID *ContractID
	gas        int64
	amount     int64
	parameters []byte
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	tx := ContractExecuteTransaction{
		transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(2))
	tx.e = &tx

	return &tx
}

func _ContractExecuteTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *ContractExecuteTransaction {
	resultTx := &ContractExecuteTransaction{
		transaction: tx,
		contractID:  _ContractIDFromProtobuf(pb.GetContractCall().GetContractID()),
		gas:         pb.GetContractCall().GetGas(),
		amount:      pb.GetContractCall().GetAmount(),
		parameters:  pb.GetContractCall().GetFunctionParameters(),
	}
	resultTx.e = resultTx
	return resultTx
}

// SetContractID sets the contract instance to call.
func (tx *ContractExecuteTransaction) SetContractID(contractID ContractID) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// GetContractID returns the contract instance to call.
func (tx *ContractExecuteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// SetGas sets the maximum amount of gas to use for the call.
func (tx *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.gas = int64(gas)
	return tx
}

// GetGas returns the maximum amount of gas to use for the call.
func (tx *ContractExecuteTransaction) GetGas() uint64 {
	return uint64(tx.gas)
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if tx is nonzero)
func (tx *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.amount = amount.AsTinybar()
	return tx
}

// GetPayableAmount returns the amount of Hbar sent (the function must be payable if tx is nonzero)
func (tx ContractExecuteTransaction) GetPayableAmount() Hbar {
	return HbarFromTinybar(tx.amount)
}

// SetFunctionParameters sets the function parameters
func (tx *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.parameters = params
	return tx
}

// GetFunctionParameters returns the function parameters
func (tx *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return tx.parameters
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (tx *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	if params == nil {
		params = NewContractFunctionParameters()
	}

	tx.parameters = params._Build(&name)
	return tx
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *ContractExecuteTransaction) Sign(
	privateKey PrivateKey,
) *ContractExecuteTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *ContractExecuteTransaction) SignWithOperator(
	client *Client,
) (*ContractExecuteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *ContractExecuteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractExecuteTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *ContractExecuteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractExecuteTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when tx deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *ContractExecuteTransaction) SetGrpcDeadline(deadline *time.Duration) *ContractExecuteTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *ContractExecuteTransaction) Freeze() (*ContractExecuteTransaction, error) {
	_, err := tx.transaction.Freeze()
	return tx, err
}

func (tx *ContractExecuteTransaction) FreezeWith(client *Client) (*ContractExecuteTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *ContractExecuteTransaction) SetMaxTransactionFee(fee Hbar) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *ContractExecuteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx ContractExecuteTransaction.
func (tx *ContractExecuteTransaction) SetTransactionMemo(memo string) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx ContractExecuteTransaction.
func (tx *ContractExecuteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx ContractExecuteTransaction.
func (tx *ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) *ContractExecuteTransaction {
	tx._RequireNotFrozen()

	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for tx ContractExecuteTransaction.
func (tx *ContractExecuteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractExecuteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *ContractExecuteTransaction) SetMaxRetry(count int) *ContractExecuteTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *ContractExecuteTransaction) SetMaxBackoff(max time.Duration) *ContractExecuteTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *ContractExecuteTransaction) SetMinBackoff(min time.Duration) *ContractExecuteTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *ContractExecuteTransaction) SetLogLevel(level LogLevel) *ContractExecuteTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *ContractExecuteTransaction) getName() string {
	return "ContractExecuteTransaction"
}
func (tx *ContractExecuteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *ContractExecuteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractCall{
			ContractCall: tx.buildProtoBody(),
		},
	}
}

func (tx *ContractExecuteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractCall{
			ContractCall: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *ContractExecuteTransaction) buildProtoBody() *services.ContractCallTransactionBody {
	body := &services.ContractCallTransactionBody{
		Gas:                tx.gas,
		Amount:             tx.amount,
		FunctionParameters: tx.parameters,
	}

	if tx.contractID != nil {
		body.ContractID = tx.contractID._ToProtobuf()
	}

	return body
}

func (tx *ContractExecuteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().ContractCallMethod,
	}
}
func (tx *ContractExecuteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
