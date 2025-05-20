package main

import (
	"fmt"
	"os"
	"strings"

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

	// Generate keys for sender and receiver
	senderKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate sender private key: %v", err))
	}
	receiverKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate receiver private key: %v", err))
	}

	// Create accounts
	senderAccountResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(senderKey.PublicKey()).
		SetInitialBalance(hiero.NewHbar(10)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("Failed to execute sender account creation transaction: %v", err))
	}

	senderAccountReceipt, err := senderAccountResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("Failed to get receipt for sender account creation: %v", err))
	}
	senderId := *senderAccountReceipt.AccountID

	receiverAccountResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(receiverKey).
		SetInitialBalance(hiero.NewHbar(1)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("Failed to execute receiver account creation transaction: %v", err))
	}

	receiverAccountReceipt, err := receiverAccountResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("Failed to get receipt for receiver account creation: %v", err))
	}
	receiverId := *receiverAccountReceipt.AccountID

	// Single node transaction example
	if err := singleNodeTransactionExample(client, senderId, receiverId, senderKey); err != nil {
		panic(fmt.Sprintf("Single node transaction example failed: %v", err))
	}

	// Multi-node multi-chunk transaction example
	if err := multiNodeFileTransactionExample(client, senderId, senderKey); err != nil {
		panic(fmt.Sprintf("Multi-node file transaction example failed: %v", err))
	}
}

func singleNodeTransactionExample(client *hiero.Client, senderId, receiverId hiero.AccountID, senderKey hiero.PrivateKey) error {
	// Step 1 - Create and prepare transfer transaction
	// Get first node from network
	network := client.GetNetwork()
	var nodeAccountId hiero.AccountID
	for _, node := range network {
		nodeAccountId = node
		break
	}

	// Create transfer transaction
	transferTx := hiero.NewTransferTransaction().
		AddHbarTransfer(senderId, hiero.NewHbar(-1)).
		AddHbarTransfer(receiverId, hiero.NewHbar(1)).
		SetNodeAccountIDs([]hiero.AccountID{nodeAccountId}).
		SetTransactionID(hiero.TransactionIDGenerate(senderId))

	transferTx, err := transferTx.FreezeWith(client)
	if err != nil {
		return fmt.Errorf("failed to freeze transfer transaction: %w", err)
	}

	// Step 2 - Get signable bytes and sign with HSM
	signableList, err := transferTx.GetSignableNodeBodyBytesList()
	if err != nil {
		return fmt.Errorf("failed to get signable bytes list: %w", err)
	}

	// Sign with HSM for each entry
	for _, signable := range signableList {
		signature := hsmSign(senderKey, signable.Body)
		transferTx.AddSignatureForMultiNodeMultiChunk(senderKey.PublicKey(), signature, signable.TransactionID, signable.NodeID)
	}

	// Step 3 - Execute transaction and get receipt
	transferResponse, err := transferTx.Execute(client)
	if err != nil {
		return fmt.Errorf("failed to execute transfer transaction: %w", err)
	}

	transferReceipt, err := transferResponse.GetReceipt(client)
	if err != nil {
		return fmt.Errorf("failed to get transfer receipt: %w", err)
	}

	fmt.Printf("Single node transaction status: %v\n", transferReceipt.Status)
	return nil
}

func multiNodeFileTransactionExample(client *hiero.Client, senderId hiero.AccountID, senderKey hiero.PrivateKey) error {
	// Step 1 - Create initial file
	// Create large content for testing
	bigContents := strings.Repeat("Lorem ipsum dolor sit amet. ", 1000)

	// Create file transaction
	fileCreateTx := hiero.NewFileCreateTransaction().
		SetKeys(senderKey.PublicKey()).
		SetContents([]byte("[e2e::FileCreateTransaction]")).
		SetMaxTransactionFee(hiero.NewHbar(5))

	fileCreateTx, err := fileCreateTx.FreezeWith(client)
	if err != nil {
		return fmt.Errorf("failed to freeze file create transaction: %w", err)
	}

	fileCreateResponse, err := fileCreateTx.Sign(senderKey).Execute(client)
	if err != nil {
		return fmt.Errorf("failed to execute file create transaction: %w", err)
	}

	fileCreateReceipt, err := fileCreateResponse.GetReceipt(client)
	if err != nil {
		return fmt.Errorf("failed to get file create receipt: %w", err)
	}

	fileId := *fileCreateReceipt.FileID
	fmt.Printf("Created file with ID: %v\n", fileId)

	// Step 2 - Prepare file append transaction
	fileAppendTx := hiero.NewFileAppendTransaction().
		SetFileID(fileId).
		SetContents([]byte(bigContents)).
		SetMaxTransactionFee(hiero.NewHbar(5)).
		SetTransactionID(hiero.TransactionIDGenerate(senderId))

	fileAppendTx, err = fileAppendTx.FreezeWith(client)
	if err != nil {
		return fmt.Errorf("failed to freeze file append transaction: %w", err)
	}

	// Step 3 - Get signable bytes and sign with HSM for each node
	fmt.Printf("Signing transaction with HSM for nodes: %v\n", fileAppendTx.GetNodeAccountIDs())

	multiNodeSignableList, err := fileAppendTx.GetSignableNodeBodyBytesList()
	if err != nil {
		return fmt.Errorf("failed to get signable bytes list: %w", err)
	}

	// Sign with HSM for each entry
	for _, signable := range multiNodeSignableList {
		signature := hsmSign(senderKey, signable.Body)
		fileAppendTx.AddSignatureForMultiNodeMultiChunk(senderKey.PublicKey(), signature, signable.TransactionID, signable.NodeID)
	}

	// Step 4 - Execute transaction and verify results
	fileAppendResponse, err := fileAppendTx.Execute(client)
	if err != nil {
		return fmt.Errorf("failed to execute file append transaction: %w", err)
	}

	fileAppendReceipt, err := fileAppendResponse.GetReceipt(client)
	if err != nil {
		return fmt.Errorf("failed to get file append receipt: %w", err)
	}

	fmt.Printf("Multi-node file append transaction status: %v\n", fileAppendReceipt.Status)

	// Step 5 - Verify file contents
	contents, err := hiero.NewFileContentsQuery().
		SetFileID(fileId).
		Execute(client)
	if err != nil {
		return fmt.Errorf("failed to query file contents: %w", err)
	}

	fmt.Printf("File content length according to FileContentsQuery: %d\n", len(contents))
	return nil
}

// hsmSign simulates signing with an HSM.
// In a real implementation, this would use actual HSM SDK logic.
func hsmSign(key hiero.PrivateKey, bodyBytes []byte) []byte {
	// This is a placeholder that simulates HSM signing
	return key.Sign(bodyBytes)
}
