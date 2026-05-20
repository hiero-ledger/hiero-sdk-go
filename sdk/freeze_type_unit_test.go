//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitFreezeTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ft   FreezeType
		want string
	}{
		{name: "Unknown", ft: FreezeTypeUnknown, want: "UNKNOWN_FREEZE_TYPE"},
		{name: "FreezeOnly", ft: FreezeTypeFreezeOnly, want: "FREEZE_ONLY"},
		{name: "PrepareUpgrade", ft: FreezeTypePrepareUpgrade, want: "PREPARE_UPGRADE"},
		{name: "FreezeUpgrade", ft: FreezeTypeFreezeUpgrade, want: "FREEZE_UPGRADE"},
		{name: "FreezeAbort", ft: FreezeTypeFreezeAbort, want: "FREEZE_ABORT"},
		{name: "TelemetryUpgrade", ft: FreezeTypeTelemetryUpgrade, want: "TELEMETRY_UPGRADE"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.want, test.ft.String(), "FreezeType(%d).String()", test.ft)
		})
	}
}

func TestUnitFreezeTypeStringPanicsOnUnknownValue(t *testing.T) {
	t.Parallel()

	require.PanicsWithValue(
		t,
		"unreachable: FreezeType.String() switch statement is non-exhaustive. Status: 6",
		func() {
			_ = FreezeType(6).String()
		},
	)
}
