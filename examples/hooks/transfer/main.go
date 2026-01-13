package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2/examples/hooks/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// HIP-1195 example for Transfer Transaction with Hook
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

	// Create sender account
	fmt.Println("Creating sender account with hook")
	senderAccountId, senderPrivateKey := createAccountWithHook(client, hookContractId)
	fmt.Printf("Sender account created: %v\n", senderAccountId)

	// Create receiver account
	fmt.Println("Creating receiver account")
	receiverPrivateKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating receiver private key", err))
	}
	response, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(receiverPrivateKey.PublicKey()).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating receiver account", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receiver account create transaction receipt", err))
	}
	receiverAccountId := *receipt.AccountID
	fmt.Printf("Receiver account created: %v\n", receiverAccountId)

	// Transfer HBAR with hook from sender to receiver
	fmt.Printf("Transferring HBAR from sender account %v to receiver account %v\n", senderAccountId, receiverAccountId)
	transferHbarWithPreHook(client, senderAccountId, &receiverAccountId)
	fmt.Println("Transfer completed")

	// Transfer Fungible Token with hook from sender to receiver
	fmt.Printf("Transferring Fungible Token from sender account %v to receiver account %v\n", senderAccountId, receiverAccountId)
	transferFungibleTokenWithPreHook(client, senderAccountId, &receiverAccountId, senderPrivateKey)
	fmt.Println("Transfer completed")

	// Transfer NFT with pre hook from sender to receiver
	fmt.Printf("Transferring NFT from sender account %v to receiver account %v\n", senderAccountId, receiverAccountId)
	transferNftWithPreHook(client, senderAccountId, &receiverAccountId, senderPrivateKey)
	fmt.Println("Transfer completed")
}

func transferHbarWithPreHook(client *hiero.Client, senderAccountId *hiero.AccountID, receiverAccountId *hiero.AccountID) {
	fmt.Println("Creating HBAR transfer with hook...")
	hookCall := hiero.NewFungibleHookCall(1, *hiero.NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), hiero.PRE_HOOK)
	fmt.Println("Executing HBAR transfer transaction...")
	response, err := hiero.NewTransferTransaction().
		AddHbarTransferWithHook(*senderAccountId, hiero.NewHbar(-1), *hookCall).
		AddHbarTransfer(*receiverAccountId, hiero.NewHbar(1)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating transfer transaction", err))
	}
	_, err = response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting transfer transaction receipt", err))
	}
	fmt.Println("HBAR transfer with hook completed successfully!")
}

func transferFungibleTokenWithPreHook(client *hiero.Client, senderAccountId *hiero.AccountID, receiverAccountId *hiero.AccountID, senderPrivateKey hiero.PrivateKey) {
	fmt.Println("Creating fungible token...")
	frozenTxn, err := hiero.NewTokenCreateTransaction().
		SetTokenName("Test Token").
		SetTokenSymbol("TT").
		SetDecimals(1).
		SetInitialSupply(10).
		SetTreasuryAccountID(*senderAccountId).
		SetAdminKey(senderPrivateKey.PublicKey()).
		SetSupplyKey(senderPrivateKey.PublicKey()).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	response, err := frozenTxn.Sign(senderPrivateKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token create transaction receipt", err))
	}
	tokenId := *receipt.TokenID
	fmt.Printf("Fungible token created: %v\n", tokenId)

	fmt.Println("Creating fungible token transfer with hook...")
	hookCall := hiero.NewFungibleHookCall(1, *hiero.NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), hiero.PRE_HOOK)
	fmt.Println("Executing fungible token transfer transaction...")
	response, err = hiero.NewTransferTransaction().
		AddTokenTransferWithHook(tokenId, *senderAccountId, -1, *hookCall).
		AddTokenTransfer(tokenId, *receiverAccountId, 1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating transfer transaction", err))
	}
	_, err = response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting transfer transaction receipt", err))
	}
	fmt.Println("Fungible token transfer with hook completed successfully!")
}

func transferNftWithPreHook(client *hiero.Client, senderAccountId *hiero.AccountID, receiverAccountId *hiero.AccountID, senderPrivateKey hiero.PrivateKey) {
	fmt.Println("Creating NFT token...")
	frozenTxn, err := hiero.NewTokenCreateTransaction().
		SetTokenName("Test Token").
		SetTokenSymbol("TT").
		SetTreasuryAccountID(*senderAccountId).
		SetAdminKey(senderPrivateKey.PublicKey()).
		SetSupplyKey(senderPrivateKey.PublicKey()).
		SetTokenType(hiero.TokenTypeNonFungibleUnique).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	response, err := frozenTxn.Sign(senderPrivateKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token create transaction receipt", err))
	}
	tokenId := *receipt.TokenID
	fmt.Printf("NFT token created: %v\n", tokenId)

	fmt.Println("  Minting NFTs...")
	frozenMint, err := hiero.NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadatas([][]byte{{1}, {2}, {3}}).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating mint transaction", err))
	}
	response, err = frozenMint.Sign(senderPrivateKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating mint transaction", err))
	}
	receipt, err = response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting mint transaction receipt", err))
	}
	serialNumbers := receipt.SerialNumbers
	fmt.Printf("  NFTs minted with serial numbers: %v\n", serialNumbers)

	// Transfer NFT with pre hook of the sender
	fmt.Println("Creating NFT transfer with hook...")
	hookCall := hiero.NewNftHookCall(1, *hiero.NewEvmHookCall().SetData([]byte{}).SetGasLimit(25000), hiero.PRE_HOOK_SENDER)
	fmt.Println("Executing NFT transfer transaction...")
	response, err = hiero.NewTransferTransaction().
		AddNftTransferWitHook(tokenId.Nft(serialNumbers[0]), *senderAccountId, *receiverAccountId, hookCall, nil).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating transfer transaction", err))
	}
	_, err = response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting transfer transaction receipt", err))
	}
	fmt.Println("NFT transfer with hook completed successfully!")
}

// Helper function to create account with hook
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
		SetInitialBalance(hiero.NewHbar(10)).
		SetMaxTransactionFee(hiero.NewHbar(10)).
		SetMaxAutomaticTokenAssociations(-1).
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
