package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// EthereumEIP7702Transaction represents the EIP-1559 Ethereum transaction data.
type EthereumEIP7702Transaction struct {
	ChainId           []byte
	Nonce             []byte
	MaxPriorityGas    []byte
	MaxGas            []byte
	GasLimit          []byte
	To                []byte
	Value             []byte
	CallData          []byte
	AccessList        [][]byte
	AuthorizationList [][]byte
	RecoveryId        []byte
	R                 []byte
	S                 []byte
}

// nolint
// NewEthereumEIP7702Transaction creates a new EthereumEIP7702Transaction with the provided fields.
func NewEthereumEIP7702Transaction(
	chainId, nonce, maxPriorityGas, maxGas, gasLimit, to, value, callData, recoveryId, r, s []byte, accessList [][]byte, authorizationList [][]byte) *EthereumEIP7702Transaction {
	return &EthereumEIP7702Transaction{
		ChainId:           chainId,
		Nonce:             nonce,
		MaxPriorityGas:    maxPriorityGas,
		MaxGas:            maxGas,
		GasLimit:          gasLimit,
		To:                to,
		Value:             value,
		CallData:          callData,
		AccessList:        accessList,
		AuthorizationList: authorizationList,
		RecoveryId:        recoveryId,
		R:                 r,
		S:                 s,
	}
}

// FromBytes decodes the RLP encoded bytes into an EthereumEIP7702Transaction.
func EthereumEIP7702TransactionFromBytes(bytes []byte) (*EthereumEIP7702Transaction, error) {
	if len(bytes) == 0 || bytes[0] != 0x04 {
		return nil, errors.New("input byte array is malformed; it should start with 0x02 followed by 12 RLP-encoded elements")
	}

	// Remove the prefix byte (0x04)
	item := NewRLPItem(LIST_TYPE)
	if err := item.Read(bytes[1:]); err != nil {
		return nil, errors.Wrap(err, "failed to read RLP data")
	}

	if item.itemType != LIST_TYPE || len(item.childItems) != 12 {
		return nil, errors.New("input byte array is malformed; it should be a list of 12 RLP-encoded elements")
	}

	// Handle the access list
	var accessListValues [][]byte
	for _, child := range item.childItems[8].childItems {
		accessListValues = append(accessListValues, child.itemValue)
	}

	// Handle the authorization list
	var authorizationListValues [][]byte
	for _, child := range item.childItems[9].childItems {
		authorizationListValues = append(authorizationListValues, child.itemValue)
	}

	// Extract values from the RLP item
	return NewEthereumEIP7702Transaction(
		item.childItems[0].itemValue,
		item.childItems[1].itemValue,
		item.childItems[2].itemValue,
		item.childItems[3].itemValue,
		item.childItems[4].itemValue,
		item.childItems[5].itemValue,
		item.childItems[6].itemValue,
		item.childItems[7].itemValue,
		item.childItems[9].itemValue,
		item.childItems[10].itemValue,
		item.childItems[11].itemValue,
		accessListValues,
		authorizationListValues,
	), nil
}

// ToBytes encodes the EthereumEIP7702Transaction into RLP format.
func (txn *EthereumEIP7702Transaction) ToBytes() ([]byte, error) {
	item := NewRLPItem(LIST_TYPE)
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.ChainId))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Nonce))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.MaxPriorityGas))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.MaxGas))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.GasLimit))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.To))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Value))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.CallData))
	accessListItem := NewRLPItem(LIST_TYPE)
	for _, itemBytes := range txn.AccessList {
		accessListItem.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(itemBytes))
	}
	item.PushBack(accessListItem)
	authorizationListItem := NewRLPItem(LIST_TYPE)
	for _, itemBytes := range txn.AuthorizationList {
		authorizationListItem.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(itemBytes))
	}
	item.PushBack(authorizationListItem)
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.RecoveryId))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.R))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.S))

	transactionBytes, err := item.Write()
	if err != nil {
		return nil, err
	}
	// Append 04 byte as it is the standard for EIP7702
	combinedBytes := append([]byte{0x04}, transactionBytes...)

	return combinedBytes, nil
}

// String returns a string representation of the EthereumEIP7702Transaction.
func (txn *EthereumEIP7702Transaction) String() string {
	// Encode each element in the AccessList slice individually
	var encodedAccessList []string
	for _, entry := range txn.AccessList {
		encodedAccessList = append(encodedAccessList, hex.EncodeToString(entry))
	}

	accessListStr := "[" + strings.Join(encodedAccessList, ", ") + "]"

	// Encode each element in the AuthorizationList slice individually
	var encodedAuthorizationList []string
	for _, entry := range txn.AuthorizationList {
		encodedAuthorizationList = append(encodedAuthorizationList, hex.EncodeToString(entry))
	}

	authorizationListStr := "[" + strings.Join(encodedAuthorizationList, ", ") + "]"

	return fmt.Sprintf("ChainId: %s\nNonce: %s\nMaxPriorityGas: %s\nMaxGas: %s\nGasLimit: %s\nTo: %s\nValue: %s\nCallData: %s\nAccessList: %s\nAuthorizationList: %s\nRecoveryId: %s\nR: %s\nS: %s",
		hex.EncodeToString(txn.ChainId),
		hex.EncodeToString(txn.Nonce),
		hex.EncodeToString(txn.MaxPriorityGas),
		hex.EncodeToString(txn.MaxGas),
		hex.EncodeToString(txn.GasLimit),
		hex.EncodeToString(txn.To),
		hex.EncodeToString(txn.Value),
		hex.EncodeToString(txn.CallData),
		accessListStr,
		authorizationListStr,
		hex.EncodeToString(txn.RecoveryId),
		hex.EncodeToString(txn.R),
		hex.EncodeToString(txn.S),
	)
}
