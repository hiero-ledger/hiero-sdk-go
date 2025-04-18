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

	// ## Example
	// Create an ED25519 admin private key and ECSDA private key
	// Extract the ECDSA public key public key
	// Extract the Ethereum public address
	// Use the `AccountCreateTransaction` and populate `setAlias(evmAddress)` field with the Ethereum public address and the `setReceiverSignatureRequired` to `true`
	// Sign the `AccountCreateTransaction` transaction with both the new private key and the admin key
	// Get the `AccountInfo` on the new account and show that the account has contractAccountId

	// Create an ED25519 admin private key and ECSDA private key
	adminKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(err.Error())
	}

	privateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err.Error())
	}

	// Use the `AccountCreateTransaction` and set the EVM address field to the Ethereum public address
	frozenTxn, err := hiero.NewAccountCreateTransaction().SetReceiverSignatureRequired(true).SetInitialBalance(hiero.HbarFromTinybar(100)).
		SetKeyWithAlias(adminKey, privateKey).FreezeWith(client)
	if err != nil {
		panic(err.Error())
	}

	response, err := frozenTxn.Sign(adminKey).Sign(privateKey).Execute(client)
	if err != nil {
		panic(err.Error())
	}

	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt}", err))
	}

	newAccountId := *transactionReceipt.AccountID

	// Get the `AccountInfo` on the new account and show that the account has contractAccountId
	info, err := hiero.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err.Error())
	}
	// Verify account is created with the provided EVM address
	fmt.Println(info.ContractAccountID == privateKey.PublicKey().ToEvmAddress())
	// Verify the account Id is the same from the create account transaction
	fmt.Println(info.AccountID.String() == newAccountId.String())
}
