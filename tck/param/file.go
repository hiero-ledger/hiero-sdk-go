package param

// SPDX-License-Identifier: Apache-2.0
type CreateFileParams struct {
	Keys                    *[]string                `json:"keys"`
	Contents                *string                  `json:"contents"`
	ExpirationTime          *string                  `json:"expirationTime"`
	Memo                    *string                  `json:"memo"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}
