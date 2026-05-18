//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitFreezeTypeString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		freezeType FreezeType
		expected   string
	}{
		{"unknown", FreezeTypeUnknown, "UNKNOWN_FREEZE_TYPE"},
		{"freeze only", FreezeTypeFreezeOnly, "FREEZE_ONLY"},
		{"prepare upgrade", FreezeTypePrepareUpgrade, "PREPARE_UPGRADE"},
		{"freeze upgrade", FreezeTypeFreezeUpgrade, "FREEZE_UPGRADE"},
		{"freeze abort", FreezeTypeFreezeAbort, "FREEZE_ABORT"},
		{"telemetry upgrade", FreezeTypeTelemetryUpgrade, "TELEMETRY_UPGRADE"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.freezeType.String())
		})
	}
}

func TestUnitFreezeTypeStringPanicsForUnknownValue(t *testing.T) {
	t.Parallel()

	defer func() {
		recovered := recover()
		require.NotNil(t, recovered)
		assert.Equal(
			t,
			"unreachable: FreezeType.String() switch statement is non-exhaustive. Status: 99",
			fmt.Sprint(recovered),
		)
	}()

	_ = FreezeType(99).String()
}
