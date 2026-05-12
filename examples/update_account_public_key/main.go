package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to update an account's key.
func main() {
	fmt.Println("Update Account Public Key Example Start!")

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

	// Setting the client operator ID and key, plus a generous default tx fee.
	client.SetOperator(operatorAccountID, operatorKey)
	if err := client.SetDefaultMaxTransactionFee(hiero.NewHbar(10)); err != nil {
		panic(fmt.Sprintf("%v : error setting default max transaction fee", err))
	}

	// Step 1: Generate ED25519 key pairs.
	fmt.Println("Generating ED25519 key pairs...")
	privateKey1, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}
	privateKey2, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Step 2: Create a new account using privateKey1's public key.
	fmt.Println("Creating new account...")
	accountTxResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(privateKey1.PublicKey()).
		SetInitialBalance(hiero.NewHbar(1)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}

	accountTxReceipt, err := accountTxResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account creation receipt", err))
	}
	accountID := *accountTxReceipt.AccountID
	fmt.Printf("Created new account with ID: %v and public key: %v\n", accountID, privateKey1.PublicKey())

	// Step 3: Update the account's key to privateKey2's public key.
	// Both signatures are required: the previous key authorizes the change,
	// the new key proves possession.
	fmt.Printf("Updating public key of new account...(Setting key: %v).\n", privateKey2.PublicKey())
	accountUpdateTx, err := hiero.NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetKey(privateKey2.PublicKey()).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account update transaction", err))
	}

	accountUpdateTx.Sign(privateKey1)
	accountUpdateTx.Sign(privateKey2)

	accountUpdateTxResponse, err := accountUpdateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating account", err))
	}

	if _, err := accountUpdateTxResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error getting account update receipt", err))
	}

	// Step 4: Confirm the key was changed via AccountInfoQuery.
	info, err := hiero.NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account info query", err))
	}
	fmt.Printf("New account public key: %v\n", info.Key)

	// Cleanup: delete the created account, transferring remaining funds back to the operator.
	deleteTx, err := hiero.NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(operatorAccountID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account delete transaction", err))
	}

	deleteTx.Sign(privateKey2)

	deleteResponse, err := deleteTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account delete", err))
	}
	if _, err := deleteResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error getting account delete receipt", err))
	}

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}
	fmt.Println("Update Account Public Key Example Complete!")
}
