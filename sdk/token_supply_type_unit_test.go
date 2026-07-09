package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnitTokenSupplyTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tst  TokenSupplyType
		want string
	}{
		{name: "Infinite", tst: TokenSupplyTypeInfinite, want: "TOKEN_SUPPLY_TYPE_INFINITE"},
		{name: "Finite", tst: TokenSupplyTypeFinite, want: "TOKEN_SUPPLY_TYPE_FINITE"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.tst.String(), "TokenSupplyType(%d).String()", test.tst)
		})
	}
}

func TestUnitTokenSupplyTypeStringPanic(t *testing.T) {
	t.Parallel()

	var unknownSupplyType TokenSupplyType = 100
	assert.PanicsWithValue(
		t,
		"unreachable: TokenSupplyType.String() switch statement is non-exhaustive. Status: 100",
		func() {
			_ = unknownSupplyType.String()
		},
	)
}
