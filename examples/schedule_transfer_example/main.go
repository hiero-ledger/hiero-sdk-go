package main

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to schedule a transfer transaction that requires a second party's
// signature before executing.
//
// Bob's account has receiverSignatureRequired = true. Alice (the operator)
// proposes a transfer to Bob via a scheduled transaction. Bob's balance is
// unchanged until he signs the schedule using ScheduleSignTransaction, at
// which point the transfer executes automatically. The 30-minute schedule
// expiration window applies if Bob does not sign in time.

func signerCount(info hiero.ScheduleInfo) int {
	if info.Signatories == nil {
		return 0
	}
	return len(info.Signatories.GetKeys())
}

func main() {
	fmt.Println("Scheduled Transfer Example Start!")

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

	// Step 1: Generate Bob's ECDSA key pair.
	fmt.Println("Generating Bob's ECDSA key pair...")
	bobsKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating Bob's key", err))
	}

	// Step 2: Create Bob's account with receiverSignatureRequired=true.
	// Because receiver-sig is required, the account creation transaction itself
	// must also be signed with Bob's key.
	fmt.Println("Creating Bob's account (receiverSignatureRequired=true)...")
	bobsAccountCreate, err := hiero.NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKeyWithoutAlias(bobsKey).
		SetInitialBalance(hiero.NewHbar(1)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing account creation", err))
	}
	bobsAccountCreate.Sign(bobsKey)

	response, err := bobsAccountCreate.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating Bob's account", err))
	}

	transactionReceipt, err := response.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}
	if transactionReceipt.AccountID == nil {
		panic("missing Bob's AccountID on receipt")
	}
	bobsID := *transactionReceipt.AccountID
	fmt.Printf("Bob's account: %v\n", bobsID)

	// Step 3: Read Bob's initial balance for the before/after comparison.
	bobsInitialBalance, err := hiero.NewAccountBalanceQuery().
		SetAccountID(bobsID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting Bob's balance", err))
	}
	fmt.Printf("Bob's balance before schedule: %v\n", bobsInitialBalance.Hbars)

	// Step 4: Alice builds the transfer and wraps it in a scheduled tx.
	// No Bob signature is added yet — the schedule will sit pending.
	fmt.Println("Alice scheduling a 1 ℏ transfer to Bob...")
	transferToSchedule := hiero.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.HbarFrom(-1, hiero.HbarUnits.Hbar)).
		AddHbarTransfer(bobsID, hiero.HbarFrom(1, hiero.HbarUnits.Hbar))

	scheduleTransaction, err := transferToSchedule.Schedule()
	if err != nil {
		panic(fmt.Sprintf("%v : error wrapping transfer in schedule", err))
	}
	scheduleTransaction.SetPayerAccountID(client.GetOperatorAccountID())

	scheduleResponse, err := scheduleTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing ScheduleCreate", err))
	}
	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting schedule create receipt", err))
	}
	if scheduleReceipt.ScheduleID == nil {
		panic("missing ScheduleID on receipt")
	}
	scheduleID := *scheduleReceipt.ScheduleID
	fmt.Printf("Schedule ID: %v\n", scheduleID)

	// Step 5: Confirm Bob's balance hasn't changed — the schedule is pending
	// because Bob's signature is still missing.
	balancePending, err := hiero.NewAccountBalanceQuery().
		SetAccountID(bobsID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting Bob's balance", err))
	}
	fmt.Printf("Bob's balance while schedule pending: %v\n", balancePending.Hbars)
	if balancePending.Hbars.AsTinybar() != bobsInitialBalance.Hbars.AsTinybar() {
		panic("expected Bob's balance to be unchanged while the schedule is pending")
	}

	// Step 6: Inspect the schedule. The inner scheduled transaction should be a
	// TransferTransaction.
	infoBefore, err := hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying schedule info", err))
	}
	scheduledTx, err := infoBefore.GetScheduledTransaction()
	if err != nil {
		panic(fmt.Sprintf("%v : error decoding scheduled transaction", err))
	}
	if _, ok := scheduledTx.(hiero.TransferTransaction); !ok {
		panic(fmt.Sprintf("expected scheduled tx to be TransferTransaction; got %T", scheduledTx))
	}
	fmt.Printf("ScheduleInfo before Bob signs: executed=%v signers=%d\n",
		infoBefore.ExecutedAt, signerCount(infoBefore))

	// Step 7: Bob signs the schedule via ScheduleSignTransaction. Once Bob's
	// signature is added, the schedule executes automatically (Alice's
	// signature was already attached at creation time).
	fmt.Println("Bob signing the schedule...")
	bobsSignTx, err := hiero.NewScheduleSignTransaction().
		SetScheduleID(scheduleID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing ScheduleSign", err))
	}
	bobsSignResponse, err := bobsSignTx.Sign(bobsKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing ScheduleSign", err))
	}
	if _, err := bobsSignResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error getting ScheduleSign receipt", err))
	}

	// Step 8: Confirm Bob's balance now reflects the transfer.
	balanceAfter, err := hiero.NewAccountBalanceQuery().
		SetAccountID(bobsID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting Bob's balance", err))
	}
	fmt.Printf("Bob's balance after Bob signs: %v\n", balanceAfter.Hbars)

	// Step 9: ScheduleInfo should now show an executed timestamp.
	infoAfter, err := hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error querying schedule info", err))
	}
	fmt.Printf("ScheduleInfo after Bob signs: executed=%v signers=%d\n",
		infoAfter.ExecutedAt, signerCount(infoAfter))
	if infoAfter.ExecutedAt == nil {
		panic("expected the schedule to have executed after Bob's signature")
	}

	// Cleanup: delete Bob's account.
	deleteAccount, err := hiero.NewAccountDeleteTransaction().
		SetAccountID(bobsID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing cleanup", err))
	}
	cleanupResponse, err := deleteAccount.Sign(bobsKey).Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing cleanup", err))
	}
	if _, err := cleanupResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error getting cleanup receipt", err))
	}

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}
	fmt.Println("Scheduled Transfer Example Complete!")
}
