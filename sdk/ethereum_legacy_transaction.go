package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// EthereumLegacyTransaction represents the legacy Ethereum transaction data.
type EthereumLegacyTransaction struct {
	Nonce    []byte
	GasPrice []byte
	GasLimit []byte
	To       []byte
	Value    []byte
	CallData []byte
	V        []byte
	R        []byte
	S        []byte
}

// nolint
// NewEthereumLegacyTransaction creates a new EthereumLegacyTransaction with the provided fields.
func NewEthereumLegacyTransaction(nonce, gasPrice, gasLimit, to, value, callData, v, r, s []byte) *EthereumLegacyTransaction {
	return &EthereumLegacyTransaction{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		To:       to,
		Value:    value,
		CallData: callData,
		V:        v,
		R:        r,
		S:        s,
	}
}

// EthereumLegacyTransactionFromBytes decodes a signed legacy transaction
// (RLP list of 9 elements, no type prefix).
func EthereumLegacyTransactionFromBytes(bytes []byte) (*EthereumLegacyTransaction, error) {
	item := NewRLPItem(LIST_TYPE)
	if err := item.Read(bytes); err != nil {
		return nil, errors.Wrap(err, "failed to read RLP data")
	}

	if item.itemType != LIST_TYPE {
		return nil, errors.New("input byte array does not represent a list of RLP-encoded elements")
	}

	if len(item.childItems) != 9 {
		return nil, errors.New("input byte array does not contain 9 RLP-encoded elements")
	}

	return NewEthereumLegacyTransaction(
		item.childItems[0].itemValue,
		item.childItems[1].itemValue,
		item.childItems[2].itemValue,
		item.childItems[3].itemValue,
		item.childItems[4].itemValue,
		item.childItems[5].itemValue,
		item.childItems[6].itemValue,
		item.childItems[7].itemValue,
		item.childItems[8].itemValue,
	), nil
}

// _toUnsignedRLP builds the unsigned portion of the legacy RLP list
func (txn *EthereumLegacyTransaction) _toUnsignedRLP() *RLPItem {
	item := NewRLPItem(LIST_TYPE)
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Nonce))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.GasPrice))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.GasLimit))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.To))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Value))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.CallData))
	return item
}

// _encodeWithSignature appends V/R/S to the given unsigned RLP list
// and serializes it. Legacy has no type prefix.
func (txn *EthereumLegacyTransaction) _encodeWithSignature(item *RLPItem) ([]byte, error) {
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.V))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.R))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.S))
	return item.Write()
}

// ToBytes encodes the EthereumLegacyTransaction into RLP format.
func (txn *EthereumLegacyTransaction) ToBytes() ([]byte, error) {
	return txn._encodeWithSignature(txn._toUnsignedRLP())
}

// Sign signs the transaction with the given ECDSA key, populates V/R/S on
// the receiver, and returns the signed RLP bytes.
func (txn *EthereumLegacyTransaction) Sign(key PrivateKey) ([]byte, error) {
	item := txn._toUnsignedRLP()
	message, err := item.Write()
	if err != nil {
		return nil, err
	}

	sig := key.Sign(message)
	if len(sig) < 64 {
		return nil, errors.New("signing produced an invalid signature; expected an ECDSA key")
	}
	r := sig[0:32]
	s := sig[32:64]
	v := key.GetRecoveryId(r, s, message)
	if v < 0 {
		return nil, errors.New("unable to compute recovery id; expected an ECDSA key")
	}

	txn.R = r
	txn.S = s
	txn.V = []byte{byte(27 + v)}

	return txn._encodeWithSignature(item)
}

// String returns a string representation of the EthereumLegacyTransaction.
func (txn *EthereumLegacyTransaction) String() string {
	return fmt.Sprintf("Nonce: %s\nGasPrice: %s\nGasLimit: %s\nTo: %s\nValue: %s\nCallData: %s\nV: %s\nR: %s\nS: %s",
		hex.EncodeToString(txn.Nonce),
		hex.EncodeToString(txn.GasPrice),
		hex.EncodeToString(txn.GasLimit),
		hex.EncodeToString(txn.To),
		hex.EncodeToString(txn.Value),
		hex.EncodeToString(txn.CallData),
		hex.EncodeToString(txn.V),
		hex.EncodeToString(txn.R),
		hex.EncodeToString(txn.S),
	)
}
