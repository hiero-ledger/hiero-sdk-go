//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDecodeBytesBound(t *testing.T) {
	typ, _ := NewType("tuple(string)")
	budget := 0
	decodeTuple(typ, nil, &budget) // it should not panic
}

func TestDecodeDynamicLengthOutOfBounds(t *testing.T) {
	input := []byte("00000000000000000000000000000000\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00 00000000000000000000000000")
	typ, _ := NewType("tuple(bytes32, bytes, bytes)")

	_, err := Decode(typ, input)
	require.Error(t, err)
}

func TestDecodeArrayLongerThanInput(t *testing.T) {
	// uint8[3] needs 3 words but the input only has 2
	typ, err := NewType("uint8[3]")
	require.NoError(t, err)
	_, err = Decode(typ, make([]byte, 64))
	require.ErrorContains(t, err, "size is too big")

	// the slice has a length of 5 but only 2 words follow the length
	typ, err = NewType("uint8[]")
	require.NoError(t, err)
	input := make([]byte, 96)
	input[31] = 5
	_, err = Decode(typ, input)
	require.ErrorContains(t, err, "size is too big")
}

func TestDecodeLargeArrayFromShortInput(t *testing.T) {
	sliceOfLengthOne := make([]byte, 128)
	sliceOfLengthOne[31] = 1

	cases := []struct {
		typ   string
		input []byte
	}{
		{"uint8[100000000][1]", make([]byte, 64)},
		{"string[10000000][]", sliceOfLengthOne},
	}

	for _, c := range cases {
		typ, err := NewType(c.typ)
		require.NoError(t, err)

		// the decoded value would take hundreds of megabytes, it should fail before allocating it
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		_, err = Decode(typ, c.input)
		runtime.ReadMemStats(&after)

		require.Error(t, err)
		require.Less(t, after.TotalAlloc-before.TotalAlloc, uint64(1<<20), c.typ)
	}
}

// nestedArraysSharingData encodes levels nested dynamic arrays of n elements each, where all the offsets
// of one level point to the same array of the next level. The input has levels*(n+1) words, but
// decoding it without a limit gives n^levels values.
func nestedArraysSharingData(levels, n int) []byte {
	word := func(v int) []byte {
		b := make([]byte, 32)
		big.NewInt(int64(v)).FillBytes(b)
		return b
	}
	var input []byte
	for level := 0; level < levels; level++ {
		input = append(input, word(n)...)
		for i := 0; i < n; i++ {
			if level < levels-1 {
				// offsets count from the end of the length word, the next level starts after n words
				input = append(input, word(32*n)...)
			} else {
				input = append(input, word(1)...)
			}
		}
	}
	return input
}

func TestDecodeNestedArraysSharingData(t *testing.T) {
	typ, err := NewType("uint8[][][][][]")
	require.NoError(t, err)

	// 4,000 bytes that would decode into 24^5 = ~8 million values
	input := nestedArraysSharingData(5, 24)

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	start := time.Now()
	_, err = Decode(typ, input)
	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)

	require.ErrorIs(t, err, errABIDecodeBudgetExceeded)
	require.Less(t, after.TotalAlloc-before.TotalAlloc, uint64(1<<20))
	require.Less(t, elapsed, time.Second)
}

func TestDecodeNestedArraysWithinBudget(t *testing.T) {
	typ, err := NewType("uint8[][][][][]")
	require.NoError(t, err)

	// 2 elements per level share their data 2^5 = 32 times, far fewer values than the input has bytes
	_, err = Decode(typ, nestedArraysSharingData(5, 2))
	require.NoError(t, err)

	// a valid encoding of nested arrays decodes back to the same value
	value := [][][]uint8{{{1, 2}, {}}, {{3}}, {}}
	typ, err = NewType("uint8[][][]")
	require.NoError(t, err)
	input, err := Encode(value, typ)
	require.NoError(t, err)
	decoded, err := Decode(typ, input)
	require.NoError(t, err)
	require.Equal(t, value, decoded)
}

func TestDecodeDeeplyNestedStaticArrays(t *testing.T) {
	// one word holding a single value nested in 32 static arrays, the deepest nesting the budget allows
	typ, err := NewType("uint8" + strings.Repeat("[1]", 32))
	require.NoError(t, err)
	_, err = Decode(typ, make([]byte, 32))
	require.NoError(t, err)

	typ, err = NewType("uint8" + strings.Repeat("[1]", 33))
	require.NoError(t, err)
	_, err = Decode(typ, make([]byte, 32))
	require.ErrorIs(t, err, errABIDecodeBudgetExceeded)
}
