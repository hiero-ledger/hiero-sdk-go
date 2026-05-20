package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to create a file under a non-operator key, delete it (signing with that
// key), and verify the deletion via FileInfoQuery.
func main() {
	fmt.Println("Delete File Example Start!")

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

	// Step 1: Generate the key to be used with the new file (separate from operator).
	fmt.Println("Generating ECDSA key pair...")
	newKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Step 2: Create a file with newKey as its only key — operator pays the tx fee
	// but newKey is required to delete or modify the file.
	fmt.Println("Creating a file to delete:")
	freezeTransaction, err := hiero.NewFileCreateTransaction().
		SetContents([]byte("The quick brown fox jumps over the lazy dog")).
		SetKeys(newKey.PublicKey()).
		SetTransactionMemo("go sdk example delete_file/main.go").
		SetMaxTransactionFee(hiero.HbarFrom(8, hiero.HbarUnits.Hbar)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing transaction", err))
	}

	transactionResponse, err := freezeTransaction.Sign(newKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}

	receipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving file creation receipt", err))
	}

	newFileID := *receipt.FileID
	fmt.Printf("file = %v\n", newFileID)

	// Step 3: Delete the file — must be signed with newKey since it's the file's key.
	fmt.Println("Deleting created file...")
	deleteTransaction, err := hiero.NewFileDeleteTransaction().
		SetFileID(newFileID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing file delete transaction", err))
	}

	deleteTransactionResponse, err := deleteTransaction.Sign(newKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing file delete transaction", err))
	}

	deleteTransactionReceipt, err := deleteTransactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving file deletion receipt", err))
	}
	fmt.Printf("file delete transaction status: %v\n", deleteTransactionReceipt.Status)

	// Step 4: Querying file info on a deleted file returns isDeleted=true.
	fileInfo, err := hiero.NewFileInfoQuery().
		SetFileID(newFileID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing file info query", err))
	}
	fmt.Printf("file %v was deleted: %v\n", newFileID, fileInfo.IsDeleted)

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}

	fmt.Println("Delete File Example Complete!")
}
