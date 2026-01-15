package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// AuthorizationTuple represents a single authorization entry: [chainId, contractAddress, nonce, yParity, r, s]
type AuthorizationTuple [6][]byte

// EthereumEIP7702Transaction represents the EIP-7702 Ethereum transaction data.
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
	AuthorizationList []AuthorizationTuple
	RecoveryId        []byte
	R                 []byte
	S                 []byte
}

// nolint
// NewEthereumEIP7702Transaction creates a new EthereumEIP7702Transaction with the provided fields.
func NewEthereumEIP7702Transaction(
	chainId, nonce, maxPriorityGas, maxGas, gasLimit, to, value, callData, recoveryId, r, s []byte, accessList [][]byte, authorizationList []AuthorizationTuple) *EthereumEIP7702Transaction {
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
		return nil, errors.New("input byte array is malformed; it should start with 0x04 followed by 13 RLP-encoded elements")
	}

	// Remove the prefix byte (0x04)
	item := NewRLPItem(LIST_TYPE)
	if err := item.Read(bytes[1:]); err != nil {
		return nil, errors.Wrap(err, "failed to read RLP data")
	}

	if item.itemType != LIST_TYPE || len(item.childItems) != 13 {
		return nil, errors.New("input byte array is malformed; it should be a list of 13 RLP-encoded elements")
	}

	// Handle the access list
	var accessListValues [][]byte
	for _, child := range item.childItems[8].childItems {
		accessListValues = append(accessListValues, child.itemValue)
	}

	// Handle the authorization list: array of [chainId, contractAddress, nonce, yParity, r, s] tuples
	var authorizationListValues []AuthorizationTuple
	if item.childItems[9].itemType != LIST_TYPE {
		return nil, errors.New("authorization list must be an array")
	}
	for _, authTupleItem := range item.childItems[9].childItems {
		if authTupleItem.itemType != LIST_TYPE {
			return nil, errors.New("invalid authorization list entry: must be a list")
		}
		if len(authTupleItem.childItems) != 6 {
			return nil, errors.New("invalid authorization list entry: must be [chainId, contractAddress, nonce, yParity, r, s]")
		}
		var tuple AuthorizationTuple
		for i := 0; i < 6; i++ {
			tuple[i] = authTupleItem.childItems[i].itemValue
		}
		authorizationListValues = append(authorizationListValues, tuple)
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
		item.childItems[10].itemValue,
		item.childItems[11].itemValue,
		item.childItems[12].itemValue,
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
	for _, authTuple := range txn.AuthorizationList {
		// Each authorization entry is a tuple: [chainId, contractAddress, nonce, yParity, r, s]
		authTupleItem := NewRLPItem(LIST_TYPE)
		for i := 0; i < 6; i++ {
			authTupleItem.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(authTuple[i]))
		}
		authorizationListItem.PushBack(authTupleItem)
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

	// Encode each authorization tuple in the AuthorizationList
	var encodedAuthorizationList []string
	for _, authTuple := range txn.AuthorizationList {
		var tupleParts []string
		for i := 0; i < 6; i++ {
			tupleParts = append(tupleParts, hex.EncodeToString(authTuple[i]))
		}
		encodedAuthorizationList = append(encodedAuthorizationList, "["+strings.Join(tupleParts, ", ")+"]")
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
