package main

import (
	"fmt"
	"os"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// FeeEstimateQuery example (HIP-1261).
//
// Demonstrates how to estimate the fees for a transaction without submitting it
// to the network. The query talks to the mirror node REST endpoint
// POST /api/v1/network/fees and returns network, node, and service fee components.
func main() {
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

	fmt.Println("FeeEstimateQuery Example (HIP-1261)")

	// On a fresh solo deployment the mirror node's FeeEstimationService races
	// the importer ingesting the genesis fee schedule, so the first few calls
	// return HTTP 400 "Unknown transaction type". Poll until it answers or 10
	// minutes have elapsed.
	if err := waitForFeeEstimationService(client); err != nil {
		panic(fmt.Sprintf("%v : fee estimation service never became ready", err))
	}

	// Step 1: Create and freeze a transfer transaction. The query auto-freezes
	// if the transaction is not already frozen, but freezing up front lets us
	// reuse the same transaction across multiple estimates.
	transferTx, err := hiero.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(-1)).
		AddHbarTransfer(hiero.AccountID{Account: 3}, hiero.NewHbar(1)).
		FreezeWith(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error freezing transaction", err))
	}

	// Step 2: Estimate fees in INTRINSIC mode (default). Estimated from the
	// payload alone, ignoring state-dependent costs. Deterministic and fast.
	intrinsicEstimate, err := hiero.NewFeeEstimateQuery().
		SetMode(hiero.FeeEstimateModeIntrinsic).
		SetTransaction(transferTx).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing intrinsic fee estimate", err))
	}
	printEstimate("INTRINSIC", intrinsicEstimate)

	// Step 3: Estimate fees in STATE mode. Uses the mirror node's latest known
	// state to account for state-dependent costs (auto-creation, custom fees,
	// hooks, etc).
	stateEstimate, err := hiero.NewFeeEstimateQuery().
		SetMode(hiero.FeeEstimateModeState).
		SetTransaction(transferTx).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing state fee estimate", err))
	}
	printEstimate("STATE", stateEstimate)

	fmt.Printf("\nDifference (state - intrinsic): %d tinycents\n",
		int64(stateEstimate.Total)-int64(intrinsicEstimate.Total))

	// Step 4: Same query reached via the Transaction.EstimateFee() helper.
	// Optionally simulate high-volume pricing (HIP-1313) at 50% throttle
	// utilization (5000 basis points).
	helperEstimate, err := transferTx.EstimateFee().
		SetMode(hiero.FeeEstimateModeIntrinsic).
		SetHighVolumeThrottle(5000).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing fee estimate via helper", err))
	}
	printEstimate("INTRINSIC (high-volume 50%)", helperEstimate)
	fmt.Printf("  high_volume_multiplier: %d\n", helperEstimate.HighVolumeMultiplier)
}

func waitForFeeEstimationService(client *hiero.Client) error {
	probe := hiero.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(-1)).
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(1))

	deadline := time.Now().Add(10 * time.Minute)
	for time.Now().Before(deadline) {
		_, err := hiero.NewFeeEstimateQuery().
			SetMode(hiero.FeeEstimateModeIntrinsic).
			SetTransaction(probe).
			SetMaxAttempts(1).
			Execute(client)
		if err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("timed out after 10 minutes")
}

func printEstimate(label string, response hiero.FeeEstimateResponse) {
	fmt.Printf("\n%s estimate (tinycents):\n", label)
	fmt.Printf("  network: %d (multiplier=%d)\n", response.NetworkFee.Subtotal, response.NetworkFee.Multiplier)
	fmt.Printf("  node:    %d (base=%d, %d extras)\n", response.NodeFee.Subtotal(), response.NodeFee.Base, len(response.NodeFee.Extras))
	fmt.Printf("  service: %d (base=%d, %d extras)\n", response.ServiceFee.Subtotal(), response.ServiceFee.Base, len(response.ServiceFee.Extras))
	fmt.Printf("  total:   %d\n", response.Total)
}
