package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

const TotalMessages = 5

// How to operate with a private HCS topic.
//
// Create a new HCS topic with a single ECDSA Submit Key,
// publish a number of messages to the topic signed by the Submit Key
// and subscribe to the topic (no key required).
func main() {
	fmt.Println("Consensus Service Submit Message To The Private Topic And Subscribe Example Start!")

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
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

	// Step 1: Generate ECDSA key pair (Submit Key for the topic).
	fmt.Println("Generating ECDSA key pair...")
	submitKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Step 2: Create the topic with admin + submit keys.
	fmt.Println("Creating new HCS topic...")
	transactionResponse, err := hiero.NewTopicCreateTransaction().
		SetTopicMemo("HCS topic with Submit Key").
		SetAdminKey(client.GetOperatorPublicKey()).
		// Access control for TopicSubmitMessage. Submitters must sign with this key.
		SetSubmitKey(submitKey.PublicKey()).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating topic", err))
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic create transaction receipt", err))
	}

	topicID := *transactionReceipt.TopicID
	fmt.Printf("Created topic with ID: %v and public ECDSA submit key: %v\n", topicID, submitKey)

	// Step 3: Wait for the new topic to propagate to mirror nodes.
	fmt.Println("Wait 5 seconds (to ensure data propagated to mirror nodes) ...")
	time.Sleep(5 * time.Second)

	// Step 4: Subscribe to the topic. Each received message decrements the latch.
	fmt.Println("Setting up a mirror client...")
	done := make(chan struct{}, TotalMessages)
	_, err = hiero.NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message hiero.TopicMessage) {
			fmt.Printf("Topic message received! | Time: %v | Content: %s\n",
				message.ConsensusTimestamp, string(message.Contents))
			done <- struct{}{}
		})
	if err != nil {
		panic(fmt.Sprintf("%v : error subscribing", err))
	}

	// Step 5: Publish messages, signing each with the submit key.
	for i := 0; i < TotalMessages; i++ {
		message := "random message " + strconv.Itoa(rand.Int()) //nolint:gosec
		fmt.Printf("Publishing message to the topic: %s\n", message)

		submitTx, err := hiero.NewTopicMessageSubmitTransaction().
			SetTopicID(topicID).
			SetMessage([]byte(message)).
			FreezeWith(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error freezing topic message submit transaction", err))
		}

		// The transaction is implicitly signed by the operator (payer).
		// The topic has a submitKey requirement — sign with that key too.
		submitTx.Sign(submitKey)

		submitTxResponse, err := submitTx.Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error executing topic message submit transaction", err))
		}

		if _, err := submitTxResponse.GetReceipt(client); err != nil {
			panic(fmt.Sprintf("%v : error retrieving topic message submit transaction receipt", err))
		}

		time.Sleep(2 * time.Second)
	}

	// Wait up to 60s for all messages to arrive via the subscription, fail otherwise.
	timeout := time.After(60 * time.Second)
	for i := 0; i < TotalMessages; i++ {
		select {
		case <-done:
			// got one
		case <-timeout:
			panic("Not all topic messages were received! (Fail)")
		}
	}

	// Cleanup: delete created topic.
	deleteResponse, err := hiero.NewTopicDeleteTransaction().
		SetTopicID(topicID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing topic delete transaction", err))
	}
	if _, err := deleteResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic delete receipt", err))
	}

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}

	fmt.Println("Consensus Service Submit Message To The Private Topic And Subscribe Example Complete!")
}
