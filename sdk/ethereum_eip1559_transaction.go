package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

// EthereumEIP1559Transaction represents the EIP-1559 Ethereum transaction data.
type EthereumEIP1559Transaction struct {
	ChainId        []byte
	Nonce          []byte
	MaxPriorityGas []byte
	MaxGas         []byte
	GasLimit       []byte
	To             []byte
	Value          []byte
	CallData       []byte
	AccessList     [][]byte
	RecoveryId     []byte
	R              []byte
	S              []byte
}

// nolint
// NewEthereumEIP1559Transaction creates a new EthereumEIP1559Transaction with the provided fields.
func NewEthereumEIP1559Transaction(
	chainId, nonce, maxPriorityGas, maxGas, gasLimit, to, value, callData, recoveryId, r, s []byte, accessList [][]byte) *EthereumEIP1559Transaction {
	return &EthereumEIP1559Transaction{
		ChainId:        chainId,
		Nonce:          nonce,
		MaxPriorityGas: maxPriorityGas,
		MaxGas:         maxGas,
		GasLimit:       gasLimit,
		To:             to,
		Value:          value,
		CallData:       callData,
		AccessList:     accessList,
		RecoveryId:     recoveryId,
		R:              r,
		S:              s,
	}
}

// FromBytes decodes the RLP encoded bytes into an EthereumEIP1559Transaction.
func EthereumEIP1559TransactionFromBytes(bytes []byte) (*EthereumEIP1559Transaction, error) {
	if len(bytes) == 0 || bytes[0] != 0x02 {
		return nil, errors.New("input byte array is malformed; it should start with 0x02 followed by 12 RLP-encoded elements")
	}

	// Remove the prefix byte (0x02)
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

	// Extract values from the RLP item
	return NewEthereumEIP1559Transaction(
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
	), nil
}

// _toUnsignedRLP builds the unsigned portion of the EIP-1559 RLP list
func (txn *EthereumEIP1559Transaction) _toUnsignedRLP() *RLPItem {
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
	return item
}

// _encodeWithSignature appends RecoveryId/R/S to the given
// unsigned RLP list, serializes it, and prepends the 0x02 type prefix.
func (txn *EthereumEIP1559Transaction) _encodeWithSignature(item *RLPItem) ([]byte, error) {
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.RecoveryId))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.R))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.S))
	bytes, err := item.Write()
	if err != nil {
		return nil, err
	}

	// Append 02 byte as it is the standard for EIP1559
	return append([]byte{0x02}, bytes...), nil
}

// ToBytes encodes the EthereumEIP1559Transaction into RLP format.
func (txn *EthereumEIP1559Transaction) ToBytes() ([]byte, error) {
	return txn._encodeWithSignature(txn._toUnsignedRLP())
}

// Sign signs the unsigned transaction with the given ECDSA key, populates
// RecoveryId/R/S on the receiver, and returns the signed RLP bytes ready to
// wrap in EthereumTransaction.
func (txn *EthereumEIP1559Transaction) Sign(key PrivateKey) ([]byte, error) {
	item := txn._toUnsignedRLP()
	unsignedBytes, err := item.Write()
	if err != nil {
		return nil, err
	}
	message := append([]byte{0x02}, unsignedBytes...)

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

// String returns a string representation of the EthereumEIP1559Transaction.
func (txn *EthereumEIP1559Transaction) String() string {
	// Encode each element in the AccessList slice individually
	var encodedAccessList []string
	for _, entry := range txn.AccessList {
		encodedAccessList = append(encodedAccessList, hex.EncodeToString(entry))
	}

	accessListStr := "[" + strings.Join(encodedAccessList, ", ") + "]"

	return fmt.Sprintf("ChainId: %s\nNonce: %s\nMaxPriorityGas: %s\nMaxGas: %s\nGasLimit: %s\nTo: %s\nValue: %s\nCallData: %s\nAccessList: %s\nRecoveryId: %s\nR: %s\nS: %s",
		hex.EncodeToString(txn.ChainId),
		hex.EncodeToString(txn.Nonce),
		hex.EncodeToString(txn.MaxPriorityGas),
		hex.EncodeToString(txn.MaxGas),
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
func (txn *EthereumEIP1559Transaction) GetChainId() uint64 {
	return _ethBytesToUint64(txn.ChainId)
}

// SetChainId sets the chain id from a uint64.
func (txn *EthereumEIP1559Transaction) SetChainId(v uint64) *EthereumEIP1559Transaction {
	txn.ChainId = _uint64ToEthBytes(v)
	return txn
}

// GetChainIdBytes returns the raw canonical big-endian chain id bytes.
func (txn *EthereumEIP1559Transaction) GetChainIdBytes() []byte { return txn.ChainId }

// SetChainIdBytes sets the chain id from raw canonical big-endian bytes.
func (txn *EthereumEIP1559Transaction) SetChainIdBytes(v []byte) *EthereumEIP1559Transaction {
	txn.ChainId = v
	return txn
}

// GetNonce returns the nonce as a uint64.
func (txn *EthereumEIP1559Transaction) GetNonce() uint64 {
	return _ethBytesToUint64(txn.Nonce)
}

// SetNonce sets the nonce from a uint64.
func (txn *EthereumEIP1559Transaction) SetNonce(v uint64) *EthereumEIP1559Transaction {
	txn.Nonce = _uint64ToEthBytes(v)
	return txn
}

// GetNonceBytes returns the raw canonical big-endian nonce bytes.
func (txn *EthereumEIP1559Transaction) GetNonceBytes() []byte { return txn.Nonce }

// SetNonceBytes sets the nonce from raw canonical big-endian bytes.
func (txn *EthereumEIP1559Transaction) SetNonceBytes(v []byte) *EthereumEIP1559Transaction {
	txn.Nonce = v
	return txn
}

// GetMaxPriorityGas returns the max priority fee per gas as a *big.Int.
func (txn *EthereumEIP1559Transaction) GetMaxPriorityGas() *big.Int {
	return _ethBytesToBigInt(txn.MaxPriorityGas)
}

// SetMaxPriorityGas sets the max priority fee per gas from a *big.Int.
func (txn *EthereumEIP1559Transaction) SetMaxPriorityGas(v *big.Int) *EthereumEIP1559Transaction {
	txn.MaxPriorityGas = _bigIntToEthBytes(v)
	return txn
}

// GetMaxPriorityGasBytes returns the raw canonical big-endian max priority fee bytes.
func (txn *EthereumEIP1559Transaction) GetMaxPriorityGasBytes() []byte { return txn.MaxPriorityGas }

// SetMaxPriorityGasBytes sets the max priority fee per gas from raw bytes.
func (txn *EthereumEIP1559Transaction) SetMaxPriorityGasBytes(v []byte) *EthereumEIP1559Transaction {
	txn.MaxPriorityGas = v
	return txn
}

// GetMaxGas returns the max fee per gas as a *big.Int.
func (txn *EthereumEIP1559Transaction) GetMaxGas() *big.Int {
	return _ethBytesToBigInt(txn.MaxGas)
}

// SetMaxGas sets the max fee per gas from a *big.Int.
func (txn *EthereumEIP1559Transaction) SetMaxGas(v *big.Int) *EthereumEIP1559Transaction {
	txn.MaxGas = _bigIntToEthBytes(v)
	return txn
}

// GetMaxGasBytes returns the raw canonical big-endian max fee bytes.
func (txn *EthereumEIP1559Transaction) GetMaxGasBytes() []byte { return txn.MaxGas }

// SetMaxGasBytes sets the max fee per gas from raw bytes.
func (txn *EthereumEIP1559Transaction) SetMaxGasBytes(v []byte) *EthereumEIP1559Transaction {
	txn.MaxGas = v
	return txn
}

// GetGasLimit returns the gas limit as a uint64.
func (txn *EthereumEIP1559Transaction) GetGasLimit() uint64 {
	return _ethBytesToUint64(txn.GasLimit)
}

// SetGasLimit sets the gas limit from a uint64.
func (txn *EthereumEIP1559Transaction) SetGasLimit(v uint64) *EthereumEIP1559Transaction {
	txn.GasLimit = _uint64ToEthBytes(v)
	return txn
}

// GetGasLimitBytes returns the raw canonical big-endian gas limit bytes.
func (txn *EthereumEIP1559Transaction) GetGasLimitBytes() []byte { return txn.GasLimit }

// SetGasLimitBytes sets the gas limit from raw bytes.
func (txn *EthereumEIP1559Transaction) SetGasLimitBytes(v []byte) *EthereumEIP1559Transaction {
	txn.GasLimit = v
	return txn
}

// GetTo returns the recipient address bytes.
func (txn *EthereumEIP1559Transaction) GetTo() []byte { return txn.To }

// SetTo sets the recipient address bytes.
func (txn *EthereumEIP1559Transaction) SetTo(v []byte) *EthereumEIP1559Transaction {
	txn.To = v
	return txn
}

// GetValue returns the transaction value (wei) as a *big.Int.
func (txn *EthereumEIP1559Transaction) GetValue() *big.Int {
	return _ethBytesToBigInt(txn.Value)
}

// SetValue sets the transaction value (wei) from a *big.Int.
func (txn *EthereumEIP1559Transaction) SetValue(v *big.Int) *EthereumEIP1559Transaction {
	txn.Value = _bigIntToEthBytes(v)
	return txn
}

// GetValueBytes returns the raw canonical big-endian value bytes.
func (txn *EthereumEIP1559Transaction) GetValueBytes() []byte { return txn.Value }

// SetValueBytes sets the value from raw bytes.
func (txn *EthereumEIP1559Transaction) SetValueBytes(v []byte) *EthereumEIP1559Transaction {
	txn.Value = v
	return txn
}

// GetCallData returns the call data.
func (txn *EthereumEIP1559Transaction) GetCallData() []byte { return txn.CallData }

// SetCallData sets the call data.
func (txn *EthereumEIP1559Transaction) SetCallData(v []byte) *EthereumEIP1559Transaction {
	txn.CallData = v
	return txn
}

// GetRecoveryId returns the recovery id as an int.
func (txn *EthereumEIP1559Transaction) GetRecoveryId() int {
	if len(txn.RecoveryId) == 0 {
		return 0
	}
	return int(txn.RecoveryId[0])
}

// SetRecoveryId sets the recovery id from an int.
func (txn *EthereumEIP1559Transaction) SetRecoveryId(v int) *EthereumEIP1559Transaction {
	if v == 0 {
		txn.RecoveryId = []byte{}
	} else {
		txn.RecoveryId = []byte{byte(v)}
	}
	return txn
}

// GetRecoveryIdBytes returns the raw recovery id bytes.
func (txn *EthereumEIP1559Transaction) GetRecoveryIdBytes() []byte { return txn.RecoveryId }

// SetRecoveryIdBytes sets the recovery id from raw bytes.
func (txn *EthereumEIP1559Transaction) SetRecoveryIdBytes(v []byte) *EthereumEIP1559Transaction {
	txn.RecoveryId = v
	return txn
}

// GetR returns the R signature component.
func (txn *EthereumEIP1559Transaction) GetR() []byte { return txn.R }

// SetR sets the R signature component.
func (txn *EthereumEIP1559Transaction) SetR(v []byte) *EthereumEIP1559Transaction {
	txn.R = v
	return txn
}

// GetS returns the S signature component.
func (txn *EthereumEIP1559Transaction) GetS() []byte { return txn.S }

// SetS sets the S signature component.
func (txn *EthereumEIP1559Transaction) SetS(v []byte) *EthereumEIP1559Transaction {
	txn.S = v
	return txn
}

// GetAccessListItems returns the access list as structured AccessListItem entries.
func (txn *EthereumEIP1559Transaction) GetAccessListItems() []AccessListItem {
	return _accessListItemsFromBytes(txn.AccessList)
}

// SetAccessListItems replaces the access list from structured AccessListItem entries.
func (txn *EthereumEIP1559Transaction) SetAccessListItems(items []AccessListItem) *EthereumEIP1559Transaction {
	txn.AccessList = _accessListItemsToBytes(items)
	return txn
}

// AddAccessListItem appends a single access list entry.
func (txn *EthereumEIP1559Transaction) AddAccessListItem(item AccessListItem) *EthereumEIP1559Transaction {
	txn.AccessList = append(txn.AccessList, _accessListItemToBytes(item))
	return txn
}