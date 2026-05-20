package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to transfer Hbar between accounts.
func main() {
	fmt.Println("Transfer Crypto Example Start!")

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	recipientID := hiero.AccountID{Account: 3}

	// Step 1: Check Hbar balance of sender and recipient.
	senderBalanceBefore, err := hiero.NewAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying sender balance", err))
	}

	recipientBalanceBefore, err := hiero.NewAccountBalanceQuery().
		SetAccountID(recipientID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying recipient balance", err))
	}

	fmt.Printf("Sender (%v) balance before transfer: %v\n", operatorAccountID, senderBalanceBefore.Hbars)
	fmt.Printf("Recipient (%v) balance before transfer: %v\n", recipientID, recipientBalanceBefore.Hbars)

	// Step 2: Execute the transfer transaction to send Hbars from operator to recipient.
	fmt.Println("Executing the transfer transaction...")
	transferAmount := hiero.NewHbar(1)
	transferResponse, err := hiero.NewTransferTransaction().
		// AddHbarTransfer can be called as many times as you want as long as the total
		// sum of inputs and outputs is zero.
		AddHbarTransfer(operatorAccountID, transferAmount.Negated()).
		AddHbarTransfer(recipientID, transferAmount).
		SetTransactionMemo("Transfer example").
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing transfer", err))
	}

	record, err := transferResponse.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving transfer record", err))
	}

	fmt.Printf("Transferred %v\n", transferAmount)
	fmt.Printf("Transfer memo: %v\n", record.TransactionMemo)

	// Step 3: Check Hbar balance of sender and recipient after the transfer.
	senderBalanceAfter, err := hiero.NewAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying sender balance after", err))
	}

	recipientBalanceAfter, err := hiero.NewAccountBalanceQuery().
		SetAccountID(recipientID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying recipient balance after", err))
	}

	fmt.Printf("Sender (%v) balance after transfer: %v\n", operatorAccountID, senderBalanceAfter.Hbars)
	fmt.Printf("Recipient (%v) balance after transfer: %v\n", recipientID, recipientBalanceAfter.Hbars)

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}
	fmt.Println("Example complete!")
}
