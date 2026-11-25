//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitHbarFromTinybar(t *testing.T) {
	t.Parallel()

	tinybar := HbarUnits.Hbar._NumberOfTinybar()

	hbar := HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MaxHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)

	tinybar = MinHbar.tinybar

	hbar = HbarFromTinybar(tinybar)

	assert.Equal(t, tinybar, hbar.tinybar)
}

func TestUnitHbarUnit(t *testing.T) {
	t.Parallel()

	tinybar := HbarUnits.Kilobar._NumberOfTinybar()

	hbar := HbarFromTinybar(tinybar)

	hbar2, err := HbarFromString(hbar.ToString(HbarUnits.Kilobar))
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	tinybar = HbarUnits.Gigabar._NumberOfTinybar()

	hbar = HbarFromTinybar(tinybar)

	hbar2, err = HbarFromString(hbar.ToString(HbarUnits.Gigabar))
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	tinybar = HbarUnits.Microbar._NumberOfTinybar()

	hbar = HbarFromTinybar(tinybar)

	hbar2, err = HbarFromString(hbar.ToString(HbarUnits.Microbar))
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, hbar.tinybar)

	hbar2, err = HbarFromString("-5.123 Gℏ")
	require.NoError(t, err)
	assert.Equal(t, hbar2.tinybar, int64(-512300000000000000))

	hbar2, err = HbarFromString("5")
	require.NoError(t, err)
	assert.Equal(t, hbar2.ToString(HbarUnits.Hbar), "5 ℏ")

	hbar2, err = HbarFromString("+5.123 ℏ")
	require.NoError(t, err)
	assert.Equal(t, hbar2.ToString(HbarUnits.Millibar), "5123 mℏ")

	hbar2, err = HbarFromString("1.151 uℏ")
	assert.Error(t, err)

	hbar2, err = HbarFromString("1.151.")
	assert.Error(t, err)
}

func TestUnitHbarStringLargeAmountRoundTrip(t *testing.T) {
	t.Parallel()

	// String() used to print these with an exponent, which HbarFromString rejects,
	// and float64 cannot hold MaxHbar exactly
	for _, tc := range []struct {
		hbar     Hbar
		expected string
	}{
		{HbarFromTinybar(100_000_000_000_000), "1000000 ℏ"},
		{HbarFromTinybar(-512_300_000_000_000_000), "-5123000000 ℏ"},
		{MaxHbar, "92233720368.54775807 ℏ"},
		{MinHbar, "-92233720368.54775808 ℏ"},
	} {
		require.Equal(t, tc.expected, tc.hbar.String())

		parsed, err := HbarFromString(tc.hbar.String())
		require.NoError(t, err)
		assert.Equal(t, tc.hbar, parsed)
	}
}

func TestUnitHbarToStringLargeAmount(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "500000000 tℏ", NewHbar(5).ToString(HbarUnits.Tinybar))
	assert.Equal(t, "0.00000001 ℏ", HbarFromTinybar(1).ToString(HbarUnits.Hbar))
}

func TestUnitHbarFromStringExact(t *testing.T) {
	t.Parallel()

	// 0.29 * 1e8 is 28999999.999999996 as a float64
	hbar, err := HbarFromString("0.29 ℏ")
	require.NoError(t, err)
	assert.Equal(t, int64(29_000_000), hbar.AsTinybar())

	hbar, err = HbarFromString("0.00010001 ℏ")
	require.NoError(t, err)
	assert.Equal(t, int64(10_001), hbar.AsTinybar())

	// A fraction of a tinybar is still truncated toward zero
	hbar, err = HbarFromString("-1.5 tℏ")
	require.NoError(t, err)
	assert.Equal(t, int64(-1), hbar.AsTinybar())
}

func TestUnitHbarFromStringOutOfRange(t *testing.T) {
	t.Parallel()

	for _, s := range []string{"92233720368.54775808 ℏ", "-92233720368.54775809 ℏ", "99999999999999999999 Gℏ"} {
		_, err := HbarFromString(s)
		assert.Error(t, err, s)
	}
}
