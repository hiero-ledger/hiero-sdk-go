package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"math/big"
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

// _toUnsignedRLP builds the unsigned portion of the EIP-7702 RLP list
func (txn *EthereumEIP7702Transaction) _toUnsignedRLP() *RLPItem {
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
	return item
}

// _encodeWithSignature appends RecoveryId/R/S to the given (already-built)
// unsigned RLP list, serializes it, and prepends the 0x04 type prefix.
func (txn *EthereumEIP7702Transaction) _encodeWithSignature(item *RLPItem) ([]byte, error) {
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.RecoveryId))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.R))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.S))

	transactionBytes, err := item.Write()
	if err != nil {
		return nil, err
	}

	// Append 04 byte as it is the standard for EIP7702
	return append([]byte{0x04}, transactionBytes...), nil
}

// ToBytes encodes the EthereumEIP7702Transaction into RLP format.
func (txn *EthereumEIP7702Transaction) ToBytes() ([]byte, error) {
	return txn._encodeWithSignature(txn._toUnsignedRLP())
}

// Sign signs the unsigned transaction with the given ECDSA key, populates
// RecoveryId/R/S on the receiver, and returns the signed RLP bytes ready to
// wrap in EthereumTransaction. AuthorizationList entries are expected to
// already carry their own signatures (yParity, r, s) per the EIP-7702 spec.
func (txn *EthereumEIP7702Transaction) Sign(key PrivateKey) ([]byte, error) {
	item := txn._toUnsignedRLP()
	unsignedBytes, err := item.Write()
	if err != nil {
		return nil, err
	}
	message := append([]byte{0x04}, unsignedBytes...)

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
	txn.RecoveryId = []byte{byte(v)}

	return txn._encodeWithSignature(item)
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

// GetChainId returns the chain id as a uint64.
func (txn *EthereumEIP7702Transaction) GetChainId() uint64 {
	return _ethBytesToUint64(txn.ChainId)
}

// SetChainId sets the chain id from a uint64.
func (txn *EthereumEIP7702Transaction) SetChainId(v uint64) *EthereumEIP7702Transaction {
	txn.ChainId = _uint64ToEthBytes(v)
	return txn
}

// GetChainIdBytes returns the raw canonical big-endian chain id bytes.
func (txn *EthereumEIP7702Transaction) GetChainIdBytes() []byte { return txn.ChainId }

// SetChainIdBytes sets the chain id from raw canonical big-endian bytes.
func (txn *EthereumEIP7702Transaction) SetChainIdBytes(v []byte) *EthereumEIP7702Transaction {
	txn.ChainId = v
	return txn
}

// GetNonce returns the nonce as a uint64.
func (txn *EthereumEIP7702Transaction) GetNonce() uint64 {
	return _ethBytesToUint64(txn.Nonce)
}

// SetNonce sets the nonce from a uint64.
func (txn *EthereumEIP7702Transaction) SetNonce(v uint64) *EthereumEIP7702Transaction {
	txn.Nonce = _uint64ToEthBytes(v)
	return txn
}

// GetNonceBytes returns the raw canonical big-endian nonce bytes.
func (txn *EthereumEIP7702Transaction) GetNonceBytes() []byte { return txn.Nonce }

// SetNonceBytes sets the nonce from raw canonical big-endian bytes.
func (txn *EthereumEIP7702Transaction) SetNonceBytes(v []byte) *EthereumEIP7702Transaction {
	txn.Nonce = v
	return txn
}

// GetMaxPriorityGas returns the max priority fee per gas as a *big.Int.
func (txn *EthereumEIP7702Transaction) GetMaxPriorityGas() *big.Int {
	return _ethBytesToBigInt(txn.MaxPriorityGas)
}

// SetMaxPriorityGas sets the max priority fee per gas from a *big.Int.
func (txn *EthereumEIP7702Transaction) SetMaxPriorityGas(v *big.Int) *EthereumEIP7702Transaction {
	txn.MaxPriorityGas = _bigIntToEthBytes(v)
	return txn
}

// GetMaxPriorityGasBytes returns the raw canonical big-endian max priority fee bytes.
func (txn *EthereumEIP7702Transaction) GetMaxPriorityGasBytes() []byte { return txn.MaxPriorityGas }

// SetMaxPriorityGasBytes sets the max priority fee per gas from raw bytes.
func (txn *EthereumEIP7702Transaction) SetMaxPriorityGasBytes(v []byte) *EthereumEIP7702Transaction {
	txn.MaxPriorityGas = v
	return txn
}

// GetMaxGas returns the max fee per gas as a *big.Int.
func (txn *EthereumEIP7702Transaction) GetMaxGas() *big.Int {
	return _ethBytesToBigInt(txn.MaxGas)
}

// SetMaxGas sets the max fee per gas from a *big.Int.
func (txn *EthereumEIP7702Transaction) SetMaxGas(v *big.Int) *EthereumEIP7702Transaction {
	txn.MaxGas = _bigIntToEthBytes(v)
	return txn
}

// GetMaxGasBytes returns the raw canonical big-endian max fee bytes.
func (txn *EthereumEIP7702Transaction) GetMaxGasBytes() []byte { return txn.MaxGas }

// SetMaxGasBytes sets the max fee per gas from raw bytes.
func (txn *EthereumEIP7702Transaction) SetMaxGasBytes(v []byte) *EthereumEIP7702Transaction {
	txn.MaxGas = v
	return txn
}

// GetGasLimit returns the gas limit as a uint64.
func (txn *EthereumEIP7702Transaction) GetGasLimit() uint64 {
	return _ethBytesToUint64(txn.GasLimit)
}

// SetGasLimit sets the gas limit from a uint64.
func (txn *EthereumEIP7702Transaction) SetGasLimit(v uint64) *EthereumEIP7702Transaction {
	txn.GasLimit = _uint64ToEthBytes(v)
	return txn
}

// GetGasLimitBytes returns the raw canonical big-endian gas limit bytes.
func (txn *EthereumEIP7702Transaction) GetGasLimitBytes() []byte { return txn.GasLimit }

// SetGasLimitBytes sets the gas limit from raw bytes.
func (txn *EthereumEIP7702Transaction) SetGasLimitBytes(v []byte) *EthereumEIP7702Transaction {
	txn.GasLimit = v
	return txn
}

// GetTo returns the recipient address bytes.
func (txn *EthereumEIP7702Transaction) GetTo() []byte { return txn.To }

// SetTo sets the recipient address bytes.
func (txn *EthereumEIP7702Transaction) SetTo(v []byte) *EthereumEIP7702Transaction {
	txn.To = v
	return txn
}

// GetValue returns the transaction value (wei) as a *big.Int.
func (txn *EthereumEIP7702Transaction) GetValue() *big.Int {
	return _ethBytesToBigInt(txn.Value)
}

// SetValue sets the transaction value (wei) from a *big.Int.
func (txn *EthereumEIP7702Transaction) SetValue(v *big.Int) *EthereumEIP7702Transaction {
	txn.Value = _bigIntToEthBytes(v)
	return txn
}

// GetValueBytes returns the raw canonical big-endian value bytes.
func (txn *EthereumEIP7702Transaction) GetValueBytes() []byte { return txn.Value }

// SetValueBytes sets the value from raw bytes.
func (txn *EthereumEIP7702Transaction) SetValueBytes(v []byte) *EthereumEIP7702Transaction {
	txn.Value = v
	return txn
}

// GetCallData returns the call data.
func (txn *EthereumEIP7702Transaction) GetCallData() []byte { return txn.CallData }

// SetCallData sets the call data.
func (txn *EthereumEIP7702Transaction) SetCallData(v []byte) *EthereumEIP7702Transaction {
	txn.CallData = v
	return txn
}

// GetRecoveryId returns the recovery id as an int.
func (txn *EthereumEIP7702Transaction) GetRecoveryId() int {
	if len(txn.RecoveryId) == 0 {
		return 0
	}
	return int(txn.RecoveryId[0])
}

// SetRecoveryId sets the recovery id from an int.
func (txn *EthereumEIP7702Transaction) SetRecoveryId(v int) *EthereumEIP7702Transaction {
	if v == 0 {
		txn.RecoveryId = []byte{}
	} else {
		txn.RecoveryId = []byte{byte(v)}
	}
	return txn
}

// GetRecoveryIdBytes returns the raw recovery id bytes.
func (txn *EthereumEIP7702Transaction) GetRecoveryIdBytes() []byte { return txn.RecoveryId }

// SetRecoveryIdBytes sets the recovery id from raw bytes.
func (txn *EthereumEIP7702Transaction) SetRecoveryIdBytes(v []byte) *EthereumEIP7702Transaction {
	txn.RecoveryId = v
	return txn
}

// GetR returns the R signature component.
func (txn *EthereumEIP7702Transaction) GetR() []byte { return txn.R }

// SetR sets the R signature component.
func (txn *EthereumEIP7702Transaction) SetR(v []byte) *EthereumEIP7702Transaction {
	txn.R = v
	return txn
}

// GetS returns the S signature component.
func (txn *EthereumEIP7702Transaction) GetS() []byte { return txn.S }

// SetS sets the S signature component.
func (txn *EthereumEIP7702Transaction) SetS(v []byte) *EthereumEIP7702Transaction {
	txn.S = v
	return txn
}

// GetAccessListItems returns the access list as structured AccessListItem entries.
func (txn *EthereumEIP7702Transaction) GetAccessListItems() []AccessListItem {
	return _accessListItemsFromBytes(txn.AccessList)
}

// SetAccessListItems replaces the access list from structured AccessListItem entries.
func (txn *EthereumEIP7702Transaction) SetAccessListItems(items []AccessListItem) *EthereumEIP7702Transaction {
	txn.AccessList = _accessListItemsToBytes(items)
	return txn
}

// AddAccessListItem appends a single access list entry.
func (txn *EthereumEIP7702Transaction) AddAccessListItem(item AccessListItem) *EthereumEIP7702Transaction {
	txn.AccessList = append(txn.AccessList, _accessListItemToBytes(item))
	return txn
}

// GetAuthorizationList returns the EIP-7702 authorization list.
func (txn *EthereumEIP7702Transaction) GetAuthorizationList() []AuthorizationTuple {
	return txn.AuthorizationList
}

// SetAuthorizationList replaces the EIP-7702 authorization list.
func (txn *EthereumEIP7702Transaction) SetAuthorizationList(list []AuthorizationTuple) *EthereumEIP7702Transaction {
	txn.AuthorizationList = list
	return txn
}

// AddAuthorization appends a single authorization tuple to the list.
func (txn *EthereumEIP7702Transaction) AddAuthorization(tuple AuthorizationTuple) *EthereumEIP7702Transaction {
	txn.AuthorizationList = append(txn.AuthorizationList, tuple)
	return txn
}
