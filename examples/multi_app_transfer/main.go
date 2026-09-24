package main

import (
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// How to transfer Hbar to an account with the receiver signature enabled.
//
// Demonstrates two-party signing across a byte serialization boundary:
// the user signs locally, sends the bytes to an "exchange" (simulated in
// this example by calling TransactionFromBytes / Sign / ToBytes), and the
// operator submits the doubly-signed transaction.
func main() {
	fmt.Println("MultiApp Transfer Example Start!")

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

	// Step 1: Generate ECDSA key pairs.
	// The exchange should possess this key — we generate it here for demo only.
	exchangePrivateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating exchange PrivateKey", err))
	}

	// The user's key — the only key we should actually possess in real life.
	userPrivateKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating user PrivateKey", err))
	}

	// Step 2: Create exchange and user accounts.
	fmt.Println("Creating exchange and receiver accounts...")

	exchangeCreateTx, err := hiero.NewAccountCreateTransaction().
		SetReceiverSignatureRequired(true).
		SetKeyWithoutAlias(exchangePrivateKey.PublicKey()).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing exchange account create transaction", err))
	}
	exchangeCreateTx.Sign(exchangePrivateKey)
	exchangeCreateResponse, err := exchangeCreateTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating exchange account", err))
	}
	exchangeReceipt, err := exchangeCreateResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving exchange account creation receipt", err))
	}
	exchangeAccountID := *exchangeReceipt.AccountID

	userCreateResponse, err := hiero.NewAccountCreateTransaction().
		SetInitialBalance(hiero.NewHbar(2)).
		SetKeyWithoutAlias(userPrivateKey.PublicKey()).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error creating user account", err))
	}
	userReceipt, err := userCreateResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving user account creation receipt", err))
	}
	userAccountID := *userReceipt.AccountID

	// Balances before
	senderBalanceBefore := hbarBalance(client, userAccountID, anyBalance)
	exchangeBalanceBefore := hbarBalance(client, exchangeAccountID, anyBalance)
	fmt.Printf("User account (%v) balance: %v\n", userAccountID, senderBalanceBefore)
	fmt.Printf("Exchange account (%v) balance: %v\n", exchangeAccountID, exchangeBalanceBefore)

	// Step 3: Build the transfer; sign with the user; serialize; "send to exchange".
	transferAmount := hiero.NewHbar(1)
	transferTx, err := hiero.NewTransferTransaction().
		AddHbarTransfer(userAccountID, transferAmount.Negated()).
		AddHbarTransfer(exchangeAccountID, transferAmount).
		// The exchange-provided memo required to validate the transaction.
		SetTransactionMemo("https://some-exchange.com/user1/account1").
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing transfer transaction", err))
	}
	transferTx.Sign(userPrivateKey)

	// The exchange must sign the transaction in order for it to be accepted by
	// the network. In real life this would be a REST call to the exchange API.
	userSignedBytes, err := transferTx.ToBytes()
	if err != nil {
		panic(fmt.Sprintf("%v : error serializing transfer transaction", err))
	}

	fmt.Println("Sending user-signed transaction bytes to exchange for countersignature...")
	exchangeSignedBytes := exchangeSigningService(userSignedBytes, exchangePrivateKey)
	fmt.Println("Exchange countersigned; received bytes back.")

	// Parse the bytes returned from the exchange and execute.
	finalTxIface, err := hiero.TransactionFromBytes(exchangeSignedBytes)
	if err != nil {
		panic(fmt.Sprintf("%v : error deserializing signed transfer transaction", err))
	}

	finalTx, ok := finalTxIface.(hiero.TransferTransaction)
	if !ok {
		panic("did not receive TransferTransaction back from signed bytes")
	}

	fmt.Printf("Transferring %v from the user account to the exchange account...\n", transferAmount)
	transferResponse, err := finalTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing transfer transaction", err))
	}
	if _, err := transferResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error retrieving transfer receipt", err))
	}

	// Step 4: Confirm balances after the transfer.
	exchangeBalanceAfter := hbarBalance(client, exchangeAccountID, func(balance hiero.Hbar) bool {
		return balance == hiero.HbarFromTinybar(exchangeBalanceBefore.AsTinybar()+transferAmount.AsTinybar())
	})
	senderBalanceAfter := hbarBalance(client, userAccountID, anyBalance)
	fmt.Printf("User account (%v) balance: %v\n", userAccountID, senderBalanceAfter)
	fmt.Printf("Exchange account (%v) balance: %v\n", exchangeAccountID, exchangeBalanceAfter)

	// Cleanup: delete both accounts, returning balances to the operator.
	exchangeDeleteTx, err := hiero.NewAccountDeleteTransaction().
		SetAccountID(exchangeAccountID).
		SetTransferAccountID(operatorAccountID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing exchange account delete", err))
	}
	exchangeDeleteTx.Sign(exchangePrivateKey)
	exchangeDeleteResponse, err := exchangeDeleteTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing exchange account delete", err))
	}
	if _, err := exchangeDeleteResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error retrieving exchange account delete receipt", err))
	}

	userDeleteTx, err := hiero.NewAccountDeleteTransaction().
		SetAccountID(userAccountID).
		SetTransferAccountID(operatorAccountID).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing user account delete", err))
	}
	userDeleteTx.Sign(userPrivateKey)
	userDeleteResponse, err := userDeleteTx.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing user account delete", err))
	}
	if _, err := userDeleteResponse.GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error retrieving user account delete receipt", err))
	}

	if err := client.Close(); err != nil {
		panic(fmt.Sprintf("%v : error closing client", err))
	}

	fmt.Println("MultiApp Transfer Example Complete!")
}

// exchangeSigningService simulates an offline exchange-side signing service:
// it deserializes the transaction bytes, signs with the exchange's key, and
// returns the re-serialized bytes.
func exchangeSigningService(txBytes []byte, exchangePrivateKey hiero.PrivateKey) []byte {
	tx, err := hiero.TransactionFromBytes(txBytes)
	if err != nil {
		panic(fmt.Sprintf("%v : exchange could not parse transaction", err))
	}

	transferTx, ok := tx.(hiero.TransferTransaction)
	if !ok {
		panic("exchange expected TransferTransaction in bytes")
	}

	signed := transferTx.Sign(exchangePrivateKey)

	signedBytes, err := signed.ToBytes()
	if err != nil {
		panic(fmt.Sprintf("%v : exchange could not serialize signed transaction", err))
	}
	return signedBytes
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
