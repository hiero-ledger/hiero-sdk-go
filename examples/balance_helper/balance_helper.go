package balance_helper

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

const (
	// How long to wait for the mirror node to ingest state produced by a consensus transaction.
	waitTimeout = 30 * time.Second
	// How often to re-read the balance while waiting.
	waitInterval = 2 * time.Second
)

// WaitForMirrorBalance polls the account's balance until ready returns true, giving the mirror
// node time to ingest. On timeout the last result is returned, so the caller's error handling
// still applies.
func WaitForMirrorBalance(client *hiero.Client, accountID hiero.AccountID, ready func(hiero.AccountBalance) bool) (hiero.AccountBalance, error) {
	var balance hiero.AccountBalance
	var err error

	deadline := time.Now().Add(waitTimeout)
	for {
		balance, err = hiero.NewMirrorNodeAccountBalanceQuery().
			SetAccountID(accountID).
			Execute(client)
		if err == nil && ready(balance) {
			return balance, nil
		}
		if !time.Now().Before(deadline) {
			return balance, err
		}
		time.Sleep(waitInterval)
	}
}

// Anytime is a ready predicate for callers that only need the mirror node to answer at all,
// without waiting for a particular balance.
func Anytime(hiero.AccountBalance) bool { return true }
