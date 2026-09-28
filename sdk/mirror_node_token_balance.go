package hiero

// SPDX-License-Identifier: Apache-2.0

// MirrorNodeTokenBalance is an account's relationship with one token, as returned by MirrorNodeTokenBalanceQuery.
type MirrorNodeTokenBalance struct {
	TokenID TokenID
	// Balance is in the token's smallest unit; for an NFT it is the number of serials held.
	Balance uint64
	// Decimals is 0 for an NFT.
	Decimals uint64
	// AutomaticAssociation is nil when not reported.
	AutomaticAssociation *bool
	// CreatedTimestamp is a Unix timestamp in seconds.nanoseconds format, or empty when not reported.
	CreatedTimestamp string
	// FreezeStatus is NOT_APPLICABLE, FROZEN or UNFROZEN, or empty when not reported.
	FreezeStatus string
	// KycStatus is NOT_APPLICABLE, GRANTED or REVOKED.
	KycStatus string
}

// MirrorNodeTokenBalancePage is one page of MirrorNodeTokenBalanceQuery results.
type MirrorNodeTokenBalancePage struct {
	Tokens []MirrorNodeTokenBalance
	// Next is the cursor for MirrorNodeTokenBalanceQuery.SetNextPage, or empty on the last page.
	Next string
}
