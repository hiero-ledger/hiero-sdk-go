//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"
)

func FuzzABIDecode(f *testing.F) {
	// Valid encodings, built from values used in abi_encode_unit_test.go
	for _, c := range []struct {
		typ   string
		value interface{}
	}{
		{"bool", true},
		{"uint8", uint8(7)},
		{"uint40", big.NewInt(50)},
		{"int256", big.NewInt(-10)},
		{"address", Address{1}},
		{"bytes5", [5]byte{0x1, 0x2, 0x3, 0x4, 0x5}},
		{"bytes", []byte{0x12, 0x34, 0x56, 0x78}},
		{"string", "foobar"},
		{"uint8[]", []uint8{1, 2}},
		{"uint8[][2]", [2][]uint8{{1}, {1}}},
		{"string[2]", [2]string{"hello", "foobar"}},
		{"bytes32[][]", [][][32]uint8{{{1}, {2}}, {{3}, {4}, {5}}}},
		{"tuple(bytes[] a)", map[string]interface{}{"a": [][]byte{{0xf0, 0xf0, 0xf0}, {0xf0, 0xf0, 0xf0}}}},
		{"tuple(uint64[2] a)", map[string]interface{}{"a": [2]uint64{1, 2}}},
		{"tuple(address owner, uint256 balance, string name)", map[string]interface{}{
			"owner":   Address{1},
			"balance": big.NewInt(100),
			"name":    "hello",
		}},
	} {
		typ, err := NewType(c.typ)
		if err != nil {
			f.Fatal(err)
		}
		input, err := Encode(c.value, typ)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(c.typ, input)
	}

	// Type strings from abi_type_unit_test.go, valid and invalid
	for _, s := range []string{
		"int",
		"int32[]",
		"string[2][]",
		"tuple(int64 indexed arg0)",
		"tuple(int32 indexed arg0,tuple(int32 c) b_2)",
		"tuple()",
		"tuple((int32))",
		"int[[",
		"tuple[](a int32)",
		"int32[a]",
		"tuple(a int32",
	} {
		f.Add(s, make([]byte, 64))
	}

	// Regression inputs from abi_decode_unit_test.go
	f.Add("tuple(string)", []byte{})
	f.Add("tuple(bytes32, bytes, bytes)", []byte("00000000000000000000000000000000\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00 00000000000000000000000000"))

	f.Fuzz(func(t *testing.T, typeStr string, input []byte) {
		typ, err := NewType(typeStr)
		if err != nil {
			return
		}
		_, _ = Decode(typ, input)
	})
}
