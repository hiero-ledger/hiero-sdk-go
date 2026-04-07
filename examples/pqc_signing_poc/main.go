package main

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/cloudflare/circl/sign/dilithium/mode3"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// simulatedTransactionBody represents mock Hedera transaction body bytes.
// In production, these would come from a frozen protobuf-serialized transaction.
var simulatedTransactionBody = []byte(`{
	"transactionID": {"accountID": {"shardNum":0,"realmNum":0,"accountNum":1234}},
	"nodeAccountID": {"shardNum":0,"realmNum":0,"accountNum":3},
	"transactionFee": 100000000,
	"transactionValidDuration": {"seconds":120},
	"cryptoTransfer": {
		"transfers": {
			"accountAmounts": [
				{"accountID":{"accountNum":1234},"amount":-100000000},
				{"accountID":{"accountNum":5678},"amount":100000000}
			]
		}
	}
}`)

func main() {
	fmt.Println("=============================================================")
	fmt.Println("  Post-Quantum Cryptography (PQC) Signing POC")
	fmt.Println("  Hedera SDK for Go - CRYSTALS-Dilithium (Mode3 / ML-DSA-65)")
	fmt.Println("=============================================================")
	fmt.Println()

	msgLen := len(simulatedTransactionBody)
	fmt.Printf("Transaction body size: %d bytes\n\n", msgLen)

	// ---------------------------------------------------------------
	// 1. SDK Key Abstraction: Ed25519
	// ---------------------------------------------------------------
	fmt.Println("--- SDK PrivateKey: Ed25519 ---")
	edKey, err := hiero.PrivateKeyGenerateEd25519()
	if err != nil {
		panic(err)
	}
	edPub := edKey.PublicKey()

	start := time.Now()
	edSig := edKey.Sign(simulatedTransactionBody)
	edSignTime := time.Since(start)

	start = time.Now()
	edValid := edPub.VerifySignedMessage(simulatedTransactionBody, edSig)
	edVerifyTime := time.Since(start)

	fmt.Printf("  Public key size:  %d bytes\n", len(edPub.BytesRaw()))
	fmt.Printf("  Private key size: %d bytes\n", len(edKey.BytesRaw()))
	fmt.Printf("  Signature size:   %d bytes\n", len(edSig))
	fmt.Printf("  Sign time:        %v\n", edSignTime)
	fmt.Printf("  Verify time:      %v\n", edVerifyTime)
	fmt.Printf("  Valid:            %v\n", edValid)
	fmt.Printf("  PublicKey hex:    %s...\n\n", edPub.StringRaw()[:16])

	// ---------------------------------------------------------------
	// 2. SDK Key Abstraction: ECDSA
	// ---------------------------------------------------------------
	fmt.Println("--- SDK PrivateKey: ECDSA (secp256k1) ---")
	ecKey, err := hiero.PrivateKeyGenerateEcdsa()
	if err != nil {
		panic(err)
	}
	ecPub := ecKey.PublicKey()

	start = time.Now()
	ecSig := ecKey.Sign(simulatedTransactionBody)
	ecSignTime := time.Since(start)

	start = time.Now()
	ecValid := ecPub.VerifySignedMessage(simulatedTransactionBody, ecSig)
	ecVerifyTime := time.Since(start)

	fmt.Printf("  Public key size:  %d bytes\n", len(ecPub.BytesRaw()))
	fmt.Printf("  Private key size: %d bytes\n", len(ecKey.BytesRaw()))
	fmt.Printf("  Signature size:   %d bytes\n", len(ecSig))
	fmt.Printf("  Sign time:        %v\n", ecSignTime)
	fmt.Printf("  Verify time:      %v\n", ecVerifyTime)
	fmt.Printf("  Valid:            %v\n", ecValid)
	fmt.Printf("  PublicKey hex:    %s...\n\n", ecPub.StringRaw()[:16])

	// ---------------------------------------------------------------
	// 3. SDK Key Abstraction: Dilithium (PQC) — NEW
	// ---------------------------------------------------------------
	fmt.Println("--- SDK PrivateKey: Dilithium Mode3 / ML-DSA-65 [PQC] ---")
	dilKey, err := hiero.PrivateKeyGenerateDilithium()
	if err != nil {
		panic(err)
	}
	dilPub := dilKey.PublicKey()

	start = time.Now()
	dilSig := dilKey.Sign(simulatedTransactionBody)
	dilSignTime := time.Since(start)

	start = time.Now()
	dilValid := dilPub.VerifySignedMessage(simulatedTransactionBody, dilSig)
	dilVerifyTime := time.Since(start)

	fmt.Printf("  Public key size:  %d bytes\n", len(dilPub.BytesRaw()))
	fmt.Printf("  Private key size: %d bytes\n", len(dilKey.BytesRaw()))
	fmt.Printf("  Signature size:   %d bytes\n", len(dilSig))
	fmt.Printf("  Sign time:        %v\n", dilSignTime)
	fmt.Printf("  Verify time:      %v\n", dilVerifyTime)
	fmt.Printf("  Valid:            %v\n", dilValid)
	fmt.Printf("  PublicKey hex:    %s...\n\n", dilPub.StringRaw()[:16])

	// ---------------------------------------------------------------
	// Key serialization round-trip
	// ---------------------------------------------------------------
	fmt.Println("=============================================================")
	fmt.Println("  Key Serialization Round-Trip")
	fmt.Println("=============================================================")

	// Serialize and deserialize Dilithium keys via SDK
	dilPrivBytes := dilKey.BytesRaw()
	dilPubBytes := dilPub.BytesRaw()

	restoredPriv, err := hiero.PrivateKeyFromBytesDilithium(dilPrivBytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to restore private key: %v", err))
	}
	restoredPub, err := hiero.PublicKeyFromBytesDilithium(dilPubBytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to restore public key: %v", err))
	}

	// Sign with restored key, verify with restored public key
	restoredSig := restoredPriv.Sign(simulatedTransactionBody)
	restoredValid := restoredPub.VerifySignedMessage(simulatedTransactionBody, restoredSig)
	fmt.Printf("  Dilithium key round-trip: serialize → deserialize → sign → verify = %v\n\n", restoredValid)

	// ---------------------------------------------------------------
	// Summary comparison table
	// ---------------------------------------------------------------
	fmt.Println("=============================================================")
	fmt.Println("  Size Comparison (bytes)")
	fmt.Println("=============================================================")
	fmt.Printf("  %-20s %10s %10s %10s\n", "Algorithm", "PubKey", "PrivKey", "Signature")
	fmt.Printf("  %-20s %10d %10d %10d\n", "Ed25519", len(edPub.BytesRaw()), len(edKey.BytesRaw()), len(edSig))
	fmt.Printf("  %-20s %10d %10d %10d\n", "ECDSA secp256k1", len(ecPub.BytesRaw()), len(ecKey.BytesRaw()), len(ecSig))
	fmt.Printf("  %-20s %10d %10d %10d\n", "Dilithium3/ML-DSA-65", len(dilPub.BytesRaw()), len(dilKey.BytesRaw()), len(dilSig))
	fmt.Println()
	fmt.Println("=============================================================")
	fmt.Println("  Timing Comparison (single operation, via SDK PrivateKey)")
	fmt.Println("=============================================================")
	fmt.Printf("  %-20s %12s %12s\n", "Algorithm", "Sign", "Verify")
	fmt.Printf("  %-20s %12v %12v\n", "Ed25519", edSignTime, edVerifyTime)
	fmt.Printf("  %-20s %12v %12v\n", "ECDSA secp256k1", ecSignTime, ecVerifyTime)
	fmt.Printf("  %-20s %12v %12v\n", "Dilithium3/ML-DSA-65", dilSignTime, dilVerifyTime)
	fmt.Println()

	// ---------------------------------------------------------------
	// SDK integration demo
	// ---------------------------------------------------------------
	fmt.Println("=============================================================")
	fmt.Println("  SDK Integration Demonstration")
	fmt.Println("=============================================================")
	fmt.Println()
	fmt.Println("  All three key types use the same SDK API:")
	fmt.Println()
	fmt.Println("    key, _ := hiero.PrivateKeyGenerateDilithium()")
	fmt.Println("    pub := key.PublicKey()")
	fmt.Println("    sig := key.Sign(message)")
	fmt.Println("    valid := pub.VerifySignedMessage(message, sig)")
	fmt.Println()
	fmt.Println("  The Dilithium key type plugs into the same PrivateKey/PublicKey")
	fmt.Println("  abstraction as Ed25519 and ECDSA. It supports:")
	fmt.Println("    - PrivateKeyGenerateDilithium()")
	fmt.Println("    - PrivateKeyFromBytesDilithium(bytes)")
	fmt.Println("    - PublicKeyFromBytesDilithium(bytes)")
	fmt.Println("    - key.Sign(message)")
	fmt.Println("    - pub.VerifySignedMessage(message, sig)")
	fmt.Println("    - key.PublicKey()")
	fmt.Println("    - key.BytesRaw() / key.StringRaw()")
	fmt.Println("    - key.SignTransaction(tx)  [when network supports it]")
	fmt.Println()

	// ---------------------------------------------------------------
	// Key size impact analysis
	// ---------------------------------------------------------------
	fmt.Println("=============================================================")
	fmt.Println("  Key Size Impact on Hedera Transactions")
	fmt.Println("=============================================================")
	fmt.Printf("  Dilithium public key is ~%.0fx larger than Ed25519\n",
		float64(len(dilPub.BytesRaw()))/float64(ed25519.PublicKeySize))
	fmt.Printf("  Dilithium signature is ~%.0fx larger than Ed25519\n",
		float64(mode3.SignatureSize)/float64(ed25519.SignatureSize))
	fmt.Println()
	fmt.Println("  NOTE: The Hedera network does not yet accept PQC signatures.")
	fmt.Println("  Protobuf schema changes (Key and SignaturePair messages) are")
	fmt.Println("  needed before Dilithium keys can be used on-network.")
	fmt.Println()
}
