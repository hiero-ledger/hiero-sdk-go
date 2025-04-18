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

	// Generating key for the new account
	key1, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Generating the key to update to
	key2, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Creating new account
	accountTxResponse, err := hiero.NewAccountCreateTransaction().
		// The key that must sign each transfer out of the account. If receiverSigRequired is true, then
		// it must also sign any transfer into the account.
		// Using the public key for this, but a PrivateKey or a KeyList can also be used
		SetKeyWithoutAlias(key1.PublicKey()).
		SetInitialBalance(hiero.ZeroHbar).
		SetTransactionID(hiero.TransactionIDGenerate(client.GetOperatorAccountID())).
		SetTransactionMemo("sdk example create_account__with_manual_signing/main.go").
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}

	println("transaction ID:", accountTxResponse.TransactionID.String())

	// Get the receipt making sure transaction worked
	accountTxReceipt, err := accountTxResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account creation receipt", err))
	}

	// Retrieve the account ID out of the Receipt
	accountID := *accountTxReceipt.AccountID
	println("account =", accountID.String())
	println("key =", key1.PublicKey().String())
	println(":: update public key of account", accountID.String())
	println("set key =", key2.PublicKey().String())

	// Updating the account with the new key
	accountUpdateTx, err := hiero.NewAccountUpdateTransaction().
		SetAccountID(accountID).
		// The new key
		SetKey(key2.PublicKey()).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account update transaction", err))
	}

	// Have to sign with both keys, the initial key first
	accountUpdateTx.Sign(key1)
	accountUpdateTx.Sign(key2)

	// Executing the account update transaction
	accountUpdateTxResponse, err := accountUpdateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating account", err))
	}

	println("transaction ID:", accountUpdateTxResponse.TransactionID.String())

	// Make sure the transaction went through
	_, err = accountUpdateTxResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting the transaction receipt", err))
	}

	println(":: getAccount and check our current key")
	info, err := hiero.NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account info query", err))
	}

	// This should be same as key2
	println("key =", info.Key.String())
}
