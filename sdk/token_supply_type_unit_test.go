//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenSupplyTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		supplyType TokenSupplyType
		expected   string
	}{
		{TokenSupplyTypeInfinite, "TOKEN_SUPPLY_TYPE_INFINITE"},
		{TokenSupplyTypeFinite, "TOKEN_SUPPLY_TYPE_FINITE"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.supplyType.String())
	}
}
