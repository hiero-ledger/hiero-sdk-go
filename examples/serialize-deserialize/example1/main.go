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

	// The new account needs hbar of its own to send
	resp, err := hiero.NewAccountCreateTransaction().SetKeyWithoutAlias(newKey).SetInitialBalance(hiero.NewHbar(2)).Execute(client)
	if err != nil {
		panic(err)
	}

	receipt, err := resp.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	newAccountId := *receipt.AccountID

	// Set every field before serializing; the deserialized transaction sends the serialized body.
	bytes, err := hiero.NewTransferTransaction().
		AddHbarTransfer(newAccountId, hiero.NewHbar(-1)).
		AddHbarTransfer(operatorAccountID, hiero.NewHbar(1)).
		ToBytes()
	if err != nil {
		panic(err)
	}

	txFromBytes, err := hiero.TransactionFromBytes(bytes)
	if err != nil {
		panic(err)
	}

	// Freeze the deserialized transaction and sign it with the sender's key; the operator pays the fee.
	transaction := txFromBytes.(hiero.TransferTransaction)
	frozen, err := transaction.FreezeWith(client)
	if err != nil {
		panic(err)
	}

	executed, err := frozen.Sign(newKey).Execute(client)
	if err != nil {
		panic(err)
	}
	if _, err = executed.SetValidateStatus(true).GetReceipt(client); err != nil {
		panic(err)
	}

	// Get the `AccountInfo` on the new account and show the transfer left it 1 hbar
	info, err := hiero.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err)
	}

	fmt.Println("Balance of new account: ", info.Balance)
}
