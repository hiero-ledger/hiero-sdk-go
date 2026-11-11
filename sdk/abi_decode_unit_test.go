//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeBytesBound(t *testing.T) {
	typ, _ := NewType("tuple(string)")
	decodeTuple(typ, nil) // it should not panic
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
