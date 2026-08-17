package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type AccountBalance struct {
	Hbars Hbar
	// Deprecated: Use `AccountBalance.Tokens` instead
	Token         map[TokenID]uint64
	Tokens        TokenBalanceMap
	TokenDecimals TokenDecimalMap
}

func _AccountBalanceFromProtobuf(pb *services.CryptoGetAccountBalanceResponse) AccountBalance { //nolint
	if pb == nil {
		return AccountBalance{}
	}
	var tokens map[TokenID]uint64
	if pb.TokenBalances != nil { //nolint
		tokens = make(map[TokenID]uint64, len(pb.TokenBalances)) //nolint
		for _, token := range pb.TokenBalances {                 //nolint
			if t := _TokenIDFromProtobuf(token.TokenId); t != nil {
				tokens[*t] = token.Balance
			}
		}
	}
	return AccountBalance{
		Hbars:         HbarFromTinybar(int64(pb.Balance)),
		Token:         tokens,
		Tokens:        _TokenBalanceMapFromProtobuf(pb.TokenBalances), //nolint
		TokenDecimals: _TokenDecimalMapFromProtobuf(pb.TokenBalances), //nolint
	}
}

// _AccountBalanceFromMirrorNode builds an AccountBalance from mirror node REST data. balances and
// decimals are keyed by token ID string for the public maps, and tokens carries the same balances
// keyed by TokenID so the deprecated Token map stays in sync as the protobuf path does.
//
// The caller passes both keyings because it has already parsed each token ID; re-deriving one from
// the other here would parse every ID a second time. The maps are adopted, not copied.
func _AccountBalanceFromMirrorNode(hbars Hbar, balances map[string]uint64, decimals map[string]uint64, tokens map[TokenID]uint64) AccountBalance {
	return AccountBalance{
		Hbars:         hbars,
		Token:         tokens,
		Tokens:        TokenBalanceMap{balances: balances},
		TokenDecimals: TokenDecimalMap{decimals: decimals},
	}
}

func (balance *AccountBalance) _ToProtobuf() *services.CryptoGetAccountBalanceResponse { //nolint
	return &services.CryptoGetAccountBalanceResponse{
		Balance: uint64(balance.Hbars.AsTinybar()),
	}
}
