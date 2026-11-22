//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
)

func FuzzHbarFromString(f *testing.F) {
	// From hbar_unit_test.go
	for _, seed := range []string{"+5.123 ℏ", "-5.123 Gℏ", "1.151 uℏ", "1.151.", "5"} {
		f.Add(seed)
	}

	// Every unit, and values around where String() switches from tinybar to hbar
	for _, seed := range []string{
		"1 tℏ", "1 μℏ", "1 mℏ", "1 ℏ", "1 kℏ", "1 Mℏ", "1 Gℏ",
		"10000 tℏ", "10001 tℏ", "-10001 tℏ", "0.00010001 ℏ", "0.29 ℏ", "0.1 ℏ",
		"92233720368.54775807 ℏ", "92233720368.54775808 ℏ", "99999999999999999999 Gℏ",
	} {
		f.Add(seed)
	}

	// Malformed
	for _, seed := range []string{"", " ", "ℏ", "1 ", "1  ℏ", "1 xℏ", ".5 ℏ", "5. ℏ", "1e5 ℏ", "--1 ℏ", "0x10 ℏ"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		hbar, err := HbarFromString(s)
		if err != nil {
			return
		}

		// String() must give a string HbarFromString parses back to the same amount
		reparsed, err := HbarFromString(hbar.String())
		if err != nil {
			t.Fatalf("HbarFromString(%q) parsed to %q, which does not parse again: %v", s, hbar.String(), err)
		}
		if reparsed != hbar {
			t.Fatalf("HbarFromString(%q) parsed to %d tinybar (%q), which parses again to %d tinybar",
				s, hbar.AsTinybar(), hbar.String(), reparsed.AsTinybar())
		}
	})
}
