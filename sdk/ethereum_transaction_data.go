package hiero

import (
	"github.com/pkg/errors"
)

// SPDX-License-Identifier: Apache-2.0

// EthereumTransactionData represents the data of an Ethereum transaction.
type EthereumTransactionData struct {
	eip1559 *EthereumEIP1559Transaction
	eip2930 *EthereumEIP2930Transaction
	eip7702 *EthereumEIP7702Transaction
	legacy  *EthereumLegacyTransaction
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
