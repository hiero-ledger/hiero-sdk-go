package main

import (
	"encoding/hex"
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

	/*
	 * Demonstrate different account creation methods.
	 */
	createAccountWithAlias(client)
	createAccountWithAliasAndBothKeys(client, operatorKey)
	createAccountWithoutAlias(client)

}

func createAccountWithAlias(client *hiero.Client) {
	/**
	 * Step 1
	 *
	 * Create an account key and an ECSDA private alias key
	 */
	ecdsaPrivateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 2
	 *
	 * Use the `AccountCreateTransaction`
	 *   - Populate `SetECDSAKeyWithAlias(ecdsaPrivateKey)` field with the generated ECDSA private key
	 */
	frozenTxn, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.HbarFromTinybar(100)).
		SetECDSAKeyWithAlias(ecdsaPrivateKey).
		FreezeWith(client)
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 3
	 *
	 * Sign the `AccountCreateTransaction` transaction with the generated private key and execute it
	 */
	response, err := frozenTxn.Sign(ecdsaPrivateKey).Execute(client)
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 4
	 *
	 * Get the account ID of the newly created account
	 */
	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt}", err))
	}

	newAccountId := *transactionReceipt.AccountID

	/**
	 *
	 * Step 5
	 *
	 * Get the `AccountInfo` and examine the created account key and alias
	 */
	info, err := hiero.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Initial EVM address: %s is the same as %s \n", ecdsaPrivateKey.PublicKey().ToEvmAddress(), info.ContractAccountID)
}

func createAccountWithAliasAndBothKeys(client *hiero.Client, operatorKey hiero.PrivateKey) {
	/**
	 * Step 1
	 *
	 * Create an account key and an ECSDA private alias key
	 */
	ecdsaPrivateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 2
	 *
	 * Use the `AccountCreateTransaction`
	 *   - Populate `SetKeyWithAlias(key, ecdsaPrivateKey)` field with the generated ECDSA private key
	 */
	frozenTxn, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.HbarFromTinybar(100)).
		SetKeyWithAlias(operatorKey.PublicKey(), ecdsaPrivateKey).
		FreezeWith(client)
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 3
	 *
	 * Sign the `AccountCreateTransaction` transaction with both keys and execute.
	 */
	response, err := frozenTxn.Sign(ecdsaPrivateKey).Sign(operatorKey).Execute(client)
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 4
	 *
	 * Get the account ID of the newly created account
	 */
	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt}", err))
	}

	newAccountId := *transactionReceipt.AccountID

	/**
	 *
	 * Step 5
	 *
	 * Get the `AccountInfo` and examine the created account key and alias
	 */
	info, err := hiero.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Account's key: %s is the same as %s \n", info.Key.String(), operatorKey.PublicKey().String())
	fmt.Printf("Initial EVM address: %s is the same as %s \n", ecdsaPrivateKey.PublicKey().ToEvmAddress(), info.ContractAccountID)
}

func createAccountWithoutAlias(client *hiero.Client) {
	/**
	 * Step 1
	 *
	 * Create an account key and an ECSDA private alias key
	 */
	ecdsaPrivateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 2
	 *
	 * Use the `AccountCreateTransaction`
	 *   - Populate `SetKeyWithoutAlias(Key)` field with the generated ECDSA private key
	 */
	response, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.HbarFromTinybar(100)).
		SetKeyWithoutAlias(ecdsaPrivateKey).
		Execute(client)
	if err != nil {
		panic(err.Error())
	}

	/**
	 *
	 * Step 3
	 *
	 * Get the account ID of the newly created account
	 */
	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt}", err))
	}

	newAccountId := *transactionReceipt.AccountID

	/**
	 *
	 * Step 4
	 *
	 * Get the `AccountInfo` and examine the created account key and alias
	 */
	info, err := hiero.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Account's key: %s is the same as %s \n", info.Key.String(), ecdsaPrivateKey.PublicKey().String())
	hexBytes, err := hex.DecodeString(info.ContractAccountID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Account has no alias: %v \n", isLongZero(hexBytes))
}

func isLongZero(address []byte) bool {
	for i := 0; i < 12; i++ {
		if address[i] != 0 {
			return false
		}
	}
	return true
}
