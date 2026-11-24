//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// Keys from crypto_unit_test.go in encodings the fuzz seeds built from fixed keys do
// not cover
var keyFuzzVectors = []string{
	// Ed25519 private key as 64 raw bytes
	"db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7",
	// ECDSA private key, legacy DER and DER with the public key
	"3030020100300706052b8104000a04220420c78aa6584228422d2edb027dec01ff374eb0720e32be2da827d51ed0711d4f8c",
	"30540201010420ac318ea8ff8d991ab2f16172b4738e74dc35a56681199cfb1c0cb2e7cb560ffda00706052b8104000aa124032200036843f5cb338bbb4cdb21b0da4ea739d910951d6e8a5f703d313efe31afe788f4",
	// ECDSA public key, legacy DER and DER
	"302d300706052b8104000a0322000298c5d6efb814ead640467934b5ef9a02b81d3c483719675cb261cc5fde3edd57",
	"3036301006072a8648ce3d020106052b8104000a032200036843f5cb338bbb4cdb21b0da4ea739d910951d6e8a5f703d313efe31afe788f4",
	// ECDSA public key DER with the wrong curve and the wrong algorithm
	"3036301006072a8648ce3d020106052b8104000b03220002f3c22171a370bd5eba014c4ffa2002c4de4e0d8545ab964cda10467f23cfa2e0",
	"3036301006072a8648ce3d020206052b8104000a032200021ef6d7f710abb6f137f1c0dd87c6fd0200486c0754273d628cde6688e0bc9bf5",
	// Ed25519 public key DER with a key that is too long
	"3053300506032b6570032100b6ecf351b0d0d8fce400647f4944621b411e97ae83a249fcadb33cf165f40ca2",
	// Out of range ECDSA private keys: the curve order, and zero
	"fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141",
	"0000000000000000000000000000000000000000000000000000000000000000",
}

// keyFuzzStrings returns the given keys in all their string encodings, followed by
// keyFuzzVectors and some malformed strings.
func keyFuzzStrings(f *testing.F) []string {
	f.Helper()

	ed25519Key, err := PrivateKeyFromStringEd25519(mockPrivateKey)
	if err != nil {
		f.Fatal(err)
	}

	var seeds []string
	for _, key := range []PrivateKey{ed25519Key, privateKeyECDSA} {
		for _, s := range []string{
			key.String(), key.StringRaw(), key.StringDer(),
			key.PublicKey().String(), key.PublicKey().StringRaw(), key.PublicKey().StringDer(),
		} {
			seeds = append(seeds, s, "0x"+s, strings.ToUpper(s))
		}
	}
	seeds = append(seeds, keyFuzzVectors...)

	return append(seeds, "", "0x", "00", "0", "zz", "302e020100300506032b657004220420")
}

func FuzzPrivateKeyFromString(f *testing.F) {
	for _, seed := range keyFuzzStrings(f) {
		f.Add(seed)
	}

	parsers := map[string]func(string) (PrivateKey, error){
		"PrivateKeyFromString":        PrivateKeyFromString,
		"PrivateKeyFromStringDer":     PrivateKeyFromStringDer,
		"PrivateKeyFromStringEd25519": PrivateKeyFromStringEd25519,
		"PrivateKeyFromStringECDSA":   PrivateKeyFromStringECDSA,
	}

	f.Fuzz(func(t *testing.T, s string) {
		for name, parse := range parsers {
			key, err := parse(s)
			if err != nil {
				continue
			}

			exercisePrivateKey(t, key)

			// The input may be encoded differently, but the key it parsed to must
			// survive a round trip through its own string form
			reparsed, err := parse(key.String())
			if err != nil {
				t.Fatalf("%s(%q) parsed to %s, which does not parse again: %v", name, s, key.String(), err)
			}
			if reparsed.String() != key.String() {
				t.Fatalf("%s(%q) parsed to %s, which parses again to %s", name, s, key.String(), reparsed.String())
			}
		}
	})
}

func FuzzPublicKeyFromString(f *testing.F) {
	for _, seed := range keyFuzzStrings(f) {
		f.Add(seed)
	}

	parsers := map[string]func(string) (PublicKey, error){
		"PublicKeyFromString":        PublicKeyFromString,
		"PublicKeyFromStringEd25519": PublicKeyFromStringEd25519,
		"PublicKeyFromStringECDSA":   PublicKeyFromStringECDSA,
	}

	f.Fuzz(func(t *testing.T, s string) {
		for name, parse := range parsers {
			key, err := parse(s)
			if err != nil {
				continue
			}

			exercisePublicKey(key)

			reparsed, err := parse(key.String())
			if err != nil {
				t.Fatalf("%s(%q) parsed to %s, which does not parse again: %v", name, s, key.String(), err)
			}
			if reparsed.String() != key.String() {
				t.Fatalf("%s(%q) parsed to %s, which parses again to %s", name, s, key.String(), reparsed.String())
			}
		}
	})
}

func FuzzKeyFromBytes(f *testing.F) {
	for _, seed := range keyFuzzBytes(f) {
		f.Add(seed)
	}

	privateParsers := []func([]byte) (PrivateKey, error){
		PrivateKeyFromBytes, PrivateKeyFromBytesDer, PrivateKeyFromBytesEd25519, PrivateKeyFromBytesECDSA,
	}
	publicParsers := []func([]byte) (PublicKey, error){
		PublicKeyFromBytes, PublicKeyFromBytesDer, PublicKeyFromBytesEd25519, PublicKeyFromBytesECDSA,
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if key, err := KeyFromBytes(data); err == nil {
			_ = key.String()
			if _, err := KeyToBytes(key); err != nil {
				t.Fatalf("KeyFromBytes(%x) parsed to a key that does not serialize: %v", data, err)
			}
		}

		for _, parse := range privateParsers {
			if key, err := parse(data); err == nil {
				exercisePrivateKey(t, key)
			}
		}

		for _, parse := range publicParsers {
			if key, err := parse(data); err == nil {
				exercisePublicKey(key)
			}
		}
	})
}

// keyFuzzBytes returns the keys from keyFuzzStrings as bytes, and keys of every type
// serialized as protobuf Key messages for KeyFromBytes.
func keyFuzzBytes(f *testing.F) [][]byte {
	f.Helper()

	var seeds [][]byte
	for _, s := range keyFuzzStrings(f) {
		if b, err := PublicKeyFromString(s); err == nil {
			seeds = append(seeds, b.BytesRaw(), b.BytesDer())
		}
		if b, err := PrivateKeyFromString(s); err == nil {
			seeds = append(seeds, b.BytesRaw(), b.BytesDer())
		}
	}

	ed25519Key := privateKeyED25519.PublicKey()
	ecdsaKey := privateKeyECDSA.PublicKey()
	keyList := NewKeyList().AddAll([]Key{ed25519Key, ecdsaKey})
	threshold := KeyListWithThreshold(1).AddAll([]Key{ed25519Key, keyList})
	for _, key := range []Key{
		ed25519Key,
		ecdsaKey,
		keyList,
		threshold,
		ContractID{Contract: 5},
		DelegatableContractID{Contract: 5},
	} {
		b, err := KeyToBytes(key)
		if err != nil {
			f.Fatal(err)
		}
		seeds = append(seeds, b)
	}

	// A Key message with no key set, and one with an unknown field
	for _, msg := range []*services.Key{{}, {Key: &services.Key_RSA_3072{RSA_3072: []byte{1}}}} {
		b, err := protobuf.Marshal(msg)
		if err != nil {
			f.Fatal(err)
		}
		seeds = append(seeds, b)
	}

	return seeds
}

// exercisePrivateKey calls what a caller does with a parsed private key, which must
// not panic either, and checks that the key's signatures verify.
func exercisePrivateKey(t *testing.T, key PrivateKey) {
	t.Helper()

	_ = key.StringRaw()
	_ = key.StringDer()
	_ = key.BytesRaw()
	_ = key.BytesDer()

	message := []byte("message")
	publicKey := key.PublicKey()
	exercisePublicKey(publicKey)
	if !publicKey.VerifySignedMessage(message, key.Sign(message)) {
		t.Fatalf("signature of private key %s does not verify with its public key %s", key.String(), publicKey.String())
	}
}

// exercisePublicKey calls what a caller does with a parsed public key, which must
// not panic either.
func exercisePublicKey(key PublicKey) {
	_ = key.String()
	_ = key.StringRaw()
	_ = key.StringDer()
	_ = key.BytesRaw()
	_ = key.BytesDer()
	_ = key.ToAccountID(0, 0).String()
	if key.ecdsaPublicKey != nil {
		_ = key.ToEvmAddress()
	}
}
