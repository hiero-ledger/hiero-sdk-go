package response

// SPDX-License-Identifier: Apache-2.0

type EthereumTransactionResponse struct {
	ContractId string `json:"contractId"`
	Status     string `json:"status"`
}
