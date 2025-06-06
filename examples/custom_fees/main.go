package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

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

	// Generate new key to be used with new account
	aliceKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Create three accounts, Alice, Bob, and Charlie.  Alice will be the treasury for our example token.
	// Fees only apply in transactions not involving the treasury, so we need two other accounts.

	aliceAccountCreate, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.NewHbar(10)).
		SetKeyWithoutAlias(aliceKey).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account create for alice", err))
	}

	aliceAccountCreate.Sign(aliceKey)
	resp, err := aliceAccountCreate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create for alice", err))
	}

	receipt, err := resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for alice account create", err))
	}

	var aliceId hiero.AccountID
	if receipt.AccountID != nil {
		aliceId = *receipt.AccountID
	} else {
		panic("Receipt didn't return alice's ID")
	}

	bobKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	bobAccountCreate, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.NewHbar(10)).
		SetKeyWithoutAlias(bobKey).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account create for bob", err))
	}

	bobAccountCreate.Sign(bobKey)
	resp, err = bobAccountCreate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create for bob", err))
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for bob account create", err))
	}

	var bobId hiero.AccountID
	if receipt.AccountID != nil {
		bobId = *receipt.AccountID
	} else {
		panic("Receipt didn't return bob's ID")
	}

	charlieKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	charlieAccountCreate, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.NewHbar(10)).
		SetKeyWithoutAlias(charlieKey).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account create for charlie", err))
	}

	charlieAccountCreate.Sign(aliceKey)
	resp, err = charlieAccountCreate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing account create for charlie", err))
	}

	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for charlie account create", err))
	}

	var charlieId hiero.AccountID
	if receipt.AccountID != nil {
		charlieId = *receipt.AccountID
	} else {
		panic("Receipt didn't return charlie's ID")
	}

	println("Alice:", aliceId.String())
	println("Bob:", bobId.String())
	println("Charlie:", charlieId.String())

	// Let's start with a custom fee list of 1 fixed fee.  A custom fee list can be a list of up to
	// 10 custom fees, where each fee is a fixed fee or a fractional fee.
	// This fixed fee will mean that every time Bob transfers any number of tokens to Charlie,
	// Alice will collect 1 Hbar from each account involved in the transaction who is SENDING
	// the Token (in this case, Bob).

	customHbarFee := hiero.NewCustomFixedFee().
		SetHbarAmount(hiero.NewHbar(1)).
		SetFeeCollectorAccountID(aliceId)

	// In this example the fee is in Hbar, but you can charge a fixed fee in a token if you'd like.
	// EG, you can make it so that each time an account transfers Foo tokens,
	// they must pay a fee in Bar tokens to the fee collecting account.
	// To charge a fixed fee in tokens, instead of calling setHbarAmount(), call
	// setDenominatingTokenId(tokenForFee) and setAmount(tokenFeeAmount).

	// Setting the feeScheduleKey to Alice's key will enable Alice to change the custom
	// fees list on this token later using the TokenFeeScheduleUpdateTransaction.
	// We will create an initial supply of 100 of these tokens.

	tokenCreate, err := hiero.NewTokenCreateTransaction().
		// Token name and symbol are only things required to create a token
		SetTokenName("Example Token").
		SetTokenSymbol("EX").
		// The key which can perform update/delete operations on the token. If empty, the token can be
		// perceived as immutable (not being able to be updated/deleted)
		SetAdminKey(aliceKey).
		// The key which can change the supply of a token. The key is used to sign Token Mint/Burn
		// operations
		SetSupplyKey(aliceKey).
		// The key which can change the token's custom fee schedule; must sign a TokenFeeScheduleUpdate
		// transaction
		SetFeeScheduleKey(aliceKey).
		// The account which will act as a treasury for the token. This account
		// will receive the specified initial supply or the newly minted NFTs
		SetTreasuryAccountID(aliceId).
		// The custom fees to be assessed during a CryptoTransfer that transfers units of this token
		SetCustomFees([]hiero.Fee{*customHbarFee}).
		// Specifies the initial supply of tokens to be put in circulation. The
		// initial supply is sent to the Treasury Account.
		SetInitialSupply(100).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing token create transaction", err))
	}

	// Sign with alice's key before executing
	tokenCreate.Sign(aliceKey)
	resp, err = tokenCreate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing token create transaction", err))
	}

	// Get receipt to make sure the transaction passed through
	receipt, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for token create transaction", err))
	}

	// Get the token out of the receipt
	var tokenId hiero.TokenID
	if receipt.TokenID != nil {
		tokenId = *receipt.TokenID
	} else {
		println("Token ID missing in the receipt")
	}

	println("TokenID:", tokenId.String())

	tokenInfo1, err := hiero.NewTokenInfoQuery().
		SetTokenID(tokenId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing TokenInfoQuery", err))
	}

	fmt.Println("Custom Fees according to TokenInfoQuery:")
	for _, fee := range tokenInfo1.CustomFees {
		if fixedFee, ok := fee.(hiero.CustomFixedFee); ok {
			fmt.Println(fixedFee.String())
		}
	}

	// We must associate the token with Bob and Charlie before they can trade in it.
	tokenAssociate, err := hiero.NewTokenAssociateTransaction().
		// Account to associate token with
		SetAccountID(bobId).
		// The token to associate with
		SetTokenIDs(tokenId).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing token associate transaction for bob", err))
	}

	// Signing with bob's key
	tokenAssociate.Sign(bobKey)
	resp, err = tokenAssociate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing token associate transaction for bob", err))
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for token associate transaction for bob", err))
	}

	// Associating charlie's account with the token
	tokenAssociate, err = hiero.NewTokenAssociateTransaction().
		// Account to associate token with
		SetAccountID(charlieId).
		// The token to associate with
		SetTokenIDs(tokenId).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing token associate transaction for charlie", err))
	}

	// Signing with charlie's key
	tokenAssociate.Sign(charlieKey)
	resp, err = tokenAssociate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing token associate transaction for charlie", err))
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for token associate transaction for charlie", err))
	}

	// Give all 100 tokens to Bob
	transferTransaction, err := hiero.NewTransferTransaction().
		// The 100 tokens being given to bob
		AddTokenTransfer(tokenId, bobId, 100).
		// Have to take the 100 tokens from alice by negating the 100
		AddTokenTransfer(tokenId, aliceId, -100).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing token transfer transaction for alice", err))
	}

	// Have to sign with alice's key as we are taking alice's tokens
	transferTransaction.Sign(aliceKey)
	resp, err = transferTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing token transfer transaction for alice", err))
	}

	// Make sure the transaction passed through
	_, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for token transfer transaction for alice", err))
	}

	// Check alice's balance before Bob transfers 20 tokens to Charlie
	// This is a free query
	aliceBalance1, err := hiero.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance 1 for alice", err))
	}

	println("Alice's Hbar balance before Bob transfers 20 tokens to Charlie:", aliceBalance1.Hbars.String())

	// Transfer 20 tokens from bob to charlie
	transferTransaction, err = hiero.NewTransferTransaction().
		// Taking away 20 tokens from bob
		AddTokenTransfer(tokenId, bobId, -20).
		// Giving 20 to charlie
		AddTokenTransfer(tokenId, charlieId, 20).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing token transfer transaction for bob", err))
	}

	// As we are taking from bob, bob has to sign this.
	transferTransaction.Sign(bobKey)
	resp, err = transferTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing token transfer transaction for bob", err))
	}

	// Getting the record to show the assessed custom fees
	record1, err := resp.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting record for token transfer transaction for bob", err))
	}

	// Query to check alice's balance
	aliceBalance2, err := hiero.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance 2 for alice", err))
	}

	println("Alice's Hbar balance after Bob transfers 20 tokens to Charlie:", aliceBalance2.Hbars.String())
	println("Assessed fees according to transaction record:")
	for _, k := range record1.AssessedCustomFees {
		println(k.String())
	}

	// Let's use the TokenUpdateFeeScheduleTransaction with Alice's key to change the custom fees on our token.
	// TokenUpdateFeeScheduleTransaction will replace the list of fees that apply to the token with
	// an entirely new list.  Let's charge a 10% fractional fee.  This means that when Bob attempts to transfer
	// 20 tokens to Charlie, 10% of the tokens he attempts to transfer (2 in this case) will be transferred to
	// Alice instead.

	// Fractional fees default to FeeAssessmentMethod.INCLUSIVE, which is the behavior described above.
	// If you set the assessment method to EXCLUSIVE, then when Bob attempts to transfer 20 tokens to Charlie,
	// Charlie will receive all 20 tokens, and Bob will be charged an _additional_ 10% fee which
	// will be transferred to Alice.

	customFractionalFee := hiero.NewCustomFractionalFee().
		SetNumerator(1).
		SetDenominator(10).
		// The minimum amount to assess
		SetMin(1).
		// The maximum amount to assess (zero implies no maximum)
		SetMax(10).
		// The account to receive the custom fee
		SetFeeCollectorAccountID(aliceId)

	tokenFeeUpdate, err := hiero.NewTokenFeeScheduleUpdateTransaction().
		// The token for which the custom fee will be updated
		SetTokenID(tokenId).
		// The updated custom fee
		SetCustomFees([]hiero.Fee{*customFractionalFee}).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing token fee update", err))
	}

	// As the token is owned by alice and all keys are set to alice's key we have to sign with that
	tokenFeeUpdate.Sign(aliceKey)
	resp, err = tokenFeeUpdate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing token fee update", err))
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt for token fee update", err))
	}

	// Get token info, we can check if the custom fee is updated
	tokenInfo2, err := hiero.NewTokenInfoQuery().
		SetTokenID(tokenId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting token info 2", err))
	}

	println("Custom Fees according to TokenInfoQuery:")
	for _, fee := range tokenInfo2.CustomFees {
		if fractionalFee, ok := fee.(hiero.CustomFractionalFee); ok {
			fmt.Println(fractionalFee.String())
		}
	}

	// Another account balance query to check alice's token balance before Bob transfers 20 tokens to Charlie
	aliceBalance3, err := hiero.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance 3 for alice", err))
	}

	println("Alice's token balance before Bob transfers 20 tokens to Charlie:", aliceBalance3.Tokens.Get(tokenId))

	// Once again transfer 20 tokens from bob to charlie
	transferTransaction, err = hiero.NewTransferTransaction().
		AddTokenTransfer(tokenId, bobId, -20).
		AddTokenTransfer(tokenId, charlieId, 20).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing token transfer transaction for bob", err))
	}

	// Bob's is losing 20 tokens again. so he has to sign this transfer
	transferTransaction.Sign(bobKey)
	resp, err = transferTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing token transfer transaction for bob", err))
	}

	record2, err := resp.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting record for token transfer transaction for bob", err))
	}

	// Checking alice's token balance again
	aliceBalance4, err := hiero.NewAccountBalanceQuery().
		SetAccountID(aliceId).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting account balance 2 for alice", err))
	}

	println("Alice's token balance after Bob transfers 20 tokens to Charlie:", aliceBalance4.Tokens.Get(tokenId))
	println("Token transfers according to transaction record:")
	for token, transfer := range record2.TokenTransfers {
		tokenT := ""
		for _, t := range transfer {
			tokenT = tokenT + " " + t.String()
		}
		println("for token", token.String()+":", tokenT)
	}
	println("Assessed fees according to transaction record:")
	for _, k := range record2.AssessedCustomFees {
		println(k.String())
	}

	// Clean up

	tokenDelete, _ := hiero.NewTokenDeleteTransaction().
		SetTokenID(tokenId).
		FreezeWith(client)

	tokenDelete.Sign(aliceKey)
	resp, _ = tokenDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	accDelete, _ := hiero.NewAccountDeleteTransaction().
		SetTransactionID(hiero.TransactionIDGenerate(charlieId)).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetAccountID(charlieId).
		FreezeWith(client)

	accDelete.Sign(charlieKey)
	resp, _ = accDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	accDelete, _ = hiero.NewAccountDeleteTransaction().
		SetTransactionID(hiero.TransactionIDGenerate(bobId)).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetAccountID(bobId).
		FreezeWith(client)

	accDelete.Sign(bobKey)
	resp, _ = accDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	accDelete, _ = hiero.NewAccountDeleteTransaction().
		SetTransactionID(hiero.TransactionIDGenerate(aliceId)).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetAccountID(aliceId).
		FreezeWith(client)

	accDelete.Sign(aliceKey)
	resp, _ = accDelete.Execute(client)
	_, _ = resp.GetReceipt(client)

	_ = client.Close()
}
