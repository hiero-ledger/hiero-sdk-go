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

	// Create a small file to start with
	newFileResponse, err := hiero.NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello from hiero.")).
		SetMemo("go file update chunked example").
		SetMaxTransactionFee(hiero.NewHbar(2)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}

	receipt, err := newFileResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving file creation receipt", err))
	}
	fileID := *receipt.FileID

	// Contents larger than the ~6 KiB single-transaction limit. ExecuteAll transparently overwrites
	// the file with the first chunk and appends the remainder, so this one call replaces the whole
	// file with content that could never fit in a single FileUpdateTransaction.
	//
	// Note: each chunk is a separate network transaction charged its own fee, and each is signed
	// with the operator only. If the file needs additional keys, or you want explicit control of the
	// FileID / per-chunk fees / error recovery, use the manual FileUpdate + FileAppend two-step.
	bigContents := strings.Repeat("The quick brown fox jumps over the lazy dog. ", 400)

	responses, err := hiero.NewFileUpdateTransaction().
		SetNodeAccountIDs([]hiero.AccountID{newFileResponse.NodeID}).
		SetFileID(fileID).
		SetContents([]byte(bigContents)).
		SetMaxTransactionFee(hiero.NewHbar(5)).
		ExecuteAll(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing chunked file update", err))
	}

	fmt.Printf("Chunked update ran as %d transactions (1 update + %d appends)\n", len(responses), len(responses)-1)

	// Confirm the whole file was replaced
	info, err := hiero.NewFileInfoQuery().
		SetNodeAccountIDs([]hiero.AccountID{newFileResponse.NodeID}).
		SetFileID(fileID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing file info query", err))
	}

	fmt.Printf("Uploaded %d bytes; file size according to FileInfoQuery: %d\n", len(bigContents), info.Size)
}
