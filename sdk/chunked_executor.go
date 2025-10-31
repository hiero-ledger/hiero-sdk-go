package hiero

// SPDX-License-Identifier: Apache-2.0

// ExecuteAll executes all chunks of a transaction with proper throttling handling
func executeAll(
	tx TransactionInterface,
	client *Client,
) ([]TransactionResponse, error) {
	if client == nil || client.operator == nil {
		return []TransactionResponse{}, errNoClientProvided
	}
	baseTxn := tx.getBaseTransaction()

	if !baseTxn.IsFrozen() {
		var err error
		fileAppendTx, ok := tx.(*FileAppendTransaction)
		if ok {
			_, err = fileAppendTx.FreezeWith(client)
		}

		messageSubmitTx, ok := tx.(*TopicMessageSubmitTransaction)
		if ok {
			_, err = messageSubmitTx.FreezeWith(client)
		}

		if err != nil {
			return []TransactionResponse{}, err
		}
	}

	transactionID := baseTxn.GetTransactionID()
	accountID := AccountID{}
	if transactionID.AccountID != nil {
		accountID = *transactionID.AccountID
	}

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(accountID) {
		baseTxn.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	size := baseTxn.signedTransactions._Length() / baseTxn.nodeAccountIDs._Length()
	list := make([]TransactionResponse, size)

	for i := 0; i < size; i++ {
		resp, err := _Execute(client, tx)
		if err != nil {
			return list, err
		}

		list[i] = resp.(TransactionResponse)
		receipt, err := list[i].GetReceipt(client)
		if err != nil {
			return list, err
		}

		// Retry in case of throttle error
		for receipt.Status == StatusThrottledAtConsensus {
			baseTxn.regenerateID(client)
			resp, err := _Execute(client, tx)
			if err != nil {
				return list, err
			}
			respTx := resp.(TransactionResponse)
			receipt, err = NewTransactionReceiptQuery().
				SetTransactionID(respTx.TransactionID).
				SetNodeAccountIDs([]AccountID{respTx.NodeID}).
				SetIncludeChildren(respTx.IncludeChildReceipts).
				Execute(client)

			// If we get a non-throttled receipt, we can break out of the loop
			if err == nil && receipt.Status != StatusThrottledAtConsensus {
				list[i] = respTx
				break
			}
		}

		// Validate the receipt status
		if err = receipt.ValidateStatus(true); err != nil {
			return list, err
		}
	}

	return list, nil
}
