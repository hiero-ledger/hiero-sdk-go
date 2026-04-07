package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

// msg is a representative transaction body payload for benchmarking.
var msg = []byte(`{"transactionID":{"accountID":{"shardNum":0,"realmNum":0,"accountNum":1234}},"nodeAccountID":{"shardNum":0,"realmNum":0,"accountNum":3},"transactionFee":100000000,"transactionValidDuration":{"seconds":120},"cryptoTransfer":{"transfers":{"accountAmounts":[{"accountID":{"accountNum":1234},"amount":-100000000},{"accountID":{"accountNum":5678},"amount":100000000}]}}}`)

// ===================================================================
// Key Generation Benchmarks
// ===================================================================

func BenchmarkKeyGen_Ed25519(b *testing.B) {
	for b.Loop() {
		_, _, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKeyGen_ECDSA_P256(b *testing.B) {
	for b.Loop() {
		_, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKeyGen_Dilithium3(b *testing.B) {
	for b.Loop() {
		_, _, err := mode3.GenerateKey(rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ===================================================================
// Signing Benchmarks
// ===================================================================

func BenchmarkSign_Ed25519(b *testing.B) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		ed25519.Sign(priv, msg)
	}
}

func BenchmarkSign_ECDSA_P256(b *testing.B) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	// ECDSA signs a hash, not raw bytes. Use first 32 bytes as hash proxy.
	hash := msg[:32]
	b.ResetTimer()
	for b.Loop() {
		_, err := ecdsa.SignASN1(rand.Reader, priv, hash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSign_Dilithium3(b *testing.B) {
	_, priv, err := mode3.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	sig := make([]byte, mode3.SignatureSize)
	b.ResetTimer()
	for b.Loop() {
		mode3.SignTo(priv, msg, sig)
	}
}

// ===================================================================
// Verification Benchmarks
// ===================================================================

func BenchmarkVerify_Ed25519(b *testing.B) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	sig := ed25519.Sign(priv, msg)
	b.ResetTimer()
	for b.Loop() {
		if !ed25519.Verify(pub, msg, sig) {
			b.Fatal("verification failed")
		}
	}
}

func BenchmarkVerify_ECDSA_P256(b *testing.B) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	hash := msg[:32]
	sig, err := ecdsa.SignASN1(rand.Reader, priv, hash)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if !ecdsa.VerifyASN1(&priv.PublicKey, hash, sig) {
			b.Fatal("verification failed")
		}
	}
}

func BenchmarkVerify_Dilithium3(b *testing.B) {
	pub, priv, err := mode3.GenerateKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	sig := make([]byte, mode3.SignatureSize)
	mode3.SignTo(priv, msg, sig)
	b.ResetTimer()
	for b.Loop() {
		if !mode3.Verify(pub, msg, sig) {
			b.Fatal("verification failed")
		}
	}
}
