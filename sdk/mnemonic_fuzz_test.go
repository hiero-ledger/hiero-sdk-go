//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"
)

func FuzzMnemonicFromString(f *testing.F) {
	// 24 and 22 (legacy) words from crypto_unit_test.go, and BIP-39 test vectors with
	// 12 and 24 words
	for _, seed := range []string{
		testMnemonic,
		testMnemonic3,
		iosMnemonicString,
		androidMnemonicString,
		"jolly kidnap tom lawn drunk chick optic lust mutter mole bride galley dense member sage neural widow decide curb aboard margin manure",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		strings.Repeat("abandon ", 23) + "art",
	} {
		f.Add(seed)
	}

	// Malformed: wrong word count, unknown word, bad checksum, extra spaces
	for _, seed := range []string{
		"",
		" ",
		"abandon",
		strings.Repeat("abandon ", 11) + "abandon",
		strings.Repeat("abandon ", 11) + "zzzz",
		strings.Repeat("abandon ", 21) + "abandon",
		strings.Repeat("abandon  ", 11) + "about",
		strings.Repeat("ABANDON ", 11) + "ABOUT",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		mnemonic, err := MnemonicFromString(s)
		if err != nil {
			return
		}

		// What a caller does with a parsed mnemonic must not panic either. Errors are
		// fine: e.g. a legacy mnemonic has no standard derivation.
		_ = mnemonic.String()
		_ = mnemonic.Words()
		_, _ = mnemonic.ToLegacyPrivateKey()
		_, _ = mnemonic.ToStandardEd25519PrivateKey("", 0)
		_, _ = mnemonic.ToStandardECDSAsecp256k1PrivateKey("", 0)
	})
}
