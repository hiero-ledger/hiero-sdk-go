package param

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
