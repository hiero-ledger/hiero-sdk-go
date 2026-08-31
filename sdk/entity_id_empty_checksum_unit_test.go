//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// A trailing "-" with an empty checksum (e.g. "0.0.123-") must be rejected, not
// silently parsed as an ID carrying an empty-string checksum. An empty checksum
// pointer bypasses the `checksum == nil` guard in ValidateChecksum and turns a
// "missing checksum" case into a misleading "wrong checksum" comparison.
func TestUnitEntityIDFromStringRejectsEmptyChecksum(t *testing.T) {
	t.Parallel()

	for _, s := range []string{"0.0.123-", "0.0.0-", "1.2.3-"} {
		_, err := AccountIDFromString(s)
		require.Errorf(t, err, "AccountIDFromString(%q) should be rejected", s)

		_, err = ContractIDFromString(s)
		require.Errorf(t, err, "ContractIDFromString(%q) should be rejected", s)

		_, err = TokenIDFromString(s)
		require.Errorf(t, err, "TokenIDFromString(%q) should be rejected", s)
	}

	// A well-formed checksummed ID still parses and preserves the checksum.
	id, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	require.NotNil(t, id.GetChecksum())
	require.Equal(t, "esxsf", *id.GetChecksum())
}
