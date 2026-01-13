package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2/examples/hooks/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// HIP-1195 example for creating and updating an account with a hook
func main() {
	var client *hiero.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}
	defer client.Close()

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

	fmt.Println("Creating hook contract")
	hookContractId := utils.CreateHookContractId(client)
	fmt.Printf("Hook contract created: %v\n", hookContractId)

	fmt.Println("Creating account with hook")
	createAccountWithHook(client, hookContractId)

	fmt.Println("Updating account to add hook")
	addHookToAccount(client, hookContractId)

	fmt.Println("Updating account to delete hook")
	deleteHookFromAccount(client, hookContractId)

	fmt.Println("Example completed")
}

func createAccountWithHook(client *hiero.Client, hookContractId *hiero.ContractID) {
	fmt.Println("Creating account with hook...")

	// Create hook detail
	hookDetail := hiero.NewHookCreationDetails().
		SetExtensionPoint(hiero.ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetEvmHook(*hiero.NewEvmHook().SetContractId(hookContractId))

	// Generate account private key
	fmt.Println("Generating account private key...")
	accountPrivateKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating account private key", err))
	}
	accountPublicKey := accountPrivateKey.PublicKey()

	// Create account create transaction with hook
	fmt.Println("Executing account create transaction with hook...")
	response, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(accountPublicKey).
		AddHook(*hookDetail).
		SetMaxTransactionFee(hiero.NewHbar(10)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account create transaction receipt", err))
	}
	fmt.Printf("Account created with hook: %v\n", receipt.AccountID)
}

func addHookToAccount(client *hiero.Client, hookContractId *hiero.ContractID) {
	fmt.Println("Creating account without hook first...")

	// Generate account private key
	fmt.Println("Generating account private key...")
	accountPrivateKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating account private key", err))
	}
	accountPublicKey := accountPrivateKey.PublicKey()

	// Create account create transaction without hook
	fmt.Println("Executing account create transaction...")
	response, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(accountPublicKey).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account create transaction receipt", err))
	}
	accountId := *receipt.AccountID
	fmt.Printf("Account created: %v\n", accountId)

	// Create hook detail
	fmt.Println("  Creating hook details...")
	hookDetail := hiero.NewHookCreationDetails().
		SetExtensionPoint(hiero.ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetEvmHook(*hiero.NewEvmHook().SetContractId(hookContractId))

	// Create account update transaction to add hook
	fmt.Println("Adding hook to account via update transaction...")
	frozenUpdateTransaction, err := hiero.NewAccountUpdateTransaction().
		SetAccountID(accountId).
		AddHookToCreate(*hookDetail).
		SetMaxTransactionFee(hiero.NewHbar(15)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account update transaction", err))
	}
	response, err = frozenUpdateTransaction.Sign(accountPrivateKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account update transaction", err))
	}
	_, err = response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account update transaction receipt", err))
	}
	fmt.Println("Hook added to account successfully!")
}

func deleteHookFromAccount(client *hiero.Client, hookContractId *hiero.ContractID) {
	fmt.Println("  Creating account with hook first...")

	// Create hook detail
	hookDetail := hiero.NewHookCreationDetails().
		SetExtensionPoint(hiero.ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetEvmHook(*hiero.NewEvmHook().SetContractId(hookContractId))

	// Generate account private key
	fmt.Println("Generating account private key...")
	accountPrivateKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating account private key", err))
	}
	accountPublicKey := accountPrivateKey.PublicKey()

	// Create account create transaction with hook
	fmt.Println("Executing account create transaction with hook...")
	response, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(accountPublicKey).
		AddHook(*hookDetail).
		SetMaxTransactionFee(hiero.NewHbar(10)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account create transaction receipt", err))
	}
	accountId := *receipt.AccountID
	fmt.Printf("  Account created with hook: %v\n", accountId)

	// Create account update transaction to delete hook
	fmt.Println("Deleting hook from account via update transaction...")
	frozenUpdateTransaction, err := hiero.NewAccountUpdateTransaction().
		SetAccountID(accountId).
		AddHookToDelete(1).
		SetMaxTransactionFee(hiero.NewHbar(15)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating account", err))
	}
	response, err = frozenUpdateTransaction.Sign(accountPrivateKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account update transaction", err))
	}
	_, err = response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account update transaction receipt", err))
	}
	fmt.Println("Hook deleted from account successfully!")
}
