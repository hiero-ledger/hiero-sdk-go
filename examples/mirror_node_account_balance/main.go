package main

import (
	"fmt"
	"os"
	"time"

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

	// MirrorNodeAccountBalanceQuery replaces the deprecated AccountBalanceQuery, which the
	// consensus node stops serving in release 0.77. It reads from the mirror node REST API, so it
	// is free and needs no query payment.
	balance, err := hiero.NewMirrorNodeAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing mirror node account balance query", err))
	}

	fmt.Printf("balance = %v\n", balance.Hbars.String())

	// The account can also be addressed by anything the mirror node resolves: an EVM address, a
	// public key alias, or a contract ID. A contract goes through SetAccountID as well -- there is
	// no SetContractID, because the balances endpoint takes a contract through the same parameter.
	//
	//	contractAccountID := hiero.AccountID{
	//		Shard: contractID.Shard, Realm: contractID.Realm, Account: contractID.Contract,
	//	}

	// The mirror node ingests consensus state asynchronously and trails the network by a few
	// seconds, so a balance read immediately after a transfer may still show the old value. Poll
	// until the expected value appears rather than trusting a single read.
	transfer, err := hiero.NewTransferTransaction().
		AddHbarTransfer(operatorAccountID, hiero.NewHbar(-1)).
		AddHbarTransfer(hiero.AccountID{Account: 3}, hiero.NewHbar(1)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing transfer transaction", err))
	}
	if _, err = transfer.SetValidateStatus(true).GetReceipt(client); err != nil {
		panic(fmt.Sprintf("%v : error getting transfer receipt", err))
	}

	spent := balance.Hbars.AsTinybar()
	for attempt := 0; attempt < 10; attempt++ {
		if attempt > 0 {
			time.Sleep(2 * time.Second)
		}

		updated, err := hiero.NewMirrorNodeAccountBalanceQuery().
			SetAccountID(operatorAccountID).
			Execute(client)
		if err != nil {
			panic(fmt.Sprintf("%v : error executing mirror node account balance query", err))
		}

		if updated.Hbars.AsTinybar() < spent {
			fmt.Printf("balance after transfer = %v\n", updated.Hbars.String())
			return
		}
	}

	fmt.Println("mirror node had not yet ingested the transfer")
}
