package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

/**
 * @summary HIP-991 hips.hedera.com/hip/hip-991
 * @description Revenue-generating topics
 */
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
	 * Create account - alice
	 */
	fmt.Println("Creating account - alice")
	alicePrivateKey, _ := hiero.PrivateKeyGenerateEd25519()
	transactionResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(alicePrivateKey).
		SetInitialBalance(hiero.NewHbar(10)).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	receipt, err := transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	alice := receipt.AccountID
	fmt.Println("Alice account id: ", alice)

	/*
	 * Step 2:
	 * Create a topic with hbar custom fee
	 */
	fmt.Println("Creating a topic with hbar custom fee")
	customFee := hiero.NewCustomFixedFee().
		SetAmount(hiero.HbarFrom(1, hiero.HbarUnits.Hbar).AsTinybar()).
		SetFeeCollectorAccountID(operatorAccountID)

	transactionResponse, err = hiero.NewTopicCreateTransaction().
		SetTransactionMemo("go sdk example revenue-generating topic").
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFeeScheduleKey(client.GetOperatorPublicKey()).
		AddCustomFee(customFee).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating topic", err))
	}
	transactionReceipt, err := transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting topic create receipt", err))
	}

	topicID := *transactionReceipt.TopicID
	fmt.Println("Created topic with ID: ", topicID)

	/*
	 * Step 3:
	 * Submit a message to that topic, paid for by alice, specifying max custom fee amount bigger than the topic’s amount.
	 */
	accountBalanceBefore, err := hiero.NewAccountBalanceQuery().
		SetAccountID(*alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	feeCollectorBalanceBefore, err := hiero.NewAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	fmt.Println("Submitting a message as alice to the topic")
	customFeeLimit := hiero.NewCustomFeeLimit().
		SetPayerId(*alice).
		AddCustomFee(hiero.NewCustomFixedFee().
			SetAmount(hiero.HbarFrom(2, hiero.HbarUnits.Hbar).AsTinybar()))

	client.SetOperator(*alice, alicePrivateKey)
	transactionResponse, err = hiero.NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage([]byte("message")).
		AddCustomFeeLimit(customFeeLimit).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting topic", err))
	}

	_, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting topic", err))
	}
	fmt.Println("Message submitted successfully")
	client.SetOperator(operatorAccountID, operatorKey)

	/*
	 * Step 4:
	 * Verify alice was debited the fee amount and the fee collector account was credited the amount.
	 */
	accountBalanceAfter, err := hiero.NewAccountBalanceQuery().
		SetAccountID(*alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	feeCollectorBalanceAfter, err := hiero.NewAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	fmt.Println("Alice account Hbar balance before: ", accountBalanceBefore.Hbars.String())
	fmt.Println("Alice account Hbar balance after: ", accountBalanceAfter.Hbars.String())
	fmt.Println("Fee collector account Hbar balance before: ", feeCollectorBalanceBefore.Hbars.String())
	fmt.Println("Fee collector account Hbar balance after: ", feeCollectorBalanceAfter.Hbars.String())

	/*
	 * Step 5:
	 * Create a fungible token and transfer some tokens to alice
	 */
	fmt.Println("Creating a token")
	transactionResponse, err = hiero.NewTokenCreateTransaction().
		SetTokenName("revenue-generating token").
		SetTokenSymbol("RGT").
		SetDecimals(8).
		SetInitialSupply(100).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating token", err))
	}

	receipt, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating token", err))
	}

	tokenID := *receipt.TokenID
	fmt.Println("Created token with ID: ", tokenID)

	// transfer token to alice
	fmt.Println("Transferring the token to alice")
	transactionResponse, err = hiero.NewTransferTransaction().
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -1).
		AddTokenTransfer(tokenID, *alice, 1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error transferring token", err))
	}
	_, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error transferring token", err))
	}

	/*
	 * Step 6:
	 * Update the topic to have a fee of the token.
	 */
	fmt.Println("Updating the topic to have a custom fee of the token")
	customFee = hiero.NewCustomFixedFee().
		SetAmount(1).
		SetFeeCollectorAccountID(operatorAccountID).
		SetDenominatingTokenID(tokenID)

	transactionResponse, err = hiero.NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetCustomFees([]*hiero.CustomFixedFee{customFee}).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating topic", err))
	}

	_, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating topic", err))
	}

	/*
	 * Step 7:
	 * Submit another message to that topic, paid by alice, without specifying max custom fee amount.
	 */
	accountBalanceBefore, err = hiero.NewAccountBalanceQuery().
		SetAccountID(*alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	feeCollectorBalanceBefore, err = hiero.NewAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	fmt.Println("Submitting a message as alice to the topic")
	client.SetOperator(*alice, alicePrivateKey)
	transactionResponse, err = hiero.NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage([]byte("message")).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting topic", err))
	}

	_, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting topic", err))
	}
	fmt.Println("Message submitted successfully")
	client.SetOperator(operatorAccountID, operatorKey)

	feeCollectorBalanceAfter, err = hiero.NewAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	/*
	 * Step 8:
	 * Verify alice was debited the new fee amount and the fee collector account was credited the amount.
	 */
	accountBalanceAfter, err = hiero.NewAccountBalanceQuery().
		SetAccountID(*alice).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	fmt.Println("Alice account Hbar balance before: ", accountBalanceBefore.Hbars.String())
	fmt.Println("Alice account Token balance before: ", accountBalanceBefore.Tokens.Get(tokenID))
	fmt.Println("Alice account Hbar balance after: ", accountBalanceAfter.Hbars.String())
	fmt.Println("Alice account Token balance after: ", accountBalanceAfter.Tokens.Get(tokenID))
	fmt.Println("Fee collector account Token balance before: ", feeCollectorBalanceBefore.Tokens.Get(tokenID))
	fmt.Println("Fee collector account Token balance after: ", feeCollectorBalanceAfter.Tokens.Get(tokenID))

	/*
	 * Step 9:
	 * Create account - bob
	 */
	fmt.Println("Creating account - bob")
	bobPrivateKey, _ := hiero.PrivateKeyGenerateEd25519()
	transactionResponse, err = hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(bobPrivateKey).
		SetInitialBalance(hiero.NewHbar(10)).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	receipt, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v error creating account", err))
	}
	bob := receipt.AccountID
	fmt.Println("Bob account id: ", bob)

	/*
	 * Step 10:
	 * Update the topic’s fee exempt keys and add bob’s public key.
	 */
	fmt.Println("Updating the topic’s fee exempt keys and add bob’s public key")
	transactionResponse, err = hiero.NewTopicUpdateTransaction().
		SetTopicID(topicID).
		AddFeeExemptKey(bobPrivateKey.PublicKey()).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating topic", err))
	}

	_, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating topic", err))
	}

	/*
	 * Step 11:
	 * Submit another message to that topic, paid with bob, without specifying max custom fee amount.
	 */
	accountBalanceBefore, err = hiero.NewAccountBalanceQuery().
		SetAccountID(*bob).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	fmt.Println("Submitting a message as bob to the topic")
	client.SetOperator(*bob, bobPrivateKey)

	transactionResponse, err = hiero.NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage([]byte("message")).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting topic", err))
	}

	_, err = transactionResponse.SetValidateStatus(true).GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting topic", err))
	}
	fmt.Println("Message submitted successfully")
	client.SetOperator(operatorAccountID, operatorKey)

	/*
	 * Step 12:
	 * Verify bob was not debited the fee amount.
	 */
	accountBalanceAfter, err = hiero.NewAccountBalanceQuery().
		SetAccountID(*bob).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v error getting account balance", err))
	}

	fmt.Println("Bob account Hbar balance before: ", accountBalanceBefore.Hbars.String())
	fmt.Println("Bob account Token balance after: ", accountBalanceAfter.Tokens.Get(tokenID))
	fmt.Println("Bob account Hbar balance after: ", accountBalanceAfter.Hbars.String())
	fmt.Println("Bob account Token balance after: ", accountBalanceAfter.Tokens.Get(tokenID))
}
