package main

import (
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

	fmt.Println("Example Start!")

	/*
	 * Step 1:
	 * Create account - alice
	 */
	fmt.Println("Create account - alice")
	alicePrivateKey, _ := hiero.PrivateKeyGenerateEd25519()
	accountCreateResp, err := hiero.NewAccountCreateTransaction().
		SetKey(alicePrivateKey).
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
	fmt.Println("Alice account id: ", alice)

	/*
	 * Step 2:
	 * Create a topic with hbar custom fee
	 */
	fmt.Println("Create a topic with hbar custom fee")
	customFee := hiero.NewCustomFixedFee().
		SetAmount(hiero.HbarFrom(1, hiero.HbarUnits.Hbar).AsTinybar()).
		SetFeeCollectorAccountID(client.GetOperatorAccountID())

	transactionResponse, err := hiero.NewTopicCreateTransaction().
		SetTransactionMemo("go sdk example revenue-generating topic").
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFeeScheduleKey(client.GetOperatorPublicKey()).
		AddCustomFee(customFee).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating topic", err))
	}
	transactionReceipt, err := transactionResponse.GetReceipt(client)
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

	fmt.Println("Submitting a message as alice to the topic")
	customFeeLimit := hiero.NewCustomFeeLimit().
		SetPayerId(*alice).
		AddCustomFee(hiero.NewCustomFixedFee().SetAmount(hiero.HbarFrom(2, hiero.HbarUnits.Hbar).AsTinybar()))

	client.SetOperator(*alice, alicePrivateKey)
	transactionResponse, err = hiero.NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage([]byte("message")).
		AddCustomFeeLimit(customFeeLimit).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error submitting topic", err))
	}

	_, err = transactionResponse.GetReceipt(client)
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

	fmt.Println("Alice account balance before: ", accountBalanceBefore.Hbars.String())
	fmt.Println("Alice account balance after: ", accountBalanceAfter.Hbars.String())
}
