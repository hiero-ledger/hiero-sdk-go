// Example demonstrating how to create and delete an account using the Hiero Go SDK.
//
// This program performs the full lifecycle of creating a temporary account
// and then deleting it:
//
//  1. Initialize a client using environment variables (HEDERA_NETWORK, OPERATOR_ID, OPERATOR_KEY)
//  2. Generate a new key pair for the temporary account
//  3. Create a new account with an initial balance of 2 HBAR
//  4. Build an AccountDeleteTransaction specifying the account to delete
//     and a transfer account to receive the remaining balance
//  5. Freeze the transaction to prepare it for signing
//  6. Sign the transaction with the private key of the account being deleted
//     (required because only the account owner can authorize its deletion)
//  7. Execute the transaction and verify successful deletion via the receipt
//
// When an account is deleted, any remaining balance must be transferred to another account.
// The account being deleted must sign the transaction to authorize the deletion.
//
// Required environment variables:
//   - HEDERA_NETWORK: The network to connect to (e.g., "testnet", "mainnet")
//   - OPERATOR_ID: The operator account ID
//   - OPERATOR_KEY: The operator private key
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
