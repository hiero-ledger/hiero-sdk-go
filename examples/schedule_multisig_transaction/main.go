package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to schedule a transaction with a multi-sig account.
//
// Creates an account with a 3-key KeyList, schedules a transfer pre-signed
// by one of the three keys, then later adds a second signature via
// ScheduleSignTransaction. The schedule never executes (3-of-3 KeyList only
// has 2 sigs) — the example demonstrates the progressive-signing pattern,
// not full execution.
func main() {
	fmt.Println("Scheduled Transaction Multi-Sig Example Start!")

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

	// Step 1: generate three ED25519 private keys.
	fmt.Println("Generating ED25519 private keys...")
	privateKey1, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}
	privateKey2, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}
	privateKey3, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// Step 2: build a KeyList from the three public keys.
	// No threshold → default is "all keys required".
	fmt.Println("Creating a Key List...")
	keyList := hiero.NewKeyList().
		AddAllPublicKeys([]hiero.PublicKey{
			privateKey1.PublicKey(),
			privateKey2.PublicKey(),
			privateKey3.PublicKey(),
		})
	fmt.Printf("Created a Key List: %v\n", keyList)

	// Step 3: create a new account with the KeyList as its key.
	fmt.Println("Creating new account...")
	accountCreateResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(keyList).
		SetInitialBalance(hiero.NewHbar(2)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating account", err))
	}
	accountReceipt, err := accountCreateResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving account creation receipt", err))
	}
	accountID := *accountReceipt.AccountID
	fmt.Printf("Created new account with ID: %v\n", accountID)

	// Step 4: build a transfer, schedule it, pre-sign with privateKey2.
	fmt.Println("Creating a token transfer transaction...")
	transferTx := hiero.NewTransferTransaction().
		AddHbarTransfer(accountID, hiero.NewHbar(1).Negated()).
		AddHbarTransfer(operatorAccountID, hiero.NewHbar(1))

	fmt.Println("Scheduling the token transfer transaction...")
	scheduledTx, err := transferTx.Schedule()
	if err != nil {
		panic(fmt.Sprintf("%v : error creating scheduled transaction", err))
	}
	scheduledTx.
		SetPayerAccountID(operatorAccountID).
		SetAdminKey(operatorKey.PublicKey())

	frozenSchedule, err := scheduledTx.FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing schedule create transaction", err))
	}
	frozenSchedule.Sign(privateKey2)

	scheduleResponse, err := frozenSchedule.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule create transaction", err))
	}
	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving schedule create receipt", err))
	}
	scheduleID := *scheduleReceipt.ScheduleID
	fmt.Printf("Schedule ID: %v\n", scheduleID)

	// Step 5: query schedule info (should now have 1 signature from privateKey2).
	infoBefore, err := hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying schedule info", err))
	}
	fmt.Printf("Schedule info (before last signature): %v\n", infoBefore)

	// Step 6: send a ScheduleSignTransaction signed with privateKey3.
	// This is the 2nd of the 3 required signatures (privateKey1 is never
	// added — the schedule needs all 3 for the transfer to execute, so it
	// remains pending. That's the intended demonstration: progressive
	// signing without full execution.)
	fmt.Println("Appending private key #3 signature to a schedule transaction...")
	signTx, err := hiero.NewScheduleSignTransaction().
		SetScheduleID(scheduleID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing schedule sign transaction", err))
	}
	signTx.Sign(privateKey3)

	signResponse, err := signTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing schedule sign transaction", err))
	}
	signReceipt, err := signResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving schedule sign receipt", err))
	}
	fmt.Printf("A transaction that appends signature to a schedule transaction (private key #3) was complete with status: %v\n", signReceipt.Status)

	// Step 7: query schedule info again (should now have 2 signatures).
	infoAfter, err := hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying schedule info after sign", err))
	}
	fmt.Printf("Schedule info (after second signature): %v\n", infoAfter)

	// Cleanup: delete the account. Requires all 3 keys to sign.
	deleteFrozen, err := hiero.NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(operatorAccountID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account delete", err))
	}
	deleteFrozen.Sign(privateKey1)
	deleteFrozen.Sign(privateKey2)
	deleteFrozen.Sign(privateKey3)
	if _, err := deleteFrozen.Execute(client); err != nil {
		panic(fmt.Sprintf("%v : error executing account delete", err))
	}

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}
	fmt.Println("Scheduled Transaction Multi-Sig Example Complete!")
}
