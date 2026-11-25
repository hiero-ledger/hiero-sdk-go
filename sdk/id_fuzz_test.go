//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
)

func FuzzEntityIDFromString(f *testing.F) {
	// From the *_id_unit_test.go files
	for _, seed := range []string{
		"0.0.3", "0.0.123", "1.2.3", "0.0.5005",
		"0.0.123-rmkyk", "0.0.123-esxsf", "0.0.123-wrong", "0.0.100-quros",
		"0.0.0011223344556677889900112233445566778899",
		"0x0011223344556677889900112233445566778899",
		"0011223344556677889900112233445566778899",
	} {
		f.Add(seed)
	}

	// Aliases: a public key and an EVM address
	f.Add("0.0." + privateKeyED25519.PublicKey().String())
	f.Add("0.0." + privateKeyECDSA.PublicKey().String())
	f.Add("0.0." + privateKeyECDSA.PublicKey().ToEvmAddress())

	// Malformed
	for _, seed := range []string{
		"", ".", "..", "0.0", "0.0.", "0.0.-", "0.0.3-", "-0.0.3", "0.0.+3", "0.0.-3",
		"0.0.3.4", "a.b.c", "0.0.ab", "0x", "0.0.9223372036854775807", "0.0.9223372036854775808",
		"18446744073709551615.0.0",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		// The checksum is not part of String(), so the IDs are compared by their
		// String() instead of struct equality
		checkIDRoundTrip(t, "AccountIDFromString", s, AccountIDFromString)
		checkIDRoundTrip(t, "ContractIDFromString", s, ContractIDFromString)
		checkIDRoundTrip(t, "FileIDFromString", s, FileIDFromString)
		checkIDRoundTrip(t, "TokenIDFromString", s, TokenIDFromString)
		checkIDRoundTrip(t, "TopicIDFromString", s, TopicIDFromString)
		checkIDRoundTrip(t, "ScheduleIDFromString", s, ScheduleIDFromString)
	})
}

// checkIDRoundTrip checks that if parse accepts s, the String() of the ID it parsed to
// parses again to an ID with the same String().
func checkIDRoundTrip[T interface{ String() string }](t *testing.T, name string, s string, parse func(string) (T, error)) {
	t.Helper()

	id, err := parse(s)
	if err != nil {
		return
	}

	reparsed, err := parse(id.String())
	if err != nil {
		t.Fatalf("%s(%q) parsed to %q, which does not parse again: %v", name, s, id.String(), err)
	}
	if reparsed.String() != id.String() {
		t.Fatalf("%s(%q) parsed to %q, which parses again to %q", name, s, id.String(), reparsed.String())
	}
}
