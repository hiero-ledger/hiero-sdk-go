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

	// Generate new key to use with new account
	newKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	// HIP-1313: opt in to high-volume throttles by setting the high-volume flag.
	// During busy periods, dynamic pricing applies and the transaction fee may be
	// multiplied — SetMaxTransactionFee caps the price the operator is willing to pay.
	transactionResponse, err := hiero.NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetInitialBalance(hiero.NewHbar(1)).
		SetHighVolume(true).
		SetMaxTransactionFee(hiero.NewHbar(5)).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing high-volume account create", err))
	}

	receipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}

	fmt.Printf("account = %v\n", *receipt.AccountID)

	// The high-volume pricing multiplier is reported on the transaction record.
	// Value is divided by 1000 to get the actual multiplier (e.g. 1000 = 1.000x).
	record, err := transactionResponse.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting record", err))
	}

	fmt.Printf("transaction fee = %v\n", record.TransactionFee)
	fmt.Printf("high-volume pricing multiplier = %s\n", formatMultiplier(record.HighVolumePricingMultiplier))
}

// formatMultiplier renders the high-volume pricing multiplier from the record.
func formatMultiplier(m *uint64) string {
	if m == nil {
		return "(not set)"
	}
	return fmt.Sprintf("%.3fx", float64(*m)/1000)
}
