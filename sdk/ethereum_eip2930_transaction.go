package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
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

// EthereumEIP2930TransactionFromBytes decodes the RLP encoded bytes into an EthereumEIP2930Transaction.
func EthereumEIP2930TransactionFromBytes(bytes []byte) (*EthereumEIP2930Transaction, error) {
	if len(bytes) == 0 || bytes[0] != 0x01 {
		return nil, errors.New("input byte array is malformed; it should start with 0x01 followed by 11 RLP-encoded elements")
	}

	// Remove the prefix byte (0x01)
	item := NewRLPItem(LIST_TYPE)
	if err := item.Read(bytes[1:]); err != nil {
		return nil, errors.Wrap(err, "failed to read RLP data")
	}

	if item.itemType != LIST_TYPE || len(item.childItems) != 11 {
		return nil, errors.New("input byte array is malformed; it should be a list of 11 RLP-encoded elements")
	}

	// Handle the access list
	var accessListValues [][]byte
	for _, child := range item.childItems[7].childItems {
		accessListValues = append(accessListValues, child.itemValue)
	}

	// Extract values from the RLP item
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

// ToBytes encodes the EthereumEIP2930Transaction into RLP format.
func (txn *EthereumEIP2930Transaction) ToBytes() ([]byte, error) {
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
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.RecoveryId))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.R))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.S))

	transactionBytes, err := item.Write()
	if err != nil {
		return nil, err
	}
	// Append 01 byte as it is the standard for EIP-2930
	combinedBytes := append([]byte{0x01}, transactionBytes...)

	return combinedBytes, nil
}

// String returns a string representation of the EthereumEIP2930Transaction.
func (txn *EthereumEIP2930Transaction) String() string {
	// Encode each element in the AccessList slice individually
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
