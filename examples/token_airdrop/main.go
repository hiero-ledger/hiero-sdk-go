package main

import (
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

/**
 * @summary HIP-904 https://hips.hedera.com/hip/hip-904
 * @description Airdrop fungible and non fungible tokens to an account
 */

//nolint:gocyclo
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

	fmt.Println("Example Start!")

	/*
	 * Step 1:
	 * Create 4 accounts
	 */
	privateKey1, _ := hiero.PrivateKeyGenerateEd25519()
	accountCreateResp, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(privateKey1).
		SetInitialBalance(hiero.NewHbar(10)).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	receipt, err := accountCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	alice := receipt.AccountID

	privateKey2, _ := hiero.PrivateKeyGenerateEd25519()
	accountCreateResp, err = hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(privateKey2).
		SetInitialBalance(hiero.NewHbar(10)).
		SetMaxAutomaticTokenAssociations(1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	receipt, err = accountCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	bob := receipt.AccountID

	privateKey3, _ := hiero.PrivateKeyGenerateEd25519()
	accountCreateResp, err = hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(privateKey3).
		SetInitialBalance(hiero.NewHbar(10)).
		SetMaxAutomaticTokenAssociations(0).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	receipt, err = accountCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	carol := receipt.AccountID

	treasuryKey, _ := hiero.PrivateKeyGenerateEd25519()
	accountCreateResp, err = hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(treasuryKey).
		SetInitialBalance(hiero.NewHbar(10)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}

	receipt, err = accountCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	treasury := receipt.AccountID

	/*
	 * Step 2:
	 * Create FT and NFT and mint
	 */
	tokenCreateTxn, _ := hiero.NewTokenCreateTransaction().
		SetTokenName("Fungible Token").
		SetTokenSymbol("TFT").
		SetTokenMemo("Example memo").
		SetDecimals(3).
		SetInitialSupply(1000).
		SetMaxSupply(1000).
		SetTreasuryAccountID(*treasury).
		SetSupplyType(hiero.TokenSupplyTypeFinite).
		SetAdminKey(operatorKey).
		SetFreezeKey(operatorKey).
		SetSupplyKey(operatorKey).
		SetMetadataKey(operatorKey).
		SetPauseKey(operatorKey).
		FreezeWith(client)

	tokenCreateResp, err := tokenCreateTxn.Sign(treasuryKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating token", err))
	}

	receipt, err = tokenCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	tokenID := receipt.TokenID

	nftCreateTransaction, _ := hiero.NewTokenCreateTransaction().
		SetTokenName("Example NFT").
		SetTokenSymbol("ENFT").
		SetTokenType(hiero.TokenTypeNonFungibleUnique).
		SetDecimals(0).
		SetInitialSupply(0).
		SetMaxSupply(10).
		SetTreasuryAccountID(*treasury).
		SetSupplyType(hiero.TokenSupplyTypeFinite).
		SetAdminKey(operatorKey).
		SetFreezeKey(operatorKey).
		SetSupplyKey(operatorKey).
		FreezeWith(client)
	tokenCreateResp, err = nftCreateTransaction.Sign(treasuryKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	nftCreateReceipt, err := tokenCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating token", err))
	}
	nftID := *nftCreateReceipt.TokenID
	var initialMetadataList = [][]byte{{2, 1}, {1, 2}, {1, 5}}

	mintTransaction, _ := hiero.NewTokenMintTransaction().
		SetTokenID(nftID).
		SetMetadatas(initialMetadataList).
		FreezeWith(client)

	mintTransactionSubmit, err := mintTransaction.Sign(operatorKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error minting NFT", err))
	}
	receipt, err = mintTransactionSubmit.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error minting NFT", err))
	}

	/*
	 * Step 3:
	 * Airdrop fungible tokens to all 3 accounts
	 */

	airdropTx, _ := hiero.NewTokenAirdropTransaction().
		AddTokenTransfer(*tokenID, *alice, 100).
		AddTokenTransfer(*tokenID, *bob, 100).
		AddTokenTransfer(*tokenID, *carol, 100).
		AddTokenTransfer(*tokenID, *treasury, -300).
		FreezeWith(client)
	airdropResponse, err := airdropTx.Sign(treasuryKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error airdropping tokens", err))
	}
	record, err := airdropResponse.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error airdropping tokens", err))
	}

	fmt.Println("Pending airdrops length: ", len(record.PendingAirdropRecords))
	fmt.Println("Pending airdrops: ", record.PendingAirdropRecords[0].String())

	/*
	 * Step 5:
	 * Query to verify alice and bob received the airdrops and carol did not
	 */
	fmt.Println("Alice ft balance after airdrop: ", tokenBalance(client, *alice, *tokenID, is(100)))
	fmt.Println("Bob ft balance after airdrop: ", tokenBalance(client, *bob, *tokenID, is(100)))
	fmt.Println("Carol ft balance after airdrop: ", tokenBalance(client, *carol, *tokenID, is(0)))

	/*
	 * Step 6:
	 * Claim the airdrop for carol
	 */
	fmt.Println("Claiming ft with carol")

	claimTx, _ := hiero.NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(client)

	_, err = claimTx.Sign(privateKey3).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error claiming tokens", err))
	}
	fmt.Println("Carol ft balance after claim: ", tokenBalance(client, *carol, *tokenID, is(100)))

	/*
	 * Step 7:
	 * Airdrop the NFTs to all three accounts
	 */

	airdropTx, _ = hiero.NewTokenAirdropTransaction().
		AddNftTransfer(nftID.Nft(1), *treasury, *alice).
		AddNftTransfer(nftID.Nft(2), *treasury, *bob).
		AddNftTransfer(nftID.Nft(3), *treasury, *carol).
		FreezeWith(client)
	airdropResponse, err = airdropTx.Sign(treasuryKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error airdropping tokens", err))
	}
	record, err = airdropResponse.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error airdropping tokens", err))
	}

	/*
	 * Step 8:
	 * Get the transaction record and verify two pending airdrops (for bob & carol)
	 */

	fmt.Println("Pending airdrops length: ", len(record.PendingAirdropRecords))
	fmt.Println("Pending airdrops for Bob: ", record.PendingAirdropRecords[0].String())
	fmt.Println("Pending airdrops for Carol: ", record.PendingAirdropRecords[1].String())

	/*
	 * Step 9:
	 * Query to verify alice received the airdrop and bob and carol did not
	 */

	fmt.Println("Alice nft balance after airdrop: ", tokenBalance(client, *alice, nftID, is(1)))
	fmt.Println("Bob nft balance after airdrop: ", tokenBalance(client, *bob, nftID, is(0)))
	fmt.Println("Carol nft balance after airdrop: ", tokenBalance(client, *carol, nftID, is(0)))

	/*
	 * Step 10:
	 * Claim the airdrop for bob
	 */
	fmt.Println("Claiming nft with Bob")
	claimTx, _ = hiero.NewTokenClaimAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[0].GetPendingAirdropId()).
		FreezeWith(client)

	claimResp, err := claimTx.Sign(privateKey2).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error claiming tokens", err))
	}
	_, err = claimResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error claiming tokens", err))
	}
	fmt.Println("Bob nft balance after claim: ", tokenBalance(client, *bob, nftID, is(1)))

	/*
	 * Step 11:
	 * Cancel the airdrop for carol
	 */
	fmt.Println("Canceling nft with Carol")
	cancelTx, _ := hiero.NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(record.PendingAirdropRecords[1].GetPendingAirdropId()).
		FreezeWith(client)

	cancelResp, err := cancelTx.Sign(treasuryKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error canceling tokens", err))
	}
	_, err = cancelResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error canceling tokens", err))
	}
	fmt.Println("Carol nft balance after cancel: ", tokenBalance(client, *carol, nftID, is(0)))

	/*
	 * Step 12:
	 * Reject the NFT for bob
	 */
	fmt.Println("Rejecting nft with Bob")

	rejectTxn, _ := hiero.NewTokenRejectTransaction().
		AddNftID(nftID.Nft(2)).
		SetOwnerID(*bob).
		FreezeWith(client)

	rejectResp, err := rejectTxn.Sign(privateKey2).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error rejecting tokens", err))
	}
	_, err = rejectResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error rejecting tokens", err))
	}

	/*
	 * Step 13:
	 * Query to verify bob no longer has the NFT
	 */
	fmt.Println("Bob nft balance after reject: ", tokenBalance(client, *bob, nftID, is(0)))

	/*
	 * Step 13:
	 * Query to verify the NFT was returned to the Treasury
	 */
	fmt.Println("Treasury nft balance after reject: ", tokenBalance(client, *treasury, nftID, is(2)))

	/*
	 * Step 14:
	 * Reject the fungible tokens for Carol
	 */

	rejectTxn, _ = hiero.NewTokenRejectTransaction().
		AddTokenID(*tokenID).
		SetOwnerID(*carol).
		FreezeWith(client)

	rejectResp, err = rejectTxn.Sign(privateKey3).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error rejecting tokens", err))
	}
	_, err = rejectResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error rejecting tokens", err))
	}

	/*
	 * Step 14:
	 * Query to verify carol no longer has the fungible tokens
	 */
	fmt.Println("Carol ft balance after reject: ", tokenBalance(client, *carol, *tokenID, is(0)))

	/*
	 * Step 15:
	 * Query to verify Treasury received the rejected fungible tokens
	 */
	fmt.Println("Treasury ft balance after reject: ", tokenBalance(client, *treasury, *tokenID, is(800)))

	/*
	 * Clean up:
	 */
	client.Close()

	fmt.Println("Example Complete!")
}

// tokenBalance polls the mirror node until ready accepts accountID's balance of tokenID; no relationship reads as 0.
func tokenBalance(client *hiero.Client, accountID hiero.AccountID, tokenID hiero.TokenID, ready func(uint64) bool) uint64 {
	for attempt := 1; ; attempt++ {
		page, err := hiero.NewMirrorNodeTokenBalanceQuery().SetAccountID(accountID).SetTokenID(tokenID).Execute(client)
		if err == nil {
			var balance uint64
			if len(page.Tokens) == 1 {
				balance = page.Tokens[0].Balance
			}
			if ready(balance) {
				return balance
			}
		}
		if attempt == 20 {
			panic(fmt.Sprintf("mirror node did not show the expected %v balance for %v (last error: %v)", tokenID, accountID, err))
		}
		time.Sleep(time.Second)
	}
}

func is(want uint64) func(uint64) bool {
	return func(balance uint64) bool { return balance == want }
}
