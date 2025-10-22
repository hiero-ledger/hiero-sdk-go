package main

import (
	"fmt"
	"os"

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
	hookContractId := createHookContractId(client)
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
		SetLambdaEvmHook(*hiero.NewLambdaEvmHook().SetContractId(hookContractId))

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

// Helper function to create hook contract ID
func createHookContractId(client *hiero.Client) *hiero.ContractID {
	fileId := createBytecodeFile(client)

	response, err := hiero.NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(1000000).
		SetBytecodeFileID(fileId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating contract", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting contract create transaction receipt", err))
	}
	return receipt.ContractID
}

// Helper function to create bytecode file
func createBytecodeFile(client *hiero.Client) hiero.FileID {
	const smartContractBytecode = "6080604052348015600e575f5ffd5b506107d18061001c5f395ff3fe608060405260043610610033575f3560e01c8063124d8b301461003757806394112e2f14610067578063bd0dd0b614610097575b5f5ffd5b610051600480360381019061004c91906106f2565b6100c7565b60405161005e9190610782565b60405180910390f35b610081600480360381019061007c91906106f2565b6100d2565b60405161008e9190610782565b60405180910390f35b6100b160048036038101906100ac91906106f2565b6100dd565b6040516100be9190610782565b60405180910390f35b5f6001905092915050565b5f6001905092915050565b5f6001905092915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f60a08284031215610112576101116100f9565b5b81905092915050565b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6101658261011f565b810181811067ffffffffffffffff821117156101845761018361012f565b5b80604052505050565b5f6101966100e8565b90506101a2828261015c565b919050565b5f5ffd5b5f5ffd5b5f67ffffffffffffffff8211156101c9576101c861012f565b5b602082029050602081019050919050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610207826101de565b9050919050565b610217816101fd565b8114610221575f5ffd5b50565b5f813590506102328161020e565b92915050565b5f8160070b9050919050565b61024d81610238565b8114610257575f5ffd5b50565b5f8135905061026881610244565b92915050565b5f604082840312156102835761028261011b565b5b61028d604061018d565b90505f61029c84828501610224565b5f8301525060206102af8482850161025a565b60208301525092915050565b5f6102cd6102c8846101af565b61018d565b905080838252602082019050604084028301858111156102f0576102ef6101da565b5b835b818110156103195780610305888261026e565b8452602084019350506040810190506102f2565b5050509392505050565b5f82601f830112610337576103366101ab565b5b81356103478482602086016102bb565b91505092915050565b5f67ffffffffffffffff82111561036a5761036961012f565b5b602082029050602081019050919050565b5f67ffffffffffffffff8211156103955761039461012f565b5b602082029050602081019050919050565b5f606082840312156103bb576103ba61011b565b5b6103c5606061018d565b90505f6103d484828501610224565b5f8301525060206103e784828501610224565b60208301525060406103fb8482850161025a565b60408301525092915050565b5f6104196104148461037b565b61018d565b9050808382526020820190506060840283018581111561043c5761043b6101da565b5b835b81811015610465578061045188826103a6565b84526020840193505060608101905061043e565b5050509392505050565b5f82601f830112610483576104826101ab565b5b8135610493848260208601610407565b91505092915050565b5f606082840312156104b1576104b061011b565b5b6104bb606061018d565b90505f6104ca84828501610224565b5f83015250602082013567ffffffffffffffff8111156104ed576104ec6101a7565b5b6104f984828501610323565b602083015250604082013567ffffffffffffffff81111561051d5761051c6101a7565b5b6105298482850161046f565b60408301525092915050565b5f61054761054284610350565b61018d565b9050808382526020820190506020840283018581111561056a576105696101da565b5b835b818110156105b157803567ffffffffffffffff81111561058f5761058e6101ab565b5b80860161059c898261049c565b8552602085019450505060208101905061056c565b5050509392505050565b5f82601f8301126105cf576105ce6101ab565b5b81356105df848260208601610535565b91505092915050565b5f604082840312156105fd576105fc61011b565b5b610607604061018d565b90505f82013567ffffffffffffffff811115610626576106256101a7565b5b61063284828501610323565b5f83015250602082013567ffffffffffffffff811115610655576106546101a7565b5b610661848285016105bb565b60208301525092915050565b5f604082840312156106825761068161011b565b5b61068c604061018d565b90505f82013567ffffffffffffffff8111156106ab576106aa6101a7565b5b6106b7848285016105e8565b5f83015250602082013567ffffffffffffffff8111156106da576106d96101a7565b5b6106e6848285016105e8565b60208301525092915050565b5f5f60408385031215610708576107076100f1565b5b5f83013567ffffffffffffffff811115610725576107246100f5565b5b610731858286016100fd565b925050602083013567ffffffffffffffff811115610752576107516100f5565b5b61075e8582860161066d565b9150509250929050565b5f8115159050919050565b61077c81610768565b82525050565b5f6020820190506107955f830184610773565b9291505056fea26469706673582212207dfe7723f6d6869419b1cb0619758b439da0cf4ffd9520997c40a3946299d4dc64736f6c634300081e0033"
	response, err := hiero.NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte(smartContractBytecode)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting file create transaction receipt", err))
	}
	return *receipt.FileID
}
