package param

// SPDX-License-Identifier: Apache-2.0
type CreateFileParams struct {
	Keys                    *[]string                `json:"keys"`
	Contents                *string                  `json:"contents"`
	ExpirationTime          *string                  `json:"expirationTime"`
	Memo                    *string                  `json:"memo"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}

type UpdateFileParams struct {
	FileId                  *string                  `json:"fileId"`
	Keys                    *[]string                `json:"keys"`
	Contents                *string                  `json:"contents"`
	ExpirationTime          *string                  `json:"expirationTime"`
	Memo                    *string                  `json:"memo"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}

type DeleteFileParams struct {
	FileId                  *string                  `json:"fileId"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}

type AppendFileParams struct {
	FileId                  *string                  `json:"fileId"`
	Contents                *string                  `json:"contents"`
	MaxChunks               *int                     `json:"maxChunks"`
	ChunkSize               *int                     `json:"chunkSize"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}
