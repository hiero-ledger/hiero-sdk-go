package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2/examples/hooks/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// HIP-1195 example for HookStoreTransaction
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
	accountId, accountPrivateKey := createAccountWithHook(client, hookContractId)
	fmt.Printf("Account created: %v\n", accountId)
	fmt.Printf("Account private key: %v\n", accountPrivateKey.String())

	fmt.Println("Updating hook storage")
	updateHookStorage(client, accountId, accountPrivateKey)

	fmt.Println("Example completed")
}

func createAccountWithHook(client *hiero.Client, hookContractId *hiero.ContractID) (*hiero.AccountID, hiero.PrivateKey) {
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
	return receipt.AccountID, accountPrivateKey
}

func updateHookStorage(client *hiero.Client, accountId *hiero.AccountID, accountPrivateKey hiero.PrivateKey) {
	fmt.Println("Updating hook storage...")

	// Create hook entity ID and hook ID
	fmt.Println("Creating hook entity ID and hook ID...")
	entityId := hiero.NewHookEntityIdWithAccountId(*accountId)
	hookId := hiero.NewHookId(*entityId, 1)

	// Create storage slot update
	fmt.Println("Creating storage slot update...")
	storageSlot := hiero.NewEvmHookStorageSlot().
		SetKey([]byte{0x01, 0x02}).
		SetValue([]byte{0x03, 0x04})

	// Create HookStore transaction
	fmt.Println("Creating and freezing HookStore transaction...")
	frozenTxn, err := hiero.NewHookStoreTransaction().
		SetHookId(*hookId).
		AddStorageUpdate(storageSlot).
		SetMaxTransactionFee(hiero.NewHbar(5)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing HookStore transaction", err))
	}

	fmt.Println("Signing and executing HookStore transaction...")
	response, err := frozenTxn.Sign(accountPrivateKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing HookStore transaction", err))
	}

	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting HookStore transaction receipt", err))
	}

	fmt.Printf("Hook storage updated successfully! Status: %v\n", receipt.Status)
}
