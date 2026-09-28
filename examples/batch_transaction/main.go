package main

import (
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func main() {
	fmt.Println("Batch Transaction Example Start!")

	client, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	operatorId, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorId, operatorKey)

	fmt.Println("Showcasing manual batch transaction preparation")
	executeBatchWithManualInnerTransactionFreeze(client)

	fmt.Println("\nShowcasing automatic batch transaction preparation using batchify")
	executeBatchWithBatchify(client)

	client.Close()
	fmt.Println("Batch Transaction Example Complete!")
}

func executeBatchWithManualInnerTransactionFreeze(client *hiero.Client) {
	// Step 1: Create batch keys
	batchKey1, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate batch key 1: %v", err))
	}
	batchKey2, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate batch key 2: %v", err))
	}
	batchKey3, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate batch key 3: %v", err))
	}

	// Step 2: Create 3 accounts and prepare transfers for batching
	fmt.Println("Creating three accounts and preparing batched transfers...")

	// Create Alice's account
	aliceKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate Alice's key: %v", err))
	}

	aliceCreateTx := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(aliceKey).
		SetInitialBalance(hiero.NewHbar(2))

	aliceCreateResp, err := aliceCreateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("failed to execute Alice's account creation: %v", err))
	}

	aliceReceipt, err := aliceCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("failed to get Alice's receipt: %v", err))
	}

	alice := aliceReceipt.AccountID
	if alice == nil {
		panic("failed to get Alice's account ID")
	}

	// Prepare Alice's transfer
	aliceBatchedTransfer, err := hiero.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(1)).
		AddHbarTransfer(*alice, hiero.NewHbar(-1)).
		SetTransactionID(hiero.TransactionIDGenerate(*alice)).
		SetBatchKey(batchKey1).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("failed to prepare Alice's transfer: %v", err))
	}
	aliceBatchedTransfer.Sign(aliceKey)

	fmt.Printf("Created first account (Alice): %v\n", alice.String())

	// Create Bob's account
	bobKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate Bob's key: %v", err))
	}

	bobCreateTx := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(bobKey).
		SetInitialBalance(hiero.NewHbar(2))

	bobCreateResp, err := bobCreateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("failed to execute Bob's account creation: %v", err))
	}

	bobReceipt, err := bobCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("failed to get Bob's receipt: %v", err))
	}

	bob := bobReceipt.AccountID
	if bob == nil {
		panic("failed to get Bob's account ID")
	}

	// Prepare Bob's transfer
	bobBatchedTransfer, err := hiero.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(1)).
		AddHbarTransfer(*bob, hiero.NewHbar(-1)).
		SetTransactionID(hiero.TransactionIDGenerate(*bob)).
		SetBatchKey(batchKey2).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("failed to prepare Bob's transfer: %v", err))
	}
	bobBatchedTransfer.Sign(bobKey)

	fmt.Printf("Created second account (Bob): %v\n", bob.String())

	// Create Carol's account
	carolKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate Carol's key: %v", err))
	}

	carolCreateTx := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(carolKey).
		SetInitialBalance(hiero.NewHbar(2))

	carolCreateResp, err := carolCreateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("failed to execute Carol's account creation: %v", err))
	}

	carolReceipt, err := carolCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("failed to get Carol's receipt: %v", err))
	}

	carol := carolReceipt.AccountID
	if carol == nil {
		panic("failed to get Carol's account ID")
	}

	// Prepare Carol's transfer
	carolBatchedTransfer, err := hiero.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(1)).
		AddHbarTransfer(*carol, hiero.NewHbar(-1)).
		SetTransactionID(hiero.TransactionIDGenerate(*carol)).
		SetBatchKey(batchKey3).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("failed to prepare Carol's transfer: %v", err))
	}
	carolBatchedTransfer.Sign(carolKey)

	fmt.Printf("Created third account (Carol): %v\n", carol.String())

	// Step 3: Get initial balances
	aliceBalanceBefore := hbarBalance(client, *alice, anyBalance)

	bobBalanceBefore := hbarBalance(client, *bob, anyBalance)

	carolBalanceBefore := hbarBalance(client, *carol, anyBalance)

	operatorBalanceBefore := hbarBalance(client, client.GetOperatorAccountID(), anyBalance)

	// Step 4: Execute the batch
	fmt.Println("Executing batch transaction...")
	batchTx, err := hiero.NewBatchTransaction().
		AddInnerTransaction(aliceBatchedTransfer).
		AddInnerTransaction(bobBatchedTransfer).
		AddInnerTransaction(carolBatchedTransfer).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("failed to prepare batch transaction: %v", err))
	}

	batchTx.Sign(batchKey1)
	batchTx.Sign(batchKey2)
	batchTx.Sign(batchKey3)

	batchResp, err := batchTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("failed to execute batch transaction: %v", err))
	}

	batchReceipt, err := batchResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("failed to get batch receipt: %v", err))
	}

	fmt.Printf("Batch transaction executed with status: %v\n", batchReceipt.Status)

	// Step 5: Verify new balances
	fmt.Println("Verifying the balances after batch execution...")
	aliceBalanceAfter := hbarBalance(client, *alice, func(balance hiero.Hbar) bool {
		return balance != aliceBalanceBefore
	})

	bobBalanceAfter := hbarBalance(client, *bob, anyBalance)

	carolBalanceAfter := hbarBalance(client, *carol, anyBalance)

	operatorBalanceAfter := hbarBalance(client, client.GetOperatorAccountID(), anyBalance)

	fmt.Printf("Alice's initial balance: %v, after: %v\n", aliceBalanceBefore, aliceBalanceAfter)
	fmt.Printf("Bob's initial balance: %v, after: %v\n", bobBalanceBefore, bobBalanceAfter)
	fmt.Printf("Carol's initial balance: %v, after: %v\n", carolBalanceBefore, carolBalanceAfter)
	fmt.Printf("Operator's initial balance: %v, after: %v\n", operatorBalanceBefore, operatorBalanceAfter)
}

func executeBatchWithBatchify(client *hiero.Client) {
	// Step 1: Create batch key
	batchKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate batch key: %v", err))
	}

	// Step 2: Create Alice's account
	fmt.Println("Creating account and preparing batched transfer...")
	aliceKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("failed to generate Alice's key: %v", err))
	}

	aliceCreateTx := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(aliceKey).
		SetInitialBalance(hiero.NewHbar(2))

	aliceCreateResp, err := aliceCreateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("failed to execute Alice's account creation: %v", err))
	}

	aliceReceipt, err := aliceCreateResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("failed to get Alice's receipt: %v", err))
	}

	alice := aliceReceipt.AccountID
	if alice == nil {
		panic("failed to get Alice's account ID")
	}

	fmt.Printf("Created Alice: %v\n", alice.String())

	// Step 3: Create client for Alice
	aliceClient, err := hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("failed to create Alice's client: %v", err))
	}
	aliceClient.SetOperator(*alice, aliceKey)

	// Step 4: Batchify a transfer transaction
	aliceBatchedTransfer, err := hiero.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(1)).
		AddHbarTransfer(*alice, hiero.NewHbar(-1)).
		Batchify(aliceClient, batchKey)
	if err != nil {
		panic(fmt.Sprintf("failed to batchify Alice's transfer: %v", err))
	}

	// Step 5: Get initial balances
	aliceBalanceBefore := hbarBalance(client, *alice, anyBalance)

	operatorBalanceBefore := hbarBalance(client, client.GetOperatorAccountID(), anyBalance)

	// Step 6: Execute the batch
	fmt.Println("Executing batch transaction...")
	batchTx, err := hiero.NewBatchTransaction().
		AddInnerTransaction(aliceBatchedTransfer).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("failed to prepare batch transaction: %v", err))
	}

	batchTx.Sign(batchKey)

	batchResp, err := batchTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("failed to execute batch transaction: %v", err))
	}

	batchReceipt, err := batchResp.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("failed to get batch receipt: %v", err))
	}

	fmt.Printf("Batch transaction executed with status: %v\n", batchReceipt.Status)

	// Step 7: Verify new balances
	fmt.Println("Verifying the balances after batch execution...")
	aliceBalanceAfter := hbarBalance(client, *alice, func(balance hiero.Hbar) bool {
		return balance != aliceBalanceBefore
	})

	operatorBalanceAfter := hbarBalance(client, client.GetOperatorAccountID(), anyBalance)

	fmt.Printf("Alice's initial balance: %v, after: %v\n", aliceBalanceBefore, aliceBalanceAfter)
	fmt.Printf("Operator's initial balance: %v, after: %v\n", operatorBalanceBefore, operatorBalanceAfter)
}

// hbarBalance polls the mirror node until ready accepts accountID's balance, absorbing mirror node lag.
func hbarBalance(client *hiero.Client, accountID hiero.AccountID, ready func(hiero.Hbar) bool) hiero.Hbar {
	for attempt := 1; ; attempt++ {
		balance, err := hiero.NewMirrorNodeAccountBalanceQuery().SetAccountID(accountID).Execute(client)
		if err == nil && ready(balance.Hbars) {
			return balance.Hbars
		}
		if attempt == 20 {
			panic(fmt.Sprintf("mirror node did not show the expected balance for %v (last error: %v)", accountID, err))
		}
		time.Sleep(time.Second)
	}
}

func anyBalance(hiero.Hbar) bool { return true }
