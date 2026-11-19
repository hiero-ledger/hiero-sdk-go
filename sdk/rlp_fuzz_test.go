//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func FuzzRLPItemRead(f *testing.F) {
	// Valid vectors from rlp_unit_test.go
	f.Add([]byte{})
	f.Add([]byte{42})
	f.Add([]byte{0x80})
	f.Add([]byte{0x85, 'h', 'e', 'l', 'l', 'o'})
	f.Add([]byte{0x83, 0x01, 0x02, 0x03})
	f.Add([]byte{0xc0})
	f.Add([]byte{0xc8, 0x01, 0x83, 't', 'w', 'o', 0x82, 0x03, 0x04})
	f.Add([]byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'})

	// Long string and long list
	f.Add(append([]byte{0xb8, 0x38}, bytes.Repeat([]byte{0x01}, 0x38)...))
	f.Add(append([]byte{0xf8, 0x38}, bytes.Repeat([]byte{0x01}, 0x38)...))

	// Signed Ethereum transactions from ethereum_transaction_data_unit_test.go
	// and ethereum_data_unit_test.go, without the transaction type byte
	for _, s := range []string{
		"f86e82012a0585a54f4c3c00832dc6c094000000000000000000000000000000000000041d8502540be40080c001a017138522c841d5864818ab16f5d5dc7e353009c7a9ecb951e84a37991d14e47aa04bb17a3123f800f63e881dd0fa753bbd8b41af0093d85c3a35e66d9ba66b0bee",
		"f87082012a022f2f83018000947e3a9eaf9bcc39e2ffa38eb30bf7a93feacbc181880de0b6b3a764000083123456c001a0df48f2efd10421811de2bfb125ab75b2d3c44139c4642837fb1fccce911fd479a01aaf7ae92bee896651dfc9d99ae422a296bf5d9f1ca49b2d96d82b79eb112d66",
	} {
		seed, err := hex.DecodeString(s)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(seed)
	}

	// Truncated and malformed inputs
	f.Add([]byte{0x83, 0x01})
	f.Add([]byte{0xb8})
	f.Add([]byte{0xb9, 0xff, 0xff})
	f.Add([]byte{0xbf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	f.Add([]byte{0xc3, 0x01})
	f.Add([]byte{0xc1, 0x82, 0x01, 0x02})
	f.Add([]byte{0xf8})
	f.Add([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	f.Fuzz(func(t *testing.T, data []byte) {
		item := NewRLPItem(LIST_TYPE)
		if err := item.Read(data); err != nil {
			return
		}

		encoded, err := item.Write()
		if err != nil {
			t.Fatalf("failed to encode decoded item: %v", err)
		}

		decoded := NewRLPItem(LIST_TYPE)
		if err := decoded.Read(encoded); err != nil {
			t.Fatalf("failed to decode re-encoded item %x: %v", encoded, err)
		}

		if !rlpItemsEqual(item, decoded) {
			t.Fatalf("round trip mismatch for input %x, re-encoded as %x", data, encoded)
		}
	})
}

func rlpItemsEqual(a, b *RLPItem) bool {
	if a.itemType != b.itemType || !bytes.Equal(a.itemValue, b.itemValue) || len(a.childItems) != len(b.childItems) {
		return false
	}
	for i := range a.childItems {
		if !rlpItemsEqual(a.childItems[i], b.childItems[i]) {
			return false
		}
	}
	return true
}
