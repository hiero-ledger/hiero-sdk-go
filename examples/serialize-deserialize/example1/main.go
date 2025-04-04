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

	// Generate new key to use with new account
	newKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	resp, err := hiero.NewAccountCreateTransaction().SetKeyWithoutAlias(newKey).Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err := resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	newAccountId := *receipt.AccountID

	bytes, err := hiero.NewTransferTransaction().AddHbarTransfer(operatorAccountID, hiero.NewHbar(1)).
		ToBytes()
	if err != nil {
		panic(err)
	}

	txFromBytes, err := hiero.TransactionFromBytes(bytes)
	if err != nil {
		panic(err)
	}

	transaction := txFromBytes.(hiero.TransferTransaction)
	_, err = transaction.AddHbarTransfer(newAccountId, hiero.NewHbar(-1)).SignWithOperator(client)
	if err != nil {
		panic(err)
	}

	_, err = transaction.Execute(client)
	if err != nil {
		panic(err)
	}

	// Get the `AccountInfo` on the new account and show it is a hollow account by not having a public key
	info, err := hiero.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err)
	}

	fmt.Println("Balance of new account: ", info.Balance)
}
