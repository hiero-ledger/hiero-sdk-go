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

	keys := make([]hiero.PrivateKey, 3)
	pubKeys := make([]hiero.PublicKey, 3)

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	// Loop to generate keys for the KeyList
	for i := range keys {
		newKey, err := hiero.GeneratePrivateKey()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey}", err))
		}

		fmt.Printf("Key %v:\n", i)
		fmt.Printf("private = %v\n", newKey)
		fmt.Printf("public = %v\n", newKey.PublicKey())

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	// A threshold key with a threshold of 2 and length of 3 requires
	// at least 2 of the 3 keys to sign anything modifying the account
	keyList := hiero.NewKeyList().
		AddAllPublicKeys(pubKeys)

	// We are using all of these keys, so the scheduled transaction doesn't automatically go through
	// It works perfectly fine with just one key
	createResponse, err := hiero.NewAccountCreateTransaction().
		// The key that must sign each transfer out of the account. If receiverSigRequired is true, then
		// it must also sign any transfer into the account.
		SetKeyWithoutAlias(keyList).
		SetInitialBalance(hiero.NewHbar(10)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing create account transaction", err))
	}

	// Make sure the transaction succeeded
	transactionReceipt, err := createResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	// Pre-generating transaction id for the scheduled transaction so we can track it
	transactionID := hiero.TransactionIDGenerate(client.GetOperatorAccountID())

	println("transactionId for scheduled transaction = ", transactionID.String())

	// Not really necessary as its client.GetOperatorAccountID()
	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)

	// Creating a non frozen transaction for the scheduled transaction
	// In this case its TransferTransaction
	transferTx := hiero.NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(newAccountID, hiero.HbarFrom(-1, hiero.HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.HbarFrom(1, hiero.HbarUnits.Hbar))

	// Scheduling it, this gives us hiero.ScheduleCreateTransaction
	scheduled, err := transferTx.Schedule()
	if err != nil {
		panic(fmt.Sprintf("%v : error scheduling Transfer Transaction", err))
	}

	// Executing the scheduled transaction
	scheduleResponse, err := scheduled.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule create", err))
	}

	// Make sure it executed successfully
	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule create receipt", err))
	}

	// Taking out the schedule ID
	scheduleID := *scheduleReceipt.ScheduleID

	// Using the schedule ID to get the schedule transaction info, which contains the whole scheduled transaction
	info, err := hiero.NewScheduleInfoQuery().
		SetNodeAccountIDs([]hiero.AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule info", err))
	}

	// Taking out the TransferTransaction from earlier
	transfer, err := info.GetScheduledTransaction()
	if err != nil {
		panic(fmt.Sprintf("%v : error getting transaction from schedule info", err))
	}

	// Converting it from interface to hiero.TransferTransaction() and retrieving the amount of transfers
	// to check if we have the right one, and that it's not empty
	var transfers map[hiero.AccountID]hiero.Hbar
	if tx, ok := transfer.(*hiero.TransferTransaction); ok {
		transfers = tx.GetHbarTransfers()
	}

	if len(transfers) != 2 {
		println("more transfers than expected")
	}

	// Checking if the Hbar values are correct
	if transfers[newAccountID].AsTinybar() != -hiero.NewHbar(1).AsTinybar() {
		println("transfer for ", newAccountID.String(), " is not whats is expected")
	}

	// Checking if the Hbar values are correct
	if transfers[client.GetOperatorAccountID()].AsTinybar() != hiero.NewHbar(1).AsTinybar() {
		println("transfer for ", client.GetOperatorAccountID().String(), " is not whats is expected")
	}

	println("sending schedule sign transaction")

	// Creating a scheduled sign transaction, we have to sign with all of the keys in the KeyList
	signTransaction, err := hiero.NewScheduleSignTransaction().
		SetNodeAccountIDs([]hiero.AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing sign transaction", err))
	}

	// Signing the scheduled transaction
	signTransaction.Sign(keys[0])
	signTransaction.Sign(keys[1])
	signTransaction.Sign(keys[2])

	resp, err := signTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule sign transaction", err))
	}

	// Getting the receipt to make sure the signing executed properly
	_, err = resp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule sign receipt", err))
	}

	// Making sure the scheduled transaction executed properly with schedule info query
	info, err = hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		SetNodeAccountIDs([]hiero.AccountID{createResponse.NodeID}).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving schedule info after signing", err))
	}

	// Checking if the scheduled transaction was executed and signed, and retrieving the signatories
	if !info.ExecutedAt.IsZero() {
		println("Singing success, signed at: ", info.ExecutedAt.String())
		println("Signatories: ", info.Signatories.String())
	} else {
		panic("Signing failed")
	}
}
