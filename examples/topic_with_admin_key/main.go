package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to create and manage an HCS topic using a threshold key as the adminKey
// and going through a key rotation to a new set of keys.
//
// Create a new HCS topic with a 2-of-3 threshold key for the Admin Key and
// update the HCS topic to a 3-of-4 threshold key for the adminKey.
func main() {
	fmt.Println("Topic With Admin (Threshold) Key Example Start!")

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

	client.SetOperator(operatorAccountID, operatorKey)

	// Step 1: Generate the initial admin key pairs (3 keys, 2-of-3 threshold).
	fmt.Println("Generating ECDSA key pairs...")
	initialAdminKeys := make([]hiero.PrivateKey, 3)
	for i := range initialAdminKeys {
		key, err := hiero.PrivateKeyGenerateEcdsa()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey", err))
		}
		initialAdminKeys[i] = key
	}

	// Step 2: Build the threshold key.
	fmt.Println("Creating a Key List (threshold key)...")
	keyList := hiero.KeyListWithThreshold(2)
	for _, key := range initialAdminKeys {
		keyList.Add(key.PublicKey())
	}
	fmt.Printf("Created a Key List: %v\n", keyList)

	// Step 3: Create the topic create transaction.
	fmt.Println("Creating topic create transaction...")
	topicCreateTx, err := hiero.NewTopicCreateTransaction().
		SetTopicMemo("demo topic").
		SetAdminKey(keyList).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing topic create transaction", err))
	}

	// Step 4: Sign with 2 of 3 admin keys.
	for i := 0; i < 2; i++ {
		fmt.Printf("Signing topic create transaction with key %v\n", initialAdminKeys[i])
		topicCreateTx.Sign(initialAdminKeys[i])
	}

	// Step 5: Execute.
	topicCreateResponse, err := topicCreateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating topic", err))
	}

	topicCreateReceipt, err := topicCreateResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic creation receipt", err))
	}
	topicID := *topicCreateReceipt.TopicID
	fmt.Printf("Created new topic (%v) with 2-of-3 threshold key as admin key.\n", topicID)

	// Step 6: Generate the new admin key pairs (4 keys, 3-of-4 threshold).
	fmt.Println("Generating new ECDSA key pairs...")
	newAdminKeys := make([]hiero.PrivateKey, 4)
	for i := range newAdminKeys {
		key, err := hiero.PrivateKeyGenerateEcdsa()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey", err))
		}
		newAdminKeys[i] = key
	}

	// Step 7: Build the new threshold key.
	fmt.Println("Creating new Key List (threshold key)...")
	newKeyList := hiero.KeyListWithThreshold(3)
	for _, key := range newAdminKeys {
		newKeyList.Add(key.PublicKey())
	}
	fmt.Printf("Created new Key List: %v\n", newKeyList)

	// Step 8: Create the topic update transaction.
	fmt.Println("Creating topic update transaction...")
	topicUpdateTx, err := hiero.NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetTopicMemo("This topic will be updated").
		SetAdminKey(newKeyList).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing topic update transaction", err))
	}

	// Step 9a: Sign with 2 of the OLD admin keys (authorize the change).
	for i := 0; i < 2; i++ {
		fmt.Printf("Signing topic update transaction with initial admin key %v\n", initialAdminKeys[i])
		topicUpdateTx.Sign(initialAdminKeys[i])
	}

	// Step 9b: Sign with 3 of the NEW admin keys (prove possession).
	for i := 0; i < 3; i++ {
		fmt.Printf("Signing topic update transaction with new admin key %v\n", newAdminKeys[i])
		topicUpdateTx.Sign(newAdminKeys[i])
	}

	// Step 10: Execute.
	topicUpdateResponse, err := topicUpdateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error updating topic", err))
	}
	if _, err := topicUpdateResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic update receipt", err))
	}
	fmt.Printf("Updated topic (%v) with 3-of-4 threshold key as admin key.\n", topicID)

	// Cleanup: delete the topic, signed with 3 of 4 new admin keys.
	topicDeleteTx, err := hiero.NewTopicDeleteTransaction().
		SetTopicID(topicID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing topic delete transaction", err))
	}
	for i := 0; i < 3; i++ {
		topicDeleteTx.Sign(newAdminKeys[i])
	}
	deleteResponse, err := topicDeleteTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing topic delete transaction", err))
	}
	if _, err := deleteResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error retrieving topic delete receipt", err))
	}

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}

	fmt.Println("Topic With Admin (Threshold) Key Example Complete!")
}
