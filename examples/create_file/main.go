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

	transactionResponse, err := hiero.NewFileCreateTransaction().
		// A file is not implicitly owned by anyone, even the operator
		// But we do use operator's key for this one
		SetKeys(client.GetOperatorPublicKey()).
		// Initial contents of the file
		SetContents([]byte("Hello, World")).
		// Optional memo
		SetTransactionMemo("go sdk example create_file/main.go").
		// Set max transaction fee just in case we are required to pay more
		SetMaxTransactionFee(hiero.HbarFrom(8, hiero.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}

	// Make sure the transaction went through
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving file create transaction receipt", err))
	}

	// Get and then display the file ID from the receipt
	fmt.Printf("file = %v\n", *transactionReceipt.FileID)
}
