/*
Account Deletion Example

This example demonstrates the complete lifecycle of account creation and deletion using the Hiero SDK:

Step 1: Initialize the client by loading environment variables (HEDERA_NETWORK, OPERATOR_ID, OPERATOR_KEY)
Step 2: Configure the client with operator credentials
Step 3: Generate a new keypair for the account to be created (and later deleted)
Step 4: Create a new account using the generated public key
Step 5: Build an account delete transaction specifying:
  - The account ID to delete
  - A transfer account ID to receive the remaining balance

Step 6: Freeze the transaction to prepare it for signing
Step 7: Manually sign the transaction with the private key of the account being deleted
Step 8: Execute the delete transaction
Step 9: Verify the deletion by checking the transaction receipt status
*/
package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func main() {
	var client *hiero.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
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

	// Generate the key to use with the new account
	newKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	fmt.Println("Creating an account to delete:")
	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	// First create an account
	transactionResponse, err := hiero.NewAccountCreateTransaction().
		// This key will be required to delete the account later
		SetKeyWithoutAlias(newKey.PublicKey()).
		// Initial balance
		SetInitialBalance(hiero.NewHbar(2)).
		SetTransactionMemo("go sdk example delete_account/main.go").
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account creation receipt", err))
	}

	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)
	fmt.Println("deleting created account")

	// To delete an account you must do the following:
	deleteTransaction, err := hiero.NewAccountDeleteTransaction().
		// Set the account to be deleted
		SetAccountID(newAccountID).
		// Set an account ID to transfer the balance of the deleted account to
		SetTransferAccountID(hiero.AccountID{Account: 3}).
		SetTransactionMemo("go sdk example delete_account/main.go").
		FreezeWith(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account delete transaction", err))
	}

	// Manually sign the transaction with the private key of the account to be deleted
	deleteTransaction = deleteTransaction.Sign(newKey)

	// Execute the transaction
	deleteTransactionResponse, err := deleteTransaction.Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error deleting account", err))
	}

	deleteTransactionReceipt, err := deleteTransactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account deletion receipt", err))
	}

	fmt.Printf("account delete transaction status: %v\n", deleteTransactionReceipt.Status)
}
