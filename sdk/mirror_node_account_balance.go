package hiero

// SPDX-License-Identifier: Apache-2.0

// MirrorNodeAccountBalance is the hbar balance returned by MirrorNodeAccountBalanceQuery.
// Token balances are not included; the balances endpoint does not return them.
type MirrorNodeAccountBalance struct {
	Hbars Hbar
}
