package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How identical scheduled transactions submitted by different parties are
// merged into a single schedule.
//
// Three sub-accounts each act as a separate party. Each submits the same
// scheduled transfer. The first submission creates the schedule; the next
// two return Status.IDENTICAL_SCHEDULE_ALREADY_CREATED with the same
// scheduleId, and then add their signatures via ScheduleSignTransaction.
// The schedule executes automatically when 2-of-3 signatures accumulate.
func main() {
	fmt.Println("Schedule Identical Transaction Example Start!")

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

	// Step 1: Generate 3 ECDSA key pairs and create 3 sub-accounts with 1 ℏ
	// each, each backed by its own Client (simulating 3 separate parties).
	fmt.Println("Generating 3 ECDSA key pairs and creating 3 sub-accounts...")
	keys := make([]hiero.PrivateKey, 3)
	pubKeys := make([]hiero.PublicKey, 3)
	accounts := make([]hiero.AccountID, 3)
	subClients := make([]*hiero.Client, 3)

	for i := 0; i < 3; i++ {
		newKey, err := hiero.PrivateKeyGenerateEcdsa()
		if err != nil {
			panic(fmt.Sprintf("%v : error generating PrivateKey", err))
		}
		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()

		createResponse, err := hiero.NewAccountCreateTransaction().
			SetKeyWithoutAlias(newKey).
			SetInitialBalance(hiero.NewHbar(1)).
			Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error creating sub-account", err))
		}
		receipt, err := createResponse.GetReceipt(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error getting sub-account receipt", err))
		}
		accounts[i] = *receipt.AccountID

		subClient, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
		if err != nil {
			panic(fmt.Sprintf("%v : error creating sub-client", err))
		}
		subClient.SetOperator(accounts[i], newKey)
		subClients[i] = subClient

		fmt.Printf("  Sub-account %d: %v\n", i+1, accounts[i])
	}

	// Step 2: Build a 2-of-3 threshold KeyList from the 3 public keys.
	fmt.Println("Building a 2-of-3 threshold KeyList...")
	thresholdKey := hiero.KeyListWithThreshold(2).AddAllPublicKeys(pubKeys)

	// Step 3: Create the threshold account with the KeyList as its key
	// and 10 ℏ initial balance — this is the source of the scheduled transfer.
	fmt.Println("Creating threshold account with 10 ℏ initial balance...")
	thresholdResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(thresholdKey).
		SetInitialBalance(hiero.NewHbar(10)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating threshold account", err))
	}
	thresholdReceipt, err := thresholdResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting threshold receipt", err))
	}
	thresholdAccount := *thresholdReceipt.AccountID
	fmt.Printf("Threshold account: %v\n", thresholdAccount)

	// Step 4: Each sub-client submits the SAME ScheduleCreateTransaction.
	// First creates the schedule. Subsequent identical submissions return
	// IDENTICAL_SCHEDULE_ALREADY_CREATED with the same scheduleId. After
	// detecting identical, each sub-client adds its signature via
	// ScheduleSignTransaction. The schedule executes automatically when the
	// 2-of-3 threshold is reached.
	var scheduleID hiero.ScheduleID
	for i := 0; i < 3; i++ {
		fmt.Printf("Sub-client %d submitting identical ScheduleCreate...\n", i+1)

		tx := hiero.NewTransferTransaction().
			AddHbarTransfer(thresholdAccount, hiero.NewHbar(3).Negated())
		for _, acc := range accounts {
			tx.AddHbarTransfer(acc, hiero.NewHbar(1))
		}

		frozenTx, err := tx.FreezeWith(subClients[i])
		if err != nil {
			panic(fmt.Sprintf("%v : error freezing transfer", err))
		}
		signedTx, err := frozenTx.SignWithOperator(subClients[i])
		if err != nil {
			panic(fmt.Sprintf("%v : error signing transfer", err))
		}

		scheduleTx, err := hiero.NewScheduleCreateTransaction().
			SetScheduledTransaction(signedTx)
		if err != nil {
			panic(fmt.Sprintf("%v : error building ScheduleCreate", err))
		}
		scheduleTx.SetPayerAccountID(thresholdAccount)

		createResponse, err := scheduleTx.Execute(subClients[i])
		if err != nil {
			panic(fmt.Sprintf("%v : error executing ScheduleCreate", err))
		}

		// TransactionReceiptQuery.Execute returns the receipt regardless of
		// status, so IDENTICAL_SCHEDULE_ALREADY_CREATED comes through cleanly.
		createReceipt, err := hiero.NewTransactionReceiptQuery().
			SetTransactionID(createResponse.TransactionID).
			Execute(subClients[i])
		if err != nil {
			panic(fmt.Sprintf("%v : error reading ScheduleCreate receipt", err))
		}

		if i == 0 {
			scheduleID = *createReceipt.ScheduleID
			fmt.Printf("  Schedule created with ID: %v\n", scheduleID)
		} else {
			if createReceipt.ScheduleID.String() != scheduleID.String() {
				panic(fmt.Sprintf("Schedule ID mismatch! Expected %v, got %v", scheduleID, createReceipt.ScheduleID))
			}
			fmt.Printf("  Status: %v, scheduleId: %v\n", createReceipt.Status, *createReceipt.ScheduleID)

			signTx, err := hiero.NewScheduleSignTransaction().
				SetScheduleID(scheduleID).
				FreezeWith(subClients[i])
			if err != nil {
				panic(fmt.Sprintf("%v : error freezing ScheduleSign", err))
			}
			signResponse, err := signTx.Execute(subClients[i])
			if err != nil {
				panic(fmt.Sprintf("%v : error executing ScheduleSign", err))
			}
			signReceipt, err := hiero.NewTransactionReceiptQuery().
				SetTransactionID(signResponse.TransactionID).
				Execute(subClients[i])
			if err != nil {
				panic(fmt.Sprintf("%v : error reading ScheduleSign receipt", err))
			}
			fmt.Printf("  ScheduleSign status: %v\n", signReceipt.Status)
		}
	}

	// Step 5: Query the schedule's final state. After 2 of the 3 sub-clients
	// signed, the threshold was reached and the transfer executed.
	info, err := hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying schedule info", err))
	}
	executedAt := "not yet"
	if info.ExecutedAt != nil && !info.ExecutedAt.IsZero() {
		executedAt = info.ExecutedAt.String()
	}
	signers := 0
	if info.Signatories != nil {
		signers = len(info.Signatories.GetKeys())
	}
	fmt.Printf("Final ScheduleInfo: executed=%s signers=%d\n", executedAt, signers)

	// Cleanup: delete the threshold account (requires 2-of-3 from threshold keys).
	deleteThreshold, err := hiero.NewAccountDeleteTransaction().
		SetAccountID(thresholdAccount).
		SetTransferAccountID(operatorAccountID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing threshold delete", err))
	}
	deleteThreshold.Sign(keys[0]).Sign(keys[1])
	if _, err := deleteThreshold.Execute(client); err != nil {
		panic(fmt.Sprintf("%v : error executing threshold delete", err))
	}

	// Cleanup: delete each sub-account, signing with its own key.
	for i := 0; i < 3; i++ {
		subDelete, err := hiero.NewAccountDeleteTransaction().
			SetAccountID(accounts[i]).
			SetTransferAccountID(operatorAccountID).
			FreezeWith(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error freezing sub-account delete", err))
		}
		if _, err := subDelete.Sign(keys[i]).Execute(client); err != nil {
			panic(fmt.Sprintf("%v : error executing sub-account delete", err))
		}
		if err := subClients[i].Close(); err != nil {
			panic(fmt.Sprintf("%v : error closing sub-client", err))
		}
	}

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}
	fmt.Println("Schedule Identical Transaction Example Complete!")
}
