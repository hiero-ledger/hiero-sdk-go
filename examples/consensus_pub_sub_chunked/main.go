package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

//go:embed large_message.txt
var bigContents string

// How to send a large message to a private HCS topic and how to subscribe to
// the topic to receive it.
//
// The message exceeds the per-transaction size limit, so the SDK automatically
// splits it into multiple chunks under one logical operation. The subscriber
// receives the reassembled message via the mirror node.
func main() {
	fmt.Println("Consensus Service Submit Large Message And Subscribe Example Start!")

	client, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	operatorAccountID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	client.SetOperator(operatorAccountID, operatorKey)

	// Step 1: Generate ED25519 key pair (Submit Key for the topic).
	fmt.Println("Generating ED25519 key pair...")
	submitKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating submit PrivateKey", err))
	}

	// Step 2: Create the topic with admin + submit keys.
	fmt.Println("Creating new topic...")
	topicCreateResponse, err := hiero.NewTopicCreateTransaction().
		SetTopicMemo("hedera-sdk-go/ConsensusPubSubChunkedExample").
		SetAdminKey(client.GetOperatorPublicKey()).
		SetSubmitKey(submitKey.PublicKey()).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating topic", err))
	}

	topicCreateReceipt, err := topicCreateResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic creation receipt", err))
	}
	topicID := *topicCreateReceipt.TopicID
	fmt.Printf("Created new topic with ID: %v\n", topicID)

	// Step 3: Wait for the new topic to propagate to mirror nodes.
	fmt.Println("Wait 5 seconds (to ensure data propagated to mirror nodes) ...")
	time.Sleep(5 * time.Second)

	// Step 4: Subscribe to the topic. The latch fires after the first
	// (reassembled) message arrives.
	fmt.Println("Setting up a mirror client...")
	done := make(chan struct{}, 1)
	_, err = hiero.NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message hiero.TopicMessage) {
			fmt.Printf("Topic message received! | Time: %v | Sequence No.: %d | Size: %d bytes.\n",
				message.ConsensusTimestamp, message.SequenceNumber, len(message.Contents))
			select {
			case done <- struct{}{}:
			default:
			}
		})
	if err != nil {
		panic(fmt.Sprintf("%v : error subscribing", err))
	}

	// Step 5: Build, sign with operator, serialize, deserialize, sign with
	// submit key, execute. The bytes round-trip mirrors the pattern where
	// the operator and the submit-key holder are different parties.
	topicMessageSubmitTx, err := hiero.NewTopicMessageSubmitTransaction().
		// Default is 10 chunks; increase so a large message will fit.
		SetMaxChunks(15).
		SetTopicID(topicID).
		SetMessage([]byte(bigContents)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing topic message submit transaction", err))
	}

	// The operator signs first (charged the transaction fee).
	operatorSignedTx, err := topicMessageSubmitTx.SignWithOperator(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error signing with operator", err))
	}

	// Serialize so the bytes can be signed "somewhere else" by the submit key.
	transactionBytes, err := operatorSignedTx.ToBytes()
	if err != nil {
		panic(fmt.Sprintf("%v : error serializing topic message submit transaction", err))
	}
	parsedTxIface, err := hiero.TransactionFromBytes(transactionBytes)
	if err != nil {
		panic(fmt.Sprintf("%v : error deserializing topic message submit transaction", err))
	}
	parsedTx, ok := parsedTxIface.(hiero.TopicMessageSubmitTransaction)
	if !ok {
		panic("did not receive TopicMessageSubmitTransaction back from signed bytes")
	}

	fmt.Printf("Preparing to submit a message to the created topic (size of the message: %d bytes)...\n", len(bigContents))

	// Sign with the submit key (required because the topic has a submitKey).
	signedTx := parsedTx.Sign(submitKey)

	// Submit the chunked message and wait for the receipt.
	submitResponse, err := signedTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing topic message submit transaction", err))
	}
	if _, err := submitResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic message submit receipt", err))
	}

	// Wait up to 60s for the reassembled message to arrive via the mirror.
	select {
	case <-done:
		// got it
	case <-time.After(60 * time.Second):
		panic("Large topic message was not received! (Fail)")
	}

	// Cleanup: delete the topic.
	deleteResponse, err := hiero.NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetMaxTransactionFee(hiero.NewHbar(5)).
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

	fmt.Println("Consensus Service Submit Large Message And Subscribe Example Complete!")
}
