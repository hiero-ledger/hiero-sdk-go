package param

// SPDX-License-Identifier: Apache-2.0

type EthereumTransactionParams struct {
	EthereumData            *string                  `json:"ethereumData,omitempty"`
	CallDataFileID          *string                  `json:"callDataFileId,omitempty"`
	MaxGasAllowed           *string                  `json:"maxGasAllowed,omitempty"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}
