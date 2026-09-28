package main

import (
	"fmt"
	"os"
	"time"

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
	defer client.Close()

	// Step 1: Create an account that auto-associates with any token it receives.
	accountKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating key", err))
	}
	accountCreate, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(accountKey.PublicKey()).
		SetInitialBalance(hiero.NewHbar(1)).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}
	accountReceipt, err := accountCreate.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account create receipt", err))
	}
	accountID := *accountReceipt.AccountID
	fmt.Printf("account: %v\n", accountID)

	// Step 2: Create a fungible token and an NFT collection with the operator as treasury.
	fungible := createToken(client, hiero.NewTokenCreateTransaction().
		SetTokenName("Example Fungible").
		SetTokenSymbol("EXF").
		SetTokenType(hiero.TokenTypeFungibleCommon).
		SetDecimals(2).
		SetInitialSupply(10_000).
		SetTreasuryAccountID(operatorAccountID).
		SetAdminKey(operatorKey.PublicKey()))
	nft := createToken(client, hiero.NewTokenCreateTransaction().
		SetTokenName("Example NFT").
		SetTokenSymbol("EXN").
		SetTokenType(hiero.TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(operatorAccountID).
		SetAdminKey(operatorKey.PublicKey()).
		SetSupplyKey(operatorKey.PublicKey()))

	mint, err := hiero.NewTokenMintTransaction().SetTokenID(nft).SetMetadatas([][]byte{{1}, {2}}).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error minting NFTs", err))
	}
	minted, err := mint.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting mint receipt", err))
	}

	// Step 3: Send the account 25.00 of the fungible token and both NFTs.
	transfer, err := hiero.NewTransferTransaction().
		AddTokenTransfer(fungible, operatorAccountID, -2_500).
		AddTokenTransfer(fungible, accountID, 2_500).
		AddNftTransfer(nft.Nft(minted.SerialNumbers[0]), operatorAccountID, accountID).
		AddNftTransfer(nft.Nft(minted.SerialNumbers[1]), operatorAccountID, accountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error transferring tokens", err))
	}
	if _, err = transfer.SetValidateStatus(true).GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error getting transfer receipt", err))
	}

	// Step 4: Read every token balance of the account, one page of up to 100 per Execute.
	// Retry until the mirror node has ingested the transfer.
	var all []hiero.MirrorNodeTokenBalance
	for attempt := 0; len(all) < 2; attempt++ {
		if attempt == 15 {
			panic("mirror node did not ingest the token transfer in time")
		}
		time.Sleep(time.Second)

		all = nil
		page, err := hiero.NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).Execute(client)
		for err == nil {
			all = append(all, page.Tokens...)
			if page.Next == "" {
				break
			}
			page, err = hiero.NewMirrorNodeTokenBalanceQuery().SetNextPage(page.Next).Execute(client)
		}
	}
	for _, token := range all {
		fmt.Printf("%v: balance %d, decimals %d, KYC %s, freeze %s\n",
			token.TokenID, token.Balance, token.Decimals, token.KycStatus, token.FreezeStatus)
	}

	// Step 5: Read the balance of one token. An unassociated token returns an empty page.
	one, err := hiero.NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).SetTokenID(fungible).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying one token balance", err))
	}
	if len(one.Tokens) == 0 {
		fmt.Printf("%v is not associated with %v\n", accountID, fungible)
		return
	}
	fmt.Printf("%v alone: %d (%.2f)\n", fungible, one.Tokens[0].Balance, float64(one.Tokens[0].Balance)/100)
}

func createToken(client *hiero.Client, tokenCreate *hiero.TokenCreateTransaction) hiero.TokenID {
	response, err := tokenCreate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	receipt, err := response.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token create receipt", err))
	}

	return *receipt.TokenID
}
