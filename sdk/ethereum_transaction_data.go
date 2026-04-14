package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/pkg/errors"
)

// EthereumTransactionBody is implemented by every Ethereum transaction
// variant: Legacy, EIP-1559, EIP-2930, EIP-7702.
type EthereumTransactionBody interface {
	// ToBytes returns the signed RLP, with the type prefix (0x01/0x02/0x04)
	// on typed transactions.
	ToBytes() ([]byte, error)

	// Sign ECDSA-signs the transaction, writes R/S and RecoveryId (V on
	// Legacy) back onto the receiver, and returns the signed RLP.
	Sign(key PrivateKey) ([]byte, error)

	String() string
}

// EthereumTransactionData represents the data of an Ethereum transaction.
type EthereumTransactionData struct {
	eip1559 *EthereumEIP1559Transaction
	eip2930 *EthereumEIP2930Transaction
	eip7702 *EthereumEIP7702Transaction
	legacy  *EthereumLegacyTransaction
}

// NewEthereumTransactionData wraps tx. Returns nil if tx is not one of the
// four concrete variants.
func NewEthereumTransactionData(tx EthereumTransactionBody) *EthereumTransactionData {
	data := &EthereumTransactionData{}
	switch concrete := tx.(type) {
	case *EthereumEIP1559Transaction:
		data.eip1559 = concrete
	case *EthereumEIP2930Transaction:
		data.eip2930 = concrete
	case *EthereumEIP7702Transaction:
		data.eip7702 = concrete
	case *EthereumLegacyTransaction:
		data.legacy = concrete
	default:
		return nil
	}
	return data
}

// GetTransaction returns the wrapped body. Use a type assertion to recover
// the concrete pointer (e.g. *EthereumEIP1559Transaction).
func (ethereumTxData *EthereumTransactionData) GetTransaction() EthereumTransactionBody {
	if ethereumTxData.eip1559 != nil {
		return ethereumTxData.eip1559
	}
	if ethereumTxData.eip2930 != nil {
		return ethereumTxData.eip2930
	}
	if ethereumTxData.eip7702 != nil {
		return ethereumTxData.eip7702
	}
	if ethereumTxData.legacy != nil {
		return ethereumTxData.legacy
	}
	return nil
}

// Sign delegates to the wrapped body.
func (ethereumTxData *EthereumTransactionData) Sign(key PrivateKey) ([]byte, error) {
	body := ethereumTxData.GetTransaction()
	if body == nil {
		return nil, errors.New("transaction data is empty")
	}
	return body.Sign(key)
}

// EthereumTransactionDataFromBytes constructs an EthereumTransactionData from a raw byte array.
func EthereumTransactionDataFromBytes(b []byte) (*EthereumTransactionData, error) {
	var transactionData EthereumTransactionData

	if len(b) == 0 {
		return nil, errors.New("input byte array is empty")
	}

	switch b[0] {
	case 0x02:
		eip1559, err := EthereumEIP1559TransactionFromBytes(b)
		if err != nil {
			return nil, err
		}
		transactionData.eip1559 = eip1559
		return &transactionData, nil
	case 0x01:
		eip2930, err := EthereumEIP2930TransactionFromBytes(b)
		if err != nil {
			return nil, err
		}
		transactionData.eip2930 = eip2930
		return &transactionData, nil
	case 0x04:
		eip7702, err := EthereumEIP7702TransactionFromBytes(b)
		if err != nil {
			return nil, err
		}
		transactionData.eip7702 = eip7702
		return &transactionData, nil
	default:
		legacy, err := EthereumLegacyTransactionFromBytes(b)
		if err != nil {
			return nil, err
		}
		transactionData.legacy = legacy
		return &transactionData, nil
	}
}

// ToBytes returns the raw bytes of the Ethereum transaction.
func (ethereumTxData *EthereumTransactionData) ToBytes() ([]byte, error) {
	if ethereumTxData.eip1559 != nil {
		return ethereumTxData.eip1559.ToBytes()
	}

	if ethereumTxData.eip2930 != nil {
		return ethereumTxData.eip2930.ToBytes()
	}

	if ethereumTxData.eip7702 != nil {
		return ethereumTxData.eip7702.ToBytes()
	}

	if ethereumTxData.legacy != nil {
		return ethereumTxData.legacy.ToBytes()
	}

	return nil, errors.New("transaction data is empty")
}

// GetData retrieves the CallData from the transaction.
func (ethereumTxData *EthereumTransactionData) GetData() []byte {
	if ethereumTxData.eip1559 != nil {
		return ethereumTxData.eip1559.CallData
	}
	if ethereumTxData.eip2930 != nil {
		return ethereumTxData.eip2930.CallData
	}
	if ethereumTxData.eip7702 != nil {
		return ethereumTxData.eip7702.CallData
	}
	return ethereumTxData.legacy.CallData
}

// SetData sets the CallData for the transaction.
func (ethereumTxData *EthereumTransactionData) SetData(data []byte) *EthereumTransactionData {
	if ethereumTxData.eip1559 != nil {
		ethereumTxData.eip1559.CallData = data
		return ethereumTxData
	}
	if ethereumTxData.eip2930 != nil {
		ethereumTxData.eip2930.CallData = data
		return ethereumTxData
	}
	if ethereumTxData.eip7702 != nil {
		ethereumTxData.eip7702.CallData = data
		return ethereumTxData
	}

	ethereumTxData.legacy.CallData = data
	return ethereumTxData
}
