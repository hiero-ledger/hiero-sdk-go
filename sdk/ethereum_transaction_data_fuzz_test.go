//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

func FuzzEthereumTransactionDataFromBytes(f *testing.F) {
	// Signed transactions from ethereum_transaction_data_unit_test.go and ethereum_data_unit_test.go
	for _, s := range []string{
		"01f86e82012a0585a54f4c3c00832dc6c094000000000000000000000000000000000000041d8502540be40080c001a017138522c841d5864818ab16f5d5dc7e353009c7a9ecb951e84a37991d14e47aa04bb17a3123f800f63e881dd0fa753bbd8b41af0093d85c3a35e66d9ba66b0bee",
		"01f86c800685a54f4c3c00832dc6c094b19f7459342cf8c4fbd5053f5ef6c5981afa52d68502540be40080c080a07c66e585eda1a789ff33cfbd5d4b01c42b86ef64ef82097f20f1e700b7c57685a065a34302820466199cce5ed4e66e2bd59c93032ca5c4478fc9925c10972223b3",
		"02f87082012a022f2f83018000947e3a9eaf9bcc39e2ffa38eb30bf7a93feacbc181880de0b6b3a764000083123456c001a0df48f2efd10421811de2bfb125ab75b2d3c44139c4642837fb1fccce911fd479a01aaf7ae92bee896651dfc9d99ae422a296bf5d9f1ca49b2d96d82b79eb112d66",
	} {
		seed, err := hex.DecodeString(s)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(seed)
	}

	// Signed transactions of every type, with access and authorization lists, so the
	// fuzzer starts from inputs that reach every decoder
	to, _ := hex.DecodeString("000000000000000000000000000000000000041d")
	storageKey := bytes.Repeat([]byte{0x01}, 32)
	accessList := []AccessListItem{
		NewAccessListItem(to, nil),
		NewAccessListItem(to, [][]byte{storageKey}),
		NewAccessListItem(to, [][]byte{storageKey, storageKey}),
	}
	authorization := NewAuthorization([]byte{0x01, 0x2a}, to, 1, 1, bytes.Repeat([]byte{0xAA}, 32), bytes.Repeat([]byte{0xBB}, 32))

	for _, body := range []EthereumTransactionBody{
		(&EthereumLegacyTransaction{}).SetNonce(1).SetGasPrice(big.NewInt(10)).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(1)).SetCallData([]byte{0xDE, 0xAD}),
		(&EthereumEIP2930Transaction{}).SetChainId(298).SetNonce(1).SetGasPrice(big.NewInt(10)).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(1)).SetAccessListItems(accessList),
		(&EthereumEIP1559Transaction{}).SetChainId(298).SetNonce(1).SetMaxPriorityGas(big.NewInt(1)).SetMaxGas(big.NewInt(10)).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(1)).SetAccessListItems(accessList),
		(&EthereumEIP7702Transaction{}).SetChainId(298).SetNonce(1).SetMaxPriorityGas(big.NewInt(1)).SetMaxGas(big.NewInt(10)).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(1)).SetAccessListItems(accessList).AddAuthorization(authorization),
	} {
		seed, err := body.Sign(privateKeyECDSA)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(seed)
	}

	// Empty, truncated and malformed inputs
	f.Add([]byte{})
	f.Add([]byte{0x01})
	f.Add([]byte{0x02})
	f.Add([]byte{0x04})
	f.Add([]byte{0x01, 0xc0})
	f.Add([]byte{0x02, 0xc0})
	f.Add([]byte{0x04, 0xc0})
	f.Add([]byte{0xc0})
	f.Add([]byte{0x02, 0xf8})
	f.Add([]byte{0x02, 0xf8, 0x70, 0x82, 0x01})

	f.Fuzz(func(t *testing.T, data []byte) {
		txData, err := EthereumTransactionDataFromBytes(data)
		if err != nil {
			return
		}

		// What a caller does with a parsed transaction must not panic either
		_ = txData.GetData()
		body := txData.GetTransaction()
		if body == nil {
			t.Fatalf("decoded %x without error but holds no transaction", data)
		}
		_ = body.String()

		encoded, err := txData.ToBytes()
		if err != nil {
			t.Fatalf("failed to encode decoded transaction %x: %v", data, err)
		}

		// The input may be encoded differently (e.g. non-minimal RLP), but encoding
		// what was decoded must give the same bytes every time
		redecoded, err := EthereumTransactionDataFromBytes(encoded)
		if err != nil {
			t.Fatalf("failed to decode re-encoded transaction %x: %v", encoded, err)
		}
		reencoded, err := redecoded.ToBytes()
		if err != nil {
			t.Fatalf("failed to encode re-decoded transaction %x: %v", encoded, err)
		}
		if !bytes.Equal(encoded, reencoded) {
			t.Fatalf("round trip mismatch for input %x: encoded as %x, then as %x", data, encoded, reencoded)
		}
	})
}
