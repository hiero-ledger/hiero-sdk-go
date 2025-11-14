package param

// SPDX-License-Identifier: Apache-2.0

type EthereumTransactionParams struct {
	BaseParams
	EthereumData   *string `json:"ethereumData,omitempty"`
	CallDataFileID *string `json:"callDataFileId,omitempty"`
	MaxGasAllowed  *string `json:"maxGasAllowance,omitempty"`
}
