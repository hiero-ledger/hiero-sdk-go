package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2/examples/hooks/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// HIP-1195 example for creating and updating a contract with a hook
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

	fmt.Println("Creating contract with hook")
	createContractWithHook(client, hookContractId)

	fmt.Println("Updating contract to add hook")
	addHookToContract(client, hookContractId)

	fmt.Println("Updating contract to delete hook")
	deleteHookFromContract(client, hookContractId)

	fmt.Println("Example completed")
}

func createContractWithHook(client *hiero.Client, hookContractId *hiero.ContractID) {
	fmt.Println("Creating contract with hook...")

	// Create hook detail
	hookDetail := hiero.NewHookCreationDetails().
		SetExtensionPoint(hiero.ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetEvmHook(*hiero.NewEvmHook().SetContractId(hookContractId))

	// Create bytecode file
	fmt.Println("Creating bytecode file...")
	fileId := utils.CreateBytecodeFile(client)

	// Create contract create transaction with hook
	fmt.Println("Executing contract create transaction with hook...")
	response, err := hiero.NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(1000000).
		SetBytecodeFileID(fileId).
		AddHook(*hookDetail).
		SetMaxTransactionFee(hiero.NewHbar(20)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating contract", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting contract create transaction receipt", err))
	}
	fmt.Printf("Contract created successfully: %v\n", receipt.ContractID)
}

func addHookToContract(client *hiero.Client, hookContractId *hiero.ContractID) {
	fmt.Println("Creating contract without hook first...")

	// Create bytecode file
	fmt.Println("Creating bytecode file...")
	fileId := utils.CreateBytecodeFile(client)

	// Create contract create transaction without hook
	fmt.Println("Executing contract create transaction...")
	response, err := hiero.NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(1000000).
		SetBytecodeFileID(fileId).
		SetMaxTransactionFee(hiero.NewHbar(20)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating contract", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting contract create transaction receipt", err))
	}
	contractId := *receipt.ContractID
	fmt.Printf("Contract created: %v\n", contractId)

	// Create hook detail
	fmt.Println("Creating hook details...")
	hookDetail := hiero.NewHookCreationDetails().
		SetExtensionPoint(hiero.ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetEvmHook(*hiero.NewEvmHook().SetContractId(hookContractId))

	// Create contract update transaction to add hook
	fmt.Println("Adding hook to contract via update transaction...")
	frozenUpdateTransaction, err := hiero.NewContractUpdateTransaction().
		SetContractID(contractId).
		AddHookToCreate(*hookDetail).
		SetMaxTransactionFee(hiero.NewHbar(25)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing contract update transaction", err))
	}
	response, err = frozenUpdateTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing contract update transaction", err))
	}
	_, err = response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting contract update transaction receipt", err))
	}
	fmt.Println("Hook added to contract successfully!")
}

func deleteHookFromContract(client *hiero.Client, hookContractId *hiero.ContractID) {
	fmt.Println("Creating contract with hook first...")

	// Create hook detail
	hookDetail := hiero.NewHookCreationDetails().
		SetExtensionPoint(hiero.ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetEvmHook(*hiero.NewEvmHook().SetContractId(hookContractId))

	// Create bytecode file
	fmt.Println("Creating bytecode file...")
	fileId := utils.CreateBytecodeFile(client)

	// Create contract create transaction with hook
	fmt.Println("Executing contract create transaction with hook...")
	response, err := hiero.NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(1000000).
		SetBytecodeFileID(fileId).
		AddHook(*hookDetail).
		SetMaxTransactionFee(hiero.NewHbar(20)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating contract", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting contract create transaction receipt", err))
	}
	contractId := *receipt.ContractID
	fmt.Printf("Contract created with hook: %v\n", contractId)

	// Create contract update transaction to delete hook
	fmt.Println("Deleting hook from contract via update transaction...")
	frozenUpdateTransaction, err := hiero.NewContractUpdateTransaction().
		SetContractID(contractId).
		AddHookToDelete(1).
		SetMaxTransactionFee(hiero.NewHbar(25)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating contract", err))
	}
	response, err = frozenUpdateTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing contract update transaction", err))
	}
	_, err = response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting contract update transaction receipt", err))
	}
	fmt.Println("Hook deleted from contract successfully!")
}
