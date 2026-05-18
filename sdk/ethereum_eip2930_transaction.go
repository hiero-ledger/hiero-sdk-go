package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

// EthereumEIP2930Transaction represents the EIP-2930 Ethereum transaction data.
type EthereumEIP2930Transaction struct {
	ChainId    []byte
	Nonce      []byte
	GasPrice   []byte
	GasLimit   []byte
	To         []byte
	Value      []byte
	CallData   []byte
	AccessList [][]byte
	RecoveryId []byte
	R          []byte
	S          []byte
}

// NewEthereumEIP2930Transaction creates a new EthereumEIP2930Transaction with the provided fields.
func NewEthereumEIP2930Transaction(
	chainId, nonce, gasPrice, gasLimit, to, value, callData, recoveryId, r, s []byte, accessList [][]byte) *EthereumEIP2930Transaction {
	return &EthereumEIP2930Transaction{
		ChainId:    chainId,
		Nonce:      nonce,
		GasPrice:   gasPrice,
		GasLimit:   gasLimit,
		To:         to,
		Value:      value,
		CallData:   callData,
		AccessList: accessList,
		RecoveryId: recoveryId,
		R:          r,
		S:          s,
	}
}

// EthereumEIP2930TransactionFromBytes decodes signed EIP-2930 RLP bytes
// (leading 0x01 prefix + list of 11 elements) into a transaction.
func EthereumEIP2930TransactionFromBytes(bytes []byte) (*EthereumEIP2930Transaction, error) {
	if len(bytes) == 0 || bytes[0] != 0x01 {
		return nil, errors.New("input byte array is malformed; it should start with 0x01 followed by 11 RLP-encoded elements")
	}

	item := NewRLPItem(LIST_TYPE)
	if err := item.Read(bytes[1:]); err != nil {
		return nil, errors.Wrap(err, "failed to read RLP data")
	}

	if item.itemType != LIST_TYPE || len(item.childItems) != 11 {
		return nil, errors.New("input byte array is malformed; it should be a list of 11 RLP-encoded elements")
	}

	var accessListValues [][]byte
	for _, child := range item.childItems[7].childItems {
		accessListValues = append(accessListValues, child.itemValue)
	}

	return NewEthereumEIP2930Transaction(
		item.childItems[0].itemValue,
		item.childItems[1].itemValue,
		item.childItems[2].itemValue,
		item.childItems[3].itemValue,
		item.childItems[4].itemValue,
		item.childItems[5].itemValue,
		item.childItems[6].itemValue,
		item.childItems[8].itemValue,
		item.childItems[9].itemValue,
		item.childItems[10].itemValue,
		accessListValues,
	), nil
}

// _toUnsignedRLP builds the unsigned portion of the EIP-2930 RLP list
func (txn *EthereumEIP2930Transaction) _toUnsignedRLP() *RLPItem {
	item := NewRLPItem(LIST_TYPE)
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.ChainId))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Nonce))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.GasPrice))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.GasLimit))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.To))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Value))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.CallData))
	accessListItem := NewRLPItem(LIST_TYPE)
	for _, itemBytes := range txn.AccessList {
		accessListItem.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(itemBytes))
	}
	item.PushBack(accessListItem)
	return item
}

// _encodeWithSignature appends the signature to the unsigned RLP list and
// prepends the 0x01 EIP-2930 type prefix.
func (txn *EthereumEIP2930Transaction) _encodeWithSignature(item *RLPItem) ([]byte, error) {
	return _encodeTypedWithSignature(item, 0x01, txn.RecoveryId, txn.R, txn.S)
}

// ToBytes encodes the EthereumEIP2930Transaction into RLP format.
func (txn *EthereumEIP2930Transaction) ToBytes() ([]byte, error) {
	return txn._encodeWithSignature(txn._toUnsignedRLP())
}

// Sign signs the unsigned transaction with the given ECDSA key, populates
// RecoveryId/R/S on the receiver, and returns the signed RLP bytes ready to
// wrap in EthereumTransaction.
func (txn *EthereumEIP2930Transaction) Sign(key PrivateKey) ([]byte, error) {
	item := txn._toUnsignedRLP()
	r, s, v, err := _signTypedTransaction(item, 0x01, key)
	if err != nil {
		return nil, err
	}
	txn.R = r
	txn.S = s
	txn.RecoveryId = []byte{byte(v)}
	return txn._encodeWithSignature(item)
}

// String returns a string representation of the EthereumEIP2930Transaction.
func (txn *EthereumEIP2930Transaction) String() string {
	var encodedAccessList []string
	for _, entry := range txn.AccessList {
		encodedAccessList = append(encodedAccessList, hex.EncodeToString(entry))
	}

	accessListStr := "[" + strings.Join(encodedAccessList, ", ") + "]"

	return fmt.Sprintf("ChainId: %s\nNonce: %s\nGasPrice: %s\nGasLimit: %s\nTo: %s\nValue: %s\nCallData: %s\nAccessList: %s\nRecoveryId: %s\nR: %s\nS: %s",
		hex.EncodeToString(txn.ChainId),
		hex.EncodeToString(txn.Nonce),
		hex.EncodeToString(txn.GasPrice),
		hex.EncodeToString(txn.GasLimit),
		hex.EncodeToString(txn.To),
		hex.EncodeToString(txn.Value),
		hex.EncodeToString(txn.CallData),
		accessListStr,
		hex.EncodeToString(txn.RecoveryId),
		hex.EncodeToString(txn.R),
		hex.EncodeToString(txn.S),
	)
}

// GetChainId returns the chain id as a uint64.
func (txn *EthereumEIP2930Transaction) GetChainId() uint64 {
	return _ethBytesToUint64(txn.ChainId)
}

// SetChainId sets the chain id from a uint64.
func (txn *EthereumEIP2930Transaction) SetChainId(v uint64) *EthereumEIP2930Transaction {
	txn.ChainId = _uint64ToEthBytes(v)
	return txn
}

// GetChainIdBytes returns the raw canonical big-endian chain id bytes.
func (txn *EthereumEIP2930Transaction) GetChainIdBytes() []byte { return txn.ChainId }

// SetChainIdBytes sets the chain id from raw canonical big-endian bytes.
func (txn *EthereumEIP2930Transaction) SetChainIdBytes(v []byte) *EthereumEIP2930Transaction {
	txn.ChainId = v
	return txn
}

// GetNonce returns the nonce as a uint64.
func (txn *EthereumEIP2930Transaction) GetNonce() uint64 {
	return _ethBytesToUint64(txn.Nonce)
}

// SetNonce sets the nonce from a uint64.
func (txn *EthereumEIP2930Transaction) SetNonce(v uint64) *EthereumEIP2930Transaction {
	txn.Nonce = _uint64ToEthBytes(v)
	return txn
}

// GetNonceBytes returns the raw canonical big-endian nonce bytes.
func (txn *EthereumEIP2930Transaction) GetNonceBytes() []byte { return txn.Nonce }

// SetNonceBytes sets the nonce from raw canonical big-endian bytes.
func (txn *EthereumEIP2930Transaction) SetNonceBytes(v []byte) *EthereumEIP2930Transaction {
	txn.Nonce = v
	return txn
}

// GetGasPrice returns the gas price as a *big.Int.
func (txn *EthereumEIP2930Transaction) GetGasPrice() *big.Int {
	return _ethBytesToBigInt(txn.GasPrice)
}

// SetGasPrice sets the gas price from a *big.Int.
func (txn *EthereumEIP2930Transaction) SetGasPrice(v *big.Int) *EthereumEIP2930Transaction {
	txn.GasPrice = _bigIntToEthBytes(v)
	return txn
}

// GetGasPriceBytes returns the raw canonical big-endian gas price bytes.
func (txn *EthereumEIP2930Transaction) GetGasPriceBytes() []byte { return txn.GasPrice }

// SetGasPriceBytes sets the gas price from raw bytes.
func (txn *EthereumEIP2930Transaction) SetGasPriceBytes(v []byte) *EthereumEIP2930Transaction {
	txn.GasPrice = v
	return txn
}

// GetGasLimit returns the gas limit as a uint64.
func (txn *EthereumEIP2930Transaction) GetGasLimit() uint64 {
	return _ethBytesToUint64(txn.GasLimit)
}

// SetGasLimit sets the gas limit from a uint64.
func (txn *EthereumEIP2930Transaction) SetGasLimit(v uint64) *EthereumEIP2930Transaction {
	txn.GasLimit = _uint64ToEthBytes(v)
	return txn
}

// GetGasLimitBytes returns the raw canonical big-endian gas limit bytes.
func (txn *EthereumEIP2930Transaction) GetGasLimitBytes() []byte { return txn.GasLimit }

// SetGasLimitBytes sets the gas limit from raw bytes.
func (txn *EthereumEIP2930Transaction) SetGasLimitBytes(v []byte) *EthereumEIP2930Transaction {
	txn.GasLimit = v
	return txn
}

// GetTo returns the recipient address bytes.
func (txn *EthereumEIP2930Transaction) GetTo() []byte { return txn.To }

// SetTo sets the recipient address bytes.
func (txn *EthereumEIP2930Transaction) SetTo(v []byte) *EthereumEIP2930Transaction {
	txn.To = v
	return txn
}

// GetValue returns the transaction value (wei) as a *big.Int.
func (txn *EthereumEIP2930Transaction) GetValue() *big.Int {
	return _ethBytesToBigInt(txn.Value)
}

// SetValue sets the transaction value (wei) from a *big.Int.
func (txn *EthereumEIP2930Transaction) SetValue(v *big.Int) *EthereumEIP2930Transaction {
	txn.Value = _bigIntToEthBytes(v)
	return txn
}

// GetValueBytes returns the raw canonical big-endian value bytes.
func (txn *EthereumEIP2930Transaction) GetValueBytes() []byte { return txn.Value }

// SetValueBytes sets the value from raw bytes.
func (txn *EthereumEIP2930Transaction) SetValueBytes(v []byte) *EthereumEIP2930Transaction {
	txn.Value = v
	return txn
}

// GetCallData returns the call data.
func (txn *EthereumEIP2930Transaction) GetCallData() []byte { return txn.CallData }

// SetCallData sets the call data.
func (txn *EthereumEIP2930Transaction) SetCallData(v []byte) *EthereumEIP2930Transaction {
	txn.CallData = v
	return txn
}

// GetRecoveryId returns the recovery id as an int.
func (txn *EthereumEIP2930Transaction) GetRecoveryId() int {
	if len(txn.RecoveryId) == 0 {
		return 0
	}
	return int(txn.RecoveryId[0])
}

// SetRecoveryId sets the recovery id from an int.
func (txn *EthereumEIP2930Transaction) SetRecoveryId(v int) *EthereumEIP2930Transaction {
	if v == 0 {
		txn.RecoveryId = []byte{}
	} else {
		txn.RecoveryId = []byte{byte(v)}
	}
	return txn
}

// GetRecoveryIdBytes returns the raw recovery id bytes.
func (txn *EthereumEIP2930Transaction) GetRecoveryIdBytes() []byte { return txn.RecoveryId }

// SetRecoveryIdBytes sets the recovery id from raw bytes.
func (txn *EthereumEIP2930Transaction) SetRecoveryIdBytes(v []byte) *EthereumEIP2930Transaction {
	txn.RecoveryId = v
	return txn
}

// GetR returns the R signature component.
func (txn *EthereumEIP2930Transaction) GetR() []byte { return txn.R }

// SetR sets the R signature component.
func (txn *EthereumEIP2930Transaction) SetR(v []byte) *EthereumEIP2930Transaction {
	txn.R = v
	return txn
}

// GetS returns the S signature component.
func (txn *EthereumEIP2930Transaction) GetS() []byte { return txn.S }

// SetS sets the S signature component.
func (txn *EthereumEIP2930Transaction) SetS(v []byte) *EthereumEIP2930Transaction {
	txn.S = v
	return txn
}

// GetAccessListItems returns the access list as structured AccessListItem entries.
func (txn *EthereumEIP2930Transaction) GetAccessListItems() []AccessListItem {
	return _accessListItemsFromBytes(txn.AccessList)
}

// SetAccessListItems replaces the access list from structured AccessListItem entries.
func (txn *EthereumEIP2930Transaction) SetAccessListItems(items []AccessListItem) *EthereumEIP2930Transaction {
	txn.AccessList = _accessListItemsToBytes(items)
	return txn
}

// AddAccessListItem appends a single access list entry.
func (txn *EthereumEIP2930Transaction) AddAccessListItem(item AccessListItem) *EthereumEIP2930Transaction {
	txn.AccessList = append(txn.AccessList, _accessListItemToBytes(item))
	return txn
}
