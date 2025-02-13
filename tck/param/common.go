package param

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type CommonTransactionParams struct {
	TransactionId            *string   `json:"transactionId"`
	MaxTransactionFee        *int64    `json:"maxTransactionFee"`
	ValidTransactionDuration *uint64   `json:"validTransactionDuration"`
	Memo                     *string   `json:"memo"`
	RegenerateTransactionId  *bool     `json:"regenerateTransactionId"`
	Signers                  *[]string `json:"signers"`
}

func (common *CommonTransactionParams) FillOutTransaction(transactionInterface hiero.TransactionInterface, client *hiero.Client) error {
	if common.TransactionId != nil {
		txId, _ := hiero.TransactionIdFromString(*common.TransactionId)
		_, err := hiero.TransactionSetTransactionID(transactionInterface, txId)
		if err != nil {
			return err
		}
	}

	if common.MaxTransactionFee != nil {
		_, err := hiero.TransactionSetMaxTransactionFee(transactionInterface, hiero.HbarFromTinybar(*common.MaxTransactionFee))
		if err != nil {
			return err
		}
	}

	if common.ValidTransactionDuration != nil {
		_, err := hiero.TransactionSetTransactionValidDuration(transactionInterface, time.Duration(*common.ValidTransactionDuration)*time.Second)
		if err != nil {
			return err
		}
	}

	if common.Memo != nil {
		_, err := hiero.TransactionSetTransactionMemo(transactionInterface, *common.Memo)
		if err != nil {
			return err
		}
	}

	if common.RegenerateTransactionId != nil {
		_, err := hiero.TransactionSetTransactionID(transactionInterface, hiero.TransactionIDGenerate(client.GetOperatorAccountID()))
		if err != nil {
			return err
		}
	}

	if common.Signers != nil {
		_, err := hiero.TransactionFreezeWith(transactionInterface, client)
		if err != nil {
			return err
		}
		for _, signer := range *common.Signers {
			s, _ := hiero.PrivateKeyFromString(signer)
			_, err = hiero.TransactionSign(transactionInterface, s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
